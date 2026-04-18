package service

import (
	"errors"
	"time"

	"daily_task/internal/model"
	"daily_task/internal/repository"
)

type TaskService struct {
	repo       *repository.TaskRepository
	checkinSvc *CheckInService
	userSvc    *UserService
}

func NewTaskService() *TaskService {
	return &TaskService{
		repo:       repository.NewTaskRepository(),
		checkinSvc: NewCheckInService(),
		userSvc:    NewUserService(),
	}
}

func (s *TaskService) Create(req *model.CreateTaskRequest) (*model.Task, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.Points <= 0 {
		req.Points = 10
	}

	task := &model.Task{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		CircleMode:  req.CircleMode,
		Points:      req.Points,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsExpired:   false,
	}

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetByID(id uint64) (*model.Task, error) {
	return s.repo.FindByID(id)
}

func (s *TaskService) GetByUserID(userID uint64, limit, offset int) ([]*model.Task, error) {
	return s.repo.FindByUserID(userID, limit, offset)
}

func (s *TaskService) Update(id uint64, req *model.UpdateTaskRequest) (*model.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.CircleMode != "" {
		task.CircleMode = req.CircleMode
	}
	if req.Points > 0 {
		task.Points = req.Points
	}
	task.IsExpired = req.IsExpired
	task.UpdatedAt = time.Now()

	if err := s.repo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

func (s *TaskService) CheckIn(taskID, userID uint64) (*model.CheckIn, error) {
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}
	if task.UserID != userID {
		return nil, errors.New("task does not belong to user")
	}

	existing, err := s.checkinSvc.GetTodayByTaskID(taskID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("already checked in today")
	}

	checkin, err := s.checkinSvc.Create(taskID, userID, task.Points)
	if err != nil {
		return nil, err
	}

	user, err := s.userSvc.GetByID(userID)
	if err != nil {
		return nil, err
	}
	newPoints := user.Points + task.Points
	if err := s.userSvc.UpdatePoints(userID, newPoints); err != nil {
		return nil, err
	}

	return checkin, nil
}

func (s *TaskService) GetActiveTasks() ([]*model.Task, error) {
	return s.repo.FindActiveTasks()
}