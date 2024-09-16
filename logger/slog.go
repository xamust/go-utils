package logger

import (
	"context"
	"fmt"
	"github.com/xamust/go-utils/logger/custom_handler"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/xamust/go-utils/metadata/request_id"

	"golang.org/x/exp/slog"
)

const (
	slogName = "slog"
)

type slogLogger struct {
	slog *slog.Logger
	opts Options

	sync.RWMutex
}

func NewSlogLogger(o ...Option) Logger {
	l := &slogLogger{
		opts: NewOptions(o...),
	}

	return l
}

func (s *slogLogger) Init(opts ...Option) error {
	for _, o := range opts {
		o(&s.opts)
	}

	if logger, ok := s.opts.context.Value(logKey{}).(*slog.Logger); ok {
		s.slog = logger
		return nil
	}

	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: s.renameAttr,
		Level:       loggerToSlogLevel(s.opts.lvl),
		AddSource:   s.opts.addSource,
	}

	attr := fieldsToAttr(s.opts.fields)

	handler := custom_handler.NewJSONHandler(s.opts.out, handleOpt).WithAttrs(attr)

	s.slog = slog.New(handler)
	return nil
}

func (s *slogLogger) Log(ctx context.Context, lvl Level, msg string, args ...any) {
	slvl := loggerToSlogLevel(lvl)

	if !s.slog.Enabled(ctx, slvl) {
		return
	}

	var pc uintptr
	if s.opts.addSource {
		var pcs [1]uintptr
		runtime.Callers(s.opts.callerSkipCount, pcs[:])
		pc = pcs[0]
	}

	if len(msg) > s.opts.maxBytesMessage {
		msg = fmt.Sprintf("too much content: have [%d], expected [<%d]", len(msg), s.opts.maxBytesMessage)
	}

	r := slog.NewRecord(time.Now(), slvl, msg, pc)

	r.Add(getProcKey(lvl))
	r.Add(args...)

	if uuid, ok := request_id.FromContext(ctx); ok {
		r.Add(RequestID(uuid))
	}

	if ctx == nil {
		ctx = context.Background()
	} else {
		r.AddAttrs(getEvent(ctx)...)
	}

	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Info(ctx context.Context, msg string, args ...any) {
	s.Log(ctx, InfoLevel, msg, args...)
}

func (s *slogLogger) Debug(ctx context.Context, msg string, args ...any) {
	s.Log(ctx, DebugLevel, msg, args...)
}

func (s *slogLogger) Warn(ctx context.Context, msg string, args ...any) {
	args = append(args, errorText(msg))
	s.Log(ctx, WarnLevel, "", args...)
}

func (s *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	args = append(args, errorText(msg))
	s.Log(ctx, ErrorLevel, "", args...)
}

func (s *slogLogger) Fatal(ctx context.Context, msg string, args ...any) {
	args = append(args, errorText(msg))
	DefaultLogger.Log(ctx, FatalLevel, "", args...)
	os.Exit(1)
}

func (s *slogLogger) String() string {
	return slogName
}

func (s *slogLogger) Fields(fields map[string]interface{}) Logger {
	attr := make([]any, 0, len(fields))

	for k, v := range fields {
		if k == keyUID || k == keyOperation {
			attr = append([]any{slog.Any(k, v)}, attr...)
			continue
		}
		attr = append(attr, slog.Any(k, v))
	}

	cloneLog := new(slogLogger)
	s.Lock()
	cloneLog.slog = s.slog.With(attr...)
	cloneLog.opts = s.opts
	s.Unlock()

	return cloneLog
}
