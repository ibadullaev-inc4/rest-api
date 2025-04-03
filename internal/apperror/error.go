package apperror

import "errors"

var (
	ErrNotFound              = errors.New("resource not found")
	ErrUnauthorized          = errors.New("unauthorized access")
	ErrInternalServer        = errors.New("internal server error")
	ErrMissingRequiredFields = errors.New("missing required fields")
	ErrInvalidUuidFormat     = errors.New("invalid UUID format")
)

func NewError(text string) error {
	return errors.New(text)
}
