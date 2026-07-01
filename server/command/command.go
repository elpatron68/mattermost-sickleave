package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi"

	"github.com/medisoftware/mattermost-sickleave/server/dialog"
	"github.com/medisoftware/mattermost-sickleave/server/i18n"
	"github.com/medisoftware/mattermost-sickleave/server/sickleave"
)

type Settings struct {
	HRChannelID     string
	DefaultLocale   string
	MaxBackdateDays int
	ReportHashtag   string
	CommandTrigger  string
}

type SettingsProvider func() Settings

type Handler struct {
	client            *pluginapi.Client
	dialogAPI         dialog.DialogAPI
	store             sickleave.Store
	settings          SettingsProvider
	bundle            *i18n.Bundle
	pluginID          string
	siteURL           func() (string, error)
	botUserID         string
	registeredTrigger string
}

type Command interface {
	Handle(args *model.CommandArgs) (*model.CommandResponse, error)
	SubmitDialog(request *model.SubmitDialogRequest) (*model.SubmitDialogResponse, error)
	End(userID, channelID string) (*model.CommandResponse, error)
	EnsureSlashCommandRegistered() error
}

type HandlerConfig struct {
	Client    *pluginapi.Client
	DialogAPI dialog.DialogAPI
	Store     sickleave.Store
	Settings  SettingsProvider
	Bundle    *i18n.Bundle
	PluginID  string
	SiteURL   func() (string, error)
	BotUserID string
}

func NewCommandHandler(cfg HandlerConfig) Command {
	return &Handler{
		client:    cfg.Client,
		dialogAPI: cfg.DialogAPI,
		store:     cfg.Store,
		settings:  cfg.Settings,
		bundle:    cfg.Bundle,
		pluginID:  cfg.PluginID,
		siteURL:   cfg.SiteURL,
		botUserID: cfg.BotUserID,
	}
}

func (h *Handler) EnsureSlashCommandRegistered() error {
	trigger := h.commandTrigger()
	if trigger == h.registeredTrigger {
		return nil
	}

	if h.registeredTrigger != "" {
		if err := h.client.SlashCommand.Unregister("", h.registeredTrigger); err != nil {
			h.client.Log.Warn("Failed to unregister previous slash command", "trigger", h.registeredTrigger, "error", err)
		}
	}

	if err := h.client.SlashCommand.Register(&model.Command{
		Trigger:          trigger,
		AutoComplete:     true,
		AutoCompleteDesc: "Report and manage sick leave",
		AutoCompleteHint: "[start|update|extend|end|status|help]",
		AutocompleteData: buildAutocompleteData(trigger),
	}); err != nil {
		return err
	}

	h.registeredTrigger = trigger
	return nil
}

func (h *Handler) commandTrigger() string {
	return NormalizeCommandTrigger(h.settings().CommandTrigger)
}

func buildAutocompleteData(trigger string) *model.AutocompleteData {
	root := model.NewAutocompleteData(trigger, "[start|update|extend|end|status|help]", "Sick leave commands")
	root.AddCommand(model.NewAutocompleteData("start", "", "Report the first sick day"))
	root.AddCommand(model.NewAutocompleteData("update", "", "Provide expected return date and AU status"))
	root.AddCommand(model.NewAutocompleteData("extend", "", "Extend the expected return date"))
	root.AddCommand(model.NewAutocompleteData("end", "", "Close your active sick leave case"))
	root.AddCommand(model.NewAutocompleteData("status", "", "Show active sick leave"))
	root.AddCommand(model.NewAutocompleteData("help", "", "Show help"))
	return root
}

func (h *Handler) Handle(args *model.CommandArgs) (*model.CommandResponse, error) {
	locale := h.localeForUser(args.UserId)
	fields := strings.Fields(args.Command)

	subcommand := "help"
	if len(fields) > 1 {
		subcommand = strings.ToLower(fields[1])
	}

	switch subcommand {
	case "start":
		return h.handleStart(args, locale)
	case "update":
		return h.handleUpdate(args, locale)
	case "extend":
		return h.handleExtend(args, locale)
	case "end":
		return h.End(args.UserId, args.ChannelId)
	case "status":
		return h.handleStatus(args, locale), nil
	case "help":
		return h.handleHelp(locale), nil
	default:
		return ephemeral(h.bundle.T(locale, "command.error.unknown_subcommand", subcommand, h.commandTrigger())), nil
	}
}

