package handler

import (
	"net/http"

	"daily_task/internal/service"
)

type PointHandler struct {
	checkinSvc *service.CheckInService
	walletSvc  *service.WalletService
}

func NewPointHandler() *PointHandler {
	return &PointHandler{
		checkinSvc: service.NewCheckInService(),
		walletSvc:  service.NewWalletService(),
	}
}

func (h *PointHandler) GetCheckIns(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromPath(r, "/api/checkin/user/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	limit := getQueryParam(r, "limit", 10)
	offset := getQueryParam(r, "offset", 0)

	checkins, err := h.checkinSvc.GetByUserID(userID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, checkins)
}

func (h *PointHandler) GetPointHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromPath(r, "/api/points/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	limit := getQueryParam(r, "limit", 10)
	offset := getQueryParam(r, "offset", 0)

	wallets, err := h.walletSvc.GetByUserID(userID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, wallets)
}