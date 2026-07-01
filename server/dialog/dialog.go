package dialog

import (
	"time"

	"github.com/mattermost/mattermost/server/public/model"

	"github.com/medisoftware/mattermost-sickleave/server/i18n"
)

const (
	CallbackStart  = "sickleave_start"
	CallbackUpdate = "sickleave_update"
	CallbackExtend = "sickleave_extend"
	StateVariantA  = "variant=A"
	StateVariantB  = "variant=B"
	StateVariantC  = "variant=C"
	dateLayout     = "2006-01-02"
)

type DialogAPI interface {
	OpenInteractiveDialog(request model.OpenDialogRequest) *model.AppError
}

type StartDialogOptions struct {
	Today           time.Time
	MaxBackdateDays int
}

type UpdateDialogOptions struct {
	StartDate string
}

type ExtendDialogOptions struct {
	CurrentExpectedEnd string
}

func formatDate(value time.Time) string {
	return value.UTC().Format(dateLayout)
}

func BuildStartDialog(locale string, bundle *i18n.Bundle, opts StartDialogOptions) model.Dialog {
	today := truncateToDate(opts.Today.UTC())
	maxBackdate := opts.MaxBackdateDays
	if maxBackdate <= 0 {
		maxBackdate = 3
	}
	minDate := today.AddDate(0, 0, -maxBackdate)

	return model.Dialog{
		CallbackId:       CallbackStart,
		Title:            bundle.T(locale, "dialog.a.title"),
		IntroductionText: bundle.T(locale, "dialog.a.intro"),
		SubmitLabel:      bundle.T(locale, "dialog.submit"),
		NotifyOnCancel:   false,
		State:            StateVariantA,
		Elements: []model.DialogElement{{
			DisplayName: bundle.T(locale, "dialog.a.start_date"),
			Name:        "start_date",
			Type:        "date",
			Placeholder: bundle.T(locale, "dialog.date_placeholder"),
			HelpText:    bundle.T(locale, "dialog.a.start_date_help"),
			MinDate:     formatDate(minDate),
			MaxDate:     formatDate(today),
		}},
	}
}

func BuildUpdateDialog(locale string, bundle *i18n.Bundle, opts UpdateDialogOptions) model.Dialog {
	return model.Dialog{
		CallbackId:       CallbackUpdate,
		Title:            bundle.T(locale, "dialog.b.title"),
		IntroductionText: bundle.T(locale, "dialog.b.intro"),
		SubmitLabel:      bundle.T(locale, "dialog.submit"),
		NotifyOnCancel:   false,
		State:            StateVariantB,
		Elements: []model.DialogElement{
			{
				DisplayName: bundle.T(locale, "dialog.b.expected_end"),
				Name:        "expected_end_date",
				Type:        "date",
				Placeholder: bundle.T(locale, "dialog.date_placeholder"),
				HelpText:    bundle.T(locale, "dialog.b.expected_end_help"),
				MinDate:     opts.StartDate,
			},
			{
				DisplayName: bundle.T(locale, "dialog.b.au_certificate"),
				Name:        "au_certificate",
				Type:        "select",
				HelpText:    bundle.T(locale, "dialog.b.au_certificate_help"),
				Options: []*model.PostActionOptions{
					{Text: bundle.T(locale, "dialog.au.yes"), Value: "yes"},
					{Text: bundle.T(locale, "dialog.au.no"), Value: "no"},
					{Text: bundle.T(locale, "dialog.au.child"), Value: "child"},
				},
			},
		},
	}
}

func BuildExtendDialog(locale string, bundle *i18n.Bundle, opts ExtendDialogOptions) model.Dialog {
	minDate := opts.CurrentExpectedEnd
	if parsed, err := time.Parse(dateLayout, opts.CurrentExpectedEnd); err == nil {
		minDate = formatDate(parsed.AddDate(0, 0, 1))
	}

	return model.Dialog{
		CallbackId:       CallbackExtend,
		Title:            bundle.T(locale, "dialog.c.title"),
		IntroductionText: bundle.T(locale, "dialog.c.intro"),
		SubmitLabel:      bundle.T(locale, "dialog.submit"),
		NotifyOnCancel:   false,
		State:            StateVariantC,
		Elements: []model.DialogElement{
			{
				DisplayName: bundle.T(locale, "dialog.c.expected_end"),
				Name:        "expected_end_date",
				Type:        "date",
				Placeholder: bundle.T(locale, "dialog.date_placeholder"),
				HelpText:    bundle.T(locale, "dialog.c.expected_end_help"),
				MinDate:     minDate,
			},
			{
				DisplayName: bundle.T(locale, "dialog.c.au_certificate"),
				Name:        "au_certificate",
				Type:        "select",
				Optional:    true,
				HelpText:    bundle.T(locale, "dialog.c.au_certificate_help"),
				Options: []*model.PostActionOptions{
					{Text: bundle.T(locale, "dialog.au.unchanged"), Value: "unchanged"},
					{Text: bundle.T(locale, "dialog.au.yes"), Value: "yes"},
					{Text: bundle.T(locale, "dialog.au.no"), Value: "no"},
					{Text: bundle.T(locale, "dialog.au.child"), Value: "child"},
				},
			},
		},
	}
}

func OpenStartDialog(api DialogAPI, triggerID, submitURL, locale string, bundle *i18n.Bundle, opts StartDialogOptions) *model.AppError {
	return api.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: triggerID,
		URL:       submitURL,
		Dialog:    BuildStartDialog(locale, bundle, opts),
	})
}

func OpenUpdateDialog(api DialogAPI, triggerID, submitURL, locale string, bundle *i18n.Bundle, opts UpdateDialogOptions) *model.AppError {
	return api.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: triggerID,
		URL:       submitURL,
		Dialog:    BuildUpdateDialog(locale, bundle, opts),
	})
}

func OpenExtendDialog(api DialogAPI, triggerID, submitURL, locale string, bundle *i18n.Bundle, opts ExtendDialogOptions) *model.AppError {
	return api.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: triggerID,
		URL:       submitURL,
		Dialog:    BuildExtendDialog(locale, bundle, opts),
	})
}

func truncateToDate(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
