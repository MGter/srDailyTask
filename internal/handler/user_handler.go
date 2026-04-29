package handler

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"daily_task/internal/logger"
	"daily_task/internal/model"
	"daily_task/internal/service"
)

const maxAvatarSize = 2 << 20

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
		logger.Error("user_handler.go", 30, "Register failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info("user_handler.go", 35, "User registered: %s", user.Username)
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
		logger.Error("user_handler.go", 48, "Login failed: %v", err)
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	logger.Info("user_handler.go", 53, "User logged in: %s", user.Username)
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getExactUserIDFromPath(r.URL.Path)
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

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getExactUserIDFromPath(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req model.UpdateUserRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.UpdateProfile(id, &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	id, err := getAvatarUserIDFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAvatarSize+1024)
	if err := r.ParseMultipartForm(maxAvatarSize); err != nil {
		writeError(w, http.StatusBadRequest, "avatar is too large")
		return
	}

	existingUser, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if existingUser == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		writeError(w, http.StatusBadRequest, "avatar file is required")
		return
	}
	defer file.Close()

	if header.Size > maxAvatarSize {
		writeError(w, http.StatusBadRequest, "avatar is too large")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isAllowedAvatarExt(ext) || !isAllowedAvatarContent(file, ext) {
		writeError(w, http.StatusBadRequest, "unsupported avatar type")
		return
	}

	if err := os.MkdirAll("uploads/avatars", 0755); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	filename := buildAvatarFilename(id, ext)
	path := filepath.Join("uploads", "avatars", filename)
	dst, err := os.Create(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer dst.Close()

	if _, err := file.Seek(0, 0); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if _, err := dst.ReadFrom(file); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	avatarURL := "/uploads/avatars/" + filename
	user, err := h.svc.UpdateAvatar(id, avatarURL)
	if err != nil {
		_ = os.Remove(path)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	removeOldAvatar(existingUser.AvatarURL, avatarURL)

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

func getExactUserIDFromPath(path string) (uint64, error) {
	if !strings.HasPrefix(path, "/api/user/") {
		return 0, errors.New("invalid user path")
	}
	trimmed := strings.TrimPrefix(path, "/api/user/")
	trimmed = strings.TrimSuffix(trimmed, "/")
	if trimmed == "" || strings.Contains(trimmed, "/") {
		return 0, errors.New("invalid user path")
	}
	return strconv.ParseUint(trimmed, 10, 64)
}

func getAvatarUserIDFromPath(r *http.Request) (uint64, error) {
	if !strings.HasSuffix(r.URL.Path, "/avatar") {
		return 0, errors.New("invalid avatar path")
	}
	path := strings.TrimSuffix(r.URL.Path, "/avatar")
	return getExactUserIDFromPath(path)
}

func isAllowedAvatarExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return true
	default:
		return false
	}
}

func removeOldAvatar(oldURL, newURL string) {
	if oldURL == "" || oldURL == newURL || !strings.HasPrefix(oldURL, "/uploads/avatars/") {
		return
	}
	name := filepath.Base(oldURL)
	_ = os.Remove(filepath.Join("uploads", "avatars", name))
}

func isAllowedAvatarContent(file multipart.File, ext string) bool {
	sniff := make([]byte, 512)
	n, _ := file.Read(sniff)
	if _, err := file.Seek(0, 0); err != nil {
		return false
	}
	if n == 0 {
		return false
	}

	contentType := http.DetectContentType(sniff[:n])
	switch ext {
	case ".jpg", ".jpeg":
		return contentType == "image/jpeg"
	case ".png":
		return contentType == "image/png"
	case ".gif":
		return contentType == "image/gif"
	case ".webp":
		return isWebP(sniff[:n])
	default:
		return false
	}
}

func isWebP(data []byte) bool {
	return len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP"
}

func buildAvatarFilename(userID uint64, ext string) string {
	buf := make([]byte, 4)
	_, _ = rand.Read(buf)
	return strings.Join([]string{
		uintToString(userID),
		uintToString(uint64(time.Now().UnixNano())),
		hex.EncodeToString(buf),
	}, "_") + ext
}
