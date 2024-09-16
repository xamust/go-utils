package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/xamust/go-utils/logger"
	"github.com/xamust/go-utils/metadata"
	"github.com/xamust/go-utils/util/echo/common"

	"github.com/labstack/echo/v4"
)

func WriteResponse(c echo.Context, code int, data any) error {
	ctx := c.Request().Context()
	log := logger.FromContextLogger(ctx)
	meta, _ := metadata.FromContextHeader(ctx)
	rsp := c.Response()

	header := rsp.Header()
	if header.Get(echo.HeaderContentType) == "" {
		header.Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8)
	}
	rsp.WriteHeader(code)

	buf := bytes.Buffer{}
	_, err := buf.Write([]byte(xml.Header))
	if err != nil {
		return err
	}

	err = xml.NewEncoder(&buf).Encode(data)
	if err != nil {
		return err
	}

	ctx = logger.NewContextEvent(ctx, logger.AppSystem, meta.Header.SourceSystem)
	log.Info(ctx, buf.String(), logger.Operation(fmt.Sprintf(common.SuccessRsp, meta.Header.SourceSystem)))

	_, err = rsp.Write(buf.Bytes())
	return err
}