func (h *Handler) handleStart(args *model.CommandArgs, locale string) (*model.CommandResponse, error) {
	settings := h.settings()
	if settings.HRChannelID == "" {
		return ephemeral(h.bundle.T(locale, "command.error.hr_channel_not_configured")), nil
	}

	active, err := h.store.GetActive(args.UserId)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return ephemeral(h.bundle.T(locale, "command.error.active_exists", h.commandTrigger())), nil
	}

	submitURL := h.dialogSubmitURL()

	maxBackdate := settings.MaxBackdateDays
	if maxBackdate <= 0 {
		maxBackdate = 3
	}

	if appErr := dialog.OpenStartDialog(h.dialogAPI, args.TriggerId, submitURL, locale, h.bundle, dialog.StartDialogOptions{
		Today:           time.Now().UTC(),
		MaxBackdateDays: maxBackdate,
	}); appErr != nil {
		return nil, fmt.Errorf("open dialog: %s", appErr.Error())
	}

	return &model.CommandResponse{}, nil
}

func (h *Handler) handleUpdate(args *model.CommandArgs, locale string) (*model.CommandResponse, error) {
	settings := h.settings()
	if settings.HRChannelID == "" {
		return ephemeral(h.bundle.T(locale, "command.error.hr_channel_not_configured")), nil
	}

	active, err := h.store.GetActive(args.UserId)
	if err != nil {
		return nil, err
	}
	if active == nil {
		return ephemeral(h.bundle.T(locale, "command.error.no_active")), nil
	}
	if active.Status != sickleave.StatusReported {
		return ephemeral(h.bundle.T(locale, "command.error.update_not_allowed", h.commandTrigger())), nil
	}

	submitURL := h.dialogSubmitURL()

	if appErr := dialog.OpenUpdateDialog(h.dialogAPI, args.TriggerId, submitURL, locale, h.bundle, dialog.UpdateDialogOptions{
		StartDate: active.StartDate,
	}); appErr != nil {
		return nil, fmt.Errorf("open dialog: %s", appErr.Error())
	}

	return &model.CommandResponse{}, nil
}

func (h *Handler) handleExtend(args *model.CommandArgs, locale string) (*model.CommandResponse, error) {
	settings := h.settings()
	if settings.HRChannelID == "" {
		return ephemeral(h.bundle.T(locale, "command.error.hr_channel_not_configured")), nil
	}

	active, err := h.store.GetActive(args.UserId)
	if err != nil {
		return nil, err
	}
	if active == nil {
		return ephemeral(h.bundle.T(locale, "command.error.no_active")), nil
	}
	if active.Status != sickleave.StatusUpdated && active.Status != sickleave.StatusExtended {
		return ephemeral(h.bundle.T(locale, "command.error.extend_not_allowed", h.commandTrigger())), nil
	}

	submitURL := h.dialogSubmitURL()

	if appErr := dialog.OpenExtendDialog(h.dialogAPI, args.TriggerId, submitURL, locale, h.bundle, dialog.ExtendDialogOptions{
		CurrentExpectedEnd: active.ExpectedEndDate,
	}); appErr != nil {
		return nil, fmt.Errorf("open dialog: %s", appErr.Error())
	}

	return &model.CommandResponse{}, nil
}

func (h *Handler) handleStatus(args *model.CommandArgs, locale string) *model.CommandResponse {
	record, err := h.store.GetActive(args.UserId)
	if err != nil {
		return ephemeral(err.Error())
	}
	if record == nil {
		return ephemeral(h.bundle.T(locale, "command.error.no_active"))
	}

	statusLabel := h.statusLabel(locale, record.Status)
	if record.ExpectedEndDate == "" {
		return ephemeral(h.bundle.T(locale, "command.status.active", record.StartDate, statusLabel))
	}

	auLabel := h.bundle.T(locale, "hr.post.au.unchanged")
	if record.AUCertificate != nil {
		if *record.AUCertificate {
			auLabel = h.bundle.T(locale, "hr.post.au.yes")
		} else {
			auLabel = h.bundle.T(locale, "hr.post.au.no")
		}
	}

	return ephemeral(h.bundle.T(
		locale,
		"command.status.active_with_end",
		record.StartDate,
		record.ExpectedEndDate,
		auLabel,
		statusLabel,
	))
}

func (h *Handler) handleHelp(locale string) *model.CommandResponse {
	trigger := h.commandTrigger()
	text := strings.Join([]string{
		h.bundle.T(locale, "command.help.header"),
		h.bundle.T(locale, "command.help.start", trigger),
		h.bundle.T(locale, "command.help.update", trigger),
		h.bundle.T(locale, "command.help.extend", trigger),
		h.bundle.T(locale, "command.help.end", trigger),
		h.bundle.T(locale, "command.help.status", trigger),
		h.bundle.T(locale, "command.help.help", trigger),
	}, "\n")
	return ephemeral(text)
}

