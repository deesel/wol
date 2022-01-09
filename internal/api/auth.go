package api

import (
	"net/http"
	"strings"
)

type APIKey string

func (a *API) isAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authorization"]

		if authHeader == nil {
			panic(ErrTokenNotFound)
		}

		if len(authHeader) != 1 {
			panic(ErrMultipleTokens)
		}

		if !strings.HasPrefix(authHeader[0], "Bearer ") {
			panic(ErrInvalidAuthHeader)
		}

		apiKey := APIKey(authHeader[0][len("Bearer "):])
		found := false

		for _, key := range a.APIKeys {
			if apiKey == key {
				found = true
			}
		}

		if !found {
			panic(ErrForbidden)
		}

		next.ServeHTTP(w, r)
	})
}
