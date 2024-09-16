package dual_writer

import "io"

type Duplicator interface {
	Write(p []byte) (n int, err error)
}

type duplicator struct {
	writers []io.Writer
}

func NewDuplicator(writers ...io.Writer) Duplicator {
	return &duplicator{
		writers: writers,
	}
}

func (d *duplicator) Write(p []byte) (n int, err error) {
	for _, writer := range d.writers {
		n, err = writer.Write(p)
		if err != nil {
			return n, err
		}
	}
	return n, nil
}
