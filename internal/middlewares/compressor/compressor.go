package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// compressWriter is a struct of writer with zip compress options.
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewCompressWriter is a function for creating writer object with compress options.
func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header is a method of `compressWriter` object for getting header.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write is a method of `compressWriter` object for writing compressed data.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader is a method of `compressWriter` object for setting header with `statusCode`.
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close is a method of `compressWriter` for closing writer.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader is a struct of reader with zip compress options.
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewCompressReader is a function for creating reader with compress options.
func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read is a method of `compressReader` object for reading data.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close is a method of `compressReader` for closing reader.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// Compress is a function for compressing metrics data for sending to the server.
func Compress(h http.Handler) http.Handler {
	zipFn := func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			cw := NewCompressWriter(w)
			ow = cw
			defer cw.Close()

			cr, err := NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		} else if supportsGzip && !sendsGzip {
			cw := NewCompressWriter(w)
			ow = cw
			ow.Header().Set("Content-Encoding", "gzip")
			defer cw.Close()

		}
		h.ServeHTTP(ow, r)
	}
	return http.HandlerFunc(zipFn)
}
