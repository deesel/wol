package api

import (
	e "errors"
	"net/http"
)

// ErrTokenNotFound represents required auth token not found error
var ErrTokenNotFound = e.New("required auth token not found")

// ErrMultipleTokens represents multiple auth tokens found error
var ErrMultipleTokens = e.New("multiple auth tokens found")

// ErrInvalidAuthHeader represents invalid authorization header format error
var ErrInvalidAuthHeader = e.New("invalid authorization header format")

// ErrForbidden represents forbidden access error
var ErrForbidden = e.New("forbidden access")

// ErrUnknownWOLType represents unknown wol type error
var ErrUnknownWOLType = e.New("unknown wol type")

// ErrUnknown represents unknown error
var ErrUnknown = e.New("unknown error")

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	} else if isUnauthorized(err) {
		return http.StatusUnauthorized
	} else if isBadRequest(err) {
		return http.StatusBadRequest
	} else if isForbidden(err) {
		return http.StatusForbidden
	} else {
		return http.StatusInternalServerError
	}
}

func isUnauthorized(err error) bool {
	return e.Is(err, ErrTokenNotFound)
}

func isBadRequest(err error) bool {
	return e.Is(err, ErrMultipleTokens) || e.Is(err, ErrInvalidAuthHeader) || e.Is(err, ErrUnknownWOLType)
}

func isForbidden(err error) bool {
	return e.Is(err, ErrForbidden)
}
