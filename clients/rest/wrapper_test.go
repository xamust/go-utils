package rest

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xamust/go-utils/logger"
	"github.com/xamust/go-utils/server"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_NewWrapper(t *testing.T) {
	ctx := context.Background()
	bufLogg := bytes.NewBuffer(nil)
	logger.DefaultLogger = logger.NewLogger(nil, logger.WithSource(), logger.WithOutput(bufLogg))
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodConnect, r.Method)

		assert.Subset(t, r.Header, testHeaders)
		assert.Equal(t, r.Body, http.NoBody)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, r.Header.Get(echo.HeaderContentType))

		auth := r.Header.Get(echo.HeaderAuthorization)
		assert.NotEqual(t, 0, len(auth))

		w.Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

		fmt.Fprintln(w, `{"Msg":"Hello, client"}`)
	}))
	defer srv.Close()

	testParams := server.EndpointParams{
		Login:    "test_login",
		Password: "qwerty",
	}
	cli, err := NewWrapper(testParams)
	assert.Nil(t, err)
	assert.NotNil(t, cli)

	req := NewRequest(nil, http.MethodConnect, srv.URL)
	assert.Nil(t, err)

	rsp := struct {
		Msg string
	}{}

	err = cli.Call(ctx, req, &rsp, SetHeader(testHeaders))
	assert.Nil(t, err)

	expected := struct {
		Msg string
	}{
		"Hello, client",
	}
	assert.Equal(t, expected, rsp)
}
