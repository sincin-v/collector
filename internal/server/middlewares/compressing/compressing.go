package compress

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: *gzip.NewWriter(w),
	}
}

func (cw *compressWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	return cw.zw.Write(p)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusMultipleChoices {
		cw.w.Header().Set("Content-Encoding", "gzip")
	}
	cw.w.WriteHeader(statusCode)
}

func (cw *compressWriter) Close() error {
	return cw.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	cr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: cr,
	}, nil
}

func (cr *compressReader) Read(p []byte) (n int, err error) {
	return cr.zr.Read(p)
}

func (cr *compressReader) Close() error {
	if err := cr.r.Close(); err != nil {
		return err
	}
	return cr.zr.Close()
}

func CompressMiddleware(h http.Handler) http.Handler {
	comperssFuncion := func(w http.ResponseWriter, r *http.Request) {
		writer := w
		var err error
		acceptEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			cw := newCompressWriter(w)
			cw.Header().Set("Content-Encoding", "gzip")
			writer = cw
			defer func() {
				if errCompressWriter := cw.Close(); errCompressWriter != nil {
					err = errors.Join(err, fmt.Errorf("close compress writer error: %w", errCompressWriter))
				}
			}()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, "gzip") {
			cr, err := newCompressReader(r.Body)
			defer func() {
				if errCompressReader := cr.Close(); errCompressReader != nil {
					err = errors.Join(err, fmt.Errorf("close compress reader error: %w", errCompressReader))
				}
			}()
			r.Body = cr
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.ServeHTTP(writer, r)
	}

	return http.HandlerFunc(comperssFuncion)
}
