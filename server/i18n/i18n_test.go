package i18n

import (
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBundleTranslations(t *testing.T) {
	bundle, err := NewBundle()
	require.NoError(t, err)

	assert.Equal(t, "Report sick leave", bundle.T("en", "dialog.a.title"))
	assert.Equal(t, "Krankmeldung", bundle.T("de", "dialog.a.title"))
	assert.Equal(t, "Report sick leave", bundle.T("en-US", "dialog.a.title"))
}

func TestLocaleForUser(t *testing.T) {
	assert.Equal(t, "de", LocaleForUser(&model.User{Locale: "de"}, "en"))
	assert.Equal(t, "en", LocaleForUser(&model.User{Locale: "en"}, "de"))
	assert.Equal(t, "de", LocaleForUser(nil, "de"))
}
