package handler

import (
	"net/http"

	"daily_task/internal/logger"
	"daily_task/internal/model"
	"daily_task/internal/service"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		svc: service.NewTaskService(),
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTaskRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.svc.Create(&req)
	if err != nil {
		logger.Error("task_handler.go", 24, "Create task failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info("task_handler.go", 28, "Task created: %d", task.ID)
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getTaskIDFromPath(r)
	if err != nil || id == 0 {
		// 如果没有 id，返回用户的所有任务
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	task, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if task == nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromPath(r, "/api/task/user/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	limit := getQueryParam(r, "limit", 10)
	offset := getQueryParam(r, "offset", 0)

	tasks, err := h.svc.GetByUserID(userID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 添加 should_checkin_today 字段
	result := make([]map[string]interface{}, len(tasks))
	for i, task := range tasks {
		result[i] = map[string]interface{}{
			"id":            task.ID,
			"user_id":       task.UserID,
			"title":         task.Title,
			"description":   task.Description,
			"circle_mode":   task.CircleMode,
			"level":         task.Level,
			"points":        task.Points,
			"created_at":    task.CreatedAt,
			"updated_at":    task.UpdatedAt,
			"is_expired":    task.IsExpired,
			"should_checkin_today": h.svc.ShouldCheckinToday(task),
		}
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *TaskHandler) GetTodayTasks(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromPath(r, "/api/task/today/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	tasks, err := h.svc.GetByUserID(userID, 50, 0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 只返回今天需要打卡且未过期的任务
	todayTasks := []map[string]interface{}{}
	for _, task := range tasks {
		if !task.IsExpired && h.svc.ShouldCheckinToday(task) {
			todayTasks = append(todayTasks, map[string]interface{}{
				"id":            task.ID,
				"user_id":       task.UserID,
				"title":         task.Title,
				"description":   task.Description,
				"circle_mode":   task.CircleMode,
				"level":         task.Level,
				"points":        task.Points,
				"created_at":    task.CreatedAt,
				"updated_at":    task.UpdatedAt,
				"is_expired":    task.IsExpired,
			})
		}
	}

	writeJSON(w, http.StatusOK, todayTasks)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getTaskIDFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req model.UpdateTaskRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.svc.Update(id, &req)
	if err != nil {
		logger.Error("task_handler.go", 76, "Update task failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info("task_handler.go", 80, "Task updated: %d", task.ID)
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getTaskIDFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		logger.Error("task_handler.go", 93, "Delete task failed: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("task_handler.go", 97, "Task deleted: %d", id)
	writeJSON(w, http.StatusOK, map[string]string{"message": "task deleted"})
}

func (h *TaskHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	taskID, err := getCheckinTaskIDFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req model.CheckInRequest
	if err := readBody(r, &req); err != nil {
		req.TaskID = taskID
		req.UserID = 1 // 默认用户，后续加认证
	}

	checkin, err := h.svc.CheckIn(taskID, req.UserID)
	if err != nil {
		logger.Error("task_handler.go", 116, "Check-in failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info("task_handler.go", 120, "Check-in success: task %d, user %d", taskID, req.UserID)
	writeJSON(w, http.StatusOK, checkin)
}

func (h *TaskHandler) CancelCheckIn(w http.ResponseWriter, r *http.Request) {
	taskID, err := getCheckinTaskIDFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req model.CheckInRequest
	if err := readBody(r, &req); err != nil {
		req.UserID = 1
	}

	if err := h.svc.CancelCheckIn(taskID, req.UserID); err != nil {
		logger.Error("task_handler.go", 135, "Cancel check-in failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info("task_handler.go", 139, "Cancel check-in success: task %d, user %d", taskID, req.UserID)
	writeJSON(w, http.StatusOK, map[string]string{"message": "checkin cancelled"})
}