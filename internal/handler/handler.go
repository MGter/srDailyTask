package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func readBody(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func getQueryParam(r *http.Request, key string, defaultValue int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return result
}

func getURLParamFromPath(path string, prefix string) (uint64, error) {
	// 从路径中提取参数，例如 /api/task/123 -> 提取 123
	trimmed := strings.TrimPrefix(path, prefix)
	trimmed = strings.TrimSuffix(trimmed, "/")
	trimmed = strings.TrimSuffix(trimmed, "/balance")

	if trimmed == "" {
		return 0, errors.New("param not found")
	}

	// 如果路径有更多部分，只取第一个
	if strings.Contains(trimmed, "/") {
		trimmed = strings.Split(trimmed, "/")[0]
	}

	result, err := strconv.ParseUint(trimmed, 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func uintToString(n uint64) string {
	return strconv.FormatUint(n, 10)
}

func getIDFromPath(r *http.Request) (uint64, error) {
	return getURLParamFromPath(r.URL.Path, "/api/user/")
}

func getTaskIDFromPath(r *http.Request) (uint64, error) {
	return getURLParamFromPath(r.URL.Path, "/api/task/")
}

func getCheckinTaskIDFromPath(r *http.Request) (uint64, error) {
	return getURLParamFromPath(r.URL.Path, "/api/checkin/")
}

func getUserIDFromPath(r *http.Request, prefix string) (uint64, error) {
	return getURLParamFromPath(r.URL.Path, prefix)
}
