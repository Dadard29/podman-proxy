package models

import (
	"bytes"
	"net/http"
)

// https://stackoverflow.com/questions/64243247/go-gorilla-log-each-request-duration-and-status-code
type LogResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Buf        bytes.Buffer
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w}
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	w.Buf.Write(body)
	return w.ResponseWriter.Write(body)
}
