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

func (s *CheckInService) Create(taskID, userID uint64, points int, taskTitle string) (*model.CheckIn, error) {
	now := time.Now()
	checkin := &model.CheckIn{
		TaskID:    taskID,
		UserID:    userID,
		Points:    points,
		CheckTime: now,
	}

	if err := s.repo.Create(checkin); err != nil {
		return nil, err
	}

    if points > 0 {
        desc := "打卡: " + taskTitle
        wallet := &model.Wallet{
            UserID:      userID,
            CheckinID:   checkin.ID, // 关联打卡记录
            Balance:     points,
            Type:        model.WalletEarn,
            Amount:      points,
            Description: desc,
            CreatedAt:   now,
            RecordTime:  now,
        }
        if err := s.walletSvc.Create(wallet); err != nil {
            logger.Error("checkin_service.go", 42, "Failed to create wallet record: %v", err)
        }
    }

    if points > 0 {
        logger.Info("checkin_service.go", 46, "User %d checked in task %d, earned %d points", userID, taskID, points)
    } else {
        logger.Info("checkin_service.go", 46, "User %d skipped task %d", userID, taskID)
    }
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

func (s *CheckInService) GetTodayCheckedTaskIDs(userID uint64) ([]uint64, error) {
	return s.repo.FindTodayCheckedTaskIDs(userID)
}

func (s *CheckInService) Delete(id uint64) error {
	return s.repo.Delete(id)
}