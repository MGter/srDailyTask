package service

import (
	"errors"
	"time"

	"daily_task/internal/model"
	"daily_task/internal/repository"
	"daily_task/internal/logger"
)

type WalletService struct {
	repo    *repository.WalletRepository
	userSvc *UserService
}

func NewWalletService() *WalletService {
	return &WalletService{
		repo:    repository.NewWalletRepository(),
		userSvc: NewUserService(),
	}
}

func (s *WalletService) Create(wallet *model.Wallet) error {
	if wallet.RecordTime.IsZero() {
		wallet.RecordTime = time.Now()
	}
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

	now := time.Now()
	wallet := &model.Wallet{
		UserID:      userID,
		Balance:     balance - amount,
		Type:        model.WalletSpend,
		Amount:      amount,
		Description: description,
		CreatedAt:   now,
		RecordTime:  now,
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

	logger.Info("wallet_service.go", 58, "User %d spent %d points for: %s", userID, amount, description)
	return wallet, nil
}

func (s *WalletService) AddRecord(req *model.AddRecordRequest) (*model.Wallet, error) {
	if req.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if req.Type != model.WalletEarn && req.Type != model.WalletSpend {
		return nil, errors.New("invalid type")
	}

	now := time.Now()
	recordTime := req.RecordTime
	if recordTime.IsZero() {
		recordTime = now
	}

	balance, _ := s.repo.GetBalance(req.UserID)
	if req.Type == model.WalletSpend {
		balance -= req.Amount
	} else {
		balance += req.Amount
	}

	wallet := &model.Wallet{
		UserID:      req.UserID,
		Balance:     balance,
		Type:        req.Type,
		Amount:      req.Amount,
		Description: req.Description,
		CreatedAt:   now,
		RecordTime:  recordTime,
	}

	if err := s.repo.Create(wallet); err != nil {
		return nil, err
	}

	// 更新用户积分
	user, err := s.userSvc.GetByID(req.UserID)
	if err != nil {
		return nil, err
	}
	var newPoints int
	if req.Type == model.WalletEarn {
		newPoints = user.Points + req.Amount
	} else {
		newPoints = user.Points - req.Amount
	}
	if err := s.userSvc.UpdatePoints(req.UserID, newPoints); err != nil {
		return nil, err
	}

	logger.Info("wallet_service.go", 100, "User %d added record: %s %d", req.UserID, req.Type, req.Amount)
	return wallet, nil
}

func (s *WalletService) Delete(id uint64, userID uint64) error {
	// 先获取记录信息以更新用户积分
	wallet, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if wallet == nil || wallet.UserID != userID {
		return errors.New("record not found or not owned by user")
	}

	if err := s.repo.Delete(id, userID); err != nil {
		return err
	}

	// 更新用户积分（反向操作）
	user, err := s.userSvc.GetByID(userID)
	if err != nil {
		return nil
	}
	var newPoints int
	if wallet.Type == model.WalletEarn {
		newPoints = user.Points - wallet.Amount
	} else {
		newPoints = user.Points + wallet.Amount
	}
	if err := s.userSvc.UpdatePoints(userID, newPoints); err != nil {
		return nil
	}

	logger.Info("wallet_service.go", 125, "User %d deleted record %d", userID, id)
	return nil
}

func (s *WalletService) List(limit, offset int) ([]*model.Wallet, error) {
	return s.repo.List(limit, offset)
}