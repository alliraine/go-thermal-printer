package common

import "net/http"

type AppError interface {
	error
	HttpStatusCode() int
}

type InvalidAPIKeyError struct{}

func (e *InvalidAPIKeyError) Error() string {
	return "invalid API key"
}

func (e *InvalidAPIKeyError) HttpStatusCode() int {
	return http.StatusUnauthorized
}
