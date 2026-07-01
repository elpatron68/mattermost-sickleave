package sickleave

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/medisoftware/mattermost-sickleave/server/i18n"
)

type AUCertificate string

const (
	AUYes   AUCertificate = "yes"
	AUNo    AUCertificate = "no"
	AUChild AUCertificate = "child"
)

func ParseAUCertificate(value string) (AUCertificate, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(AUYes):
		return AUYes, true
	case string(AUNo):
		return AUNo, true
	case string(AUChild):
		return AUChild, true
	default:
		return "", false
	}
}

func (a AUCertificate) Format(locale string, bundle *i18n.Bundle) string {
	switch a {
	case AUYes:
		return bundle.T(locale, "hr.post.au.yes")
	case AUNo:
		return bundle.T(locale, "hr.post.au.no")
	case AUChild:
		return bundle.T(locale, "hr.post.au.child")
	default:
		return bundle.T(locale, "hr.post.au.unchanged")
	}
}

func (a *AUCertificate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*a = ""
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		if asString == "" {
			*a = ""
			return nil
		}
		parsed, ok := ParseAUCertificate(asString)
		if !ok {
			return fmt.Errorf("invalid au_certificate value: %s", asString)
		}
		*a = parsed
		return nil
	}

	var asBool bool
	if err := json.Unmarshal(data, &asBool); err == nil {
		if asBool {
			*a = AUYes
		} else {
			*a = AUNo
		}
		return nil
	}

	return fmt.Errorf("invalid au_certificate value: %s", string(data))
}

func (a AUCertificate) MarshalJSON() ([]byte, error) {
	if a == "" {
		return []byte("null"), nil
	}
	return json.Marshal(string(a))
}
