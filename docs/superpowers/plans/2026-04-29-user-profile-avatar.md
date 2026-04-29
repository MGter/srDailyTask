# 用户资料和头像上传 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 新增用户资料卡片、基础资料修改、密码修改和本地头像上传功能。

**Architecture:** 后端在现有用户模型和接口上扩展 `avatar_url`、资料更新接口和头像上传接口；头像文件保存到项目根目录 `uploads/avatars`，数据库只保存访问 URL。前端在 `Tasks.vue` 左侧栏新增资料卡片和编辑弹窗，复用现有用户加载逻辑并同步更新 `localStorage`。

**Tech Stack:** Go net/http, database/sql, MySQL, Vue 3 Composition API, Axios, Vite.

---

## File Structure

- Create: `migrations/003_add_user_avatar.sql` — 给 `users` 表新增 `avatar_url` 字段。
- Modify: `internal/model/user.go` — 新增用户头像字段和资料更新请求结构。
- Modify: `internal/repository/user_repo.go` — 查询/创建用户时处理 `avatar_url`，新增资料更新、头像更新、用户名排重方法。
- Modify: `internal/service/user_service.go` — 实现资料更新、旧密码校验、邮箱格式校验、头像 URL 更新。
- Modify: `internal/handler/user_handler.go` — 新增 `Update` 和 `UploadAvatar` handler。
- Modify: `internal/handler/router.go` — 用户路由支持 `PUT /api/user/{id}`、`POST /api/user/{id}/avatar`，并开放 `/uploads/` 静态目录。
- Modify: `web/src/api/index.js` — 新增 `updateUser`、`uploadAvatar`。
- Modify: `web/src/views/Tasks.vue` — 新增资料卡片、编辑资料弹窗、头像预览、前端校验和保存逻辑。
- Modify generated: `web/dist/*` — 前端构建产物。

---

### Task 1: Add database migration and user model fields

**Files:**
- Create: `migrations/003_add_user_avatar.sql`
- Modify: `internal/model/user.go`

- [ ] **Step 1: Create migration**

Create `/home/mgter/srDailyTask/migrations/003_add_user_avatar.sql`:

```sql
-- migrations/003_add_user_avatar.sql
-- Add avatar URL for user profile images

ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255);
```

- [ ] **Step 2: Update user model**

Replace `/home/mgter/srDailyTask/internal/model/user.go` with:

```go
package model

import "time"

type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	Points    int       `json:"points"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
```

- [ ] **Step 3: Build check**

Run:

```bash
cd /home/mgter/srDailyTask && go test ./...
```

Expected: fails because repository scans still do not include `avatar_url`. This confirms the model change requires repository updates.

---

### Task 2: Update user repository

**Files:**
- Modify: `internal/repository/user_repo.go`

- [ ] **Step 1: Replace repository implementation**

Replace `/home/mgter/srDailyTask/internal/repository/user_repo.go` with:

```go
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
```

- [ ] **Step 2: Build check**

Run:

```bash
cd /home/mgter/srDailyTask && go test ./...
```

Expected: succeeds if database connection is not required by tests; otherwise repository compiles and any failure should be unrelated to syntax.

---

### Task 3: Add profile service logic

**Files:**
- Modify: `internal/service/user_service.go`

- [ ] **Step 1: Replace service implementation**

Replace `/home/mgter/srDailyTask/internal/service/user_service.go` with:

```go
package service

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"daily_task/internal/model"
	"daily_task/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: repository.NewUserRepository(),
	}
}

