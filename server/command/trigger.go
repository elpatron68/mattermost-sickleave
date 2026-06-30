package command

import (
	"strings"
	"unicode"
)

const DefaultCommandTrigger = "sick-leave"

func NormalizeCommandTrigger(trigger string) string {
	trigger = strings.TrimSpace(trigger)
	trigger = strings.TrimPrefix(trigger, "/")

	var b strings.Builder
	for _, r := range trigger {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			b.WriteRune(unicode.ToLower(r))
		}
	}

	normalized := b.String()
	if normalized == "" {
		return DefaultCommandTrigger
	}

	return normalized
}
