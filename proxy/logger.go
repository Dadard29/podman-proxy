package proxy

import (
	"net/http"
	"time"

	"github.com/Dadard29/podman-proxy/models"
)

func (p *Proxy) dbLoggingMiddleware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		responseLog := models.NewLogResponseWriter(w)
		beganAt := time.Now()

		next.ServeHTTP(responseLog, r)

		netLog, err := p.api.NewNetworkLog(beganAt, r, responseLog)
		if err != nil {
			return
		}

		p.logger.Printf("%d %s %s %s %s", netLog.ResponseStatusCode, r.Method, r.Host, r.URL.Host, r.URL.Path)
	}

	return http.HandlerFunc(fn)
}
