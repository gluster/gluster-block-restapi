package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gluster/gluster-block-restapi/pkg/utils"

	"github.com/dgrijalva/jwt-go"
)

var (
	requiredClaims = []string{"iss", "exp", "qsh"}
)

// JWTAuth maintains the rest configuration
type JWTAuth struct {
	user        string
	secret      string
	authEnabled bool
}

func (a JWTAuth) getAuthSecret(issuer string) string {
	if issuer == a.user {
		return a.secret
	}
	return ""
}

// NewJWTAuth returns JWTAuth
func NewJWTAuth(user, secret string, authEnabled bool) *JWTAuth {
	return &JWTAuth{
		user:        user,
		secret:      secret,
		authEnabled: authEnabled,
	}
}

// Auth is a middleware which authenticates HTTP requests
func (a *JWTAuth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If Auth disabled Return as is
		if !a.authEnabled {
			next.ServeHTTP(w, r)
			return
		}
		// Verify if Authorization header exists or not
		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		if authHeader == "" {
			utils.SendHTTPError(w, http.StatusUnauthorized, errors.New("'Authorization' header is required"))
			return
		}

		// Verify the Authorization header format "Bearer <TOKEN>"
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			utils.SendHTTPError(w, http.StatusUnauthorized, errors.New("'Authorization' header must be of the format - Bearer <TOKEN>"))
			return
		}

		// Verify JWT token with additional validations for Claims
		token, err := jwt.Parse(authHeaderParts[1], func(token *jwt.Token) (interface{}, error) {
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, fmt.Errorf("unable to parse Token claims")
			}

			// Error if required claims are not sent by Client
			for _, claimName := range requiredClaims {
				if _, claimOk := claims[claimName]; !claimOk {
					return nil, fmt.Errorf("token missing %s Claim", claimName)
				}
			}

			// Validate the JWT Signing Algo
			if _, tokenOk := token.Method.(*jwt.SigningMethodHMAC); !tokenOk {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			secret := a.getAuthSecret(claims["iss"].(string))
			if secret == "" {
				return nil, fmt.Errorf("invalid App ID: %s", claims["iss"])
			}

			qsh, err := utils.GenerateQsh(r)
			if err != nil {
				return nil, err
			}
			// Check qsh claim
			if claims["qsh"] != qsh {
				return nil, errors.New("invalid qsh claim in token")
			}
			// All checks GOOD, return the Secret to validate
			return []byte(secret), nil
		})

		// Check if token is Valid
		if err != nil {
			utils.SendHTTPError(w, http.StatusUnauthorized, err)
			return
		}
		if !token.Valid {
			utils.SendHTTPError(w, http.StatusUnauthorized, errors.New("invalid token specified in 'Authorization' header"))
			return
		}
		// Authentication is successful, continue serving the request
		next.ServeHTTP(w, r)
	})
}
