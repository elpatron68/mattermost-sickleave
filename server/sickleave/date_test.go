package sickleave

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDateForLocaleGerman(t *testing.T) {
	t.Parallel()

	parsed, err := ParseDateForLocale("01.07.2026", "de")
	require.NoError(t, err)
	assert.Equal(t, time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC), parsed)

	parsed, err = ParseDateForLocale("2026-07-01", "de")
	require.NoError(t, err)
	assert.Equal(t, 1, parsed.Day())
}

func TestParseDateForLocaleEnglish(t *testing.T) {
	t.Parallel()

	parsed, err := ParseDateForLocale("07/01/2026", "en")
	require.NoError(t, err)
	assert.Equal(t, time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC), parsed)
}

func TestFormatDateForLocale(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "01.07.2026", FormatDateForLocale("2026-07-01", "de"))
	assert.Equal(t, "07/01/2026", FormatDateForLocale("2026-07-01", "en"))
}
