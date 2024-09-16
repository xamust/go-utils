package rest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/xamust/go-utils/encoder/yml"
	"github.com/xamust/go-utils/logger"
	"github.com/xamust/go-utils/metadata/middleware"
	"github.com/xamust/go-utils/models_bcon"
	"github.com/xamust/go-utils/server"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	testHeaders = http.Header{
		"Account": []string{"test_Acc", "test_adc2"},
		"User":    []string{"admin", "zero"},
	}
	headerInbody = []byte(`
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
)

func Test_CreateRequest(t *testing.T) {
	logger.DefaultLogger = logger.NewLogger(nil, logger.WithSource())
	srv := httptest.NewServer(middleware.ExtractHeader(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodConnect, r.Method)

		assert.Subset(t, r.Header, testHeaders)

		all, _ := io.ReadAll(r.Body)
		assert.Equal(t, headerInbody, all)
		//logger.FromContextLogger(r.Context()).Info(r.Context(), "start test")

		w.Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

		fmt.Fprintln(w, `{"Msg":"Hello, client"}`)
	}))
	defer srv.Close()

	cli, err := NewClient()
	assert.Nil(t, err)

	req := NewRequest(headerInbody, http.MethodConnect, srv.URL)

	assert.Nil(t, err)

	rsp := struct {
		Msg string
	}{}

	err = cli.Call(context.Background(), req, &rsp, SetHeader(testHeaders), BasicAuth("test", "123"))
	assert.Nil(t, err)

	expected := struct {
		Msg string
	}{
		"Hello, client",
	}
	assert.Equal(t, expected, rsp)
}

func Test_LoggingNilBody(t *testing.T) {
	bufLogg := bytes.NewBuffer(nil)
	logger.DefaultLogger = logger.NewLogger(nil, logger.WithSource(), logger.WithOutput(bufLogg))
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodConnect, r.Method)

		assert.Subset(t, r.Header, testHeaders)
		assert.Equal(t, r.Body, http.NoBody)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, r.Header.Get(echo.HeaderContentType))

		w.Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

		fmt.Fprintln(w, `{"Msg":"Hello, client"}`)
	}))
	defer srv.Close()

	cli, err := NewClient()
	assert.Nil(t, err)

	req := NewRequest(nil, http.MethodConnect, srv.URL)
	assert.Nil(t, err)

	rsp := struct {
		Msg string
	}{}

	err = cli.Call(context.Background(), req, &rsp, BasicAuth("test", "123"), SetHeader(testHeaders))
	assert.Nil(t, err)

	expected := struct {
		Msg string
	}{
		"Hello, client",
	}
	assert.Equal(t, expected, rsp)
}

func Test_CodecCType(t *testing.T) {
	bufLogg := bytes.NewBuffer(nil)
	logger.DefaultLogger = logger.NewLogger(nil, logger.WithSource(), logger.WithOutput(bufLogg))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodConnect, r.Method)

		assert.Subset(t, r.Header, testHeaders)
		assert.Equal(t, "application/yaml", r.Header.Get(echo.HeaderContentType))

		w.Header().Add(echo.HeaderContentType, "application/yaml")

		fmt.Fprintln(w, "server:\n  addr: \":5201\"\n  timeout_conn: 120s")
	}))
	defer srv.Close()

	cli, err := NewClient()
	assert.Nil(t, err)

	req := NewRequest(nil, http.MethodConnect, srv.URL)
	assert.Nil(t, err)

	rsp := struct {
		Server server.Config `yaml:"server"`
	}{}

	err = cli.Call(context.Background(), req, &rsp, BasicAuth("test", "123"), SetHeader(testHeaders), CallCodec(yml.NewCodec()))
	assert.Nil(t, err)

	expected := server.Config{
		Addr:        ":5201",
		TimeoutConn: models_bcon.Duration{Duration: 120 * time.Second},
	}
	assert.Equal(t, expected, rsp.Server)
}

func Test_CreateRequestWithDebugOption(t *testing.T) {
	logger.DefaultLogger = logger.NewLogger(nil, logger.WithSource())
	srv := httptest.NewServer(middleware.ExtractHeader(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		assert.Subset(t, r.Header, testHeaders)

		all, _ := io.ReadAll(r.Body)
		assert.Equal(t, headerInbody, all)
		//logger.FromContextLogger(r.Context()).Info(r.Context(), "start test")

		w.Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

		fmt.Fprintln(w, `{"Msg":"Hello, client"}`)
	}))
	defer srv.Close()
	dbgSrv := httptest.NewServer(middleware.ExtractHeader(func(w http.ResponseWriter, r *http.Request) {
		dUrl := r.Header.Get("Debugurl")
		_, err := url.Parse(dUrl)
		assert.Nil(t, err)

		newRequest, err := http.NewRequest(r.Method, dUrl, r.Body)
		assert.Nil(t, err)
		for key, values := range r.Header {
			if key == "Isproxy" || key == "Debugurl" {
				continue
			}
			for _, value := range values {
				newRequest.Header.Add(key, value)
			}
		}
		cli := &http.Client{}
		rq, err := cli.Do(newRequest)
		assert.Nil(t, err)

		body := rq.Body
		defer body.Close()
		buf := new(strings.Builder)
		_, err = io.Copy(buf, body)
		assert.Nil(t, err)

		w.Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

		fmt.Fprintln(w, buf.String())
	}))
	defer dbgSrv.Close()

	cli, err := NewClient()
	assert.Nil(t, err)

	req := NewRequest(headerInbody, http.MethodPost, srv.URL)

	assert.Nil(t, err)

	rsp := struct {
		Msg string
	}{}

	err = cli.Call(context.Background(), req, &rsp,
		SetHeader(testHeaders),
		BasicAuth("test", "123"),
		EnableDebugProxy(dbgSrv.URL, false),
	)
	assert.Nil(t, err)

	expected := struct {
		Msg string
	}{
		"Hello, client",
	}
	assert.Equal(t, expected, rsp)
}
