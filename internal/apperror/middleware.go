package apperror

import (
	"log"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func ErrorMiddleware(next appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err != nil {
			log.Println("Error:", err)

			switch err {
			case ErrMissingRequiredFields:
				http.Error(w, err.Error(), http.StatusNotFound)
			case ErrInvalidUuidFormat:
				http.Error(w, err.Error(), http.StatusNotFound)
			case ErrNotFound:
				http.Error(w, err.Error(), http.StatusNotFound)
			case ErrUnauthorized:
				http.Error(w, err.Error(), http.StatusUnauthorized)
			default:
				http.Error(w, ErrInternalServer.Error(), http.StatusInternalServerError)
			}
		}
	}
}