func (h *Handler) localeForUser(userID string) string {
	user, err := h.client.User.Get(userID)
	if err != nil {
		h.client.Log.Warn("Failed to load user for locale", "user_id", userID, "error", err)
		return i18n.LocaleForUser(nil, h.settings().DefaultLocale)
	}
	return i18n.LocaleForUser(user, h.settings().DefaultLocale)
}

func (h *Handler) dialogSubmitURL() string {
	// Use a relative plugin path so Mattermost routes dialog submissions locally
	// (DoLocalRequest) instead of making an outbound HTTP call to SiteURL.
	return fmt.Sprintf("/plugins/%s/api/v1/dialog/submit", h.pluginID)
}

func (h *Handler) statusLabel(locale string, status sickleave.Status) string {
	switch status {
	case sickleave.StatusUpdated:
		return h.bundle.T(locale, "command.status.updated")
	case sickleave.StatusExtended:
		return h.bundle.T(locale, "command.status.extended")
	case sickleave.StatusClosed:
		return h.bundle.T(locale, "command.status.closed")
	default:
		return h.bundle.T(locale, "command.status.reported")
	}
}

func ephemeral(text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         text,
	}
}

func (h *Handler) SubmitDialog(request *model.SubmitDialogRequest) (*model.SubmitDialogResponse, error) {
	switch request.CallbackId {
	case dialog.CallbackStart:
		return h.SubmitStartDialog(request)
	case dialog.CallbackUpdate:
		return h.SubmitUpdateDialog(request)
	case dialog.CallbackExtend:
		return h.SubmitExtendDialog(request)
	default:
		return &model.SubmitDialogResponse{
			Errors: map[string]string{
				"": "unknown dialog callback",
			},
		}, nil
	}
}

func (h *Handler) SubmitStartDialog(request *model.SubmitDialogRequest) (*model.SubmitDialogResponse, error) {
	locale := h.localeForUser(request.UserId)
	settings := h.settings()

	if settings.HRChannelID == "" {
		return dialogError(locale, h.bundle, "command.error.hr_channel_not_configured"), nil
	}

	active, err := h.store.GetActive(request.UserId)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return dialogError(locale, h.bundle, "command.error.active_exists", h.commandTrigger()), nil
	}

	rawStartDate, ok := request.Submission["start_date"].(string)
	if !ok || strings.TrimSpace(rawStartDate) == "" {
		return dialogFieldError(locale, h.bundle, "start_date", "dialog.error.start_date_required"), nil
	}

	startDate, err := sickleave.ParseDate(rawStartDate)
	if err != nil {
		return dialogFieldError(locale, h.bundle, "start_date", "dialog.error.start_date_invalid"), nil
	}

	maxBackdate := settings.MaxBackdateDays
	if maxBackdate <= 0 {
		maxBackdate = 3
	}

	if err = sickleave.ValidateStartDate(startDate, time.Now().UTC(), maxBackdate); err != nil {
		switch err.Error() {
		case "start date is in the future":
			return dialogFieldError(locale, h.bundle, "start_date", "dialog.error.start_date_future"), nil
		case "start date is too far in the past":
			return dialogFieldError(locale, h.bundle, "start_date", "dialog.error.start_date_backdate", maxBackdate), nil
		default:
			return dialogFieldError(locale, h.bundle, "start_date", "dialog.error.start_date_invalid"), nil
		}
	}

	user, err := h.client.User.Get(request.UserId)
	if err != nil {
		return nil, err
	}

	record := &sickleave.Record{
		ID:          uuid.NewString(),
		UserID:      request.UserId,
		TeamID:      request.TeamId,
		StartDate:   rawStartDate,
		Status:      sickleave.StatusReported,
		HRChannelID: settings.HRChannelID,
		Hashtag:     sickleave.NormalizeHashtag(settings.ReportHashtag),
		History: []sickleave.HistoryEntry{{
			Variant:   "A",
			Timestamp: time.Now().UTC(),
			Data: map[string]string{
				"start_date": rawStartDate,
			},
		}},
	}

	post := &model.Post{
		UserId:    h.botUserID,
		ChannelId: settings.HRChannelID,
		Message:   sickleave.FormatInitialHRPost(record, user, locale, h.bundle),
	}

	if err := h.client.Post.CreatePost(post); err != nil {
		return nil, err
	}

	record.HRPostID = post.Id
	if err := h.persistRecord(request.UserId, record); err != nil {
		return nil, err
	}

	h.sendEphemeralSuccess(request, locale, "dialog.success.start")
	return &model.SubmitDialogResponse{Errors: map[string]string{}}, nil
}

