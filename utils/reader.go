package utils

import "io"

type Reader struct {
	Target string
	R      io.Reader
	OnEof  func(tar string)
	OnRead func(size int)
}

func (c *Reader) Read(p []byte) (int, error) {
	n, err := c.R.Read(p)
	if err == io.EOF && c.OnEof != nil {
		c.OnEof(c.Target)
	}
	if n > 0 && c.OnRead != nil {
		c.OnRead(n)
	}
	return n, err
}
