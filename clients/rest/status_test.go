package rest

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xamust/go-utils/encoder/json"
	"github.com/xamust/go-utils/logger"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_Noop(t *testing.T) {
	bufLogg := bytes.NewBuffer(nil)
	logger.DefaultLogger = logger.NewLogger(nil, logger.WithSource(), logger.WithOutput(bufLogg))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodConnect, r.Method)

		assert.Subset(t, r.Header, testHeaders)
		assert.Contains(t, r.Header.Get(echo.HeaderContentType), "application/json")

		w.Header().Add(echo.HeaderContentType, "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"Msg":"Hello, client"}`))
	}))
	defer srv.Close()

	cli, err := NewClient()
	assert.Nil(t, err)

	req := NewRequest(nil, http.MethodConnect, srv.URL)
	assert.Nil(t, err)

	rsp := StatusNoop{}

	err = cli.Call(context.Background(), req, &rsp, BasicAuth("test", "123"), SetHeader(testHeaders), CallCodec(json.NewCodec()))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusAccepted, rsp.Status())
}
