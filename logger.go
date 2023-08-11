package logrus2slog

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	*Entry
}

func New(ctx context.Context, delegate *slog.Logger) *Logger {
	return &Logger{Entry: entryOf(ctx, delegate)}
}

func (logger *Logger) LogFn(level Level, fn LogFunction) {
	if logger.IsLevelEnabled(level) {
		logger.Log(level, fn()...)
	}
}

func (logger *Logger) TraceFn(fn LogFunction) {
	logger.LogFn(TraceLevel, fn)
}

func (logger *Logger) DebugFn(fn LogFunction) {
	logger.LogFn(DebugLevel, fn)
}

func (logger *Logger) InfoFn(fn LogFunction) {
	logger.LogFn(InfoLevel, fn)
}

func (logger *Logger) PrintFn(fn LogFunction) {
	logger.InfoFn(fn)
}

func (logger *Logger) WarnFn(fn LogFunction) {
	logger.LogFn(WarnLevel, fn)
}

func (logger *Logger) WarningFn(fn LogFunction) {
	logger.WarnFn(fn)
}

func (logger *Logger) ErrorFn(fn LogFunction) {
	logger.LogFn(ErrorLevel, fn)
}

func (logger *Logger) FatalFn(fn LogFunction) {
	logger.LogFn(FatalLevel, fn)
	os.Exit(1)
}

func (logger *Logger) PanicFn(fn LogFunction) {
	logger.LogFn(PanicLevel, fn)
	panic(fn)
}
