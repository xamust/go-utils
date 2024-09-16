package echo

import (
	"net/http"
	"strings"

	"github.com/xamust/go-utils/util/echo/json"
	"github.com/xamust/go-utils/util/echo/xml"

	"github.com/labstack/echo/v4"
)

type binder struct{}

func (b binder) Bind(i interface{}, c echo.Context) error {
	def := echo.DefaultBinder{}
	if err := def.BindPathParams(c, i); err != nil {
		return err
	}
	// Only bind query parameters for GET/DELETE/HEAD to avoid unexpected behavior with destination struct binding from body.
	// For example a request URL `&id=1&lang=en` with body `{"id":100,"lang":"de"}` would lead to precedence issues.
	// The HTTP method check restores pre-v4.1.11 behavior to avoid these problems (see issue #1670)
	method := c.Request().Method
	if method == http.MethodGet || method == http.MethodDelete || method == http.MethodHead {
		if err := def.BindQueryParams(c, i); err != nil {
			return err
		}
	}
	return b.BindBody(c, i)
}

func (b binder) BindBody(c echo.Context, i interface{}) (err error) {
	req := c.Request()
	if req.ContentLength == 0 {
		return
	}

	var deserializer func(c echo.Context, i interface{}) error

	ctype := req.Header.Get(echo.HeaderContentType)

	switch {
	case strings.HasPrefix(ctype, echo.MIMEApplicationJSON):
		deserializer = json.NewSerializeBCON().Deserialize
	case strings.HasPrefix(ctype, echo.MIMEApplicationXML), strings.HasPrefix(ctype, echo.MIMETextXML):
		deserializer = xml.NewSerializeBCON().Deserialize
	case strings.HasPrefix(ctype, echo.MIMEApplicationForm), strings.HasPrefix(ctype, echo.MIMEMultipartForm):
		def := echo.DefaultBinder{}
		if err = def.BindBody(c, i); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
		}
	default:
		return echo.ErrUnsupportedMediaType
	}

	if err = deserializer(c, i); err != nil {
		switch err.(type) {
		case *echo.HTTPError:
			return err
		default:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
		}
	}
	return nil
}

func NewBinderBCON() echo.Binder {
	return &binder{}
}
