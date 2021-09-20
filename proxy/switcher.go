package proxy

import (
	"net/http"
)

func (p *Proxy) switcher(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)

	// w.Write(bytes.NewBufferString("salut ca va").Bytes())
}
