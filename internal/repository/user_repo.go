package repository

import (
	"database/sql"

	"daily_task/internal/model"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *model.User) error {
	query := `INSERT INTO users (username, password, email, avatar_url, points, created_at)
	          VALUES (?, ?, ?, ?, ?, ?)`
	result, err := DB.Exec(query, user.Username, user.Password, user.Email, user.AvatarURL, user.Points, user.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = uint64(id)
	return nil
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, password, email, COALESCE(avatar_url, ''), points, created_at
	          FROM users WHERE username = ?`
	err := DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.AvatarURL, &user.Points, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(id uint64) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, password, email, COALESCE(avatar_url, ''), points, created_at
	          FROM users WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.AvatarURL, &user.Points, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByUsernameExceptID(username string, id uint64) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, password, email, COALESCE(avatar_url, ''), points, created_at
	          FROM users WHERE username = ? AND id != ?`
	err := DB.QueryRow(query, username, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.AvatarURL, &user.Points, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateProfile(id uint64, username, email, password string, updatePassword bool) error {
	if updatePassword {
		query := `UPDATE users SET username = ?, email = ?, password = ? WHERE id = ?`
		_, err := DB.Exec(query, username, email, password, id)
		return err
	}
	query := `UPDATE users SET username = ?, email = ? WHERE id = ?`
	_, err := DB.Exec(query, username, email, id)
	return err
}

func (r *UserRepository) UpdateAvatar(id uint64, avatarURL string) error {
	query := `UPDATE users SET avatar_url = ? WHERE id = ?`
	_, err := DB.Exec(query, avatarURL, id)
	return err
}

func (r *UserRepository) UpdatePoints(id uint64, points int) error {
	query := `UPDATE users SET points = ? WHERE id = ?`
	_, err := DB.Exec(query, points, id)
	return err
}

func (r *UserRepository) List(limit, offset int) ([]*model.User, error) {
	query := `SELECT id, username, password, email, COALESCE(avatar_url, ''), points, created_at
	          FROM users ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*model.User{}
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.AvatarURL, &user.Points, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
