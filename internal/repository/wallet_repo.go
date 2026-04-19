package repository

import (
	"daily_task/internal/logger"
	"daily_task/internal/model"
	"time"
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

// DailyStats 每日积分统计
type DailyStats struct {
	Date    string `json:"date"`
	Earn    int    `json:"earn"`
	Spend   int    `json:"spend"`
	Balance int    `json:"balance"`
}

// GetDailyStats 获取最近 N 天的积分统计
// 包含 wallet 表和 checkins 表的数据
func (r *WalletRepository) GetDailyStats(userID uint64, days int) ([]*DailyStats, error) {
	// 生成最近 N 天的日期列表
	stats := []*DailyStats{}
	now := time.Now()

	// 查询 checkins 表按日期聚合
	checkinQuery := `
		SELECT DATE(check_time) as date, SUM(points) as earn
		FROM checkins WHERE user_id = ?
		AND check_time >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
		GROUP BY DATE(check_time)`
	checkinRows, err := DB.Query(checkinQuery, userID, days)
	if err != nil {
		logger.Error("wallet_repo.go", 130, "Checkin query error: %v", err)
		return nil, err
	}
	logger.Info("wallet_repo.go", 133, "Checkin query executed for user %d, days %d", userID, days)

	walletMap := map[string]*DailyStats{}
	for checkinRows.Next() {
		var dateStr string
		var earn int
		err := checkinRows.Scan(&dateStr, &earn)
		if err != nil {
			checkinRows.Close()
			return nil, err
		}
		// 解析日期字符串，可能是 time.Time 格式或字符串
		var date time.Time
		if len(dateStr) > 10 {
			// 可能是 time.Time 格式的字符串
			date, err = time.Parse("2006-01-02T15:04:05Z07:00", dateStr)
			if err != nil {
				// 尝试其他格式
				date, err = time.Parse("2006-01-02", dateStr[:10])
				if err != nil {
					checkinRows.Close()
					return nil, err
				}
			}
			dateStr = date.Format("2006-01-02")
		}
		logger.Info("wallet_repo.go", 165, "Found checkin: date=%s, earn=%d", dateStr, earn)
		walletMap[dateStr] = &DailyStats{Date: dateStr, Earn: earn, Spend: 0}
	}
	checkinRows.Close()
	logger.Info("wallet_repo.go", 149, "Total checkin records: %d", len(walletMap))

	// 查询 wallet 表按日期聚合
	walletQuery := `
		SELECT DATE(record_time) as date,
		       SUM(CASE WHEN type = 'earn' THEN amount ELSE 0 END) as earn,
		       SUM(CASE WHEN type = 'spend' THEN amount ELSE 0 END) as spend
		FROM wallet WHERE user_id = ?
		AND record_time >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
		GROUP BY DATE(record_time)`
	walletRows, err := DB.Query(walletQuery, userID, days)
	if err != nil {
		return nil, err
	}
	defer walletRows.Close()

	for walletRows.Next() {
		s := &DailyStats{}
		var dateStr string
		err := walletRows.Scan(&dateStr, &s.Earn, &s.Spend)
		if err != nil {
			return nil, err
		}
		// 解析日期格式，可能是 time.Time 格式
		var date time.Time
		if len(dateStr) > 10 {
			date, err = time.Parse("2006-01-02T15:04:05Z07:00", dateStr)
			if err != nil {
				date, err = time.Parse("2006-01-02", dateStr[:10])
				if err != nil {
					return nil, err
				}
			}
			dateStr = date.Format("2006-01-02")
		}
		s.Date = dateStr
		logger.Info("wallet_repo.go", 185, "Found wallet: date=%s, earn=%d, spend=%d", s.Date, s.Earn, s.Spend)
		if walletMap[s.Date] != nil {
			walletMap[s.Date].Earn += s.Earn
			walletMap[s.Date].Spend += s.Spend
		} else {
			walletMap[s.Date] = s
		}
	}

	// 构建日期列表，按日期正序
	dateList := []string{}
	for i := days - 1; i >= 0; i-- {
		d := now.AddDate(0, 0, -i).Format("2006-01-02")
		dateList = append(dateList, d)
	}

	// 计算每日累计余额（从最早日期开始累加）
	balanceMap := map[string]int{}
	cumulative := 0
	for _, d := range dateList {
		s := walletMap[d]
		if s != nil {
			cumulative += s.Earn - s.Spend
		}
		balanceMap[d] = cumulative
	}

	// 构建最终结果（按日期正序）
	for _, d := range dateList {
		s := walletMap[d]
		if s == nil {
			s = &DailyStats{Date: d, Earn: 0, Spend: 0}
		}
		s.Balance = balanceMap[d]
		stats = append(stats, s)
	}

	return stats, nil
}