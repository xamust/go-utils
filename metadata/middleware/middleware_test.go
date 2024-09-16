package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/xamust/go-utils/errors"
	"github.com/xamust/go-utils/logger"
	"github.com/xamust/go-utils/metadata/request_id"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var body = []byte(`
{
    "header": {
        "operUid": "{{ReqUID}}",
        "rqUid": "{{ReqUID}}",
        "rqTm": "{{$isoTimestamp}}",
        "service": "CashMGAcctSumm",
        "receiverSystem": "CashManagement",
        "sourceSystem": "QPragma"
    }
}
`)

func Test_ReadBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/test", bytes.NewReader(body))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := ExtractHeader(testHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, body, rr.Body.Bytes())
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	all, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(all)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Test_ErrorLogging(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := logger.NewSlogLogger(logger.WithOutput(buf))
	err := l.Init()
	assert.Nil(t, err)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return errors.NewInternalErrorRsp("dial tcp")
	})

	e.Use(ExtractHeaderEcho(), logger.InjectLoggerEcho(l))

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{
    "header": {
        "operUid": "{{ReqUID}}",
        "rqUid": "{{ReqUID}}",
        "rqTm": "{{$isoTimestamp}}",
        "service": "{{service}}",
        "platform": "{{platform}}",
        "sourceSystem": "{{sourceSystem}}"
    }}`))
	res := httptest.NewRecorder()

	e.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Contains(t, res.Body.String(), "dial tcp")

	assert.NotNil(t, buf)
	assert.Contains(t, buf.String(), `"level":"error"`)
	assert.Contains(t, buf.String(), `"ErrorText":"dial tcp"`)
	assert.Contains(t, buf.String(), `"EventReceiver":"{{sourceSystem}}"`)
	assert.Contains(t, buf.String(), `"EventSource":"B-Connect"`)
	assert.Contains(t, buf.String(), `"StackTrace":"error stacktrace: `)
}

func Test_BaseErrorLogging(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := logger.NewSlogLogger(logger.WithOutput(buf))
	err := l.Init()
	assert.Nil(t, err)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return fmt.Errorf("dial tcp")
	})

	e.Use(ExtractHeaderEcho(), logger.InjectLoggerEcho(l))

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{
    "header": {
        "operUid": "{{ReqUID}}",
        "rqUid": "{{ReqUID}}",
        "rqTm": "{{$isoTimestamp}}",
        "service": "{{service}}",
        "platform": "{{platform}}",
        "sourceSystem": "{{sourceSystem}}"
    }}`))
	res := httptest.NewRecorder()

	e.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Contains(t, res.Body.String(), `{"message":"Internal Server Error"}`)

	assert.NotNil(t, buf)
	assert.Contains(t, buf.String(), `"level":"error"`)
	assert.Contains(t, buf.String(), `"ErrorText":"dial tcp"`)
	assert.Contains(t, buf.String(), `"EventReceiver":"{{sourceSystem}}"`)
	assert.Contains(t, buf.String(), `"EventSource":"B-Connect"`)
}

func Test_GroupMiddleware(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := logger.NewSlogLogger(logger.WithOutput(buf))
	err := l.Init()
	assert.Nil(t, err)

	e := echo.New()
	api := e.Group("/api", ExtractHeaderEcho(), logger.InjectLoggerEcho(l))
	api.GET("/v1", func(c echo.Context) error {
		return fmt.Errorf("dial tcp")
	})

	e.GET("/api/v2", func(c echo.Context) error { return fmt.Errorf("dial tcp") })

	req := httptest.NewRequest(http.MethodGet, "/api/v2", strings.NewReader(`{
    "header": {
        "operUid": "{{ReqUID}}",
        "rqUid": "{{ReqUID}}",
        "rqTm": "{{$isoTimestamp}}",
        "service": "{{service}}",
        "platform": "{{platform}}",
        "sourceSystem": "{{sourceSystem}}"
    }}`))
	res := httptest.NewRecorder()

	e.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Contains(t, res.Body.String(), `{"message":"Internal Server Error"}`)

	assert.NotNil(t, buf)
	assert.Equal(t, buf.String(), ``)
}

func Benchmark_RequestID(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	l := logger.NewSlogLogger(logger.WithOutput(buf), logger.WithSource())
	if err := l.Init(); err != nil {
		b.Fatal(err.Error())
	}

	e := echo.New()

	e.GET("/api/v2", func(c echo.Context) error {
		req, ok := request_id.FromContext(c.Request().Context())
		if !ok {
			b.Error(req)
		}

		if req != "{{ReqUID}}" {
			b.Error(req)
		}

		return c.String(http.StatusOK, "")
	}, ExtractHeaderEcho(), logger.InjectLoggerEcho(l))

	b.ReportAllocs()
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v2", strings.NewReader(`{
    "header": {
        "operUid": "{{ReqUID}}",
        "rqUid": "{{ReqUID}}",
        "rqTm": "{{$isoTimestamp}}",
        "service": "{{service}}",
        "platform": "{{platform}}",
        "sourceSystem": "{{sourceSystem}}"
    }}`))
			e.ServeHTTP(res, req)
		}
	})
}
