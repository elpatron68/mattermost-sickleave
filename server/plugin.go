package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/pkg/errors"

	"github.com/elpatron68/mattermost-sickleave/server/command"
	"github.com/elpatron68/mattermost-sickleave/server/i18n"
	"github.com/elpatron68/mattermost-sickleave/server/sickleave"
)

type Plugin struct {
	plugin.MattermostPlugin

	kvstore sickleave.Store
	client  *pluginapi.Client
	command command.Command
	router  *mux.Router
	bundle  *i18n.Bundle
	botID   string

	configurationLock sync.RWMutex
	configuration     *configuration
}

func (p *Plugin) OnActivate() error {
	p.client = pluginapi.NewClient(p.API, p.Driver)
	p.kvstore = sickleave.NewStore(p.client)

	bundle, err := i18n.NewBundle()
	if err != nil {
		return errors.Wrap(err, "failed to load translations")
	}
	p.bundle = bundle

	bot := &model.Bot{
		Username:    "sickleave",
		Description: "Posts sick leave reports to HR channels.",
	}
	botID, err := p.client.Bot.EnsureBot(bot, pluginapi.ProfileImageBytes(pluginIcon))
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot user")
	}
	p.botID = botID

	p.command = command.NewCommandHandler(command.HandlerConfig{
		Client:    p.client,
		DialogAPI: p.API,
		Store:     p.kvstore,
		Settings: func() command.Settings {
			settings := p.settingsFromConfig()
			return command.Settings{
				HRChannelID:     settings.HRChannelID,
				DefaultLocale:   settings.DefaultLocale,
				MaxBackdateDays: settings.MaxBackdateDays,
			}
		},
		Bundle:    p.bundle,
		PluginID:  manifest.Id,
		SiteURL:   p.getSiteURL,
		BotUserID: p.botID,
	})

	p.router = p.initRouter()

	return nil
}

func (p *Plugin) getSiteURL() (string, error) {
	config := p.API.GetConfig()
	if config == nil || config.ServiceSettings.SiteURL == nil || *config.ServiceSettings.SiteURL == "" {
		return "", errors.New("site URL is not configured")
	}
	return *config.ServiceSettings.SiteURL, nil
}

func (p *Plugin) ExecuteCommand(_ *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	response, err := p.command.Handle(args)
	if err != nil {
		return nil, model.NewAppError("ExecuteCommand", "plugin.command.execute_command.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return response, nil
}

var pluginIcon = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
	0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
	0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
}
