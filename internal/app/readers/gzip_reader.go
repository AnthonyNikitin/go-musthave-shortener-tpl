package readers

import (
	"compress/gzip"
	"io"
)

type GzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewGzipReader(r io.ReadCloser) (*GzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &GzipReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c GzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *GzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
