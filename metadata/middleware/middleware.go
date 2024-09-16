package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/xamust/go-utils/encoder/json"
	"github.com/xamust/go-utils/errors"
	"github.com/xamust/go-utils/logger"
	"github.com/xamust/go-utils/metadata"
	"github.com/xamust/go-utils/metadata/request_id"
	"github.com/xamust/go-utils/util/slice"
	"github.com/xamust/go-utils/validate"

	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
)

var (
	encodingError = "error encoding header request: %v"
	validateError = "error validate request: %v"
)

func ExtractHeader(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContextLogger(ctx)

		var h metadata.HeaderSource
		codec := json.NewCodec()
		reader, err := codec.ReadHeader(r.Body, &h)
		if err != nil {
			log.Error(ctx, fmt.Sprintf(encodingError, err))
			http.Error(rw, fmt.Sprintf(encodingError, err), http.StatusBadRequest)
			return
		}

		if err := validate.Validate(h); err != nil {
			log.Error(ctx, fmt.Sprintf(validateError, err))
			http.Error(rw, fmt.Sprintf(validateError, err), http.StatusBadRequest)
			return
		}

		ctx = metadata.NewContextHeader(ctx, h)
		r = r.WithContext(ctx)
		r.Body = io.NopCloser(reader)

		r.Header.Set(echo.HeaderXRequestID, h.Header.RqUid)

		next(rw, r)
	}
}

func ExtractHeaderMux() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return ExtractHeader(next.ServeHTTP)
	}
}

func ExtractHeaderEcho() echo.MiddlewareFunc {
	return ParsingHeaderEcho(&metadata.HeaderSource{})
}

func ParsingHeaderEcho(parser metadata.ParserHeaderSource) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			if slice.StringContains(req.URL.String(), []string{"/ready", "/live", "/swagger", "/health", "/healthz"}) {
				return next(c)
			}

			ctx := req.Context()
			log := logger.FromContextLogger(ctx)

			all, err := io.ReadAll(req.Body)
			if err != nil {
				log.Error(ctx, err.Error())
				return errors.NewInternalErrorRsp(err.Error())
			}
			defer func(r io.Closer) {
				_ = r.Close()
			}(req.Body)

			h, err := parser.ExtractHeaderFromBytes(all)
			if err != nil {
				errMsg := fmt.Sprintf(encodingError, err)
				log.Error(ctx, errMsg)
				return errors.NewBadRequestErrorRsp(errMsg)
			}

			reqID := h.RequestID()
			c.Response().Header().Set(echo.HeaderXRequestID, reqID)

			ctx = metadata.NewContextHeader(ctx, h)
			ctx = request_id.NewContext(ctx, reqID)

			if err := parser.Validate(h); err != nil {
				errMsg := fmt.Sprintf(validateError, err)
				log.Error(ctx, errMsg)
				return errors.NewBadRequestErrorRsp(errMsg)
			}

			req = req.WithContext(ctx)
			req.Header.Set(echo.HeaderXRequestID, reqID)
			req.Body = io.NopCloser(bytes.NewReader(all))

			c.SetRequest(req)

			return next(c)
		}
	}
}
