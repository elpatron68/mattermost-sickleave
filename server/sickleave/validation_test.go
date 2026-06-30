package sickleave

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDate(t *testing.T) {
	t.Parallel()

	parsed, err := ParseDate("2026-06-28")
	require.NoError(t, err)
	assert.Equal(t, 2026, parsed.Year())
	assert.Equal(t, time.June, parsed.Month())
	assert.Equal(t, 28, parsed.Day())
}

func TestValidateStartDate(t *testing.T) {
	t.Parallel()

	today := time.Date(2026, 6, 30, 15, 0, 0, 0, time.UTC)

	t.Run("today is valid", func(t *testing.T) {
		t.Parallel()
		err := ValidateStartDate(today, today, 3)
		assert.NoError(t, err)
	})

	t.Run("future date is invalid", func(t *testing.T) {
		t.Parallel()
		err := ValidateStartDate(today.AddDate(0, 0, 1), today, 3)
		assert.Error(t, err)
	})

	t.Run("within backdate window is valid", func(t *testing.T) {
		t.Parallel()
		err := ValidateStartDate(today.AddDate(0, 0, -3), today, 3)
		assert.NoError(t, err)
	})

	t.Run("outside backdate window is invalid", func(t *testing.T) {
		t.Parallel()
		err := ValidateStartDate(today.AddDate(0, 0, -4), today, 3)
		assert.Error(t, err)
	})
}

func TestValidateExpectedEndDate(t *testing.T) {
	t.Parallel()

	start := time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC)

	t.Run("same day as start is valid", func(t *testing.T) {
		t.Parallel()
		err := ValidateExpectedEndDate(start, start)
		assert.NoError(t, err)
	})

	t.Run("after start is valid", func(t *testing.T) {
		t.Parallel()
		err := ValidateExpectedEndDate(start, start.AddDate(0, 0, 5))
		assert.NoError(t, err)
	})

	t.Run("before start is invalid", func(t *testing.T) {
		t.Parallel()
		err := ValidateExpectedEndDate(start, start.AddDate(0, 0, -1))
		assert.Error(t, err)
	})
}

func TestValidateExtensionEndDate(t *testing.T) {
	t.Parallel()

	current := time.Date(2026, 6, 25, 0, 0, 0, 0, time.UTC)

	t.Run("after current end is valid", func(t *testing.T) {
		t.Parallel()
		err := ValidateExtensionEndDate(current, current.AddDate(0, 0, 2))
		assert.NoError(t, err)
	})

	t.Run("same day as current end is invalid", func(t *testing.T) {
		t.Parallel()
		err := ValidateExtensionEndDate(current, current)
		assert.Error(t, err)
	})

	t.Run("before current end is invalid", func(t *testing.T) {
		t.Parallel()
		err := ValidateExtensionEndDate(current, current.AddDate(0, 0, -1))
		assert.Error(t, err)
	})
}
