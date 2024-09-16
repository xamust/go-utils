package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/xamust/go-utils/logger"
	"github.com/xamust/go-utils/metadata"
	"github.com/xamust/go-utils/util/echo/common"

	"github.com/labstack/echo/v4"
)

type Serializer interface {
	Deserialize(c echo.Context, i interface{}) error
}

func NewSerializeBCON() Serializer {
	return &serializeEcho{}
}

type serializeEcho struct{}

func (s *serializeEcho) Deserialize(c echo.Context, i interface{}) error {
	ctx := c.Request().Context()
	log := logger.FromContextLogger(ctx)
	h, _ := metadata.FromContextHeader(ctx)

	reqBuf, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	ctx = logger.NewContextEvent(ctx, h.Header.SourceSystem, logger.AppSystem)
	log.Info(ctx, string(reqBuf), logger.Operation(fmt.Sprintf(common.CatchReq, h.Header.SourceSystem)), logger.HttpHeaders(c.Request().Header))

	err = xml.NewDecoder(bytes.NewReader(reqBuf)).Decode(i)
	if ute, ok := err.(*xml.UnsupportedTypeError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unsupported type error: type=%v, error=%v", ute.Type, ute.Error())).SetInternal(err)
	} else if se, ok := err.(*xml.SyntaxError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: line=%v, error=%v", se.Line, se.Error())).SetInternal(err)
	}
	return err
}
