package custom_handler

import (
	"context"
	"golang.org/x/exp/slog"
	"io"
)

type JSONHandler struct {
	*commonHandler
}

func (h *JSONHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.enabled(level)
}

func (h *JSONHandler) Handle(_ context.Context, r slog.Record) error {
	return h.handle(r)
}

func (h *JSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &JSONHandler{commonHandler: h.commonHandler.withAttrs(attrs)}
}

func (h *JSONHandler) WithGroup(name string) slog.Handler {
	return &JSONHandler{commonHandler: h.commonHandler.withGroup(name)}
}

func NewJSONHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return &JSONHandler{
		commonHandler: &commonHandler{
			w:    w,
			opts: *opts,
		},
	}
}
