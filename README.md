# srDailyTask - 日常打卡积分系统

## 项目概述

srDailyTask 是一个日常任务打卡系统，用户可以创建周期性任务，完成打卡获得积分奖励，积分可用于兑换或记录日常支出。

## 功能特性

- **用户管理**: 注册、登录、用户信息展示
- **任务管理**: 创建、查看、删除任务
- **任务级别**: 低/中/高三级，用颜色区分（绿/橙/红）
- **周期模式**: 单次、每周一、工作日、周末、每天
- **打卡系统**: 完成任务打卡获得积分，支持取消打卡、跳过任务（不计积分）
- **分页显示**: 任务列表和钱包记录每页最多5个，手机端按钮文字竖排
- **积分钱包**: 查看积分余额、手动添加收入/支出记录
- **积分统计**: 柱状图展示每日积累、消耗、结余
- **定时提醒**: 每分钟检查任务状态发送提醒
- **Web 界面**: Vue 3 前端，苹果风格设计，全屏布局

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + Vite + Chart.js + vue-chartjs |
| 后端 | Go 1.26 + 标准库 net/http |
| 数据库 | MySQL |
| 定时任务 | robfig/cron |

## 项目结构

```
srDailyTask/
├── cmd/
│   └── main.go                 # 程序入口，初始化数据库、路由、调度器
│
├── config/
│   └── config.yaml             # 配置文件（端口、数据库、日志）
│
├── internal/
│   ├── config/
│   │   └── config.go           # 配置加载和解析
│   │
│   ├── handler/                # HTTP 处理器层
│   │   ├── router.go           # 路由配置，静态文件托管
│   │   ├── handler.go          # 通用工具函数（JSON响应、路径参数解析）
│   │   ├── user_handler.go     # 用户注册、登录、查询
│   │   ├── task_handler.go     # 任务 CRUD、打卡、取消打卡
│   │   ├── point_handler.go    # 积分历史、每日统计、今日已打卡任务
│   │   └── wallet_handler.go   # 钱包余额、添加记录、删除记录
│   │
│   ├── service/                # 业务逻辑层
│   │   ├── user_service.go     # 用户 CRUD、积分更新
│   │   ├── task_service.go     # 任务 CRUD、打卡逻辑、周期判断
│   │   ├── checkin_service.go  # 打卡记录创建、删除
│   │   └── wallet_service.go   # 钱包记录 CRUD、每日统计
│   │
│   ├── repository/             # 数据访问层
│   │   ├── mysql.go            # 数据库连接初始化
│   │   ├── user_repo.go        # 用户表操作
│   │   ├── task_repo.go        # 任务表操作
│   │   ├── checkin_repo.go     # 打卡表操作
│   │   └── wallet_repo.go      # 钱包表操作、每日统计查询
│   │
│   ├── model/                  # 数据模型定义
│   │   ├── user.go             # User 结构体
│   │   ├── task.go             # Task、CircleMode、TaskLevel 定义
│   │   ├── checkin.go          # CheckIn 结构体
│   │   └── wallet.go           # Wallet、WalletType 定义
│   │
│   ├── scheduler/
│   │   └── scheduler.go        # 定时提醒调度器
│   │
│   └── logger/
│       └── logger.go           # 日志工具
│
├── web/                        # Vue 3 前端
│   ├── src/
│   │   ├── views/
│   │   │   ├── Login.vue       # 登录页面
│   │   │   ├── Register.vue    # 注册页面
│   │   │   └── Tasks.vue       # 主页面（打卡+钱包+统计）
│   │   ├── api/
│   │   │   └── index.js        # API 调用封装
│   │   ├── router/
│   │   │   └── index.js        # 路由配置
│   │   ├── App.vue             # 根组件
│   │   ├── main.js             # 入口文件
│   │   └── style.css           # 全局样式
│   ├── dist/                   # 打包后的静态文件
│   ├── package.json
│   └ vite.config.js
│   └── index.html              # HTML 模板（标题: srDailyTask）
│
├── migrations/
│   ├── 001_init.sql            # 初始化数据库表
│   └── 002_add_checkin_id.sql  # 添加 checkin_id 字段
│
├── pics/
│   ├── kita.png                # 原始背景图（4K PNG）
│   └── kita.webp               # 优化后的背景图（1080P WebP, 214KB）
│
├── go.mod                      # Go 模块定义
├── go.sum                      # Go 依赖锁定
└── Makefile                    # 构建脚本
```

