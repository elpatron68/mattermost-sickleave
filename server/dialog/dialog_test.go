package dialog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/medisoftware/mattermost-sickleave/server/i18n"
)

func TestBuildStartDialogUsesTextFieldWithLocaleDefault(t *testing.T) {
	t.Parallel()

	bundle, err := i18n.NewBundle()
	require.NoError(t, err)

	today := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	d := BuildStartDialog("en", bundle, StartDialogOptions{
		Today:           today,
		MaxBackdateDays: 3,
	})

	require.Len(t, d.Elements, 1)
	assert.Equal(t, "text", d.Elements[0].Type)
	assert.Equal(t, "07/01/2026", d.Elements[0].Default)
	assert.Equal(t, "MM/DD/YYYY", d.Elements[0].Placeholder)
	assert.Empty(t, d.Elements[0].MinDate)
	assert.Empty(t, d.Elements[0].MaxDate)

	d = BuildStartDialog("de", bundle, StartDialogOptions{Today: today, MaxBackdateDays: 3})
	assert.Equal(t, "01.07.2026", d.Elements[0].Default)
	assert.Equal(t, "TT.MM.JJJJ", d.Elements[0].Placeholder)
}

func TestBuildUpdateDialogUsesTextDateField(t *testing.T) {
	t.Parallel()

	bundle, err := i18n.NewBundle()
	require.NoError(t, err)

	d := BuildUpdateDialog("en", bundle, UpdateDialogOptions{StartDate: "2026-06-20"})
	require.Len(t, d.Elements, 2)
	assert.Equal(t, "text", d.Elements[0].Type)
	assert.Equal(t, "expected_end_date", d.Elements[0].Name)
	assert.Empty(t, d.Elements[0].MinDate)
}

func TestBuildExtendDialogUsesTextDateField(t *testing.T) {
	t.Parallel()

	bundle, err := i18n.NewBundle()
	require.NoError(t, err)

	d := BuildExtendDialog("en", bundle, ExtendDialogOptions{CurrentExpectedEnd: "2026-06-25"})
	require.Len(t, d.Elements, 2)
	assert.Equal(t, "text", d.Elements[0].Type)
	assert.Equal(t, "expected_end_date", d.Elements[0].Name)
	assert.Empty(t, d.Elements[0].MinDate)
}