func (s *UserService) Register(req *model.CreateUserRequest) (*model.User, error) {
	existing, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	user := &model.User{
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		AvatarURL: "",
		Points:    0,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(req *model.LoginRequest) (*model.User, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if user.Password != req.Password {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

func (s *UserService) GetByID(id uint64) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) UpdateProfile(id uint64, req *model.UpdateUserRequest) (*model.User, error) {
	username := strings.TrimSpace(req.Username)
	email := strings.TrimSpace(req.Email)
	if username == "" {
		return nil, errors.New("username is required")
	}
	if email != "" && !isValidEmail(email) {
		return nil, errors.New("invalid email")
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	existing, err := s.repo.FindByUsernameExceptID(username, id)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	password := user.Password
	updatePassword := req.NewPassword != ""
	if updatePassword {
		if req.OldPassword == "" {
			return nil, errors.New("old password is required")
		}
		if req.OldPassword != user.Password {
			return nil, errors.New("old password is incorrect")
		}
		password = req.NewPassword
	}

	if err := s.repo.UpdateProfile(id, username, email, password, updatePassword); err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

func (s *UserService) UpdateAvatar(id uint64, avatarURL string) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if err := s.repo.UpdateAvatar(id, avatarURL); err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

func (s *UserService) UpdatePoints(id uint64, points int) error {
	return s.repo.UpdatePoints(id, points)
}

func (s *UserService) List(limit, offset int) ([]*model.User, error) {
	return s.repo.List(limit, offset)
}

func isValidEmail(email string) bool {
	pattern := `^[^@\s]+@[^@\s]+\.[^@\s]+$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}
```

- [ ] **Step 2: Build check**

Run:

```bash
cd /home/mgter/srDailyTask && go test ./...
```

Expected: succeeds if handlers are not yet requiring missing methods; otherwise compile error points to handler/router work in the next task.

---

### Task 4: Add user handlers and routes

**Files:**
- Modify: `internal/handler/user_handler.go`
- Modify: `internal/handler/router.go`

- [ ] **Step 1: Replace user handler**

Replace `/home/mgter/srDailyTask/internal/handler/user_handler.go` with:

```go
package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"daily_task/internal/logger"
	"daily_task/internal/model"
	"daily_task/internal/service"
)

const maxAvatarSize = 2 << 20

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		svc: service.NewUserService(),
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.Register(&req)
	if err != nil {
		logger.Error("user_handler.go", 30, "Register failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info("user_handler.go", 35, "User registered: %s", user.Username)
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.Login(&req)
	if err != nil {
		logger.Error("user_handler.go", 48, "Login failed: %v", err)
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	logger.Info("user_handler.go", 53, "User logged in: %s", user.Username)
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req model.UpdateUserRequest
	if err := readBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.UpdateProfile(id, &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	id, err := getAvatarUserIDFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAvatarSize+1024)
	if err := r.ParseMultipartForm(maxAvatarSize); err != nil {
		writeError(w, http.StatusBadRequest, "avatar is too large")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		writeError(w, http.StatusBadRequest, "avatar file is required")
		return
	}
	defer file.Close()

	if header.Size > maxAvatarSize {
		writeError(w, http.StatusBadRequest, "avatar is too large")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isAllowedAvatarExt(ext) {
		writeError(w, http.StatusBadRequest, "unsupported avatar type")
		return
	}

	if err := os.MkdirAll("uploads/avatars", 0755); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	filename := buildAvatarFilename(id, ext)
	path := filepath.Join("uploads", "avatars", filename)
	dst, err := os.Create(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer dst.Close()

	if _, err := file.Seek(0, 0); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if _, err := dst.ReadFrom(file); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	avatarURL := "/uploads/avatars/" + filename
	user, err := h.svc.UpdateAvatar(id, avatarURL)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := getQueryParam(r, "limit", 10)
	offset := getQueryParam(r, "offset", 0)

	users, err := h.svc.List(limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func getAvatarUserIDFromPath(r *http.Request) (uint64, error) {
	path := strings.TrimSuffix(r.URL.Path, "/avatar")
	return getURLParamFromPath(path, "/api/user/")
}

func isAllowedAvatarExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return true
	default:
		return false
	}
}

func buildAvatarFilename(userID uint64, ext string) string {
	buf := make([]byte, 4)
	_, _ = rand.Read(buf)
	return strings.Join([]string{
		uintToString(userID),
		uintToString(uint64(time.Now().UnixNano())),
		hex.EncodeToString(buf),
	}, "_") + ext
}
```

- [ ] **Step 2: Add uint helper if missing**

Check `/home/mgter/srDailyTask/internal/handler/handler.go`. If no `uintToString` helper exists, add this near the other helpers:

```go
func uintToString(n uint64) string {
	return strconv.FormatUint(n, 10)
}
```

Also add `strconv` to the import list in that file if not already present.

- [ ] **Step 3: Update routes**

In `/home/mgter/srDailyTask/internal/handler/router.go`, replace the user routes block:

```go
// User routes
mux.HandleFunc("/api/user/register", methodHandler("POST", userHandler.Register))
mux.HandleFunc("/api/user/login", methodHandler("POST", userHandler.Login))
mux.HandleFunc("/api/user/", userHandler.GetByID)
mux.HandleFunc("/api/users", methodHandler("GET", userHandler.List))
```

with:

```go
// User routes
mux.HandleFunc("/api/user/register", methodHandler("POST", userHandler.Register))
mux.HandleFunc("/api/user/login", methodHandler("POST", userHandler.Login))
mux.HandleFunc("/api/user/", func(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/avatar") {
		if r.Method == "POST" {
			userHandler.UploadAvatar(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
		return
	}
	switch r.Method {
	case "GET":
		userHandler.GetByID(w, r)
	case "PUT":
		userHandler.Update(w, r)
	default:
		http.Error(w, "Method not allowed", 405)
	}
})
mux.HandleFunc("/api/users", methodHandler("GET", userHandler.List))
```

- [ ] **Step 4: Serve uploads**

In `/home/mgter/srDailyTask/internal/handler/router.go`, inside `SetupServer` after:

```go
fs := http.FileServer(http.Dir(distDir))
mux.Handle("/assets/", fs)
```

add:

```go
uploadsFS := http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads")))
mux.Handle("/uploads/", uploadsFS)
```

- [ ] **Step 5: Build check**

Run:

```bash
cd /home/mgter/srDailyTask && gofmt -w internal/model/user.go internal/repository/user_repo.go internal/service/user_service.go internal/handler/user_handler.go internal/handler/router.go internal/handler/handler.go && go test ./...
```

Expected: Go code compiles and tests pass.

---

### Task 5: Apply database migration locally

**Files:**
- Runtime database only

- [ ] **Step 1: Inspect DB config**

Run:

```bash
cd /home/mgter/srDailyTask && grep -n "database:" -A10 config/config.yaml
```

Expected: shows MySQL host, port, database, username, password.

- [ ] **Step 2: Apply migration**

Run the equivalent MySQL command using the config values:

```bash
cd /home/mgter/srDailyTask && mysql -h <host> -P <port> -u <username> -p<password> <database> < migrations/003_add_user_avatar.sql
```

Expected: command exits with status 0. If the column already exists, stop and inspect before changing the migration.

- [ ] **Step 3: Verify column exists**

Run:

```bash
mysql -h <host> -P <port> -u <username> -p<password> <database> -e "SHOW COLUMNS FROM users LIKE 'avatar_url';"
```

Expected: output includes `avatar_url`.

---

### Task 6: Add frontend API methods

**Files:**
- Modify: `web/src/api/index.js`

- [ ] **Step 1: Replace userApi block**

In `/home/mgter/srDailyTask/web/src/api/index.js`, replace:

```javascript
export const userApi = {
  register: (data) => api.post('/user/register', data),
  login: (data) => api.post('/user/login', data),
  getUser: (id) => api.get(`/user/${id}`),
  getUsers: (params) => api.get('/users', { params })
}
```

with:

```javascript
export const userApi = {
  register: (data) => api.post('/user/register', data),
  login: (data) => api.post('/user/login', data),
  getUser: (id) => api.get(`/user/${id}`),
  updateUser: (id, data) => api.put(`/user/${id}`, data),
  uploadAvatar: (id, data) => api.post(`/user/${id}/avatar`, data, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }),
  getUsers: (params) => api.get('/users', { params })
}
```

- [ ] **Step 2: Frontend build check**

Run:

```bash
cd /home/mgter/srDailyTask/web && npm run build
```

Expected: build succeeds.

---

### Task 7: Add profile card and edit modal UI

**Files:**
- Modify: `web/src/views/Tasks.vue`

- [ ] **Step 1: Insert profile card**

In `/home/mgter/srDailyTask/web/src/views/Tasks.vue`, inside:

```vue
<div class="left-column">
  <!-- 打卡任务区域 -->
```

insert before the task section:

```vue
        <div class="section profile-section">
          <div class="profile-card">
            <div class="profile-avatar">
              <img v-if="user?.avatar_url" :src="user.avatar_url" alt="头像" />
              <span v-else>{{ userInitial }}</span>
            </div>
            <div class="profile-details">
              <h3>{{ user?.username }}</h3>
              <p>{{ user?.email || '未设置邮箱' }}</p>
              <span class="profile-points">总积分 {{ user?.points || 0 }}</span>
            </div>
            <button class="edit-profile-btn" @click="openProfileModal">编辑资料</button>
          </div>
        </div>
```

- [ ] **Step 2: Insert profile modal**

In the same file, before the existing edit task modal block, insert:

```vue
      <div v-if="showProfileModal" class="modal-overlay" @click="closeProfileModal">
        <div class="modal profile-modal" @click.stop>
          <h3>编辑资料</h3>
          <form @submit.prevent="saveProfile">
            <div class="avatar-edit">
              <div class="profile-avatar large">
                <img v-if="avatarPreview || profileForm.avatar_url" :src="avatarPreview || profileForm.avatar_url" alt="头像预览" />
                <span v-else>{{ profileInitial }}</span>
              </div>
              <label class="avatar-upload-btn">
                选择头像
                <input type="file" accept="image/jpeg,image/png,image/webp,image/gif" @change="handleAvatarChange" />
              </label>
            </div>
            <div class="form-group">
              <label>用户名</label>
              <input v-model.trim="profileForm.username" type="text" required />
            </div>
            <div class="form-group">
              <label>邮箱</label>
              <input v-model.trim="profileForm.email" type="email" placeholder="可选" />
            </div>
            <div class="form-group">
              <label>旧密码</label>
              <input v-model="profileForm.old_password" type="password" placeholder="修改密码时必填" />
            </div>
            <div class="form-group">
              <label>新密码</label>
              <input v-model="profileForm.new_password" type="password" placeholder="不修改可留空" />
            </div>
            <div class="form-group">
              <label>确认新密码</label>
              <input v-model="profileForm.confirm_password" type="password" placeholder="再次输入新密码" />
            </div>
            <p class="error" v-if="profileError">{{ profileError }}</p>
            <p class="success" v-if="profileSuccess">{{ profileSuccess }}</p>
            <div class="modal-actions">
              <button type="button" class="cancel-btn" @click="closeProfileModal">取消</button>
              <button type="submit" class="save-btn" :disabled="savingProfile">{{ savingProfile ? '保存中...' : '保存' }}</button>
            </div>
          </form>
        </div>
      </div>
```

- [ ] **Step 3: Build check**

Run:

```bash
cd /home/mgter/srDailyTask/web && npm run build
```

Expected: fails because script variables/functions are not yet defined. This confirms the template is wired to new state.

---

### Task 8: Add profile modal script logic

**Files:**
- Modify: `web/src/views/Tasks.vue`

- [ ] **Step 1: Add state after `addForm`**

After the existing `addForm` ref, add:

```javascript
const showProfileModal = ref(false)
const savingProfile = ref(false)
const profileError = ref('')
const profileSuccess = ref('')
const avatarFile = ref(null)
const avatarPreview = ref('')
const profileForm = ref({
  username: '',
  email: '',
  avatar_url: '',
  old_password: '',
  new_password: '',
  confirm_password: ''
})
```

- [ ] **Step 2: Add computed values after `todayStats`**

Add:

```javascript
const userInitial = computed(() => user.value?.username?.slice(0, 1).toUpperCase() || '?')
const profileInitial = computed(() => profileForm.value.username?.slice(0, 1).toUpperCase() || userInitial.value)
```

- [ ] **Step 3: Add profile functions before `loadUser`**

Add:

```javascript
const openProfileModal = () => {
  profileForm.value = {
    username: user.value?.username || '',
    email: user.value?.email || '',
    avatar_url: user.value?.avatar_url || '',
    old_password: '',
    new_password: '',
    confirm_password: ''
  }
  avatarFile.value = null
  avatarPreview.value = ''
  profileError.value = ''
  profileSuccess.value = ''
  showProfileModal.value = true
}

const closeProfileModal = () => {
  showProfileModal.value = false
  avatarFile.value = null
  avatarPreview.value = ''
}

const handleAvatarChange = (event) => {
  const file = event.target.files?.[0]
  if (!file) return

  const allowedTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    profileError.value = '头像只支持 jpg、png、webp、gif'
    return
  }
  if (file.size > 2 * 1024 * 1024) {
    profileError.value = '头像不能超过 2MB'
    return
  }

  avatarFile.value = file
  avatarPreview.value = URL.createObjectURL(file)
  profileError.value = ''
}

const saveProfile = async () => {
  profileError.value = ''
  profileSuccess.value = ''
  if (profileForm.value.new_password !== profileForm.value.confirm_password) {
    profileError.value = '两次输入的新密码不一致'
    return
  }

  savingProfile.value = true
  try {
    const payload = {
      username: profileForm.value.username,
      email: profileForm.value.email,
      old_password: profileForm.value.old_password,
      new_password: profileForm.value.new_password
    }
    let res = await userApi.updateUser(userId, payload)

    if (avatarFile.value) {
      const formData = new FormData()
      formData.append('avatar', avatarFile.value)
      res = await userApi.uploadAvatar(userId, formData)
    }

    user.value = res.data
    localStorage.setItem('user', JSON.stringify(res.data))
    profileSuccess.value = '资料已更新'
    setTimeout(() => {
      closeProfileModal()
    }, 600)
  } catch (e) {
    profileError.value = e.response?.data?.error || '保存失败'
  } finally {
    savingProfile.value = false
  }
}
```

- [ ] **Step 4: Update loadUser localStorage sync**

In `loadUser`, after:

```javascript
user.value = res.data
```

add:

```javascript
localStorage.setItem('user', JSON.stringify(res.data))
```

- [ ] **Step 5: Build check**

Run:

```bash
cd /home/mgter/srDailyTask/web && npm run build
```

Expected: build succeeds or only CSS remains to refine.

---

### Task 9: Add profile styles

**Files:**
- Modify: `web/src/views/Tasks.vue`

- [ ] **Step 1: Add profile CSS before task styles**

Before the existing task list styles, add:

```css
.profile-section {
  padding: 16px;
}

.profile-card {
  display: flex;
  align-items: center;
  gap: 12px;
}

.profile-avatar {
  width: 54px;
  height: 54px;
  flex: 0 0 auto;
  border-radius: 50%;
  overflow: hidden;
  background: #e8f2ff;
  color: #007aff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  font-weight: 700;
}

.profile-avatar.large {
  width: 86px;
  height: 86px;
  font-size: 32px;
}

.profile-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.profile-details {
  flex: 1;
  min-width: 0;
}

.profile-details h3 {
  margin: 0 0 4px;
  color: #1d1d1f;
  font-size: 16px;
}

.profile-details p {
  margin: 0 0 5px;
  color: #86868b;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.profile-points {
  color: #007aff;
  font-size: 12px;
  font-weight: 600;
}

.edit-profile-btn {
  padding: 7px 10px;
  border: none;
  border-radius: 8px;
  background: #007aff;
  color: #fff;
  font-size: 12px;
  cursor: pointer;
}

.profile-modal {
  max-width: 420px;
}

.avatar-edit {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 18px;
}

.avatar-upload-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 8px 12px;
  border-radius: 8px;
  background: #f5f5f7;
  color: #007aff;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}

.avatar-upload-btn input {
  display: none;
}

.error {
  color: #ff3b30;
  font-size: 13px;
}

.success {
  color: #34c759;
  font-size: 13px;
}
```

- [ ] **Step 2: Add mobile profile CSS**

Inside the existing `@media (max-width: 600px)` block, add:

```css
.profile-card {
  gap: 10px;
}

.profile-avatar {
  width: 46px;
  height: 46px;
  font-size: 18px;
}

.edit-profile-btn {
  padding: 6px 8px;
  font-size: 11px;
}

.avatar-edit {
  gap: 12px;
}
```

- [ ] **Step 3: Build check**

Run:

```bash
cd /home/mgter/srDailyTask/web && npm run build
```

Expected: build succeeds and `web/dist/*` updates.

---

### Task 10: Final verification, restart, commit, push

**Files:**
- Create: `migrations/003_add_user_avatar.sql`
- Modify: `internal/model/user.go`
- Modify: `internal/repository/user_repo.go`
- Modify: `internal/service/user_service.go`
- Modify: `internal/handler/user_handler.go`
- Modify: `internal/handler/router.go`
- Modify: `internal/handler/handler.go` if helper was needed
- Modify: `web/src/api/index.js`
- Modify: `web/src/views/Tasks.vue`
- Modify generated: `web/dist/*`
- Create: `docs/superpowers/plans/2026-04-29-user-profile-avatar.md`

- [ ] **Step 1: Run backend and frontend checks**

Run:

```bash
cd /home/mgter/srDailyTask && go test ./... && cd web && npm run build
```

Expected: both commands succeed.

- [ ] **Step 2: Restart service and health check**

Run:

```bash
systemctl --user restart srdailytask && curl -s http://localhost:18888/health
```

Expected:

```json
{"status":"ok"}
```

- [ ] **Step 3: Manual browser check**

Open `http://localhost:18888` and verify:
- 左侧栏顶部显示资料卡片。
- 没头像时显示用户名首字符。
- 修改用户名和邮箱后刷新仍生效。
- 不填新密码时只更新资料。
- 旧密码错误时无法改密码。
- 新密码和确认密码不一致时前端提示。
- 上传小于 2MB 的 jpg/png/webp/gif 后头像显示。
- 上传非图片或超过 2MB 文件会提示错误。

- [ ] **Step 4: Check git status**

Run:

```bash
cd /home/mgter/srDailyTask && git status --short
```

Expected: includes only files from this plan plus any pre-existing unrelated files. Do not stage `docs/err.png` or `docs/web_page.png` unless the user explicitly asks.

- [ ] **Step 5: Commit relevant files**

Run:

```bash
cd /home/mgter/srDailyTask && git add migrations/003_add_user_avatar.sql internal/model/user.go internal/repository/user_repo.go internal/service/user_service.go internal/handler/user_handler.go internal/handler/router.go internal/handler/handler.go web/src/api/index.js web/src/views/Tasks.vue web/dist docs/superpowers/plans/2026-04-29-user-profile-avatar.md && git commit -m "$(cat <<'EOF'
新增用户资料编辑和头像上传

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

Expected: commit succeeds. If `internal/handler/handler.go` was not modified, omit it from `git add`.

- [ ] **Step 6: Push**

Run:

```bash
cd /home/mgter/srDailyTask && git push
```

Expected: push succeeds.

---

## Self-Review

- Spec coverage: covers profile card, editable username/email/password, old password requirement, avatar upload to local directory, 2MB/type limits, static serving, localStorage sync, migration, verification.
- Placeholder scan: no TBD/TODO/placeholders.
- Type consistency: uses `avatar_url`, `UpdateUserRequest`, `updateUser`, `uploadAvatar`, and `/api/user/{id}/avatar` consistently across backend and frontend.
