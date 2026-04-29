package handler

import (
	"net/http"
	"strings"

	"daily_task/internal/model"
	"daily_task/internal/service"
)

type LongTermItemHandler struct {
	svc *service.LongTermItemService
}

func NewLongTermItemHandler() *LongTermItemHandler {
	return &LongTermItemHandler{svc: service.NewLongTermItemService()}
}

func (h *LongTermItemHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromPath(r, "/api/long-term-items/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	items, summary, err := h.svc.GetByUserID(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, &model.LongTermItemListResponse{Items: items, Summary: summary})
}

func (h *LongTermItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateLongTermItemRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.svc.Create(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *LongTermItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getLongTermItemID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var req model.UpdateLongTermItemRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Update(id, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "item updated"})
}

func (h *LongTermItemHandler) Scrap(w http.ResponseWriter, r *http.Request) {
	id, err := getLongTermItemID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var req model.ScrapLongTermItemRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Scrap(id, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "item scrapped"})
}

func (h *LongTermItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getLongTermItemID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var req model.DeleteLongTermItemRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Delete(id, req.UserID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "item deleted"})
}

func getLongTermItemID(path string) (uint64, error) {
	path = strings.TrimSuffix(path, "/scrap")
	return getURLParamFromPath(path, "/api/long-term-items/")
}
