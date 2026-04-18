package handler

import (
	"net/http"
	"os"
	"strings"

	"daily_task/internal/logger"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	userHandler := NewUserHandler()
	taskHandler := NewTaskHandler()
	walletHandler := NewWalletHandler()
	pointHandler := NewPointHandler()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// User routes
	mux.HandleFunc("/api/user/register", methodHandler("POST", userHandler.Register))
	mux.HandleFunc("/api/user/login", methodHandler("POST", userHandler.Login))
	mux.HandleFunc("/api/user/", userHandler.GetByID)
	mux.HandleFunc("/api/users", methodHandler("GET", userHandler.List))

	// Task routes
	mux.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			taskHandler.Create(w, r)
		case "GET":
			taskHandler.GetByID(w, r)
		default:
			http.Error(w, "Method not allowed", 405)
		}
	})
	mux.HandleFunc("/api/task/user/", methodHandler("GET", taskHandler.GetByUserID))
	mux.HandleFunc("/api/task/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			taskHandler.GetByID(w, r)
		case "PUT":
			taskHandler.Update(w, r)
		case "DELETE":
			taskHandler.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", 405)
		}
	})

	// Check-in routes
	mux.HandleFunc("/api/checkin/", methodHandler("POST", taskHandler.CheckIn))
	mux.HandleFunc("/api/checkin/user/", methodHandler("GET", pointHandler.GetCheckIns))

	// Wallet routes
	mux.HandleFunc("/api/wallet/spend", methodHandler("POST", walletHandler.Spend))
	mux.HandleFunc("/api/wallet/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/balance") {
			if r.Method == "GET" {
				walletHandler.GetBalance(w, r)
			} else {
				http.Error(w, "Method not allowed", 405)
			}
			return
		}
		if r.Method == "GET" {
			walletHandler.GetByUserID(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	})

	// Point routes
	mux.HandleFunc("/api/points/", methodHandler("GET", pointHandler.GetPointHistory))

	logger.Info("router.go", 56, "Router initialized")
	return mux
}

// SetupServer 设置完整的 HTTP 服务器，包括 API 和静态文件
func SetupServer(distDir string) http.Handler {
	mux := NewRouter()

	// 静态文件服务
	if distDir != "" && dirExists(distDir) {
		fs := http.FileServer(http.Dir(distDir))
		mux.Handle("/assets/", fs)

		// 处理 SPA 路由 - 对于非 API 请求返回 index.html
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// API 路由已经在上面的路由中处理
			// 其他路由返回 index.html (SPA fallback)
			path := distDir + r.URL.Path
			if fileExists(path) {
				http.ServeFile(w, r, path)
			} else {
				http.ServeFile(w, r, distDir+"/index.html")
			}
		})
		logger.Info("router.go", 80, "Static files served from: %s", distDir)
	} else {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, http.StatusOK, map[string]string{"message": "daily_task API server"})
		})
		logger.Info("router.go", 85, "No static files, running as API only")
	}

	return mux
}

func methodHandler(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", 405)
			return
		}
		handler(w, r)
	}
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}