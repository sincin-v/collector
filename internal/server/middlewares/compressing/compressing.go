package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	w  http.ResponseWriter
	cw gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		cw: *gzip.NewWriter(w),
	}
}

func (cw *compressWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	return cw.cw.Write(p)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		cw.w.Header().Set("Content-Encoding", "gzip")
	}
	cw.w.WriteHeader(statusCode)
}

func (cw *compressWriter) Close() error {
	return cw.cw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	cr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	cr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		cr: cr,
	}, nil
}

func (cr *compressReader) Read(p []byte) (n int, err error) {
	return cr.cr.Read(p)
}

func (cr *compressReader) Close() error {
	if err := cr.r.Close(); err != nil {
		return err
	}
	return cr.cr.Close()
}

func CompressMiddleware(h http.Handler) http.Handler {
	comperssFuncion := func(w http.ResponseWriter, r *http.Request) {
		writer := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			cw := newCompressWriter(w)
			cw.Header().Set("Content-Encoding", "gzip")
			writer = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, "gzip") {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.cr.Close()
		}

		h.ServeHTTP(writer, r)
	}

	return http.HandlerFunc(comperssFuncion)
}
