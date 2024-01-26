package handler

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

func (h Handler) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		proxyWriter := &responseWriterProxy{ResponseWriter: w}

		next.ServeHTTP(proxyWriter, r)

		h.logger.Infof("status: %d, path: %s, method: %s, params: %v, request body: %s, response body: %s, duration: %s",
			proxyWriter.statusCode, r.URL.Path, r.Method, r.URL.Query(), requestBody, proxyWriter.body.Bytes(), time.Since(startTime))
	})
}

type responseWriterProxy struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (p *responseWriterProxy) Write(b []byte) (int, error) {
	p.body.Write(b)
	return p.ResponseWriter.Write(b)
}

func (p *responseWriterProxy) WriteHeader(code int) {
	p.statusCode = code
	p.ResponseWriter.WriteHeader(code)
}
