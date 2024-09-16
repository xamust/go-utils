package echo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xamust/go-utils/logger"
	"github.com/xamust/go-utils/metadata/middleware"
	json2 "github.com/xamust/go-utils/util/echo/json"
	"github.com/xamust/go-utils/util/echo/xml"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_HandlerEchoSJSON(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := logger.NewSlogLogger(
		logger.WithLevel(logger.DebugLevel),
		logger.WithOutput(buf),
		logger.WithSource(),
	)
	err := l.Init()
	assert.Nil(t, err)

	e := echo.New()
	e.Binder = NewBinderBCON()
	e.JSONSerializer = json2.NewSerializeBCON()
	e.HideBanner = true

	e.POST("/", testHandleJSON, middleware.ExtractHeaderEcho(), logger.InjectLoggerEcho(l))

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBodyJSON))
	assert.Nil(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()

	e.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	var tReq testReq
	err = json.Unmarshal(requestBodyJSON, &tReq)
	assert.Nil(t, err)

	logStr := buf.String()
	assert.NotEqual(t, "", logStr)

	assert.Contains(t, logStr, "{\\\"msg\\\": \\\"RequestMessage\\\"}}")
	assert.Contains(t, logStr, `"message":"TestMessage"`)
	//	assert.Contains(t, logStr, fmt.Sprintf("%+v", &tReq))
	assert.Contains(t, logStr, fmt.Sprintf(`"HTTPHeaders":"%s"`, req.Header))
	assert.NotEqual(t, testReq{}, tReq)
}

var requestBodyJSON = []byte(`{"header":{"operUid": "{{ReqUID}}","rqUid": "{{ReqUID}}","rqTm": "{{$isoTimestamp}}","service": "{{service}}","platform": "{{platform}}","sourceSystem": "{{sourceSystem}}"},"rqParms": {"msg": "RequestMessage"}}`)

func testHandleJSON(c echo.Context) error {
	ctx := c.Request().Context()
	log := logger.FromContextLogger(ctx)

	reqModel := &testReq{}
	if errDec := c.Bind(&reqModel); errDec != nil {
		return errDec
	}

	log.Info(ctx, "TestMessage")

	return xml.WriteResponse(c, http.StatusOK, reqModel)
}
