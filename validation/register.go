package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/rizalarfiyan/be-revend/constants"
)

type reValidation struct{}

func (re reValidation) SetTagName() {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			return field.Name
		}

		if name == "-" {
			return ""
		}

		return fmt.Sprintf("{{%s}}", name)
	})
}

func (re reValidation) SetPhoneNumber() {
	validate.RegisterValidation("phone_number", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(constants.RegexPhoneNumber).MatchString(fl.Field().String())
	})

	validate.RegisterTranslation("phone_number", trans, func(ut ut.Translator) error {
		return ut.Add("phone_number", "{0} invalid indonesian phone number format.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone_number", fe.Field())
		return t
	})
}
