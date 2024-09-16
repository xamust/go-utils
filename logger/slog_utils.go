package logger

import (
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/exp/slog"
)

var (
	valDebug = slog.StringValue("debug")
	valInfo  = slog.StringValue("info")
	valWarn  = slog.StringValue("warn")
	valErr   = slog.StringValue("error")
	valFatal = slog.StringValue("fatal")

	valUppSuc = slog.StringValue(StatusSuccess)
	valUppErr = slog.StringValue(StatusError)
)

type ReplacerAttribute func([]string, slog.Attr) slog.Attr

func Unpack(funcs ...ReplacerAttribute) ReplacerAttribute {
	return func(args []string, attr slog.Attr) slog.Attr {
		for _, f := range funcs {
			attr = f(args, attr)
		}
		return attr
	}
}

func renameTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		a.Key = "@timestamp"
	}

	return a
}

func renameMessage(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.MessageKey {
		a.Key = "message"
	}
	return a
}

func (s *slogLogger) renameAttr(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.SourceKey:
		source := a.Value.Any().(*slog.Source)
		a.Value = slog.StringValue(source.File + ":" + strconv.Itoa(source.Line))
		a.Key = s.opts.keys.sourceKey
	case slog.TimeKey:
		a.Key = s.opts.keys.timeKey
	case slog.MessageKey:
		a.Key = s.opts.keys.messageKey
	case slog.LevelKey:
		a.Key = s.opts.keys.levelKey
		lvl := slogToLoggerLevel(a.Value.Any().(slog.Level))
		switch {
		case lvl < InfoLevel:
			a.Value = valDebug
		case lvl < WarnLevel:
			a.Value = valInfo
		case lvl < ErrorLevel:
			a.Value = valWarn
		case lvl < FatalLevel:
			a.Value = valErr
		case lvl >= FatalLevel:
			a.Value = valFatal
		default:
			a.Value = valInfo
		}
	}

	return a
}

func fieldsToAttr(m map[string]any) []slog.Attr {
	data := make([]slog.Attr, 0, len(m))
	for k, v := range m {
		data = append(data, slog.Any(k, v))
	}

	return data
}

func loggerToSlogLevel(level Level) slog.Level {
	switch level {
	case TraceLevel, DebugLevel:
		return slog.LevelDebug
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	case FatalLevel:
		return slog.LevelError + 1
	default:
		return slog.LevelInfo
	}
}

func slogToLoggerLevel(level slog.Level) Level {
	switch level {
	case slog.LevelDebug:
		return DebugLevel
	case slog.LevelWarn:
		return WarnLevel
	case slog.LevelError:
		return ErrorLevel
	case slog.LevelError + 1:
		return FatalLevel
	default:
		return InfoLevel
	}
}

func successStatus() slog.Attr {
	return slog.Any(keyStatus, StatusSuccess)
}
func errorStatus() slog.Attr {
	return slog.Any(keyStatus, StatusError)
}

func errorText(msg string) slog.Attr {
	return slog.Any("ErrorText", msg)
}

func getProcKey(lvl Level) slog.Attr {
	if lvl < ErrorLevel {
		return slog.Any(keyStatus, valUppSuc)
	}
	return slog.Any(keyStatus, valUppErr)
}

func Operation(msg string) slog.Attr {
	return slog.Any(keyOperation, msg)
}

func StackTrace(msg string) slog.Attr {
	return slog.Any(keyStackTrace, msg)
}

func StatusCode(code int) slog.Attr {
	return slog.Any(keyStatusCode, code)
}

func HttpHeaders(h http.Header) slog.Attr {
	return slog.Any(keyHTTPHeaders, fmt.Sprintf("%+v", h))
}

func RequestID(uuid string) slog.Attr {
	return slog.Any(keyRequestID, uuid)
}
