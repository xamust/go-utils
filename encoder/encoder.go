package encoder

import (
	"io"
)

type Codec interface {
	Marshal(data any) ([]byte, error)
	Unmarshal(data []byte, desc any) error

	ReadHeader(r io.Reader, desc any) (io.Reader, error)
	ReadBody(r io.Reader, desc any) error
	ReadBodyWorker(r io.Reader, desc any, w func(reader io.Reader) ([]byte, error)) error

	Write(w io.Writer, msg any) error

	String() string
	ContentType() string
}
