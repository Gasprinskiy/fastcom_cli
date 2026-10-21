package utils

import "io"

type Reader struct {
	R      io.Reader
	OnEof  func()
	OnRead func(size int)
}

func (c *Reader) Read(p []byte) (int, error) {
	n, err := c.R.Read(p)
	if err == io.EOF && c.OnEof != nil {
		c.OnEof()
	}
	if n > 0 && c.OnRead != nil {
		c.OnRead(n)
	}
	return n, err
}