func (h *Handler) SubmitUpdateDialog(request *model.SubmitDialogRequest) (*model.SubmitDialogResponse, error) {
	locale := h.localeForUser(request.UserId)
	settings := h.settings()

	if settings.HRChannelID == "" {
		return dialogError(locale, h.bundle, "command.error.hr_channel_not_configured"), nil
	}

	active, err := h.store.GetActive(request.UserId)
	if err != nil {
		return nil, err
	}
	if active == nil {
		return dialogError(locale, h.bundle, "command.error.no_active"), nil
	}
	if active.Status != sickleave.StatusReported {
		return dialogError(locale, h.bundle, "command.error.update_not_allowed", h.commandTrigger()), nil
	}

	rawExpectedEnd, ok := request.Submission["expected_end_date"].(string)
	if !ok || strings.TrimSpace(rawExpectedEnd) == "" {
		return dialogFieldError(locale, h.bundle, "expected_end_date", "dialog.error.expected_end_required"), nil
	}

	expectedEnd, err := sickleave.ParseDate(rawExpectedEnd)
	if err != nil {
		return dialogFieldError(locale, h.bundle, "expected_end_date", "dialog.error.expected_end_invalid"), nil
	}

	startDate, err := sickleave.ParseDate(active.StartDate)
	if err != nil {
		return nil, fmt.Errorf("parse stored start date: %w", err)
	}

	if err := sickleave.ValidateExpectedEndDate(startDate, expectedEnd); err != nil {
		return dialogFieldError(locale, h.bundle, "expected_end_date", "dialog.error.expected_end_before_start"), nil
	}

	auValue, ok := request.Submission["au_certificate"].(string)
	if !ok || strings.TrimSpace(auValue) == "" {
		return dialogFieldError(locale, h.bundle, "au_certificate", "dialog.error.au_certificate_required"), nil
	}

	auCertificate, valid := parseAUCertificate(auValue)
	if !valid {
		return dialogFieldError(locale, h.bundle, "au_certificate", "dialog.error.au_certificate_invalid"), nil
	}

	active.ExpectedEndDate = rawExpectedEnd
	active.AUCertificate = &auCertificate
	active.Status = sickleave.StatusUpdated
	active.History = append(active.History, sickleave.HistoryEntry{
		Variant:   "B",
		Timestamp: time.Now().UTC(),
		Data: map[string]any{
			"expected_end_date": rawExpectedEnd,
			"au_certificate":    auCertificate,
		},
	})

	if err := h.postHRThreadReply(active, sickleave.FormatUpdateHRPost(active, rawExpectedEnd, auCertificate, locale, h.bundle)); err != nil {
		return nil, err
	}

	if err := h.persistRecord(request.UserId, active); err != nil {
		return nil, err
	}

	h.sendEphemeralSuccess(request, locale, "dialog.success.update")
	return &model.SubmitDialogResponse{Errors: map[string]string{}}, nil
}

