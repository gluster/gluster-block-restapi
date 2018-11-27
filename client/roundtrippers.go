package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gluster/gluster-block-restapi/pkg/utils"

	"github.com/dgrijalva/jwt-go"
)

const (
	expirationDuration = time.Second * 120
)

// NewBearerAuthRoundTrippers applies a Authorization header to a request with value set as bearer token
// formed using given username and password.
// It wraps a http Roundtripper
func NewBearerAuthRoundTrippers(username string, pass string, nextRT http.RoundTripper) http.RoundTripper {
	return &bearerAuthRoundTrippers{username, pass, nextRT}
}

type bearerAuthRoundTrippers struct {
	username string
	password string
	nextRT   http.RoundTripper
}

// RoundTrip will set Authorization header in request and execute the http req
func (b *bearerAuthRoundTrippers) RoundTrip(req *http.Request) (*http.Response, error) {
	if b.username == "" && b.password == "" {
		return b.nextRT.RoundTrip(req)
	}

	token, err := b.bearerToken(req)
	if err != nil {
		return b.nextRT.RoundTrip(req)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return b.nextRT.RoundTrip(req)
}

func (b *bearerAuthRoundTrippers) bearerToken(req *http.Request) (string, error) {
	qsh, err := utils.GenerateQsh(req)
	if err != nil {
		return "", err
	}
	// Create Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// Set issuer
		"iss": b.username,
		// Set expiration
		"exp": time.Now().Add(expirationDuration).Unix(),
		// Set qsh
		"qsh": qsh,
	})
	// Sign the token
	return token.SignedString([]byte(b.password))

}
