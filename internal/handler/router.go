package handler

import (
	"net/http"

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
	mux.HandleFunc("/api/user/", userHandler.GetByID) // GET with id in path
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
	mux.HandleFunc("/api/wallet/", func(w http.ResponseWriter, r *http.Request) {
		// Check if path ends with /balance
		if r.URL.Path == "/api/wallet/"+r.PathValue("user_id")+"/balance" {
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
	mux.HandleFunc("/api/wallet/spend", methodHandler("POST", walletHandler.Spend))

	// Point routes
	mux.HandleFunc("/api/points/", methodHandler("GET", pointHandler.GetPointHistory))

	logger.Info("router.go", 46, "Router initialized")
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