func (h *Handler) SubmitExtendDialog(request *model.SubmitDialogRequest) (*model.SubmitDialogResponse, error) {
	locale := h.localeForUser(request.UserId)
	settings := h.settings()

	if settings.HRChannelID == "" {
		return dialogError(locale, h.bundle, "command.error.hr_channel_not_configured"), nil
	}

	active, err := h.store.GetActive(request.UserId)
	if err != nil {
		return nil, err
	}
	if active == nil {
		return dialogError(locale, h.bundle, "command.error.no_active"), nil
	}
	if active.Status != sickleave.StatusUpdated && active.Status != sickleave.StatusExtended {
		return dialogError(locale, h.bundle, "command.error.extend_not_allowed", h.commandTrigger()), nil
	}
	if active.ExpectedEndDate == "" {
		return dialogError(locale, h.bundle, "command.error.extend_not_allowed", h.commandTrigger()), nil
	}

	rawExpectedEnd, ok := request.Submission["expected_end_date"].(string)
	if !ok || strings.TrimSpace(rawExpectedEnd) == "" {
		return dialogFieldError(locale, h.bundle, "expected_end_date", "dialog.error.expected_end_required"), nil
	}

	newExpectedEnd, err := sickleave.ParseDate(rawExpectedEnd)
	if err != nil {
		return dialogFieldError(locale, h.bundle, "expected_end_date", "dialog.error.expected_end_invalid"), nil
	}

	currentExpectedEnd, err := sickleave.ParseDate(active.ExpectedEndDate)
	if err != nil {
		return nil, fmt.Errorf("parse stored expected end date: %w", err)
	}

	if err := sickleave.ValidateExtensionEndDate(currentExpectedEnd, newExpectedEnd); err != nil {
		return dialogFieldError(locale, h.bundle, "expected_end_date", "dialog.error.expected_end_not_after_current"), nil
	}

	var auUpdate *bool
	if auValue, ok := request.Submission["au_certificate"].(string); ok && strings.TrimSpace(auValue) != "" && auValue != "unchanged" {
		auCertificate, valid := parseAUCertificate(auValue)
		if !valid {
			return dialogFieldError(locale, h.bundle, "au_certificate", "dialog.error.au_certificate_invalid"), nil
		}
		auUpdate = &auCertificate
		active.AUCertificate = auUpdate
	}

	active.ExpectedEndDate = rawExpectedEnd
	active.Status = sickleave.StatusExtended
	historyData := map[string]any{
		"expected_end_date": rawExpectedEnd,
	}
	if auUpdate != nil {
		historyData["au_certificate"] = *auUpdate
	}
	active.History = append(active.History, sickleave.HistoryEntry{
		Variant:   "C",
		Timestamp: time.Now().UTC(),
		Data:      historyData,
	})

	if err := h.postHRThreadReply(active, sickleave.FormatExtendHRPost(active, rawExpectedEnd, auUpdate, locale, h.bundle)); err != nil {
		return nil, err
	}

	if err := h.persistRecord(request.UserId, active); err != nil {
		return nil, err
	}

	h.sendEphemeralSuccess(request, locale, "dialog.success.extend")
	return &model.SubmitDialogResponse{Errors: map[string]string{}}, nil
}

func (h *Handler) End(userID, channelID string) (*model.CommandResponse, error) {
	locale := h.localeForUser(userID)
	return h.closeActiveCase(userID, locale)
}

func (h *Handler) closeActiveCase(userID, locale string) (*model.CommandResponse, error) {
	active, err := h.store.GetActive(userID)
	if err != nil {
		return nil, err
	}
	if active == nil {
		return ephemeral(h.bundle.T(locale, "command.error.no_active")), nil
	}

	active.Status = sickleave.StatusClosed
	active.History = append(active.History, sickleave.HistoryEntry{
		Variant:   "end",
		Timestamp: time.Now().UTC(),
		Data:      map[string]string{},
	})

	if err := h.postHRThreadReply(active, sickleave.FormatCloseHRPost(active, locale, h.bundle)); err != nil {
		return nil, err
	}

	if err := h.store.SaveRecord(active); err != nil {
		return nil, err
	}
	if err := h.store.ClearActive(userID); err != nil {
		return nil, err
	}

	return ephemeral(h.bundle.T(locale, "command.success.end")), nil
}

func (h *Handler) persistRecord(userID string, record *sickleave.Record) error {
	if err := h.store.SetActive(userID, record); err != nil {
		return err
	}
	return h.store.SaveRecord(record)
}

func (h *Handler) postHRThreadReply(record *sickleave.Record, message string) error {
	post := &model.Post{
		UserId:    h.botUserID,
		ChannelId: record.HRChannelID,
		RootId:    record.HRPostID,
		Message:   message,
	}
	return h.client.Post.CreatePost(post)
}

func (h *Handler) sendEphemeralSuccess(request *model.SubmitDialogRequest, locale, key string) {
	ephemeralPost := &model.Post{
		UserId:    h.botUserID,
		ChannelId: request.ChannelId,
		Message:   h.bundle.T(locale, key),
	}
	h.client.Post.SendEphemeralPost(request.UserId, ephemeralPost)
}

func parseAUCertificate(value string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "yes":
		return true, true
	case "no":
		return false, true
	default:
		return false, false
	}
}

func dialogError(locale string, bundle *i18n.Bundle, key string, args ...any) *model.SubmitDialogResponse {
	return &model.SubmitDialogResponse{
		Errors: map[string]string{
			"": bundle.T(locale, key, args...),
		},
	}
}

func dialogFieldError(locale string, bundle *i18n.Bundle, field, key string, args ...any) *model.SubmitDialogResponse {
	return &model.SubmitDialogResponse{
		Errors: map[string]string{
			field: bundle.T(locale, key, args...),
		},
	}
}
