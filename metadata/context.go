package metadata

import (
	"context"
	"runtime"
	"strings"
	"time"
)

type headerKey struct{}
type mockTime struct{}

func NewContextHeader(ctx context.Context, h HeaderSource) context.Context {
	return context.WithValue(ctx, headerKey{}, &h)
}

func NewContextMockTime(ctx context.Context, f func() time.Time) context.Context {
	return context.WithValue(ctx, mockTime{}, f)
}

func FromContextHeader(ctx context.Context) (HeaderSource, bool) {
	if ctx == nil {
		return HeaderSource{}, false
	}
	h, ok := ctx.Value(headerKey{}).(*HeaderSource)
	if !ok || h == nil {
		return HeaderSource{}, ok
	}
	return *h, ok
}

func SetTimeFormatContextHeader(ctx context.Context, layout string) bool {
	if ctx == nil {
		return false
	}

	fTime, ok := ctx.Value(mockTime{}).(func() time.Time)
	if !ok {
		fTime = time.Now
	}

	h, ok := ctx.Value(headerKey{}).(*HeaderSource)
	if ok {
		h.Header.RsTm = fTime().Local().Format(layout)
		return true
	}
	return false
}

func SetTimeContextHeader(ctx context.Context) bool {
	return SetTimeFormatContextHeader(ctx, time.RFC3339)
}

func SetRsUidContextHeader(ctx context.Context, uid string) bool {
	if ctx == nil {
		return false
	}

	h, ok := ctx.Value(headerKey{}).(*HeaderSource)
	if ok {
		h.Header.RsUid = uid
		return true
	}
	return false
}

func CopyRqToRsUid(ctx context.Context) bool {
	if ctx == nil {
		return false
	}

	h, ok := ctx.Value(headerKey{}).(*HeaderSource)
	if ok {
		h.Header.RsUid = h.Header.RqUid
		return true
	}
	return false
}

func SetServiceCaller(ctx context.Context) bool {
	if ctx == nil {
		return false
	}

	h, ok := ctx.Value(headerKey{}).(*HeaderSource)
	if !ok || len(h.Header.Service) != 0 {
		return false
	}

	pc, _, _, _ := runtime.Caller(1)
	fullMethod := runtime.FuncForPC(pc).Name()

	parts := strings.Split(fullMethod, ".")
	nameFunc := ""
	if len(parts) >= 2 {
		nameFunc = parts[len(parts)-1]
	}

	h.Header.Service = nameFunc

	return true
}
