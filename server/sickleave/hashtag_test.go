package sickleave

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeHashtag(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "#krankmeldung", NormalizeHashtag(""))
	assert.Equal(t, "#krankmeldung", NormalizeHashtag("  "))
	assert.Equal(t, "#krankmeldung", NormalizeHashtag("krankmeldung"))
	assert.Equal(t, "#krankmeldung", NormalizeHashtag("#krankmeldung"))
	assert.Equal(t, "#sick_leave", NormalizeHashtag(" sick_leave "))
	assert.Equal(t, "#krankmeldung", NormalizeHashtag("###"))
}
