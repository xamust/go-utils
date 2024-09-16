package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/xamust/go-utils/logger"
	"github.com/xamust/go-utils/metadata"
	"github.com/xamust/go-utils/util/echo/common"

	"github.com/labstack/echo/v4"
)

func NewSerializeBCON() echo.JSONSerializer {
	return &serializeEcho{}
}

type serializeEcho struct{}

func (s *serializeEcho) Serialize(c echo.Context, i interface{}, indent string) (err error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	if indent != "" {
		enc.SetIndent("", indent)
	}

	if err = enc.Encode(i); err != nil {
		return err
	}

	ctx := c.Request().Context()
	log := logger.FromContextLogger(ctx)
	h, _ := metadata.FromContextHeader(ctx)

	ctx = logger.NewContextEvent(ctx, logger.AppSystem, h.Header.SourceSystem)
	log.Info(ctx, buf.String(), logger.Operation(fmt.Sprintf(common.SuccessRsp, h.Header.SourceSystem)))

	_, err = c.Response().Write(buf.Bytes())
	return err
}

func (s *serializeEcho) Deserialize(c echo.Context, i interface{}) error {
	ctx := c.Request().Context()
	log := logger.FromContextLogger(ctx)
	h, _ := metadata.FromContextHeader(ctx)

	reqBuf, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	ctx = logger.NewContextEvent(ctx, h.Header.SourceSystem, logger.AppSystem)

	log.Info(ctx, string(reqBuf), logger.Operation(fmt.Sprintf(common.CatchReq, h.Header.SourceSystem)), logger.HttpHeaders(cloneHeaders(c.Request().Header)))

	err = json.NewDecoder(bytes.NewReader(reqBuf)).Decode(i)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	} else if se, ok := err.(*json.SyntaxError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
	}
	return err
}

func cloneHeaders(headers http.Header) http.Header {
	clone := make(http.Header, len(headers))
	for k, values := range headers {
		clone[k] = make([]string, len(values))
		copy(clone[k], values)
	}
	return clone
}
