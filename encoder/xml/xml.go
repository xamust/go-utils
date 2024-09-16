package xml

import (
	"bytes"
	"encoding/xml"
	"io"

	enc "github.com/xamust/go-utils/encoder"

	"github.com/labstack/echo/v4"
)

type xmlEnc struct {
}

func NewCodec() enc.Codec {
	return &xmlEnc{}
}

func (e *xmlEnc) String() string {
	return "xml"
}

func (e *xmlEnc) ContentType() string {
	return echo.MIMEApplicationXMLCharsetUTF8
}

func (e *xmlEnc) Marshal(data any) (result []byte, err error) {
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
		result, err = xml.Marshal(data)
	}
	return
}

func (e *xmlEnc) Unmarshal(data []byte, desc any) (err error) {
	if desc == nil {
		return nil
	}

	switch dType := desc.(type) {
	case *string:
		*dType = string(data)
	case *[]byte:
		*dType = append(*dType, data...)
	default:
		decoder := xml.NewDecoder(bytes.NewBuffer(data))
		decoder.CharsetReader = identReader
		err = decoder.Decode(desc)
	}

	return err
}

func (e *xmlEnc) ReadHeader(r io.Reader, desc any) (io.Reader, error) {
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
		decoder := xml.NewDecoder(bytes.NewBuffer(buf))
		decoder.CharsetReader = identReader
		err = decoder.Decode(desc)
	}

	return bytes.NewReader(buf), err
}

func identReader(_ string, input io.Reader) (io.Reader, error) {
	return input, nil
}

func (e *xmlEnc) ReadBody(r io.Reader, desc any) error {
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
		err = xml.Unmarshal(buf, desc)
	}

	return err
}

func (e *xmlEnc) ReadBodyWorker(r io.Reader, desc any, w func(reader io.Reader) ([]byte, error)) error {
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
		err = xml.Unmarshal(buf, desc)
	}

	return err
}

func (e *xmlEnc) Write(w io.Writer, msg any) (err error) {
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
		v, err = xml.Marshal(msgType)
	}

	if err == nil {
		_, err = w.Write(v)
	}
	return err
}
