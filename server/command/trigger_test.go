package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeCommandTrigger(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "sick-leave", NormalizeCommandTrigger(""))
	assert.Equal(t, "krankmeldung", NormalizeCommandTrigger("/krankmeldung"))
	assert.Equal(t, "krankmeldung", NormalizeCommandTrigger(" Krankmeldung "))
	assert.Equal(t, "sick_leave", NormalizeCommandTrigger("sick_leave"))
	assert.Equal(t, "sick-leave", NormalizeCommandTrigger("///"))
}
