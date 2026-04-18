package handler

import (
	"net/http"

	"daily_task/internal/model"
	"daily_task/internal/service"
)

type WalletHandler struct {
	svc *service.WalletService
}

func NewWalletHandler() *WalletHandler {
	return &WalletHandler{
		svc: service.NewWalletService(),
	}
}

func (h *WalletHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromPath(r, "/api/wallet/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	limit := getQueryParam(r, "limit", 10)
	offset := getQueryParam(r, "offset", 0)

	wallets, err := h.svc.GetByUserID(userID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, wallets)
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromPath(r, "/api/wallet/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	balance, err := h.svc.GetBalance(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"balance": balance})
}

func (h *WalletHandler) Spend(w http.ResponseWriter, r *http.Request) {
	var req model.SpendRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	wallet, err := h.svc.Spend(req.UserID, req.Amount, req.Description)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, wallet)
}

func (h *WalletHandler) AddRecord(w http.ResponseWriter, r *http.Request) {
	var req model.AddRecordRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	wallet, err := h.svc.AddRecord(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, wallet)
}

func (h *WalletHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// 从路径获取 id
	id, err := getURLParamFromPath(r.URL.Path, "/api/wallet/delete/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid record id")
		return
	}

	var req model.DeleteRecordRequest
	if err := readBody(r, &req); err != nil {
		req.ID = id
		req.UserID = 1 // 默认用户
	}
	req.ID = id

	if err := h.svc.Delete(req.ID, req.UserID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "record deleted"})
}