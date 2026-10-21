package main

import "io"

type Reader struct {
	r      io.Reader
	onEof  func()
	onRead func(size int)
}

func (c *Reader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	if err == io.EOF && c.onEof != nil {
		c.onEof()
	}
	if n > 0 && c.onRead != nil {
		c.onRead(n)
	}
	return n, err
}
