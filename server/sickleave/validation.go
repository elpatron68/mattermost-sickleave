package sickleave

import (
	"fmt"
	"time"
)

var dateLayout = isoDateLayout

func ParseDate(value string) (time.Time, error) {
	return ParseDateForLocale(value, "en")
}

func ValidateStartDate(startDate time.Time, today time.Time, maxBackdateDays int) error {
	today = truncateToDate(today)
	startDate = truncateToDate(startDate)

	if startDate.After(today) {
		return fmt.Errorf("start date is in the future")
	}

	earliest := today.AddDate(0, 0, -maxBackdateDays)
	if startDate.Before(earliest) {
		return fmt.Errorf("start date is too far in the past")
	}

	return nil
}

func ValidateExpectedEndDate(startDate, expectedEnd time.Time) error {
	startDate = truncateToDate(startDate)
	expectedEnd = truncateToDate(expectedEnd)

	if expectedEnd.Before(startDate) {
		return fmt.Errorf("expected end date is before start date")
	}

	return nil
}

func ValidateExtensionEndDate(currentEnd, newEnd time.Time) error {
	currentEnd = truncateToDate(currentEnd)
	newEnd = truncateToDate(newEnd)

	if !newEnd.After(currentEnd) {
		return fmt.Errorf("new expected end date must be after current expected end date")
	}

	return nil
}

func truncateToDate(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}
