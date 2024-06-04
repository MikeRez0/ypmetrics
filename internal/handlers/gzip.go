package handlers

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinCompress(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ow := ctx.Writer

		acceptEncoding := ctx.Request.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(ow)

			ow = cw

			defer func() { _ = cw.Close() }()
		}

		contentEncoding := ctx.Request.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(ctx.Request.Body)
			if err != nil {
				err = ctx.AbortWithError(http.StatusInternalServerError,
					fmt.Errorf("error decompressing request: %w", err))
				log.Error("error creating compress reader", zap.Error(err))
				return
			}
			ctx.Request.Body = cr
			defer func() { _ = cr.Close() }()
		}

		ctx.Writer = ow
		ctx.Next()
	}
}

type compressWriter struct {
	gin.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w gin.ResponseWriter) *compressWriter {
	return &compressWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p) //nolint:wrapcheck //error checked in base handler
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 { //nolint:gomnd //Yes, it's a magic number
		c.Header().Set("Content-Encoding", "gzip")
	}
	c.ResponseWriter.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close() //nolint:wrapcheck //error checked in base handler
}

type compressReader struct {
	io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err //nolint:wrapcheck //error checked in base handler
	}

	return &compressReader{
		ReadCloser: r,
		zr:         zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p) //nolint:wrapcheck //error checked in base handler
}

func (c *compressReader) Close() error {
	if err := c.ReadCloser.Close(); err != nil {
		return err //nolint:wrapcheck //error checked in base handler
	}
	return c.zr.Close() //nolint:wrapcheck //error checked in base handler
}
