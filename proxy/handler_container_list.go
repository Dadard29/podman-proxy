package proxy

import "net/http"

func (p *Proxy) containersGetHandler(w http.ResponseWriter, r *http.Request) {
	containersList, err := p.api.ContainersGet(w, r)
	if err != nil {
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	p.WriteJson(w, containersList)
}

// Main entrypoint for containers list management
func (p *Proxy) containerListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.containersGetHandler(w, r)
	}
}
