package service

import (
	"errors"
	"regexp"
	"strings"
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
		AvatarURL: "",
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

func (s *UserService) UpdateProfile(id uint64, req *model.UpdateUserRequest) (*model.User, error) {
	username := strings.TrimSpace(req.Username)
	email := strings.TrimSpace(req.Email)
	if username == "" {
		return nil, errors.New("username is required")
	}
	if email != "" && !isValidEmail(email) {
		return nil, errors.New("invalid email")
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	existing, err := s.repo.FindByUsernameExceptID(username, id)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	password := user.Password
	updatePassword := req.NewPassword != ""
	if updatePassword {
		if req.OldPassword == "" {
			return nil, errors.New("old password is required")
		}
		if req.OldPassword != user.Password {
			return nil, errors.New("old password is incorrect")
		}
		password = req.NewPassword
	}

	if err := s.repo.UpdateProfile(id, username, email, password, updatePassword); err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

func (s *UserService) UpdateAvatar(id uint64, avatarURL string) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if err := s.repo.UpdateAvatar(id, avatarURL); err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

func (s *UserService) UpdatePoints(id uint64, points int) error {
	return s.repo.UpdatePoints(id, points)
}

func (s *UserService) List(limit, offset int) ([]*model.User, error) {
	return s.repo.List(limit, offset)
}

func isValidEmail(email string) bool {
	pattern := `^[^@\s]+@[^@\s]+\.[^@\s]+$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}
