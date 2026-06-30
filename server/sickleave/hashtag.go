package sickleave

import (
	"strings"
	"unicode"
)

const DefaultReportHashtag = "#krankmeldung"

func NormalizeHashtag(tag string) string {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return DefaultReportHashtag
	}

	tag = strings.TrimPrefix(tag, "#")
	var b strings.Builder
	for _, r := range tag {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			b.WriteRune(r)
		}
	}

	normalized := b.String()
	if normalized == "" {
		return DefaultReportHashtag
	}

	return "#" + normalized
}
