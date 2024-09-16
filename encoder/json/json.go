package json

import (
	"bytes"
	"encoding/json"
	"io"

	enc "github.com/xamust/go-utils/encoder"

	"github.com/labstack/echo/v4"
)

type jsonEnc struct {
}

func NewCodec() enc.Codec {
	return &jsonEnc{}
}

func (e *jsonEnc) String() string {
	return "json"
}

func (e *jsonEnc) ContentType() string {
	return echo.MIMEApplicationJSONCharsetUTF8
}

func (e *jsonEnc) Marshal(data any) (result []byte, err error) {
	if data == nil {
		return nil, nil
	}

	switch dType := data.(type) {
	case string:
		result = []byte(dType)
	case *string:
		result = []byte(*dType)
	case []byte:
		result = dType
	case *[]byte:
		result = *dType
	default:
		result, err = json.Marshal(data)
	}
	return
}

func (e *jsonEnc) Unmarshal(data []byte, desc any) (err error) {
	if desc == nil {
		return nil
	}

	switch dType := desc.(type) {
	case *string:
		*dType = string(data)
	case *[]byte:
		*dType = append(*dType, data...)
	default:
		err = json.Unmarshal(data, desc)
	}

	return err
}

func (e *jsonEnc) ReadHeader(r io.Reader, desc any) (io.Reader, error) {
	if desc == nil {
		return r, nil
	}
	buf, err := io.ReadAll(r)
	if err != nil {
		return r, err
	}

	switch descType := desc.(type) {
	case *string:
		*descType = string(buf)
	case *[]byte:
		*descType = buf
	default:
		err = json.Unmarshal(buf, desc)
	}

	return bytes.NewReader(buf), err
}

func (e *jsonEnc) ReadBody(r io.Reader, desc any) error {
	if desc == nil {
		return nil
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	switch descType := desc.(type) {
	case *string:
		*descType = string(buf)
	case *[]byte:
		*descType = buf
	default:
		err = json.Unmarshal(buf, desc)
	}

	return err
}

func (e *jsonEnc) ReadBodyWorker(r io.Reader, desc any, w func(reader io.Reader) ([]byte, error)) error {
	if desc == nil {
		return nil
	}

	buf, err := w(r)
	if err != nil {
		return err
	}

	switch descType := desc.(type) {
	case *string:
		*descType = string(buf)
	case *[]byte:
		*descType = buf
	default:
		err = json.Unmarshal(buf, desc)
	}

	return err
}

func (e *jsonEnc) Write(w io.Writer, msg any) (err error) {
	if msg == nil {
		return err
	}

	var v []byte
	switch msgType := msg.(type) {
	case []byte:
		v = msgType
	case *[]byte:
		v = *msgType
	case string:
		v = []byte(msgType)
	case *string:
		v = []byte(*msgType)
	default:
		v, err = json.Marshal(msgType)
	}

	if err == nil {
		_, err = w.Write(v)
	}
	return err
}
