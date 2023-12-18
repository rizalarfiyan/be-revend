package utils

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rizalarfiyan/be-revend/constants"
)

func FixedPhoneNumber(phoneNumber string) string {
	pattern := `^(\+62|62)?[\s-]?(0)?8[1-9]{1}\d{1}[\s-]?\d{4}[\s-]?\d{2,5}$`
	regex := regexp.MustCompile(pattern)
	match := regex.FindStringSubmatch(phoneNumber)
	if match == nil || len(match) < 3 {
		return phoneNumber
	}

	prefix := match[1]
	if prefix == "" && match[2] != "" {
		prefix = match[2]
	}

	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")
	return strings.Replace(phoneNumber, prefix, "62", 1)
}

func WhatsappPhoneNumber(phoneNumber string) string {
	if strings.HasSuffix(phoneNumber, constants.WhatsappNumberSuffix) {
		return phoneNumber
	}
	return FixedPhoneNumber(phoneNumber) + constants.WhatsappNumberSuffix
}

func ParseTextTemplate(rawTemplate string, data interface{}) (string, error) {
	tmpl, err := template.New("tmpl").Parse(rawTemplate)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func FullName(firstName string, lastName pgtype.Text) string {
	name := firstName
	if lastName.Valid {
		name += " " + lastName.String
	}
	return strings.TrimSpace(name)
}

func CensorPhoneNumber(phoneNumber string) string {
	re := regexp.MustCompile(`(\d{5})\d+(\d{2})`)
	matches := re.FindStringSubmatch(phoneNumber)
	if len(matches) > 2 {
		return phoneNumber[:len(matches[1])] + strings.Repeat("*", len(phoneNumber)-len(matches[1])-len(matches[2])) + phoneNumber[len(phoneNumber)-len(matches[2]):]
	}
	return re.ReplaceAllString(phoneNumber, "$1*****$2")
}
