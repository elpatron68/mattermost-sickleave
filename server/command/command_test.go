package command

import (
	"testing"

	"github.com/elpatron68/mattermost-sickleave/server/dialog"
	"github.com/elpatron68/mattermost-sickleave/server/i18n"
	"github.com/elpatron68/mattermost-sickleave/server/sickleave"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type env struct {
	client *pluginapi.Client
	api    *plugintest.API
}

type memoryStore struct {
	active  map[string]*sickleave.Record
	records map[string]*sickleave.Record
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		active:  map[string]*sickleave.Record{},
		records: map[string]*sickleave.Record{},
	}
}

func (s *memoryStore) GetActive(userID string) (*sickleave.Record, error) {
	record, ok := s.active[userID]
	if !ok {
		return nil, nil
	}
	return record, nil
}

func (s *memoryStore) SetActive(userID string, record *sickleave.Record) error {
	s.active[userID] = record
	return nil
}

func (s *memoryStore) ClearActive(userID string) error {
	delete(s.active, userID)
	return nil
}

func (s *memoryStore) SaveRecord(record *sickleave.Record) error {
	s.records[record.ID] = record
	return nil
}

func (s *memoryStore) GetRecord(recordID string) (*sickleave.Record, error) {
	record, ok := s.records[recordID]
	if !ok {
		return nil, nil
	}
	return record, nil
}

func setupTest(t *testing.T) (*env, *i18n.Bundle) {
	t.Helper()

	api := &plugintest.API{}
	driver := &plugintest.Driver{}
	client := pluginapi.NewClient(api, driver)

	bundle, err := i18n.NewBundle()
	require.NoError(t, err)

	return &env{
		client: client,
		api:    api,
	}, bundle
}

func expectCommandRegistration(api *plugintest.API, trigger string) {
	if trigger == "" {
		trigger = DefaultCommandTrigger
	}
	api.On("RegisterCommand", &model.Command{
		Trigger:          trigger,
		AutoComplete:     true,
		AutoCompleteDesc: "Report and manage sick leave",
		AutoCompleteHint: "[start|update|extend|end|status|help]",
		AutocompleteData: buildAutocompleteData(trigger),
	}).Return(nil)
}

func newTestHandler(t *testing.T, store sickleave.Store, settings Settings, userID string) (Command, *env) {
	t.Helper()

	testEnv, bundle := setupTest(t)
	if settings.CommandTrigger == "" {
		settings.CommandTrigger = DefaultCommandTrigger
	}
	expectCommandRegistration(testEnv.api, settings.CommandTrigger)
	testEnv.api.On("GetUser", userID).Return(&model.User{Id: userID, Locale: "en"}, nil)

	handler := NewCommandHandler(HandlerConfig{
		Client: testEnv.client,
		Store:  store,
		Settings: func() Settings {
			return settings
		},
		Bundle:    bundle,
		PluginID:  "com.elpatron68.mattermost-sickleave",
		BotUserID: "bot-1",
	})

	require.NoError(t, handler.EnsureSlashCommandRegistered())

	return handler, testEnv
}

func TestHelpCommand(t *testing.T) {
	env, bundle := setupTest(t)
	expectCommandRegistration(env.api, DefaultCommandTrigger)
	env.api.On("GetUser", "").Return(&model.User{Locale: "en"}, nil)

	handler := NewCommandHandler(HandlerConfig{
		Client: env.client,
		Settings: func() Settings {
			return Settings{DefaultLocale: "en", CommandTrigger: DefaultCommandTrigger}
		},
		Bundle: bundle,
	})
	require.NoError(t, handler.EnsureSlashCommandRegistered())

	response, err := handler.Handle(&model.CommandArgs{Command: "/sick-leave help"})
	require.NoError(t, err)
	assert.Equal(t, model.CommandResponseTypeEphemeral, response.ResponseType)
	assert.Contains(t, response.Text, "Sick Leave")
	assert.Contains(t, response.Text, "/sick-leave update")
}

func TestStartWithoutHRChannel(t *testing.T) {
	env, bundle := setupTest(t)
	expectCommandRegistration(env.api, DefaultCommandTrigger)
	env.api.On("GetUser", "user-1").Return(&model.User{
		Id:     "user-1",
		Locale: "en",
	}, nil)

	handler := NewCommandHandler(HandlerConfig{
		Client: env.client,
		Settings: func() Settings {
			return Settings{DefaultLocale: "en", CommandTrigger: DefaultCommandTrigger}
		},
		Bundle: bundle,
	})
	require.NoError(t, handler.EnsureSlashCommandRegistered())

	response, err := handler.Handle(&model.CommandArgs{
		Command: "/sick-leave start",
		UserId:  "user-1",
	})
	require.NoError(t, err)
	assert.Contains(t, response.Text, "not configured")
}

