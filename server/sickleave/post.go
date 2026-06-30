package sickleave

import (
	"fmt"
	"strings"

	"github.com/elpatron68/mattermost-sickleave/server/i18n"
	"github.com/mattermost/mattermost/server/public/model"
)

func recordHashtag(record *Record) string {
	if record == nil || record.Hashtag == "" {
		return DefaultReportHashtag
	}
	return record.Hashtag
}

func FormatInitialHRPost(record *Record, user *model.User, locale string, bundle *i18n.Bundle) string {
	username := "unknown"
	if user != nil {
		username = fmt.Sprintf("@%s", user.Username)
	}

	return formatFieldValuePost(
		bundle.T(locale, "hr.post.a.title"),
		bundle.T(locale, "hr.post.table.field"),
		bundle.T(locale, "hr.post.table.value"),
		recordHashtag(record),
		[][2]string{
		{bundle.T(locale, "hr.post.field.employee"), username},
		{bundle.T(locale, "hr.post.field.first_sick_day"), record.StartDate},
		{bundle.T(locale, "hr.post.field.status"), bundle.T(locale, "command.status.reported")},
	})
}

func FormatUpdateHRPost(record *Record, expectedEnd string, auCertificate bool, locale string, bundle *i18n.Bundle) string {
	return formatFieldValuePost(
		bundle.T(locale, "hr.post.b.title"),
		bundle.T(locale, "hr.post.table.field"),
		bundle.T(locale, "hr.post.table.value"),
		recordHashtag(record),
		[][2]string{
		{bundle.T(locale, "hr.post.field.expected_end"), expectedEnd},
		{bundle.T(locale, "hr.post.field.au_certificate"), formatAUCertificate(auCertificate, locale, bundle)},
		{bundle.T(locale, "hr.post.field.status"), bundle.T(locale, "command.status.updated")},
	})
}

func FormatExtendHRPost(record *Record, newExpectedEnd string, auCertificate *bool, locale string, bundle *i18n.Bundle) string {
	auValue := bundle.T(locale, "hr.post.au.unchanged")
	if auCertificate != nil {
		auValue = formatAUCertificate(*auCertificate, locale, bundle)
	}

	return formatFieldValuePost(
		bundle.T(locale, "hr.post.c.title"),
		bundle.T(locale, "hr.post.table.field"),
		bundle.T(locale, "hr.post.table.value"),
		recordHashtag(record),
		[][2]string{
		{bundle.T(locale, "hr.post.field.expected_end"), newExpectedEnd},
		{bundle.T(locale, "hr.post.field.au_certificate"), auValue},
		{bundle.T(locale, "hr.post.field.status"), bundle.T(locale, "command.status.extended")},
	})
}

func FormatCloseHRPost(record *Record, locale string, bundle *i18n.Bundle) string {
	rows := [][2]string{
		{bundle.T(locale, "hr.post.field.first_sick_day"), record.StartDate},
		{bundle.T(locale, "hr.post.field.status"), bundle.T(locale, "command.status.closed")},
	}
	if record.ExpectedEndDate != "" {
		rows = append([][2]string{{bundle.T(locale, "hr.post.field.expected_end"), record.ExpectedEndDate}}, rows...)
	}

	return formatFieldValuePost(
		bundle.T(locale, "hr.post.d.title"),
		bundle.T(locale, "hr.post.table.field"),
		bundle.T(locale, "hr.post.table.value"),
		recordHashtag(record),
		rows,
	)
}

func formatFieldValuePost(title, fieldHeader, valueHeader, hashtag string, rows [][2]string) string {
	var b strings.Builder
	b.WriteString("**")
	b.WriteString(title)
	b.WriteString("**\n\n| ")
	b.WriteString(fieldHeader)
	b.WriteString(" | ")
	b.WriteString(valueHeader)
	b.WriteString(" |\n| :-- | :-- |\n")
	for _, row := range rows {
		fmt.Fprintf(&b, "| %s | %s |\n", row[0], row[1])
	}
	if hashtag != "" {
		b.WriteString("\n\n")
		b.WriteString(hashtag)
	}
	return strings.TrimRight(b.String(), "\n")
}

func formatAUCertificate(value bool, locale string, bundle *i18n.Bundle) string {
	if value {
		return bundle.T(locale, "hr.post.au.yes")
	}
	return bundle.T(locale, "hr.post.au.no")
}
