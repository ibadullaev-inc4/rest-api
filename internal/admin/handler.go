package admin

import (
	"net/http"
	"rest-api/internal/apperror"
	"rest-api/internal/handlers"

	"github.com/julienschmidt/httprouter"
)

var _ handlers.Handler = &handler{}

const (
	usersURL = "/admins"
	userURL  = "/admins/:uuid"
)

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, apperror.ErrorMiddleware(h.GetList))
	router.HandlerFunc(http.MethodPost, usersURL, apperror.ErrorMiddleware(h.CreateUser))
	router.HandlerFunc(http.MethodGet, userURL, apperror.ErrorMiddleware(h.GetUserByUUID))
	router.HandlerFunc(http.MethodPut, userURL, apperror.ErrorMiddleware(h.UpdateUser))
	router.HandlerFunc(http.MethodPatch, userURL, apperror.ErrorMiddleware(h.PartiallyUpdateUser))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.ErrorMiddleware(h.DeleteUser))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	return apperror.ErrNotFound
}
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(201)
	w.Write([]byte("this is create the admin"))
	return nil
}
func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(200)
	w.Write([]byte("this is get the admin"))
	return nil
}
func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(204)
	w.Write([]byte("this is update the admin"))
	return nil
}
func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(204)
	w.Write([]byte("this is partially update the admin"))
	return nil
}
func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(204)
	w.Write([]byte("this is delete the admin"))
	return nil
}