func TestUpdateRequiresReportedStatus(t *testing.T) {
	store := newMemoryStore()
	store.active["user-1"] = &sickleave.Record{
		ID:        "rec-1",
		UserID:    "user-1",
		StartDate: "2026-06-20",
		Status:    sickleave.StatusUpdated,
	}

	handler, _ := newTestHandler(t, store, Settings{
		HRChannelID:   "channel-hr",
		DefaultLocale: "en",
	}, "user-1")

	response, err := handler.Handle(&model.CommandArgs{
		Command: "/sick-leave update",
		UserId:  "user-1",
	})
	require.NoError(t, err)
	assert.Contains(t, response.Text, "only possible after the initial")
}

func TestExtendRequiresUpdatedStatus(t *testing.T) {
	store := newMemoryStore()
	store.active["user-1"] = &sickleave.Record{
		ID:        "rec-1",
		UserID:    "user-1",
		StartDate: "2026-06-20",
		Status:    sickleave.StatusReported,
	}

	handler, _ := newTestHandler(t, store, Settings{
		HRChannelID:   "channel-hr",
		DefaultLocale: "en",
	}, "user-1")

	response, err := handler.Handle(&model.CommandArgs{
		Command: "/sick-leave extend",
		UserId:  "user-1",
	})
	require.NoError(t, err)
	assert.Contains(t, response.Text, "only possible after providing")
}

func TestSubmitUpdateDialogTransitionsToUpdated(t *testing.T) {
	store := newMemoryStore()
	store.active["user-1"] = &sickleave.Record{
		ID:          "rec-1",
		UserID:      "user-1",
		StartDate:   "2026-06-20",
		Status:      sickleave.StatusReported,
		HRChannelID: "channel-hr",
		HRPostID:    "post-root",
	}

	env, bundle := setupTest(t)
	expectCommandRegistration(env.api, DefaultCommandTrigger)
	env.api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Locale: "en"}, nil)
	env.api.On("CreatePost", mockMatchedHRThreadPost("channel-hr", "post-root")).Return(&model.Post{Id: "post-update"}, nil).Once()
	env.api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(&model.Post{})

	handler := NewCommandHandler(HandlerConfig{
		Client: env.client,
		Store:  store,
		Settings: func() Settings {
			return Settings{HRChannelID: "channel-hr", DefaultLocale: "en"}
		},
		Bundle:    bundle,
		BotUserID: "bot-1",
	})

	response, err := handler.SubmitDialog(&model.SubmitDialogRequest{
		CallbackId: dialog.CallbackUpdate,
		UserId:     "user-1",
		ChannelId:  "channel-1",
		Submission: map[string]any{
			"expected_end_date": "2026-06-25",
			"au_certificate":    "yes",
		},
	})
	require.NoError(t, err)
	assert.Empty(t, response.Errors)

	record := store.active["user-1"]
	require.NotNil(t, record)
	assert.Equal(t, sickleave.StatusUpdated, record.Status)
	assert.Equal(t, "2026-06-25", record.ExpectedEndDate)
	require.NotNil(t, record.AUCertificate)
	assert.True(t, *record.AUCertificate)
}

func TestSubmitExtendDialogTransitionsToExtended(t *testing.T) {
	au := true
	store := newMemoryStore()
	store.active["user-1"] = &sickleave.Record{
		ID:              "rec-1",
		UserID:          "user-1",
		StartDate:       "2026-06-20",
		ExpectedEndDate: "2026-06-25",
		AUCertificate:   &au,
		Status:          sickleave.StatusUpdated,
		HRChannelID:     "channel-hr",
		HRPostID:        "post-root",
	}

	env, bundle := setupTest(t)
	expectCommandRegistration(env.api, DefaultCommandTrigger)
	env.api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Locale: "en"}, nil)
	env.api.On("CreatePost", mockMatchedHRThreadPost("channel-hr", "post-root")).Return(&model.Post{Id: "post-extend"}, nil).Once()
	env.api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(&model.Post{})

	handler := NewCommandHandler(HandlerConfig{
		Client: env.client,
		Store:  store,
		Settings: func() Settings {
			return Settings{HRChannelID: "channel-hr", DefaultLocale: "en"}
		},
		Bundle:    bundle,
		BotUserID: "bot-1",
	})

	response, err := handler.SubmitDialog(&model.SubmitDialogRequest{
		CallbackId: dialog.CallbackExtend,
		UserId:     "user-1",
		ChannelId:  "channel-1",
		Submission: map[string]any{
			"expected_end_date": "2026-06-30",
			"au_certificate":    "no",
		},
	})
	require.NoError(t, err)
	assert.Empty(t, response.Errors)

	record := store.active["user-1"]
	require.NotNil(t, record)
	assert.Equal(t, sickleave.StatusExtended, record.Status)
	assert.Equal(t, "2026-06-30", record.ExpectedEndDate)
	require.NotNil(t, record.AUCertificate)
	assert.False(t, *record.AUCertificate)
}

