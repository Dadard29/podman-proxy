package proxy

import (
	"bytes"
	"net/http"
	"time"

	"github.com/Dadard29/podman-proxy/models"
)

// https://stackoverflow.com/questions/64243247/go-gorilla-log-each-request-duration-and-status-code
type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w}
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	w.buf.Write(body)
	return w.ResponseWriter.Write(body)
}

func (p *Proxy) dbLoggingMiddleware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		responseLog := NewLogResponseWriter(w)
		beganAt := time.Now()

		time.Sleep(23 * time.Millisecond)

		next.ServeHTTP(responseLog, r)

		duration := time.Since(beganAt)
		responseBody := responseLog.buf.Bytes()
		netLog, err := models.NewNetworkLogFromRequest(r, responseLog.statusCode, duration, responseBody)
		if err != nil {
			p.logger.Println(err)
		}

		err = p.db.InsertNetworkLog(netLog)
		if err != nil {
			p.logger.Println(err)
		}
	}

	return http.HandlerFunc(fn)
}
