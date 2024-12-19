package authn

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sphinx-camfield/utils/stackable"
	"net/http"
	"strings"
)

type UnAuthResponder interface {
	RespondEmptyHeader(w http.ResponseWriter, r *http.Request)
	RespondInvalidJwt(w http.ResponseWriter, r *http.Request)
	RespondInvalidClaim(w http.ResponseWriter, r *http.Request, err error)
	RespondInvalidSubject(w http.ResponseWriter, r *http.Request, err error)
}

type JwtAuthnOptions struct {
	RespondEmptyHeader func(w http.ResponseWriter, r *http.Request)
	RespondInvalidJwt  func(w http.ResponseWriter, r *http.Request)
}

type JwtKeyProvider interface {
	GetKey(token *jwt.Token) (interface{}, error)
}

type JwtAuthOption func(*JwtAuthnOptions)

func WithRespondEmptyHeader(fn func(w http.ResponseWriter, r *http.Request)) JwtAuthOption {
	return func(o *JwtAuthnOptions) {
		o.RespondEmptyHeader = fn
	}
}

func WithRespondInvalidJwt(fn func(w http.ResponseWriter, r *http.Request)) JwtAuthOption {
	return func(o *JwtAuthnOptions) {
		o.RespondInvalidJwt = fn
	}
}

func JwtAuthn(provider JwtKeyProvider, opts ...JwtAuthOption) stackable.Stackable {

	compiledOpts := &JwtAuthnOptions{
		RespondEmptyHeader: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		},
		RespondInvalidJwt: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		},
	}

	for _, opt := range opts {
		opt(compiledOpts)
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				compiledOpts.RespondEmptyHeader(w, r)
				return
			}

			segments := strings.Split(authHeader, " ")

			if len(segments) != 2 || strings.ToLower(segments[0]) != "bearer" {
				compiledOpts.RespondInvalidJwt(w, r)
				return
			}

			tokenStr := segments[1]

			token, err := jwt.ParseWithClaims(tokenStr, &jwt.MapClaims{}, provider.GetKey)
			if err != nil {
				compiledOpts.RespondInvalidJwt(w, r)
				return
			}

			next(w, r.WithContext(
				context.WithValue(r.Context(), "token", token)),
			)
		}
	}
}
