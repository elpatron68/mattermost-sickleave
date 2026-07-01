package sickleave

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medisoftware/mattermost-sickleave/server/i18n"
)

func TestParseAUCertificate(t *testing.T) {
	t.Parallel()

	value, ok := ParseAUCertificate("yes")
	assert.True(t, ok)
	assert.Equal(t, AUYes, value)

	value, ok = ParseAUCertificate("no")
	assert.True(t, ok)
	assert.Equal(t, AUNo, value)

	value, ok = ParseAUCertificate("child")
	assert.True(t, ok)
	assert.Equal(t, AUChild, value)

	_, ok = ParseAUCertificate("maybe")
	assert.False(t, ok)
}

func TestAUCertificateUnmarshalLegacyBoolean(t *testing.T) {
	t.Parallel()

	var record Record
	require.NoError(t, json.Unmarshal([]byte(`{"au_certificate":true}`), &record))
	assert.Equal(t, AUYes, record.AUCertificate)

	require.NoError(t, json.Unmarshal([]byte(`{"au_certificate":false}`), &record))
	assert.Equal(t, AUNo, record.AUCertificate)
}

func TestAUCertificateFormat(t *testing.T) {
	t.Parallel()

	bundle, err := i18n.NewBundle()
	require.NoError(t, err)

	assert.Equal(t, "AU-Kind", AUChild.Format("de", bundle))
	assert.Equal(t, "Child sickness certificate", AUChild.Format("en", bundle))
}
