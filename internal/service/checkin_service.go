package service

import (
	"time"

	"daily_task/internal/model"
	"daily_task/internal/repository"
	"daily_task/internal/logger"
)

type CheckInService struct {
	repo      *repository.CheckInRepository
	walletSvc *WalletService
}

func NewCheckInService() *CheckInService {
	return &CheckInService{
		repo:      repository.NewCheckInRepository(),
		walletSvc: NewWalletService(),
	}
}

func (s *CheckInService) Create(taskID, userID uint64, points int) (*model.CheckIn, error) {
	checkin := &model.CheckIn{
		TaskID:    taskID,
		UserID:    userID,
		Points:    points,
		CheckTime: time.Now(),
	}

	if err := s.repo.Create(checkin); err != nil {
		return nil, err
	}

	wallet := &model.Wallet{
		UserID:      userID,
		Balance:     points,
		Type:        model.WalletEarn,
		Amount:      points,
		Description: "Check-in reward",
		CreatedAt:   time.Now(),
	}
	if err := s.walletSvc.Create(wallet); err != nil {
		logger.Error("checkin_service.go", 36, "Failed to create wallet record: %v", err)
	}

	logger.Info("checkin_service.go", 39, "User %d checked in task %d, earned %d points", userID, taskID, points)
	return checkin, nil
}

func (s *CheckInService) GetByUserID(userID uint64, limit, offset int) ([]*model.CheckIn, error) {
	return s.repo.FindByUserID(userID, limit, offset)
}

func (s *CheckInService) GetByTaskID(taskID uint64, limit, offset int) ([]*model.CheckIn, error) {
	return s.repo.FindByTaskID(taskID, limit, offset)
}

func (s *CheckInService) GetTodayByTaskID(taskID uint64) (*model.CheckIn, error) {
	return s.repo.FindTodayByTaskID(taskID)
}