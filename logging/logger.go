package logging

import (
	"os"

	"log/slog"

	"github.com/deveusss/evergram-core/common"
)

func NewDefaultStdOutTextLogger() *slog.Logger {
	return NewStdOutTextLogger(false, false)
}
func NewStdOutTextLogger(enableSource bool, enableDebug bool) *slog.Logger {
	opt := &slog.HandlerOptions{
		AddSource: enableSource,
		Level:     common.When[slog.Level](enableDebug).Then(slog.LevelDebug).Else(slog.LevelInfo),
	}
	logHandler := slog.NewTextHandler(os.Stdout, opt)
	logger := slog.New(logHandler)
	return logger
}
