package domain

import "errors"

var (
	ErrValidation          = errors.New("validation error")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrNotFound            = errors.New("not found")
	ErrUpstreamUnavailable = errors.New("upstream unavailable")
	ErrInternal            = errors.New("internal")
)
