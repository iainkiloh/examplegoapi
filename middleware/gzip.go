package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type GzipMiddleware struct {
	Next http.Handler
}

func NewGzipMiddleware(handlerToWrap http.Handler) *GzipMiddleware {
	return &GzipMiddleware{Next: handlerToWrap}
}

func (gm *GzipMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if gm.Next == nil {
		gm.Next = http.DefaultServeMux
	}

	encodings := r.Header.Get("Accept-Encoding")
	if !strings.Contains(encodings, "gzip") {
		gm.Next.ServeHTTP(w, r)
		return
	}

	w.Header().Add("Content-Encoding", "gzip")
	gzipWriter := gzip.NewWriter(w)
	defer gzipWriter.Close()

	//pass on a gzip response writer to the next handler
	gzrw := gzipResponseWriter{
		ResponseWriter: w,
		Writer:         gzipWriter,
	}
	gm.Next.ServeHTTP(gzrw, r)
}

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (gzrw gzipResponseWriter) Write(data []byte) (int, error) {
	return gzrw.Writer.Write(data)
}
