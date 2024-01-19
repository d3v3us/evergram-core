package logging

import (
	"os"

	"log/slog"

	"github.com/deveusss/evergram-core/common"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewDefaultStdOutTextLogger() *SlogLogger {
	return NewStdOutTextLogger(false, false)
}
func NewStdOutTextLogger(enableSource bool, enableDebug bool) *SlogLogger {
	opt := &slog.HandlerOptions{
		AddSource: enableSource,
		Level:     common.When[slog.Level](enableDebug).Then(slog.LevelDebug).Else(slog.LevelInfo),
	}
	logHandler := slog.NewTextHandler(os.Stdout, opt)
	logger := slog.New(logHandler)
	return &SlogLogger{
		logger: logger,
	}
}
