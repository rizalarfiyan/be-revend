package utils

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/validation"
)

func DefaultErrorData(isList bool) interface{} {
	if isList {
		return []any{}
	}
	return nil
}

func PanicIfError(err error, isList bool) {
	if err != nil {
		data := DefaultErrorData(isList)
		panic(response.NewErrorMessage(http.StatusInternalServerError, err.Error(), data))
	}
}

func PanicIfErrorWithoutNoSqlResult(err error, isList bool) {
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		data := DefaultErrorData(isList)
		panic(response.NewErrorMessage(http.StatusInternalServerError, err.Error(), data))
	}
}

func responseIfNull(value interface{}, isList bool, callback func(isList bool, data interface{}) *response.BaseResponse) {
	if IsNil(value) || IsZeroLength(value) {
		data := DefaultErrorData(isList)
		panic(callback(isList, data))
	}
}

func IsNotFoundMessage(value interface{}, message string, isList bool) {
	responseIfNull(value, isList, func(isList bool, data interface{}) *response.BaseResponse {
		if isList {
			return response.NewErrorMessage(http.StatusOK, message, data)
		}
		return response.NewErrorMessage(http.StatusNotFound, message, data)
	})
}

func IsNotFound(value interface{}, isList bool) {
	IsNotFoundMessage(value, "Data not found", isList)
}

func ValidateStruct(dataSet interface{}, isList bool) {
	err := validation.Get().Struct(dataSet)
	if err == nil {
		return
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		validate := make(map[string]string, len(ve))
		reflected := reflect.ValueOf(dataSet)
		for _, fe := range ve {
			field, _ := reflected.Type().FieldByName(fe.StructField())

			tagName := field.Tag.Get("json")
			if tagName == "" {
				tagName = field.Name
			}

			tagField := field.Tag.Get("field")
			if tagField == "" {
				tagField = tagName
			}

			searchTag := fmt.Sprintf("{{%s}}", tagName)
			message := fe.Translate(validation.GetTranslator())
			validate[tagName] = strings.Replace(message, searchTag, tagField, -1)
		}
		panic(response.NewErrorMessage(http.StatusBadRequest, "Error Validation", validate))
	}

	data := DefaultErrorData(isList)
	panic(response.NewErrorMessage(http.StatusUnprocessableEntity, err.Error(), data))
}

func ErrorManualValidation(key, message string) {
	panic(response.NewErrorMessage(http.StatusBadRequest, "Error Validation", fiber.Map{
		key: message,
	}))
}

func ErrorManualValidationErr(err error, key, message string) {
	if err != nil {
		ErrorManualValidation(key, message)
	}
}

func IsNotProcessMessage(value interface{}, message string, isList bool) {
	responseIfNull(value, isList, func(isList bool, data interface{}) *response.BaseResponse {
		return response.NewErrorMessage(http.StatusUnprocessableEntity, message, data)
	})
}

func IsNotProcess(value interface{}, isList bool) {
	IsNotFoundMessage(value, "Something wrong, please try again later!", isList)
}

func IsNotProcessError(err error, isList bool) {
	if err != nil {
		IsNotProcess(nil, isList)
	}
}

func IsNotProcessErrorMessage(err error, message string, isList bool) {
	if err != nil {
		IsNotProcessMessage(nil, message, isList)
	}
}

func IsNotProcessRawMessage(message string, isList bool) {
	data := DefaultErrorData(isList)
	panic(response.NewErrorMessage(http.StatusUnprocessableEntity, message, data))
}

func IsNotProcessData(message string, data interface{}) {
	panic(response.NewErrorMessage(http.StatusUnprocessableEntity, message, data))
}
