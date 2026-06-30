package sickleave

import "time"

type Status string

const (
	StatusReported Status = "reported"
	StatusUpdated  Status = "updated"
	StatusExtended Status = "extended"
	StatusClosed   Status = "closed"
)

type HistoryEntry struct {
	Variant   string    `json:"variant"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data,omitempty"`
}

type Record struct {
	ID              string         `json:"id"`
	UserID          string         `json:"user_id"`
	TeamID          string         `json:"team_id"`
	StartDate       string         `json:"start_date"`
	ExpectedEndDate string         `json:"expected_end_date,omitempty"`
	AUCertificate   *bool          `json:"au_certificate,omitempty"`
	Status          Status         `json:"status"`
	HRPostID        string         `json:"hr_post_id"`
	HRChannelID     string         `json:"hr_channel_id"`
	Hashtag           string         `json:"hashtag,omitempty"`
	History           []HistoryEntry `json:"history"`
}
