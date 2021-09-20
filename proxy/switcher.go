package proxy

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (p *Proxy) switcher(w http.ResponseWriter, r *http.Request) {
	if r.Host == p.config.proxyHost {
		// uses the built-in proxy routes
		var matcher mux.RouteMatch
		if check := p.router.Match(r, &matcher); check {

			newReq := mux.SetURLVars(r, matcher.Vars)
			matcher.Handler.ServeHTTP(w, newReq)

		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		// establish a reverse proxy with the container
		// todo
		w.WriteHeader(http.StatusBadGateway)
	}
}
