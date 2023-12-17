package logger

import "github.com/rs/zerolog"

type Logger interface {
	Logger() zerolog.Logger
	Logs() *zerolog.Logger
	Debug(msg string)
	Debugf(format string, v ...interface{})
	Info(msg string)
	Infof(format string, v ...interface{})
	Warn(msg string)
	Warnf(format string, v ...interface{})
	Error(err error, msg string)
	Errorf(err error, format string, v ...interface{})
	Fatal(err error, msg string)
	Fatalf(err error, format string, v ...interface{})
	Panic(err error, msg string)
	Panicf(err error, format string, v ...interface{})
}

type logger struct {
	log zerolog.Logger
}

func Get(types string) Logger {
	return &logger{
		log: log.With().Str("type", types).Logger(),
	}
}

func (l *logger) Logger() zerolog.Logger {
	return l.log
}

func (l *logger) Logs() *zerolog.Logger {
	return &l.log
}

func (l *logger) Debug(msg string) {
	l.log.Debug().Msg(msg)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.log.Debug().Msgf(format, v...)
}

func (l *logger) Info(msg string) {
	l.log.Info().Msg(msg)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.log.Info().Msgf(format, v...)
}

func (l *logger) Warn(msg string) {
	l.log.Warn().Msg(msg)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.log.Warn().Msgf(format, v...)
}

func (l *logger) Error(err error, msg string) {
	l.log.Err(err).Msg(msg)
}

func (l *logger) Errorf(err error, format string, v ...interface{}) {
	l.log.Err(err).Msgf(format, v...)
}

func (l *logger) Fatal(err error, msg string) {
	l.log.Fatal().Err(err).Msg(msg)
}

func (l *logger) Fatalf(err error, format string, v ...interface{}) {
	l.log.Fatal().Err(err).Msgf(format, v...)
}

func (l *logger) Panic(err error, msg string) {
	l.log.Panic().Err(err).Msg(msg)
}

func (l *logger) Panicf(err error, format string, v ...interface{}) {
	l.log.Panic().Err(err).Msgf(format, v...)
}
