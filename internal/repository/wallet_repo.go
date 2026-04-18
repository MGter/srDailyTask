package repository

import (
	"daily_task/internal/model"
)

type WalletRepository struct{}

func NewWalletRepository() *WalletRepository {
	return &WalletRepository{}
}

func (r *WalletRepository) Create(wallet *model.Wallet) error {
	query := `INSERT INTO wallet (user_id, balance, type, amount, description, created_at, record_time)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := DB.Exec(query, wallet.UserID, wallet.Balance, wallet.Type,
		wallet.Amount, wallet.Description, wallet.CreatedAt, wallet.RecordTime)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	wallet.ID = uint64(id)
	return nil
}

func (r *WalletRepository) FindByUserID(userID uint64, limit, offset int) ([]*model.Wallet, error) {
	query := `SELECT id, user_id, balance, type, amount, description, created_at, record_time
	          FROM wallet WHERE user_id = ? ORDER BY record_time DESC LIMIT ? OFFSET ?`
	rows, err := DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wallets := []*model.Wallet{}
	for rows.Next() {
		w := &model.Wallet{}
		err := rows.Scan(&w.ID, &w.UserID, &w.Balance, &w.Type, &w.Amount, &w.Description, &w.CreatedAt, &w.RecordTime)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}
	return wallets, nil
}

func (r *WalletRepository) FindByID(id uint64) (*model.Wallet, error) {
	w := &model.Wallet{}
	query := `SELECT id, user_id, balance, type, amount, description, created_at, record_time
	          FROM wallet WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&w.ID, &w.UserID, &w.Balance, &w.Type, &w.Amount, &w.Description, &w.CreatedAt, &w.RecordTime)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *WalletRepository) Delete(id uint64, userID uint64) error {
	query := `DELETE FROM wallet WHERE id = ? AND user_id = ?`
	_, err := DB.Exec(query, id, userID)
	return err
}

func (r *WalletRepository) GetBalance(userID uint64) (int, error) {
	var balance int
	query := `SELECT COALESCE(SUM(CASE WHEN type = 'earn' THEN amount ELSE -amount END), 0)
	          FROM wallet WHERE user_id = ?`
	err := DB.QueryRow(query, userID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *WalletRepository) List(limit, offset int) ([]*model.Wallet, error) {
	query := `SELECT id, user_id, balance, type, amount, description, created_at, record_time
	          FROM wallet ORDER BY record_time DESC LIMIT ? OFFSET ?`
	rows, err := DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wallets := []*model.Wallet{}
	for rows.Next() {
		w := &model.Wallet{}
		err := rows.Scan(&w.ID, &w.UserID, &w.Balance, &w.Type, &w.Amount, &w.Description, &w.CreatedAt, &w.RecordTime)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}
	return wallets, nil
}