func TestSubmitExtendDialogRejectsNonExtensionDate(t *testing.T) {
	au := true
	store := newMemoryStore()
	store.active["user-1"] = &sickleave.Record{
		ID:              "rec-1",
		UserID:          "user-1",
		StartDate:       "2026-06-20",
		ExpectedEndDate: "2026-06-25",
		AUCertificate:   &au,
		Status:          sickleave.StatusUpdated,
		HRChannelID:     "channel-hr",
		HRPostID:        "post-root",
	}

	env, bundle := setupTest(t)
	expectCommandRegistration(env.api, DefaultCommandTrigger)
	env.api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Locale: "en"}, nil)

	handler := NewCommandHandler(HandlerConfig{
		Client: env.client,
		Store:  store,
		Settings: func() Settings {
			return Settings{HRChannelID: "channel-hr", DefaultLocale: "en"}
		},
		Bundle:    bundle,
		BotUserID: "bot-1",
	})

	response, err := handler.SubmitDialog(&model.SubmitDialogRequest{
		CallbackId: dialog.CallbackExtend,
		UserId:     "user-1",
		ChannelId:  "channel-1",
		Submission: map[string]any{
			"expected_end_date": "2026-06-25",
		},
	})
	require.NoError(t, err)
	assert.Contains(t, response.Errors["expected_end_date"], "after your current expected return date")
}

func TestHandleEndClosesActiveCase(t *testing.T) {
	store := newMemoryStore()
	store.active["user-1"] = &sickleave.Record{
		ID:              "rec-1",
		UserID:          "user-1",
		StartDate:       "2026-06-20",
		ExpectedEndDate: "2026-06-25",
		Status:          sickleave.StatusUpdated,
		HRChannelID:     "channel-hr",
		HRPostID:        "post-root",
	}

	env, bundle := setupTest(t)
	expectCommandRegistration(env.api, DefaultCommandTrigger)
	env.api.On("GetUser", "user-1").Return(&model.User{Id: "user-1", Locale: "en"}, nil)
	env.api.On("CreatePost", mockMatchedHRThreadPost("channel-hr", "post-root")).Return(&model.Post{Id: "post-close"}, nil).Once()

	handler := NewCommandHandler(HandlerConfig{
		Client: env.client,
		Store:  store,
		Settings: func() Settings {
			return Settings{HRChannelID: "channel-hr", DefaultLocale: "en"}
		},
		Bundle:    bundle,
		BotUserID: "bot-1",
	})

	response, err := handler.End("user-1", "channel-1")
	require.NoError(t, err)
	assert.Contains(t, response.Text, "closed")

	_, hasActive := store.active["user-1"]
	assert.False(t, hasActive)

	closed := store.records["rec-1"]
	require.NotNil(t, closed)
	assert.Equal(t, sickleave.StatusClosed, closed.Status)
}

func TestHandleEndWithoutActiveCase(t *testing.T) {
	store := newMemoryStore()
	handler, _ := newTestHandler(t, store, Settings{DefaultLocale: "en"}, "user-1")

	response, err := handler.End("user-1", "channel-1")
	require.NoError(t, err)
	assert.Contains(t, response.Text, "do not have an active")
}

func TestParseAUCertificate(t *testing.T) {
	t.Parallel()

	value, ok := parseAUCertificate("yes")
	assert.True(t, ok)
	assert.True(t, value)

	value, ok = parseAUCertificate("no")
	assert.True(t, ok)
	assert.False(t, value)

	_, ok = parseAUCertificate("maybe")
	assert.False(t, ok)
}

func mockMatchedHRThreadPost(channelID, rootID string) any {
	return mock.MatchedBy(func(post *model.Post) bool {
		return post != nil &&
			post.UserId == "bot-1" &&
			post.ChannelId == channelID &&
			post.RootId == rootID &&
			post.Message != ""
	})
}
