package api

import (
	"fmt"
	"net/http"

	"github.com/Dadard29/podman-proxy/models"
	"golang.org/x/crypto/bcrypt"
)

// Verify an existing JWT
func (a *Api) AuthGet(w http.ResponseWriter, r *http.Request, jwtKey string) error {
	token := models.NewAccessTokenFromRequest(r)
	return token.Verify(jwtKey)
}

// Issue a new JWT
func (a *Api) AuthPost(w http.ResponseWriter, r *http.Request, jwtKey string) (*models.AccessToken, int, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		s := "invalid authentification format"
		return nil, http.StatusUnauthorized, fmt.Errorf(s)
	}

	user, err := a.db.GetUser(username)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	hashedPasswordBytes := []byte(user.HashedPassword)
	passwordBytes := []byte(password)
	err = bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	token, err := models.NewAccessToken(user, jwtKey)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &token, http.StatusOK, nil
}
