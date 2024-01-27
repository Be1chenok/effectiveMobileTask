package handler

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

func (h Handler) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				h.logger.Error(err)
				writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
			}
		}()

		startTime := time.Now()

		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Errorf("failed to read request body: %v", err)
			writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)

		h.logger.Infof("status: %d, method: %s, path: %s, request body: %s response body: %s, duration: %s",
			wrapped.status, r.Method, r.URL.EscapedPath(), requestBody, wrapped.body.Bytes(), time.Since(startTime))
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	if err != nil {
		return n, err
	}
	rw.body.Write(b)
	return n, nil
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.status != 0 {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
