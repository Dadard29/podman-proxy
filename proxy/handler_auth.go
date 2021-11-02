package proxy

import "net/http"

func (p *Proxy) authGetHandler(w http.ResponseWriter, r *http.Request) {
	err := p.api.AuthGet(w, r, p.config.jwtKey)

	if err != nil {
		p.WriteErrorJson(w, http.StatusUnauthorized, err)
		return
	}

	p.WriteMessageJson(w, "JWT is valid")
}

func (p *Proxy) authPostHandler(w http.ResponseWriter, r *http.Request) {
	token, code, err := p.api.AuthPost(w, r, p.config.jwtKey)
	if err != nil {
		p.WriteErrorJson(w, code, err)
		return
	}

	p.WriteJson(w, token)
}

func (p *Proxy) authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.authGetHandler(w, r)

	} else if r.Method == http.MethodPost {
		p.authPostHandler(w, r)
	}
}
