package utils

import (
	"reflect"

	"github.com/rizalarfiyan/be-revend/internal/sql"
)

func IsEmptyAll(i interface{}) bool {
	return (IsNil(i) || IsZeroLength(i) || IsEmpty(i))
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}

	value := reflect.ValueOf(i)
	kind := value.Kind()

	if kind == reflect.Ptr || kind == reflect.Interface {
		return value.IsNil()
	}

	return false
}

func IsZeroLength(i interface{}) bool {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Slice {
		return false
	}

	return value.Len() == 0
}

func IsEmpty(i interface{}) bool {
	return reflect.ValueOf(i).IsZero()
}

func IsValidRole(role string) bool {
	switch sql.Role(role) {
	case sql.RoleGuest, sql.RoleAdmin:
		return true
	default:
		return false
	}
}
