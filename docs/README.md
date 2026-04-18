# srDailyTask - 日常打卡积分系统

## 项目概述

srDailyTask 是一个日常任务打卡系统，用户可以创建周期性任务，完成打卡获得积分奖励，积分可用于兑换。

## 功能特性

- **用户管理**: 注册、登录、用户信息
- **任务管理**: 创建、查看、更新、删除任务
- **周期模式**: 单次、每周、工作日、周末、自定义
- **打卡系统**: 完成任务打卡，获得积分
- **积分钱包**: 查看积分余额、消费积分
- **定时提醒**: 自动检查任务状态，发送提醒

## 技术栈

- **语言**: Go 1.21
- **Web 框架**: 标准库 `net/http`
- **数据库**: MySQL
- **定时任务**: `robfig/cron`
- **配置**: YAML (`gopkg.in/yaml.v3`)

## 项目结构

```
srDailyTask/
├── cmd/
│   └── main.go              # 程序入口
├── config/
│   └── config.yaml          # 配置文件
├── internal/
│   ├── config/              # 配置加载
│   ├── handler/             # HTTP 处理器
│   ├── logger/              # 日志模块
│   ├── model/               # 数据模型
│   ├── repository/          # 数据库操作
│   ├── scheduler/           # 定时提醒
│   └── service/             # 业务逻辑
├── migrations/              # 数据库迁移
│   └ 001_init.sql
├── pkg/utils/               # 工具函数
├── go.mod
├── Makefile
└── docs/
```

## 快速开始

### 1. 安装依赖

```bash
make deps
```

### 2. 创建数据库

```bash
mysql -u root -p -e "CREATE DATABASE daily_task CHARACTER SET utf8mb4;"
mysql -u root -p daily_task < migrations/001_init.sql
```

### 3. 配置

编辑 `config/config.yaml`，修改数据库连接信息：

```yaml
database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  name: daily_task
```

### 4. 运行

```bash
make run
# 或编译后运行
make build
./daily_task
```

## API 接口

### 用户

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/user/register` | 注册用户 |
| POST | `/api/user/login` | 用户登录 |
| GET | `/api/user/{id}` | 获取用户信息 |
| GET | `/api/users` | 用户列表 |

### 任务

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/task` | 创建任务 |
| GET | `/api/task/{id}` | 获取任务详情 |
| GET | `/api/task/user/{user_id}` | 用户任务列表 |
| PUT | `/api/task/{id}` | 更新任务 |
| DELETE | `/api/task/{id}` | 删除任务 |

### 打卡

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/checkin/{task_id}` | 打卡 |
| GET | `/api/checkin/user/{user_id}` | 打卡记录 |

### 积分钱包

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/wallet/{user_id}` | 钱包流水 |
| GET | `/api/wallet/{user_id}/balance` | 积分余额 |
| POST | `/api/wallet/spend` | 消费积分 |

### 健康检查

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/health` | 服务状态 |

## 请求示例

```bash
# 注册用户
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456","email":"test@example.com"}'

# 创建任务
curl -X POST http://localhost:8080/api/task \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"title":"每日运动","circle_mode":"workday","points":10}'

# 打卡
curl -X POST http://localhost:8080/api/checkin/1 \
  -H "Content-Type: application/json" \
  -d '{"user_id":1}'

# 查看积分
curl http://localhost:8080/api/wallet/1/balance
```

## 数据模型

### User (用户)
- `id`: 用户ID
- `username`: 用户名
- `password`: 密码
- `email`: 邮箱
- `points`: 当前积分
- `created_at`: 创建时间

### Task (任务)
- `id`: 任务ID
- `user_id`: 所属用户
- `title`: 任务标题
- `description`: 任务描述
- `circle_mode`: 循环模式 (once/weekly/workday/weekend/custom)
- `points`: 完成奖励积分
- `is_expired`: 是否过期

### CheckIn (打卡)
- `id`: 打卡ID
- `task_id`: 任务ID
- `user_id`: 用户ID
- `points`: 本次获得积分
- `check_time`: 打卡时间

### Wallet (钱包)
- `id`: 记录ID
- `user_id`: 用户ID
- `balance`: 余额
- `type`: 类型 (earn/spend)
- `amount`: 金额
- `description`: 描述

## 定时提醒

系统每分钟检查活跃任务，根据循环模式判断是否需要提醒。提醒方式目前为日志输出，可扩展：
- 邮件通知
- Webhook
- 推送服务

## 开发计划

- [ ] 用户认证 (JWT)
- [ ] 任务提醒邮件/Webhook
- [ ] 积分兑换商品
- [ ] 任务统计报表
- [ ] 前端界面