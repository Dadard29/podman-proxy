package proxy

import (
	"net/http"

	"github.com/Dadard29/podman-proxy/models"
)

func (p *Proxy) authMiddleware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth" {
			token := models.NewAccessTokenFromRequest(r)
			if err := token.Verify(p.config.jwtKey); err != nil {
				p.logger.Println(err)
				p.WriteErrorJson(w, http.StatusForbidden, err)
				return
			}
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
