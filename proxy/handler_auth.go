package proxy

import (
	"fmt"
	"net/http"

	"github.com/Dadard29/podman-proxy/models"
	"golang.org/x/crypto/bcrypt"
)

// Verify an existing JWT
func (p *Proxy) authGet(w http.ResponseWriter, r *http.Request) {
	token := models.NewAccessTokenFromRequest(r)
	if err := token.Verify(p.config.jwtKey); err != nil {
		p.logger.Println(err)
		p.WriteErrorJson(w, http.StatusUnauthorized, err)
		return
	}

	p.WriteMessageJson(w, "JWT is valid")
}

// Issue a new JWT
func (p *Proxy) authPost(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		s := "invalid authentification format"
		p.logger.Println(s)
		p.WriteErrorJson(w, http.StatusUnauthorized, fmt.Errorf(s))
		return
	}

	user, err := p.db.GetUser(username)
	if err != nil {
		p.logger.Println(err)
		p.WriteErrorJson(w, http.StatusNotFound, err)
		return
	}

	hashedPasswordBytes := []byte(user.HashedPassword)
	passwordBytes := []byte(password)
	err = bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes)
	if err != nil {
		p.logger.Println(err)
		p.WriteErrorJson(w, http.StatusUnauthorized, err)
		return
	}

	token, err := models.NewAccessToken(user, p.config.jwtKey)
	if err != nil {
		p.logger.Println(err)
		p.WriteErrorJson(w, http.StatusInternalServerError, err)
		return
	}

	p.WriteJson(w, &token)
}

func (p *Proxy) authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.authGet(w, r)

	} else if r.Method == http.MethodPost {
		p.authPost(w, r)
	}
}