## 数据库表结构

### users 表
```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    points INT DEFAULT 0,           -- 当前积分余额
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### tasks 表
```sql
CREATE TABLE tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    circle_mode VARCHAR(20) NOT NULL DEFAULT 'once',  -- 周期模式
    level INT DEFAULT 1,                               -- 任务级别 1/2/3
    points INT DEFAULT 10,                             -- 打卡奖励积分
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_expired BOOLEAN DEFAULT FALSE,                  -- 单次任务打卡后过期
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### checkins 表
```sql
CREATE TABLE checkins (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    task_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    points INT NOT NULL,                 -- 本次获得的积分
    check_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### wallet 表
```sql
CREATE TABLE wallet (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    checkin_id BIGINT UNSIGNED DEFAULT 0,  -- 关联的打卡记录ID
    balance INT DEFAULT 0,
    type VARCHAR(20) NOT NULL,             -- 'earn' 或 'spend'
    amount INT NOT NULL,
    description VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    record_time DATETIME,                  -- 记录时间（可自定义）
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

## API 接口

### 用户相关

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/user/register` | 注册用户 |
| POST | `/api/user/login` | 用户登录 |
| GET | `/api/user/{id}` | 获取用户信息 |
| GET | `/api/users` | 获取所有用户列表 |

### 任务相关

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/task` | 创建任务 |
| GET | `/api/task/user/{user_id}` | 获取用户所有任务（含 should_checkin_today） |
| GET | `/api/task/today/{user_id}` | 获取今日需打卡任务 |
| DELETE | `/api/task/{id}` | 删除任务 |

### 打卡相关

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/checkin/{task_id}` | 打卡（获得积分） |
| DELETE | `/api/checkin/{task_id}` | 取消打卡（退还积分） |
| POST | `/api/checkin/{task_id}/skip` | 跳过任务（不计积分） |
| GET | `/api/checkin/user/{user_id}` | 获取用户打卡历史 |
| GET | `/api/checkin/today/{user_id}` | 获取今日已打卡任务ID列表 |

### 积分钱包相关

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/wallet/{user_id}` | 获取钱包记录列表 |
| GET | `/api/wallet/{user_id}/balance` | 获取积分余额 |
| POST | `/api/wallet/spend` | 消费积分 |
| POST | `/api/wallet/add` | 添加记录（收入/支出） |
| DELETE | `/api/wallet/delete/{id}` | 删除记录 |

### 统计相关

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/points/{user_id}` | 获取积分历史记录 |
| GET | `/api/points/daily/{user_id}?days=N` | 获取每日统计（默认7天，支持14/30/180/360） |

## 前后端逻辑详解

### 1. 打卡流程

**前端** (`Tasks.vue`):
```javascript
// 用户点击"打卡"按钮
const checkin = async (taskId) => {
  await checkinApi.checkin(taskId, { user_id: userId })
  // 更新任务状态
  task.checked = true
  // 刷新用户积分、历史记录、统计图表
  await loadUser()
  await loadHistory()
  await loadDailyStats()
}
```

**后端** (`task_service.go`):
```go
func (s *TaskService) CheckIn(taskID, userID uint64) (*model.CheckIn, error) {
  // 1. 验证任务存在且属于用户
  // 2. 检查今天是否需要打卡（ShouldCheckinToday）
  // 3. 检查今天是否已打卡
  // 4. 创建打卡记录（checkins表）
  // 5. 创建钱包记录（wallet表，描述为"打卡: 任务标题"）
  // 6. 更新用户积分
  // 7. 单次任务标记为过期
}
```

### 2. 取消打卡流程

**后端**:
```go
func (s *TaskService) CancelCheckIn(taskID, userID uint64) error {
  // 1. 找到今天的打卡记录
  // 2. 删除打卡记录（checkins表）
  // 3. 删除对应的钱包记录（通过 checkin_id 关联）
  // 4. 退还积分
}
```

### 3. 跳过任务流程

**前端** (`Tasks.vue`):
```javascript
const skipTask = async (taskId) => {
  await checkinApi.skip(taskId, { user_id: userId })
  task.checked = true  // 标记为已完成（灰色显示）
}
```

**后端**:
```go
func (s *TaskService) SkipCheckIn(taskID, userID uint64) (*model.CheckIn, error) {
  // 1. 创建打卡记录，points = 0（不计积分）
  // 2. 不创建钱包记录
  // 3. 单次任务标记为过期
}
```

### 4. 周期模式判断

**后端** (`task_service.go`):
```go
func ShouldCheckinToday(task *model.Task) bool {
  weekday := time.Now().Weekday()
  switch task.CircleMode {
  case "once":    return true                // 每天可打卡
  case "weekly":  return weekday == Monday   // 仅周一
  case "workday": return weekday <= Friday   // 周一至周五
  case "weekend": return weekday >= Saturday // 周六周日
  case "custom":  return true                // 每天
  }
}
```

### 5. 每日统计计算

**后端** (`wallet_repo.go`):
```go
func GetDailyStats(userID uint64, days int) ([]*DailyStats, error) {
  // 从 wallet 表按日期聚合：
  // - earn: 当日总收入（打卡奖励 + 手动添加的收入）
  // - spend: 当日总支出
  // - balance: 当日结余 = earn - spend
  // 不再从 checkins 表单独查询（已合并到 wallet）
}
```

### 6. 删除记录同步逻辑

删除钱包记录时，如果 `checkin_id > 0`（打卡产生的记录），会同步删除打卡历史：
```go
func (s *WalletService) Delete(id uint64, userID uint64) error {
  wallet := s.repo.FindByID(id)
  if wallet.CheckinID > 0 {
    s.checkinRepo.Delete(wallet.CheckinID)  // 同步删除打卡历史
  }
  s.repo.Delete(id, userID)
  // 反向更新用户积分
}
```

## Web 界面

### 页面结构

单页面布局（Tasks.vue），分为左侧和右侧：

**左侧**:
- 今日打卡任务列表（未打卡在前，已打卡灰色显示在后）
- 任务创建表单（标题、周期、级别、积分）
- 积分钱包（添加收入/支出表单 + 历史记录列表）
- 其他任务（今天不需要打卡的任务）

**右侧**:
- 积分统计图表（柱状图）
- 天数切换（7/14/30/180/360天）
- 三条数据：当日积累（绿）、当日消耗（红）、当日结余（橙）

### 样式特点

- 苹果风格设计：纯白卡片、浅灰输入框
- 背景图：kita.webp（214KB，1080P）
- 任务级别用颜色圆点区分（绿●/橙●/红●）
- 图例横排、字体缩小（11px）

### 交互功能

- 打卡后按钮变为橙色"取消打卡"
- 历史记录每项有绿色"重复"按钮（快速复制）
- 退出登录按钮在页面右上角

## 快速开始

### 1. 安装依赖

```bash
# Go 后端依赖
go mod download

# 前端依赖
cd web && npm install
```

### 2. 创建数据库

```bash
mysql -u root -e "CREATE DATABASE daily_task CHARACTER SET utf8mb4;"
mysql -u root daily_task < migrations/001_init.sql
mysql -u root daily_task < migrations/002_add_checkin_id.sql
```

### 3. 配置

编辑 `config/config.yaml`:
```yaml
server:
  port: 18888

database:
  host: localhost
  port: 3306
  user: root
  password: ""
  name: daily_task
```

### 4. 运行

```bash
# 构建前端
cd web && npm run build && cp ../pics/kita.webp dist/assets/kita.webp

# 构建后端
cd .. && go build -o daily_task cmd/main.go

# 运行
./daily_task

# 访问 http://localhost:18888
```

## 部署

单二进制文件部署，后端托管前端静态文件：

```bash
# 部署所需文件
daily_task              # 二进制文件
config/config.yaml      # 配置文件
web/dist/               # 前端静态文件
pics/kita.webp          # 背景图（需复制到 web/dist/assets/）
```

## 更新历史

- 2026-04-18: 项目从 C++ 迁移到 Go
- 2026-04-19: 
  - 添加任务级别和删除功能
  - 积分记录支持新增和删除
  - 实现周期模式判断逻辑
  - 合并打卡和积分页面
  - 添加统计图表（Chart.js 柱状图）
  - 添加取消打卡和重复记录功能
  - 背景图转换为 WebP（压缩95%）
  - 页面改为全屏布局
- 2026-04-20:
  - 添加跳过任务功能（不计积分）
  - 任务和钱包分页显示（每页5个）
  - 手机端按钮文字竖排显示