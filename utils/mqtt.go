package utils

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type MQTTUtils interface {
	PanicIfError(err error)
	PanicIfErrorWithoutNoSqlResult(err error)
	Recover()
	Skip(message string)
	Skipf(format string, a ...any)
}

type mqttUtils struct{}

func NewMqttUtils() MQTTUtils {
	return &mqttUtils{}
}

func (u *mqttUtils) PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func (u *mqttUtils) PanicIfErrorWithoutNoSqlResult(err error) {
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		panic(err)
	}
}

func (u *mqttUtils) Recover() {
	if r := recover(); r != nil {
		fmt.Println("Recovered in f", r)
	}
}

func (u *mqttUtils) Skip(message string) {
	panic(errors.New(message))
}

func (u *mqttUtils) Skipf(format string, a ...any) {
	panic(fmt.Errorf(format, a...))
}
