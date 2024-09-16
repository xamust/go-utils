package yml

import (
	"bytes"
	"io"

	enc "github.com/xamust/go-utils/encoder"

	"gopkg.in/yaml.v3"
)

type yamlEnc struct {
}

func NewCodec() enc.Codec {
	return &yamlEnc{}
}

func (e *yamlEnc) String() string {
	return "yaml"
}

func (e *yamlEnc) ContentType() string {
	return "application/yaml"
}

func (e *yamlEnc) Marshal(data any) (result []byte, err error) {
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
		result, err = yaml.Marshal(data)
	}
	return
}

func (e *yamlEnc) Unmarshal(data []byte, desc any) (err error) {
	if desc == nil {
		return nil
	}

	switch dType := desc.(type) {
	case *string:
		*dType = string(data)
	case *[]byte:
		*dType = append(*dType, data...)
	default:
		err = yaml.Unmarshal(data, desc)
	}

	return err
}

func (e *yamlEnc) ReadHeader(r io.Reader, desc any) (io.Reader, error) {
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
		err = yaml.Unmarshal(buf, desc)
	}

	return bytes.NewReader(buf), err
}

func (e *yamlEnc) ReadBody(r io.Reader, desc any) error {
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
		err = yaml.Unmarshal(buf, desc)
	}

	return err
}

func (e *yamlEnc) ReadBodyWorker(r io.Reader, desc any, w func(reader io.Reader) ([]byte, error)) error {
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
		err = yaml.Unmarshal(buf, desc)
	}

	return err
}

func (e *yamlEnc) Write(w io.Writer, msg any) (err error) {
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
		v, err = yaml.Marshal(msgType)
	}

	if err == nil {
		_, err = w.Write(v)
	}
	return err
}
