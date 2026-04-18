package service

import (
	"errors"
	"time"

	"daily_task/internal/model"
	"daily_task/internal/repository"
	"daily_task/internal/logger"
)

type WalletService struct {
	repo   *repository.WalletRepository
	userSvc *UserService
}

func NewWalletService() *WalletService {
	return &WalletService{
		repo:    repository.NewWalletRepository(),
		userSvc: NewUserService(),
	}
}

func (s *WalletService) Create(wallet *model.Wallet) error {
	return s.repo.Create(wallet)
}

func (s *WalletService) GetByUserID(userID uint64, limit, offset int) ([]*model.Wallet, error) {
	return s.repo.FindByUserID(userID, limit, offset)
}

func (s *WalletService) GetBalance(userID uint64) (int, error) {
	return s.repo.GetBalance(userID)
}

func (s *WalletService) Spend(userID uint64, amount int, description string) (*model.Wallet, error) {
	balance, err := s.repo.GetBalance(userID)
	if err != nil {
		return nil, err
	}
	if balance < amount {
		return nil, errors.New("insufficient balance")
	}

	wallet := &model.Wallet{
		UserID:      userID,
		Balance:     balance - amount,
		Type:        model.WalletSpend,
		Amount:      amount,
		Description: description,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(wallet); err != nil {
		return nil, err
	}

	user, err := s.userSvc.GetByID(userID)
	if err != nil {
		return nil, err
	}
	newPoints := user.Points - amount
	if err := s.userSvc.UpdatePoints(userID, newPoints); err != nil {
		return nil, err
	}

	logger.Info("wallet_service.go", 57, "User %d spent %d points for: %s", userID, amount, description)
	return wallet, nil
}

func (s *WalletService) List(limit, offset int) ([]*model.Wallet, error) {
	return s.repo.List(limit, offset)
}