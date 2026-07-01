package sickleave

import (
	"fmt"
	"strings"
	"time"
)

const (
	isoDateLayout     = "2006-01-02"
	germanDateLayout  = "02.01.2006"
	englishDateLayout = "01/02/2006"
)

func NormalizeDateLocale(locale string) string {
	locale = strings.ToLower(strings.TrimSpace(locale))
	if strings.HasPrefix(locale, "de") {
		return "de"
	}
	return "en"
}

func FormatISODate(value time.Time) string {
	year, month, day := value.Date()
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func FormatDateForLocale(isoDate, locale string) string {
	parsed, err := time.Parse(isoDateLayout, strings.TrimSpace(isoDate))
	if err != nil {
		return isoDate
	}
	return formatDateForLocale(parsed, locale)
}

func formatDateForLocale(value time.Time, locale string) string {
	switch NormalizeDateLocale(locale) {
	case "de":
		return value.Format(germanDateLayout)
	default:
		return value.Format(englishDateLayout)
	}
}

func ParseDateForLocale(value, locale string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("invalid date")
	}

	if len(value) >= 10 && value[4] == '-' {
		parsed, err := time.Parse(isoDateLayout, value[:10])
		if err == nil {
			return parsed, nil
		}
	}

	for _, layout := range dateInputLayouts(locale) {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date")
}

func dateInputLayouts(locale string) []string {
	switch NormalizeDateLocale(locale) {
	case "de":
		return []string{germanDateLayout, "2.1.2006"}
	default:
		return []string{englishDateLayout, "1/2/2006"}
	}
}
