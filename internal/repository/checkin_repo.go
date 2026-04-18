package repository

import (
	"database/sql"
	"time"

	"daily_task/internal/model"
)

type CheckInRepository struct{}

func NewCheckInRepository() *CheckInRepository {
	return &CheckInRepository{}
}

func (r *CheckInRepository) Create(checkin *model.CheckIn) error {
	query := `INSERT INTO checkins (task_id, user_id, points, check_time)
	          VALUES (?, ?, ?, ?)`
	result, err := DB.Exec(query, checkin.TaskID, checkin.UserID, checkin.Points, checkin.CheckTime)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	checkin.ID = uint64(id)
	return nil
}

func (r *CheckInRepository) FindByUserID(userID uint64, limit, offset int) ([]*model.CheckIn, error) {
	query := `SELECT id, task_id, user_id, points, check_time
	          FROM checkins WHERE user_id = ? ORDER BY check_time DESC LIMIT ? OFFSET ?`
	rows, err := DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checkins := []*model.CheckIn{}
	for rows.Next() {
		c := &model.CheckIn{}
		err := rows.Scan(&c.ID, &c.TaskID, &c.UserID, &c.Points, &c.CheckTime)
		if err != nil {
			return nil, err
		}
		checkins = append(checkins, c)
	}
	return checkins, nil
}

func (r *CheckInRepository) FindByTaskID(taskID uint64, limit, offset int) ([]*model.CheckIn, error) {
	query := `SELECT id, task_id, user_id, points, check_time
	          FROM checkins WHERE task_id = ? ORDER BY check_time DESC LIMIT ? OFFSET ?`
	rows, err := DB.Query(query, taskID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checkins := []*model.CheckIn{}
	for rows.Next() {
		c := &model.CheckIn{}
		err := rows.Scan(&c.ID, &c.TaskID, &c.UserID, &c.Points, &c.CheckTime)
		if err != nil {
			return nil, err
		}
		checkins = append(checkins, c)
	}
	return checkins, nil
}

func (r *CheckInRepository) FindTodayByTaskID(taskID uint64) (*model.CheckIn, error) {
	c := &model.CheckIn{}
	today := time.Now().Format("2006-01-02")
	query := `SELECT id, task_id, user_id, points, check_time
	          FROM checkins WHERE task_id = ? AND DATE(check_time) = ?`
	err := DB.QueryRow(query, taskID, today).Scan(
		&c.ID, &c.TaskID, &c.UserID, &c.Points, &c.CheckTime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

// FindTodayCheckedTaskIDs 返回用户今天已打卡的任务ID列表
func (r *CheckInRepository) FindTodayCheckedTaskIDs(userID uint64) ([]uint64, error) {
	today := time.Now().Format("2006-01-02")
	query := `SELECT DISTINCT task_id FROM checkins WHERE user_id = ? AND DATE(check_time) = ?`
	rows, err := DB.Query(query, userID, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	taskIDs := []uint64{}
	for rows.Next() {
		var taskID uint64
		err := rows.Scan(&taskID)
		if err != nil {
			return nil, err
		}
		taskIDs = append(taskIDs, taskID)
	}
	return taskIDs, nil
}