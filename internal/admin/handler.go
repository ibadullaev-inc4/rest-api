package admin

import (
	"encoding/json"
	"net/http"
	"rest-api/internal/apperror"
	"rest-api/internal/handlers"

	"rest-api/internal/user"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ handlers.Handler = &handler{}

const (
	usersURL = "/admins"
	userURL  = "/admins/:uuid"
)

type handler struct {
	logger  *logrus.Logger
	storage user.Storage
}

func NewHandler(logger *logrus.Logger, storage user.Storage) handlers.Handler {
	return &handler{
		logger:  logger,
		storage: storage,
	}
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
	h.logger.Info("GetList called for users")

	users, err := h.storage.GetAll(r.Context())
	if err != nil {
		h.logger.Errorf("Failed to get users: %v", err)
		return apperror.ErrInternalServer
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		h.logger.Errorf("Failed to encode users list: %v", err)
		return apperror.ErrInternalServer
	}

	return nil
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var admin Admin

	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
		return apperror.NewError("invalid request body")
	}

	if admin.Email == "" || admin.Username == "" || admin.PasswordHash == "" {
		return apperror.ErrMissingRequiredFields
	}

	id, err := h.storage.Create(r.Context(), user.User(admin))
	if err != nil {
		h.logger.Errorf("Failed to create user: %v", err)
		return apperror.ErrInternalServer
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": id}); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
		return apperror.ErrInternalServer
	}

	return nil
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("GetUserByUUID called for user")

	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("uuid")

	user, err := h.storage.FindOne(r.Context(), id)
	if err != nil {
		h.logger.Errorf("Failed to find user by ID %s: %v", id, err)
		if err == mongo.ErrNoDocuments {
			return apperror.ErrNotFound
		}
		return apperror.ErrInternalServer
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.logger.Errorf("Failed to encode user: %v", err)
		return apperror.ErrInternalServer
	}

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
