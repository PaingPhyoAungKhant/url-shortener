// Package logger provides slog based structured logger
// loggers are constructed with New
// and must be constructed explicitly into every struct and handler that need one
// no global logger variable is used
// Context is used only to carry request scoped trace IDs and not logger
// Logging methods takes in Context as argument to to get trace IDs
package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type Logger struct {
	sl *slog.Logger
}

type ctxKey string

const (
	traceIDKey ctxKey = "TRACE_ID_KEY"
)

// String returns a Attr for a string value
func String(key, value string) slog.Attr {
	return slog.String(key, value)
}

// Int returns a Attr for a int value
func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

// Int64 returns a Attr for a int64 value
func Int64(key string, value int64) slog.Attr {
	return slog.Int64(key, value)
}

// Float64 returns a Attr for a int64 value
func Float64(key string, value float64) slog.Attr {
	return slog.Float64(key, value)
}

// Bool returns a Attr for a bool value
func Bool(key string, value bool) slog.Attr {
	return slog.Bool(key, value)
}

// Time returns a Attr for a time.Time value
func Time(key string, value time.Time) slog.Attr {
	return slog.Time(key, value)
}

// Duration returns a Attr for a time.Duration value
func Duration(key string, value time.Duration) slog.Attr {
	return slog.Duration(key, value)
}

// Err returns a Field with an error value
func Err(err any) slog.Attr {
	return slog.Any("error", err)
}

// New constructs Logger with the given level and format
// level must be one of "debug", "info", "warn", "error"
// format must be one of "json", "text"
func New(level, format string) (*Logger, error) {
	var slogLevel slog.Level

	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		return nil, fmt.Errorf("logger: unknown level %q; want debug|info|warn|error", level)
	}

	var handler slog.Handler
	w := os.Stderr
	handlerOptions := &slog.HandlerOptions{
		Level: slogLevel,
	}

	switch format {
	case "json":
		handler = slog.NewJSONHandler(w, handlerOptions)
	case "text":
		handler = slog.NewTextHandler(w, handlerOptions)
	default:
		return nil, fmt.Errorf("logger: unknown format %q; want json|text", format)
	}

	sl := slog.New(handler)
	return &Logger{sl: sl}, nil
}

// WithTraceID attaches the traceID to ctx so that
// all log call after this that recives the context will include it immediately.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// TraceIDFromCtx return the the traceID from context
// that is stored by WithTraceID
// or return an empty string "" if traceID was not set
func TraceIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(traceIDKey).(string); ok {
		return v
	}
	return ""
}

// addTraceID get the traceID from the context and add it log call
func addTraceID(ctx context.Context, args []slog.Attr) []slog.Attr {
	if traceID := TraceIDFromCtx(ctx); traceID != "" {
		args = append(args, slog.String("trace_id", traceID))
	}
	return args
}

// Debug log a message at Debug level
func (l *Logger) Debug(ctx context.Context, msg string, fields ...slog.Attr) {
	l.sl.LogAttrs(ctx, slog.LevelDebug, msg, addTraceID(ctx, fields)...)
}

// Info log a message at Info level
func (l *Logger) Info(ctx context.Context, msg string, fields ...slog.Attr) {
	l.sl.LogAttrs(ctx, slog.LevelInfo, msg, addTraceID(ctx, fields)...)
}

// Warn log a message at Warn level
func (l *Logger) Warn(ctx context.Context, msg string, fields ...slog.Attr) {
	l.sl.LogAttrs(ctx, slog.LevelWarn, msg, addTraceID(ctx, fields)...)
}

// Error log a message at Error level
func (l *Logger) Error(ctx context.Context, msg string, fields ...slog.Attr) {
	l.sl.LogAttrs(ctx, slog.LevelError, msg, addTraceID(ctx, fields)...)
}
