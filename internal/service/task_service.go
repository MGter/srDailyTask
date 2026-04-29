package service

import (
	"errors"
	"time"

	"daily_task/internal/logger"
	"daily_task/internal/model"
	"daily_task/internal/repository"
)

type TaskService struct {
	repo       *repository.TaskRepository
	checkinSvc *CheckInService
	userSvc    *UserService
	walletSvc  *WalletService
}

func NewTaskService() *TaskService {
	return &TaskService{
		repo:       repository.NewTaskRepository(),
		checkinSvc: NewCheckInService(),
		userSvc:    NewUserService(),
		walletSvc:  NewWalletService(),
	}
}

func (s *TaskService) Create(req *model.CreateTaskRequest) (*model.Task, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.Points <= 0 {
		req.Points = 10
	}
	if req.Level < 1 || req.Level > 3 {
		req.Level = model.LevelLow
	}

	task := &model.Task{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		CircleMode:  req.CircleMode,
		Level:       req.Level,
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
	if req.Level >= 1 && req.Level <= 3 {
		task.Level = req.Level
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

func (s *TaskService) GetActiveTasks() ([]*model.Task, error) {
	return s.repo.FindActiveTasks()
}

// ShouldCheckinToday 判断任务今天是否需要打卡
func (s *TaskService) ShouldCheckinToday(task *model.Task) bool {
	if task.IsExpired {
		return false
	}

	now := time.Now()
	weekday := now.Weekday()

	switch task.CircleMode {
	case model.CircleOnce:
		// 单次任务：只要没过期就需要打卡（但打卡一次后就过期了）
		return true
	case model.CircleWeekly:
		// 每周：每周一打卡
		return weekday == time.Monday
	case model.CircleWorkday:
		// 工作日：周一到周五
		return weekday >= time.Monday && weekday <= time.Friday
	case model.CircleWeekend:
		// 周末：周六周日
		return weekday == time.Saturday || weekday == time.Sunday
	case model.CircleCustom:
		// 自定义：默认每天都需要（后续可扩展）
		return true
	default:
		return false
	}
}

// CheckIn 打卡，单次任务打卡后标记为过期
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

	// 检查今天是否需要打卡
	if !s.ShouldCheckinToday(task) {
		return nil, errors.New("today is not a check-in day for this task")
	}

	existing, err := s.checkinSvc.GetTodayByTaskID(taskID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("already checked in today")
	}

	checkin, err := s.checkinSvc.Create(taskID, userID, task.Points, task.Title)
	if err != nil {
		return nil, err
	}

	// 单次任务打卡后标记为过期
	if task.CircleMode == model.CircleOnce {
		task.IsExpired = true
		task.UpdatedAt = time.Now()
		s.repo.Update(task)
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

// CancelCheckIn 取消打卡，删除打卡记录和钱包记录
func (s *TaskService) CancelCheckIn(taskID, userID uint64) error {
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}
	if task.UserID != userID {
		return errors.New("task does not belong to user")
	}

	// 找到今天的打卡记录
	checkin, err := s.checkinSvc.GetTodayByTaskID(taskID)
	if err != nil {
		return err
	}
	if checkin == nil {
		return errors.New("not checked in today")
	}

	// 删除打卡记录
	if err := s.checkinSvc.Delete(checkin.ID); err != nil {
		return err
	}

	// 找到对应的钱包记录并删除
	if err := s.walletSvc.DeleteByCheckinID(checkin.ID); err != nil {
		logger.Error("task_service.go", 195, "Failed to delete wallet record: %v", err)
	}

	// 更新用户积分（只有实际打卡才扣减积分，跳过的不扣）
	if checkin.Points > 0 {
		user, err := s.userSvc.GetByID(userID)
		if err != nil {
			return err
		}
		newPoints := user.Points - checkin.Points
		if err := s.userSvc.UpdatePoints(userID, newPoints); err != nil {
			return err
		}
	}

	logger.Info("task_service.go", 205, "User %d cancelled checkin for task %d", userID, taskID)
	return nil
}

// SkipCheckIn 跳过打卡，记录跳过但不获得积分
func (s *TaskService) SkipCheckIn(taskID, userID uint64) (*model.CheckIn, error) {
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

	// 检查今天是否需要打卡
	if !s.ShouldCheckinToday(task) {
		return nil, errors.New("today is not a check-in day for this task")
	}

	existing, err := s.checkinSvc.GetTodayByTaskID(taskID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("already checked in today")
	}

	// 创建跳过记录，积分为0
	checkin, err := s.checkinSvc.Create(taskID, userID, 0, task.Title)
	if err != nil {
		return nil, err
	}

	logger.Info("task_service.go", 240, "User %d skipped task %d", userID, taskID)
	return checkin, nil
}