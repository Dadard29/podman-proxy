package models

import (
	"net/http"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type AccessToken struct {
	Token string `json:"token"`
}

const issuer = "podman-proxy-issuer"
const subject = "podman-proxy-subject"

func NewAccessToken(user User, jwtKey string) (AccessToken, error) {
	key := jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       []byte(jwtKey),
	}
	opts := (&jose.SignerOptions{}).WithType("JWT")
	sig, err := jose.NewSigner(key, opts)
	if err != nil {
		return AccessToken{}, err
	}

	duration := 30 * time.Second
	expiry := time.Now().Add(duration)

	claims := jwt.Claims{
		Subject:  subject,
		Issuer:   issuer,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Expiry:   jwt.NewNumericDate(expiry),
		Audience: jwt.Audience{user.Name},
	}

	raw, err := jwt.Signed(sig).Claims(claims).CompactSerialize()
	if err != nil {
		return AccessToken{}, err
	}

	return AccessToken{
		Token: raw,
	}, nil
}

func NewAccessTokenFromRequest(r *http.Request) AccessToken {
	auth := r.Header.Get("Authorization")
	authToken := strings.TrimPrefix(auth, "Bearer ")
	return AccessToken{
		Token: authToken,
	}
}

func (a AccessToken) Verify(jwtKey string) error {
	token, err := jwt.ParseSigned(a.Token)
	if err != nil {
		return err
	}

	out := jwt.Claims{}
	if err := token.Claims([]byte(jwtKey), &out); err != nil {
		return err
	}

	expectedClaims := jwt.Expected{
		Issuer:  issuer,
		Subject: subject,
		Time:    time.Now(),
	}

	return out.ValidateWithLeeway(expectedClaims, 0)
}
