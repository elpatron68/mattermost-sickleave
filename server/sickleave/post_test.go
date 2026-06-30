package sickleave

import (
	"strings"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medisoftware/mattermost-sickleave/server/i18n"
)

func TestFormatInitialHRPostUsesMarkdownTable(t *testing.T) {
	t.Parallel()

	bundle, err := i18n.NewBundle()
	require.NoError(t, err)

	message := FormatInitialHRPost(&Record{StartDate: "2026-06-30"}, &model.User{Username: "markus"}, "en", bundle)

	assert.True(t, strings.HasPrefix(message, "**Sick leave — Initial report**\n\n"))
	assert.Contains(t, message, "| Field | Value |")
	assert.Contains(t, message, "| Employee | @markus |")
	assert.Contains(t, message, "| First sick day | 2026-06-30 |")
	assert.NotContains(t, message, "|\n|---|")
}

func TestFormatFieldValuePostBlankLineBeforeTable(t *testing.T) {
	t.Parallel()

	message := formatFieldValuePost("Title", "Field", "Value", "#krankmeldung", [][2]string{{"Key", "Value"}})
	assert.Contains(t, message, "**Title**\n\n| Field | Value |")
	assert.Contains(t, message, "\n\n#krankmeldung")
}

func TestFormatInitialHRPostIncludesHashtag(t *testing.T) {
	t.Parallel()

	bundle, err := i18n.NewBundle()
	require.NoError(t, err)

	message := FormatInitialHRPost(&Record{
		StartDate: "2026-06-30",
		Hashtag:   "#krankmeldung",
	}, &model.User{Username: "markus"}, "en", bundle)

	assert.Contains(t, message, "\n\n#krankmeldung")
}

func TestFormatCloseHRPostUsesMarkdownTable(t *testing.T) {
	t.Parallel()

	bundle, err := i18n.NewBundle()
	require.NoError(t, err)

	message := FormatCloseHRPost(&Record{
		StartDate:       "2026-06-20",
		ExpectedEndDate: "2026-06-25",
		Status:          StatusClosed,
	}, "en", bundle)

	assert.True(t, strings.HasPrefix(message, "**Sick leave — Case closed**\n\n"))
	assert.Contains(t, message, "| Expected return | 2026-06-25 |")
	assert.Contains(t, message, "| First sick day | 2026-06-20 |")
	assert.Contains(t, message, "| Status | Case closed |")
}
