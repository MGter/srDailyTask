package handler

import (
	"net/http"

	"daily_task/internal/logger"
	"daily_task/internal/model"
	"daily_task/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		svc: service.NewUserService(),
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.Register(&req)
	if err != nil {
		logger.Error("user_handler.go", 24, "Register failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info("user_handler.go", 28, "User registered: %s", user.Username)
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.Login(&req)
	if err != nil {
		logger.Error("user_handler.go", 40, "Login failed: %v", err)
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	logger.Info("user_handler.go", 44, "User logged in: %s", user.Username)
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := getQueryParam(r, "limit", 10)
	offset := getQueryParam(r, "offset", 0)

	users, err := h.svc.List(limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, users)
}