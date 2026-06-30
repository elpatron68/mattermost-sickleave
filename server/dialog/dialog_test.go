package dialog

import (
	"testing"
	"time"

	"github.com/medisoftware/mattermost-sickleave/server/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildStartDialogUsesDatePickerWithBounds(t *testing.T) {
	t.Parallel()

	bundle, err := i18n.NewBundle()
	require.NoError(t, err)

	today := time.Date(2026, 6, 30, 12, 0, 0, 0, time.UTC)
	d := BuildStartDialog("en", bundle, StartDialogOptions{
		Today:           today,
		MaxBackdateDays: 3,
	})

	require.Len(t, d.Elements, 1)
	assert.Equal(t, "date", d.Elements[0].Type)
	assert.Equal(t, "2026-06-27", d.Elements[0].MinDate)
	assert.Equal(t, "2026-06-30", d.Elements[0].MaxDate)
}

func TestBuildUpdateDialogUsesDatePickerWithMinDate(t *testing.T) {
	t.Parallel()

	bundle, err := i18n.NewBundle()
	require.NoError(t, err)

	d := BuildUpdateDialog("en", bundle, UpdateDialogOptions{StartDate: "2026-06-20"})
	require.Len(t, d.Elements, 2)
	assert.Equal(t, "date", d.Elements[0].Type)
	assert.Equal(t, "2026-06-20", d.Elements[0].MinDate)
}

func TestBuildExtendDialogUsesDatePickerWithMinDate(t *testing.T) {
	t.Parallel()

	bundle, err := i18n.NewBundle()
	require.NoError(t, err)

	d := BuildExtendDialog("en", bundle, ExtendDialogOptions{CurrentExpectedEnd: "2026-06-25"})
	require.Len(t, d.Elements, 2)
	assert.Equal(t, "date", d.Elements[0].Type)
	assert.Equal(t, "2026-06-26", d.Elements[0].MinDate)
}
