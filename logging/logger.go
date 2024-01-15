package logging

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(str string)
	Info(str string)
	Warn(str string)
	Error(str string)
	Fatal(str string)
	Panic(str string)
	Trace(str string)
	DebugStruct(message string, fields ...Field)
	InfoStruct(message string, fields ...Field)
	WarnStruct(message string, fields ...Field)
	ErrorStruct(message string, fields ...Field)
	FatalStruct(message string, fields ...Field)
	PanicStruct(message string, fields ...Field)
	TraceStruct(message string, fields ...Field)
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Tracef(format string, args ...interface{})
	Log(level string, message string, fields ...Field)
	EnableLogging(enabled bool, options LogEntryFormattingOptions)
}

type ZerologLogger struct {
	logger zerolog.Logger
}

func NewZerologLogger() *ZerologLogger {
	return &ZerologLogger{
		logger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
	}
}

func (zl *ZerologLogger) Debug(str string) {
	zl.logger.Debug().Msg(str)
}

func (zl *ZerologLogger) Info(str string) {
	zl.logger.Info().Msg(str)
}

func (zl *ZerologLogger) Warn(str string) {
	zl.logger.Warn().Msg(str)
}

func (zl *ZerologLogger) Error(str string) {
	zl.logger.Error().Msg(str)
}

func (zl *ZerologLogger) Fatal(str string) {
	zl.logger.Fatal().Msg(str)
}

func (zl *ZerologLogger) Panic(str string) {
	zl.logger.Panic().Msg(str)
}

func (zl *ZerologLogger) Trace(str string) {
	zl.logger.Trace().Msg(str)
}

func (zl *ZerologLogger) DebugStruct(message string, fields ...Field) {
	zl.logger.Debug().Fields(convertFields(fields)).Msg(message)
}

func (zl *ZerologLogger) InfoStruct(message string, fields ...Field) {
	zl.logger.Info().Fields(convertFields(fields)).Msg(message)
}

func (zl *ZerologLogger) WarnStruct(message string, fields ...Field) {
	zl.logger.Warn().Fields(convertFields(fields)).Msg(message)
}

func (zl *ZerologLogger) ErrorStruct(message string, fields ...Field) {
	zl.logger.Error().Fields(convertFields(fields)).Msg(message)
}

func (zl *ZerologLogger) FatalStruct(message string, fields ...Field) {
	zl.logger.Fatal().Fields(convertFields(fields)).Msg(message)
}

func (zl *ZerologLogger) PanicStruct(message string, fields ...Field) {
	zl.logger.Panic().Fields(convertFields(fields)).Msg(message)
}

func (zl *ZerologLogger) TraceStruct(message string, fields ...Field) {
	zl.logger.Trace().Fields(convertFields(fields)).Msg(message)
}

func (zl *ZerologLogger) Debugf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	zl.Debug(message)
}

func (zl *ZerologLogger) Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	zl.Info(message)
}

func (zl *ZerologLogger) Warnf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	zl.Warn(message)
}

func (zl *ZerologLogger) Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	zl.Error(message)
}

func (zl *ZerologLogger) Fatalf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	zl.Fatal(message)
}

func (zl *ZerologLogger) Panicf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	zl.Panic(message)
}

func (zl *ZerologLogger) Tracef(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	zl.Trace(message)
}

func (zl *ZerologLogger) Log(level string, message string, fields ...Field) {
	switch strings.ToLower(level) {
	case "debug":
		zl.DebugStruct(message, fields...)
	case "info":
		zl.InfoStruct(message, fields...)
	case "warn", "warning":
		zl.WarnStruct(message, fields...)
	case "error":
		zl.ErrorStruct(message, fields...)
	case "fatal":
		zl.FatalStruct(message, fields...)
	case "panic":
		zl.PanicStruct(message, fields...)
	case "trace":
		zl.TraceStruct(message, fields...)
	default:
		zl.InfoStruct(message, fields...)
	}
}

func (zl *ZerologLogger) EnableLogging(enabled bool, options LogEntryFormattingOptions) {
	if enabled {
		zl.logger = configureLoggingEnhancements(zl.logger, options)
	} else {
		zl.logger = zerolog.New(zerolog.Nop()).With().Timestamp().Logger()
	}
}

func convertFields(fields []Field) map[string]interface{} {
	result := make(map[string]interface{})
	for _, field := range fields {
		result[field.Key] = field.Value
	}
	return result
}

type LogEntryFormattingOptions struct {
	DisableTimestamp bool
	DisableLevel     bool
}

type LogFormatter interface {
	Format(e *zerolog.Event, level zerolog.Level, message string) *zerolog.Event
}

type CustomLogHook struct {
	Formatter LogFormatter
}

type Field struct {
	Key   string
	Value interface{}
}

func configureLoggingEnhancements(logger zerolog.Logger, options LogEntryFormattingOptions) zerolog.Logger {
	if !options.DisableTimestamp {
		logger = logger.With().Timestamp().Logger()
	}

	if options.DisableLevel {
		logger = logger.Level(zerolog.NoLevel)
	}

	return logger
}

type LoggerFactory interface {
	CreateLogger() Logger
}

type ConfigurableLoggerFactory struct {
	Enabled                   bool
	LogLevel                  zerolog.Level
	LogEntryFormattingOptions LogEntryFormattingOptions
	CustomLogFormatter        LogFormatter
}

func NewConfigurableLoggerFactory() *ConfigurableLoggerFactory {
	return &ConfigurableLoggerFactory{}
}

func (clf *ConfigurableLoggerFactory) SetEnabled(enabled bool) {
	clf.Enabled = enabled
}

func (clf *ConfigurableLoggerFactory) SetLogLevel(level zerolog.Level) {
	clf.LogLevel = level
}

func (clf *ConfigurableLoggerFactory) SetLogEntryFormattingOptions(options LogEntryFormattingOptions) {
	clf.LogEntryFormattingOptions = options
}

func (clf *ConfigurableLoggerFactory) SetCustomLogFormatter(formatter LogFormatter) {
	clf.CustomLogFormatter = formatter
}

func (clf *ConfigurableLoggerFactory) CreateLogger() Logger {
	logger := NewZerologLogger()
	logger.EnableLogging(clf.Enabled, clf.LogEntryFormattingOptions)
	return logger
}
