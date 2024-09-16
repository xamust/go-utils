package logger

import (
	"fmt"
	"net/http"

	"github.com/xamust/go-utils/errors"
	"github.com/xamust/go-utils/metadata"

	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	giterrors "github.com/pkg/errors"
)

var (
	missHeader = "miss header"
)

type stackTracer interface {
	StackTrace() giterrors.StackTrace
}

func InjectLogger(next http.HandlerFunc, log Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		h, ok := metadata.FromContextHeader(ctx)
		if !ok {
			http.Error(rw, missHeader, http.StatusBadRequest)
			return
		}

		var fields = make(map[string]any)
		fields[keyUID] = h.Header.RqUid
		fields[keyEventDate] = h.Header.RqTm
		fields[keyService] = h.Header.Service

		ctx = NewContextLogger(ctx, log.Fields(fields))
		r = r.WithContext(ctx)

		next(rw, r)
	}
}

func InjectLoggerMux(log Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return InjectLogger(next.ServeHTTP, log)
	}
}

func InjectLoggerEcho(log Logger) echo.MiddlewareFunc {
	if log == nil {
		log = DefaultLogger
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			ctx := req.Context()

			h, ok := metadata.FromContextHeader(ctx)
			if !ok {
				log.Error(ctx, missHeader)
				return errors.NewBadRequestErrorRsp(missHeader)
			}

			var fields = make(map[string]any)

			fields[keyUID] = h.Header.RqUid
			fields[keyEventDate] = h.Header.RqTm
			fields[keyService] = h.Header.Service

			cLog := log.Fields(fields)

			ctx = NewContextLogger(ctx, cLog)
			c.SetRequest(req.WithContext(ctx))

			if err = next(c); err != nil {
				msg := err.Error()
				var args []any

				if e, ok := err.(*echo.HTTPError); ok {
					if m, ok := e.Message.(errors.ErrorRsp); ok {
						msg = m.Error()
					}

					if stErr, ok := e.Internal.(stackTracer); ok {
						args = append(args, StackTrace(fmt.Sprintf("error stacktrace: %+v", stErr.StackTrace())))
					}
				}

				ctx := NewContextEvent(ctx, AppSystem, h.Header.SourceSystem)
				cLog.Error(ctx, msg, args...)
			}
			return err
		}
	}
}
