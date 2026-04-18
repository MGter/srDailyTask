package service

import (
	"errors"
	"time"

	"daily_task/internal/model"
	"daily_task/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: repository.NewUserRepository(),
	}
}

func (s *UserService) Register(req *model.CreateUserRequest) (*model.User, error) {
	existing, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	user := &model.User{
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		Points:    0,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(req *model.LoginRequest) (*model.User, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if user.Password != req.Password {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

func (s *UserService) GetByID(id uint64) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) UpdatePoints(id uint64, points int) error {
	return s.repo.UpdatePoints(id, points)
}

func (s *UserService) List(limit, offset int) ([]*model.User, error) {
	return s.repo.List(limit, offset)
}