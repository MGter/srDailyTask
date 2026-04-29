package repository

import (
	"database/sql"
	"time"

	"daily_task/internal/model"
)

type LongTermItemRepository struct{}

func NewLongTermItemRepository() *LongTermItemRepository {
	return &LongTermItemRepository{}
}

func (r *LongTermItemRepository) Create(item *model.LongTermItem) error {
	query := `INSERT INTO long_term_items (user_id, name, price, purchase_date, status, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := DB.Exec(query, item.UserID, item.Name, item.Price, item.PurchaseDate, item.Status, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	item.ID = uint64(id)
	return nil
}

func (r *LongTermItemRepository) FindByUserID(userID uint64) ([]*model.LongTermItem, error) {
	query := `SELECT id, user_id, name, price, purchase_date, scrap_date, frozen_daily_cost, status, created_at, updated_at
	          FROM long_term_items WHERE user_id = ? ORDER BY status ASC, purchase_date DESC, id DESC`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*model.LongTermItem{}
	for rows.Next() {
		item, err := scanLongTermItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *LongTermItemRepository) FindByID(id uint64) (*model.LongTermItem, error) {
	query := `SELECT id, user_id, name, price, purchase_date, scrap_date, frozen_daily_cost, status, created_at, updated_at
	          FROM long_term_items WHERE id = ?`
	row := DB.QueryRow(query, id)
	return scanLongTermItem(row)
}

func (r *LongTermItemRepository) Update(item *model.LongTermItem) error {
	query := `UPDATE long_term_items SET name = ?, price = ?, purchase_date = ?, frozen_daily_cost = ?, updated_at = ?
	          WHERE id = ? AND user_id = ?`
	var frozenDailyCost interface{}
	if item.FrozenDailyCost != nil {
		frozenDailyCost = *item.FrozenDailyCost
	}
	_, err := DB.Exec(query, item.Name, item.Price, item.PurchaseDate, frozenDailyCost, time.Now(), item.ID, item.UserID)
	return err
}

func (r *LongTermItemRepository) Scrap(id uint64, userID uint64, scrapDate time.Time, frozenDailyCost float64) error {
	query := `UPDATE long_term_items SET scrap_date = ?, frozen_daily_cost = ?, status = ?, updated_at = ?
	          WHERE id = ? AND user_id = ?`
	_, err := DB.Exec(query, scrapDate, frozenDailyCost, model.LongTermItemScrapped, time.Now(), id, userID)
	return err
}

func (r *LongTermItemRepository) Delete(id uint64, userID uint64) error {
	query := `DELETE FROM long_term_items WHERE id = ? AND user_id = ?`
	_, err := DB.Exec(query, id, userID)
	return err
}

type longTermItemScanner interface {
	Scan(dest ...interface{}) error
}

func scanLongTermItem(scanner longTermItemScanner) (*model.LongTermItem, error) {
	item := &model.LongTermItem{}
	var scrapDate sql.NullTime
	var frozenDailyCost sql.NullFloat64
	var status string

	err := scanner.Scan(
		&item.ID,
		&item.UserID,
		&item.Name,
		&item.Price,
		&item.PurchaseDate,
		&scrapDate,
		&frozenDailyCost,
		&status,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if scrapDate.Valid {
		item.ScrapDate = &scrapDate.Time
	}
	if frozenDailyCost.Valid {
		item.FrozenDailyCost = &frozenDailyCost.Float64
	}
	item.Status = model.LongTermItemStatus(status)
	return item, nil
}
