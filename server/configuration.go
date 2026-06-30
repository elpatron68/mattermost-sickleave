package main

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/elpatron68/mattermost-sickleave/server/command"
	"github.com/elpatron68/mattermost-sickleave/server/sickleave"
)

type configuration struct {
	HRChannelID     string
	DefaultLocale   string
	MaxBackdateDays int
	ReportHashtag   string
	CommandTrigger  string
}

func (c *configuration) Clone() *configuration {
	clone := *c
	return &clone
}

func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{
			DefaultLocale:   "en",
			MaxBackdateDays: 3,
			ReportHashtag:   sickleave.DefaultReportHashtag,
			CommandTrigger:  command.DefaultCommandTrigger,
		}
	}

	return p.configuration
}

func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		if reflect.ValueOf(*configuration).NumField() == 0 {
			return
		}

		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}

func (p *Plugin) OnConfigurationChange() error {
	configuration := &configuration{
		DefaultLocale:   "en",
		MaxBackdateDays: 3,
		ReportHashtag:   sickleave.DefaultReportHashtag,
		CommandTrigger:  command.DefaultCommandTrigger,
	}

	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}

	if configuration.DefaultLocale == "" {
		configuration.DefaultLocale = "en"
	}
	if configuration.MaxBackdateDays <= 0 {
		configuration.MaxBackdateDays = 3
	}
	configuration.ReportHashtag = sickleave.NormalizeHashtag(configuration.ReportHashtag)
	configuration.CommandTrigger = command.NormalizeCommandTrigger(configuration.CommandTrigger)

	p.setConfiguration(configuration)

	if p.command != nil {
		if err := p.command.EnsureSlashCommandRegistered(); err != nil {
			return errors.Wrap(err, "failed to register slash command")
		}
	}

	return nil
}

func (p *Plugin) settingsFromConfig() commandSettings {
	config := p.getConfiguration()
	return commandSettings{
		HRChannelID:     config.HRChannelID,
		DefaultLocale:   config.DefaultLocale,
		MaxBackdateDays: config.MaxBackdateDays,
		ReportHashtag:   config.ReportHashtag,
		CommandTrigger:  config.CommandTrigger,
	}
}

type commandSettings struct {
	HRChannelID     string
	DefaultLocale   string
	MaxBackdateDays int
	ReportHashtag   string
	CommandTrigger  string
}
