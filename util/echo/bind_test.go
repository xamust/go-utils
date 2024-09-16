package echo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xamust/go-utils/logger"
	"github.com/xamust/go-utils/metadata"
	"github.com/xamust/go-utils/metadata/middleware"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type parserHeader struct{}

func (p *parserHeader) ExtractHeaderFromBytes(b []byte) (metadata.HeaderSource, error) {
	return metadata.HeaderSource{
		Header: metadata.Header{
			RqUid:        "{{ReqUID}}",
			OperUID:      "{{ReqUID}}",
			RqTm:         "{{$isoTimestamp}}",
			Service:      "{{service}}",
			Platform:     "{{platform}}",
			SourceSystem: "{{sourceSystem}}",
		},
	}, nil
}

func (p *parserHeader) Validate(s metadata.HeaderSource) error {
	return nil
}

func Test_HandlerEcho(t *testing.T) {
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
	e.HideBanner = true

	e.POST("/", testHandle, middleware.ParsingHeaderEcho(&parserHeader{}), logger.InjectLoggerEcho(l))

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))
	assert.Nil(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	res := httptest.NewRecorder()

	e.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	var tReq testReq
	err = xml.Unmarshal(requestBody, &tReq)
	assert.Nil(t, err)

	logStr := buf.String()
	assert.NotEqual(t, "", logStr)

	//assert.Contains(t, logStr, "{\\\"msg\\\": \\\"TestMessage\\\"}}")
	assert.Contains(t, logStr, `"message":"TestMessage"`)
	//assert.Contains(t, logStr, fmt.Sprintf("%+v", &tReq))
	assert.Contains(t, logStr, fmt.Sprintf(`"HTTPHeaders":"%s"`, req.Header))
}

var requestBody = []byte(`
<body>
<header>
	<operUid>{{ReqUID}}</operUid>
	<rqUid>{{ReqUID}}</rqUid>
	<rqTm>{{$isoTimestamp}}</rqTm>
	<service>{{service}}</service>
	<platform>{{platform}}</platform>
	<sourceSystem>{{sourceSystem}}</sourceSystem>
</header>
<rqParms> 
	<msg>TestMessage</msg>
</rqParms>
</body>
`)

type testReq struct {
	XMLName xml.Name           `json:"-" xml:"body"`
	Header  metadata.HeaderReq `json:"header" xml:"header"`
	RqParms struct {
		Msg string `json:"msg" xml:"msg"`
	} `json:"rqParms" xml:"rqParms"`
}

func testHandle(c echo.Context) error {
	ctx := c.Request().Context()
	log := logger.FromContextLogger(ctx)

	reqModel := &testReq{}
	if errDec := c.Bind(&reqModel); errDec != nil {
		return errDec
	}

	log.Info(ctx, "TestMessage")

	return c.JSON(http.StatusOK, reqModel)
}
