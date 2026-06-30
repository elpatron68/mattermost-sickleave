package i18n

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

//go:embed en.json
var enJSON []byte

//go:embed de.json
var deJSON []byte

type Bundle struct {
	locales map[string]map[string]string
}

func NewBundle() (*Bundle, error) {
	en := map[string]string{}
	if err := json.Unmarshal(enJSON, &en); err != nil {
		return nil, fmt.Errorf("parse en translations: %w", err)
	}

	de := map[string]string{}
	if err := json.Unmarshal(deJSON, &de); err != nil {
		return nil, fmt.Errorf("parse de translations: %w", err)
	}

	return &Bundle{
		locales: map[string]map[string]string{
			"en": en,
			"de": de,
		},
	}, nil
}

func (b *Bundle) T(locale, key string, args ...any) string {
	text := b.lookup(locale, key)
	if len(args) == 0 {
		return text
	}
	return fmt.Sprintf(text, args...)
}

func (b *Bundle) lookup(locale, key string) string {
	if messages, ok := b.locales[normalizeLocale(locale)]; ok {
		if text, ok := messages[key]; ok {
			return text
		}
	}
	if messages, ok := b.locales["en"]; ok {
		if text, ok := messages[key]; ok {
			return text
		}
	}
	return key
}

func LocaleForUser(user *model.User, defaultLocale string) string {
	if user != nil && user.Locale != "" {
		return normalizeLocale(user.Locale)
	}
	return normalizeLocale(defaultLocale)
}

func normalizeLocale(locale string) string {
	locale = strings.ToLower(strings.TrimSpace(locale))
	if strings.HasPrefix(locale, "de") {
		return "de"
	}
	return "en"
}
