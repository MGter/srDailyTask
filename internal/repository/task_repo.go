package repository

import (
	"database/sql"

	"daily_task/internal/model"
)

type TaskRepository struct{}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

func (r *TaskRepository) Create(task *model.Task) error {
	query := `INSERT INTO tasks (user_id, title, description, circle_mode, level, points, created_at, updated_at, is_expired)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := DB.Exec(query, task.UserID, task.Title, task.Description,
		task.CircleMode, task.Level, task.Points, task.CreatedAt, task.UpdatedAt, task.IsExpired)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = uint64(id)
	return nil
}

func (r *TaskRepository) FindByID(id uint64) (*model.Task, error) {
	task := &model.Task{}
	query := `SELECT id, user_id, title, description, circle_mode, level, points, created_at, updated_at, is_expired
	          FROM tasks WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(
		&task.ID, &task.UserID, &task.Title, &task.Description,
		&task.CircleMode, &task.Level, &task.Points, &task.CreatedAt, &task.UpdatedAt, &task.IsExpired)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *TaskRepository) FindByUserID(userID uint64, limit, offset int) ([]*model.Task, error) {
	query := `SELECT id, user_id, title, description, circle_mode, level, points, created_at, updated_at, is_expired
	          FROM tasks WHERE user_id = ? ORDER BY level DESC, id DESC LIMIT ? OFFSET ?`
	rows, err := DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*model.Task{}
	for rows.Next() {
		task := &model.Task{}
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description,
			&task.CircleMode, &task.Level, &task.Points, &task.CreatedAt, &task.UpdatedAt, &task.IsExpired)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) Update(task *model.Task) error {
	query := `UPDATE tasks SET title = ?, description = ?, circle_mode = ?, level = ?, points = ?, updated_at = ?, is_expired = ?
	          WHERE id = ?`
	_, err := DB.Exec(query, task.Title, task.Description, task.CircleMode,
		task.Level, task.Points, task.UpdatedAt, task.IsExpired, task.ID)
	return err
}

func (r *TaskRepository) Delete(id uint64) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := DB.Exec(query, id)
	return err
}

func (r *TaskRepository) FindActiveTasks() ([]*model.Task, error) {
	query := `SELECT id, user_id, title, description, circle_mode, level, points, created_at, updated_at, is_expired
	          FROM tasks WHERE is_expired = false`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*model.Task{}
	for rows.Next() {
		task := &model.Task{}
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description,
			&task.CircleMode, &task.Level, &task.Points, &task.CreatedAt, &task.UpdatedAt, &task.IsExpired)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}