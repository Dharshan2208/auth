package middleware

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode  int
	WroteHeader bool
}

func (rw *ResponseWriter) WriteHeader(code int) {
	if rw.WroteHeader {
		return
	}

	rw.StatusCode = code
	rw.WroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if !rw.WroteHeader {
		rw.WriteHeader(http.StatusOK)
	}

	return rw.ResponseWriter.Write(b)
}
