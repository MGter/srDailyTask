package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"daily_task/internal/config"
	"daily_task/internal/logger"
)

var DB *sql.DB

func InitDB(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error("mysql.go", 20, "Failed to open database: %v", err)
		return err
	}

	if err = db.Ping(); err != nil {
		logger.Error("mysql.go", 25, "Failed to ping database: %v", err)
		return err
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	DB = db
	logger.Info("mysql.go", 32, "Connected to database: %s", cfg.Name)
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		logger.Info("mysql.go", 38, "Database connection closed")
	}
}