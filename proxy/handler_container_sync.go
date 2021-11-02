package proxy

import "net/http"

func (p *Proxy) containerSyncGetHandler(w http.ResponseWriter, r *http.Request) {
	t, err := p.api.ContainerSyncGet(w, r)
	if err != nil {
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	p.WriteJson(w, t)
}

func (p *Proxy) containerSyncPostHandler(w http.ResponseWriter, r *http.Request) {
	t, err := p.api.ContainerSyncPost(w, r)
	if err != nil {
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	p.WriteJson(w, t)
}

func (p *Proxy) containerSyncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.containerSyncGetHandler(w, r)

	} else if r.Method == http.MethodPost {
		p.containerSyncPostHandler(w, r)

	}
}
