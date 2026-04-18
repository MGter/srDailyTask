package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"daily_task/internal/config"
	"daily_task/internal/handler"
	"daily_task/internal/logger"
	"daily_task/internal/repository"
	"daily_task/internal/scheduler"
)

func main() {
	cfgPath := "./config/config.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	if err := config.Load(cfgPath); err != nil {
		fmt.Printf("Failed to load config: %v, using defaults\n", err)
		config.AppConfig = config.GetDefault()
	}

	logger.Init(config.AppConfig.Log.Level, config.AppConfig.Log.Path)
	logger.Info("main.go", 24, "Starting daily_task service")

	if err := repository.InitDB(&config.AppConfig.Database); err != nil {
		logger.Error("main.go", 27, "Failed to init database: %v", err)
		os.Exit(1)
	}

	// 使用 SetupServer 托管静态文件
	distDir := "./web/dist"
	handler := handler.SetupServer(distDir)

	if config.AppConfig.Scheduler.Enabled {
		sched := scheduler.NewScheduler()
		sched.Start()
		logger.Info("main.go", 38, "Scheduler enabled")
	}

	go func() {
		addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
		logger.Info("main.go", 42, "Server listening on %s", addr)
		if err := http.ListenAndServe(addr, handler); err != nil {
			logger.Error("main.go", 44, "Server error: %v", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan

	logger.Info("main.go", 51, "Received signal: %v, shutting down", sig)
	repository.CloseDB()
	logger.Close()
	logger.Info("main.go", 54, "daily_task stopped")
}