package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"

	"daily_task/internal/logger"
	"daily_task/internal/model"
	"daily_task/internal/repository"
)

type Scheduler struct {
	cron    *cron.Cron
	taskRepo *repository.TaskRepository
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cron:     cron.New(),
		taskRepo: repository.NewTaskRepository(),
	}
}

func (s *Scheduler) Start() {
	s.cron.AddFunc("* * * * *", s.checkReminders)
	s.cron.Start()
	logger.Info("scheduler.go", 23, "Scheduler started")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	logger.Info("scheduler.go", 27, "Scheduler stopped")
}

func (s *Scheduler) checkReminders() {
	tasks, err := s.taskRepo.FindActiveTasks()
	if err != nil {
		logger.Error("scheduler.go", 32, "Failed to get active tasks: %v", err)
		return
	}

	now := time.Now()
	for _, task := range tasks {
		if s.shouldRemind(task, now) {
			s.sendReminder(task)
		}
	}
}

func (s *Scheduler) shouldRemind(task *model.Task, now time.Time) bool {
	switch task.CircleMode {
	case model.CircleOnce:
		return true
	case model.CircleWeekly:
		return now.Weekday() == time.Monday
	case model.CircleWorkday:
		wd := now.Weekday()
		return wd >= time.Monday && wd <= time.Friday
	case model.CircleWeekend:
		wd := now.Weekday()
		return wd == time.Saturday || wd == time.Sunday
	case model.CircleCustom:
		return true
	default:
		return false
	}
}

func (s *Scheduler) sendReminder(task *model.Task) {
	logger.Info("scheduler.go", 62, "Reminder: Task %d - %s (User %d)", task.ID, task.Title, task.UserID)
}

type ReminderNotifier interface {
	Notify(task *model.Task) error
}

type LogNotifier struct{}

func (n *LogNotifier) Notify(task *model.Task) error {
	logger.Info("scheduler.go", 72, "[NOTIFY] Task reminder: %s", task.Title)
	return nil
}