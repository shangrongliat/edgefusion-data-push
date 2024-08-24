package common

import (
	"net/http"
)

// Code code
type Code string

func (c Code) String() string {
	if msg, ok := templates[c]; ok {
		return msg
	}
	return templates[ErrUnknown]
}

// all codes
const (

	// * request
	ErrRequestAccessDenied   = "ErrRequestAccessDenied"
	ErrRequestMethodNotFound = "ErrRequestMethodNotFound"
	// * resource
	ErrResourceNotFound    = "ErrResourceNotFound"
	ErrResourceHasBeenUsed = "ErrResourceHasBeenUsed"
	// * unknown
	ErrUnknown = "UnknownError"
)

var templates = map[Code]string{}

func GetHTTPStatus(c Code) int {
	switch c {
	case ErrResourceNotFound, ErrRequestMethodNotFound:
		return http.StatusNotFound
	case ErrRequestAccessDenied:
		return http.StatusUnauthorized
	case ErrResourceHasBeenUsed:
		return http.StatusForbidden
	case ErrUnknown:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}
