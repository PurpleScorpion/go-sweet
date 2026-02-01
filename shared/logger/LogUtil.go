package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *slog.Logger

func LogInit() {
	levelStr := keqing.ValueString("${sweet.logging.level}")
	level := parseLevel(levelStr)

	adapters := keqing.ValueStringArr("${sweet.logging.adapters}")
	if len(adapters) == 0 {
		adapters = []string{"console", "file"}
	}

	var handlers []slog.Handler

	// console
	if keqing.ArrayContains(adapters, "console") {
		h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
		handlers = append(handlers, h)
	}

	// file
	if keqing.ArrayContains(adapters, "file") {
		file := keqing.ValueString("${sweet.logging.file}")
		if keqing.IsEmpty(file) {
			file = "log/go_sweet.log"
		}

		maxSize := keqing.ValueInt("${sweet.logging.maxSize}")
		if maxSize <= 0 || maxSize > 1024 {
			maxSize = 10
		}
		maxDays := keqing.ValueInt("${sweet.logging.maxDays}")
		if maxDays <= 0 {
			maxDays = 7
		}
		maxBackups := keqing.ValueInt("${sweet.logging.maxBackups}")
		if maxBackups <= 0 {
			maxBackups = 10
		}

		// ⚠ slog 不内置滚动，需要配合 lumberjack
		writer := &lumberjack.Logger{
			Filename:   file,
			MaxSize:    maxSize,
			MaxAge:     maxDays,
			MaxBackups: maxBackups,
			Compress:   false,
		}

		h := slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: level,
		})
		handlers = append(handlers, h)
	}

	// fanout
	log = slog.New(multiHandler(handlers...))

	// 设为默认 logger
	slog.SetDefault(log)
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type fanoutHandler struct {
	handlers []slog.Handler
}

func multiHandler(h ...slog.Handler) slog.Handler {
	return &fanoutHandler{handlers: h}
}

func (f *fanoutHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range f.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (f *fanoutHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range f.handlers {
		_ = h.Handle(ctx, r)
	}
	return nil
}

func (f *fanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var hs []slog.Handler
	for _, h := range f.handlers {
		hs = append(hs, h.WithAttrs(attrs))
	}
	return &fanoutHandler{handlers: hs}
}

func (f *fanoutHandler) WithGroup(name string) slog.Handler {
	var hs []slog.Handler
	for _, h := range f.handlers {
		hs = append(hs, h.WithGroup(name))
	}
	return &fanoutHandler{handlers: hs}
}

// formatMessage 智能格式化消息，支持多种格式
func formatMessage(args ...any) string {
	if len(args) == 0 {
		return ""
	}

	// 只有一个参数
	if len(args) == 1 {
		return valueToString(args[0])
	}

	// 第一个参数是 string
	if format, ok := args[0].(string); ok {

		// Python 风格 {}
		if strings.Contains(format, "{}") {
			return formatPythonStyle(format, args[1:]...)
		}

		// fmt 风格 %
		if strings.Contains(format, "%") {
			return fmt.Sprintf(format, args[1:]...)
		}

		// 普通拼接
		var b strings.Builder
		b.WriteString(format)

		for _, arg := range args[1:] {
			b.WriteByte(' ')
			b.WriteString(valueToString(arg))
		}
		return b.String()
	}

	// 第一个不是 string
	var b strings.Builder
	for i, arg := range args {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(valueToString(arg))
	}
	return b.String()
}

// formatPythonStyle 实现Python风格的{}占位符格式化
func formatPythonStyle(format string, args ...any) string {
	result := format

	for _, arg := range args {
		if !strings.Contains(result, "{}") {
			break
		}
		result = strings.Replace(result, "{}", valueToString(arg), 1)
	}

	return result
}

func valueToString(v any) string {
	if v == nil {
		return "null"
	}

	rv := reflect.ValueOf(v)

	// 解引用指针
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "null"
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct,
		reflect.Map,
		reflect.Slice,
		reflect.Array:
		if b, err := json.Marshal(v); err == nil {
			return string(b)
		}
	}

	return fmt.Sprint(v)
}

func Info(args ...interface{}) {
	log.Info(formatMessage(args...))
}

func Warn(args ...interface{}) {
	log.Warn(formatMessage(args...))
}

func Error(args ...interface{}) {
	log.Error(formatMessage(args...))
}

func Debug(args ...interface{}) {
	log.Debug(formatMessage(args...))
}
