package logrus

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// Defines the key when adding errors using WithError.
var ErrorKey = "error"

// An entry is the final or intermediate Logrus logging entry.
type Entry struct {
	Logger  *slog.Logger
	Context context.Context
	Data    []slog.Attr
}

func NewEntry(logger *Logger) *Entry {
	return logger.Entry
}

func (entry *Entry) Slog() *slog.Logger {
	return entry.Logger
}

func entryOf(ctx context.Context, logger *slog.Logger) *Entry {
	if ctx == nil {
		ctx = context.Background()
	}

	return &Entry{Logger: logger, Context: ctx}
}

func (entry *Entry) Dup() *Entry {
	var dup = *entry
	return &dup
}

// Add a context to the Entry.
func (entry *Entry) withLogger(logger *slog.Logger) *Entry {
	dup := entry.Dup()
	dup.Logger = logger
	return dup
}

func (entry *Entry) WithError(err error) *Entry {
	return entry.WithField(ErrorKey, err)
}

// Add a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	attr := slog.Any(key, value)

	e := entry.withLogger(entry.Logger.With(attr))
	e.Data = append(e.Data, attr)
	return e
}

// Add a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	e := entry.Dup()

	var attrs []any
	for key, value := range fields {
		attr := slog.Any(key, value)
		e.Data = append(e.Data, attr)
		attrs = append(attrs, attr)
	}

	e.Logger = entry.Logger.With(attrs...)
	return e
}

// Add a context to the Entry.
func (entry *Entry) WithContext(ctx context.Context) *Entry {
	dup := entry.Dup()
	dup.Context = ctx
	return dup
}

func toString(args []interface{}) string {
	return fmt.Sprint(args...)
}

func (entry *Entry) Log(level Level, args ...interface{}) {
	entry.Logger.Log(entry.Context, level.ToSlog(), toString(args))
}

func (entry *Entry) Trace(args ...interface{}) {
	entry.Log(TraceLevel, args...)
}

func (entry *Entry) Debug(args ...interface{}) {
	entry.Log(DebugLevel, args...)
}

func (entry *Entry) Print(args ...interface{}) {
	entry.Info(args...)
}

func (entry *Entry) Info(args ...interface{}) {
	entry.Log(InfoLevel, args...)
}

func (entry *Entry) Warn(args ...interface{}) {
	entry.Log(WarnLevel, args...)
}

func (entry *Entry) Warning(args ...interface{}) {
	entry.Warn(args...)
}

func (entry *Entry) Error(args ...interface{}) {
	entry.Log(ErrorLevel, args...)
}

func (entry *Entry) Fatal(args ...interface{}) {
	entry.Log(FatalLevel, args...)
	os.Exit(1)
}

func (entry *Entry) Panic(args ...interface{}) {
	entry.Log(PanicLevel, args...)
	panic(args)
}

func (entry *Entry) IsLevelEnabled(level Level) bool {
	return entry.Logger.Enabled(entry.Context, level.ToSlog())
}

// Entry Printf family functions

func (entry *Entry) Logf(level Level, format string, args ...interface{}) {
	if entry.IsLevelEnabled(level) {
		entry.Logger.Log(entry.Context, level.ToSlog(), fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Tracef(format string, args ...interface{}) {
	entry.Logf(TraceLevel, format, args...)
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.Logf(DebugLevel, format, args...)
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.Logf(InfoLevel, format, args...)
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Infof(format, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.Logf(WarnLevel, format, args...)
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.Warnf(format, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.Logf(ErrorLevel, format, args...)
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.Logf(FatalLevel, format, args...)
	os.Exit(1)
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	entry.Logf(PanicLevel, format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Entry Println family functions

func (entry *Entry) Logln(level Level, args ...interface{}) {
	if entry.IsLevelEnabled(level) {
		entry.Log(level, sprintlnn(args...))
	}
}

func (entry *Entry) Traceln(args ...interface{}) {
	entry.Logln(TraceLevel, args...)
}

func (entry *Entry) Debugln(args ...interface{}) {
	entry.Logln(DebugLevel, args...)
}

func (entry *Entry) Infoln(args ...interface{}) {
	entry.Logln(InfoLevel, args...)
}

func (entry *Entry) Println(args ...interface{}) {
	entry.Infoln(args...)
}

func (entry *Entry) Warnln(args ...interface{}) {
	entry.Logln(WarnLevel, args...)
}

func (entry *Entry) Warningln(args ...interface{}) {
	entry.Warnln(args...)
}

func (entry *Entry) Errorln(args ...interface{}) {
	entry.Logln(ErrorLevel, args...)
}

func (entry *Entry) Fatalln(args ...interface{}) {
	entry.Logln(FatalLevel, args...)
	os.Exit(1)
}

func (entry *Entry) Panicln(args ...interface{}) {
	entry.Logln(PanicLevel, args...)
}

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
