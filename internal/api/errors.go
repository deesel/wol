package api

import (
	e "errors"
	"net/http"
)

var ErrTokenNotFound = e.New("required auth token not found")
var ErrMultipleTokens = e.New("multiple auth tokens found")
var ErrInvalidAuthHeader = e.New("invalid authorization header format")
var ErrForbidden = e.New("forbidden access")
var ErrUnknownWOLType = e.New("unknown wol type")
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
