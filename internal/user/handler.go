package user

import (
	"encoding/json"
	"net/http"
	"rest-api/internal/handlers"
	"rest-api/internal/storage"
	"rest-api/pkg/metrics"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ handlers.Handler = &handler{}

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

type handler struct {
	logger  *logrus.Logger
	storage storage.Storage
}

func NewHandler(logger *logrus.Logger, storage storage.Storage) handlers.Handler {
	return &handler{
		logger:  logger,
		storage: storage,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(usersURL, metrics.PrometheusMiddleware(h.GetList, usersURL))
	router.POST(usersURL, metrics.PrometheusMiddleware(h.CreateUser, usersURL))
	router.GET(userURL, metrics.PrometheusMiddleware(h.GetUserByUUID, usersURL))
	router.PUT(userURL, h.UpdateUser)
	router.PATCH(userURL, h.PartiallyUpdateUser)
	router.DELETE(userURL, h.DeleteUser)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h.logger.Info("GetList called for users")

	users, err := h.storage.GetAll(r.Context())
	if err != nil {
		h.logger.Errorf("Failed to get users: %v", err)
		http.Error(w, "failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		h.logger.Errorf("Failed to encode users list: %v", err)
		http.Error(w, "failed to encode users", http.StatusInternalServerError)
	}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var user storage.Client
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	id, err := h.storage.Create(r.Context(), user)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h.logger.Info("GetUserByUUID called for user")

	id := params.ByName("uuid")

	user, err := h.storage.FindOne(r.Context(), id)
	if err != nil {
		h.logger.Errorf("Failed to find user by ID %s: %v", id, err)
		if err == mongo.ErrNoDocuments {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to get user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.logger.Errorf("Failed to encode user: %v", err)
		http.Error(w, "failed to encode user", http.StatusInternalServerError)
	}
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h.logger.Info("UpdateUser called for user")

	id := params.ByName("uuid")
	h.logger.Infof("Attempting to update user with id: %s", id)

	var user storage.Client
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.Errorf("Invalid request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user.ID = id

	h.logger.Infof("User data to be updated: %+v", user)

	err := h.storage.Update(r.Context(), user)
	if err != nil {
		h.logger.Errorf("Failed to update user %s: %v", id, err)
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	h.logger.Info("User updated successfully")
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h.logger.Info("PartiallyUpdateUser called for user")

	id := params.ByName("uuid")

	var user storage.Client
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.Errorf("Invalid request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user.ID = id

	err := h.storage.PartiallyUpdate(r.Context(), user)
	if err != nil {
		h.logger.Errorf("Failed to partially update user %s: %v", id, err)
		http.Error(w, "failed to partially update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("uuid")
	h.logger.Infof("Attempting to delete user with id: %s", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.logger.Errorf("Invalid UUID format: %v", err)
		http.Error(w, "invalid UUID format", http.StatusBadRequest)
		return
	}

	idString := objectID.Hex()

	err = h.storage.Delete(r.Context(), idString)
	if err != nil {
		h.logger.Errorf("Failed to delete user %s: %v", id, err)
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	h.logger.Infof("User %s deleted successfully", id)
	w.WriteHeader(http.StatusNoContent)
}
