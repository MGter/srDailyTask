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
- **Web 界面**: Vue 3 前端，响应式设计

## 技术栈

- **后端**: Go 1.21 + MySQL
- **前端**: Vue 3 + Vite
- **定时任务**: `robfig/cron`

## 项目结构

```
srDailyTask/
├── cmd/main.go              # 程序入口
├── config/config.yaml       # 配置文件
├── internal/                # 后端代码
│   ├── handler/             # HTTP 处理器
│   ├── service/             # 业务逻辑
│   ├── repository/          # 数据库操作
│   ├── model/               # 数据模型
│   ├── logger/              # 日志
│   ├── scheduler/           # 定时提醒
│   └── config/              # 配置加载
├── web/                     # Vue 3 前端
│   ├── src/
│   │   ├── views/           # 页面组件
│   │   ├── api/             # API 调用
│   │   └── router/          # 路由配置
│   └── dist/                # 打包后的静态文件
├── migrations/              # 数据库迁移
├── go.mod
├── Makefile
└── docs/
```

## 快速开始

### 1. 安装依赖

```bash
# Go 后端依赖
make deps

# 前端依赖（如需开发）
cd web && npm install
```

### 2. 创建数据库

```bash
mysql -u root -p -e "CREATE DATABASE daily_task CHARACTER SET utf8mb4;"
mysql -u root -p daily_task < migrations/001_init.sql
```

### 3. 配置

复制配置文件模板：
```bash
cp config/config.yaml.example config/config.yaml
```

修改数据库连接信息：
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
# 编译并运行（前端已打包）
make build
./daily_task

# 访问 http://localhost:8080
```

## 部署指南

### 单机部署（推荐）

一个二进制文件 + 配置即可运行，后端托管前端静态文件：

```bash
# 1. 打包前端
cd web && npm run build

# 2. 编译后端
cd .. && go build -o daily_task cmd/main.go

# 3. 部署文件
daily_task          # 二进制文件
config/config.yaml  # 配置文件
web/dist/           # 前端静态文件目录

# 4. 运行
./daily_task
```

### 开发模式

前后端分离运行：

```bash
# 后端（端口 8080）
./daily_task

# 前端开发服务器（端口 3000）
cd web && npm run dev

# 访问 http://localhost:3000
```

## API 接口

### 用户

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/user/register` | 注册用户 |
| POST | `/api/user/login` | 用户登录 |
| GET | `/api/user/{id}` | 获取用户信息 |

### 任务

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/task` | 创建任务 |
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
| GET | `/api/wallet/{user_id}/balance` | 积分余额 |
| POST | `/api/wallet/spend` | 消费积分 |

## Web 界面

访问 http://localhost:8080 即可使用 Web 界面：

- **登录/注册**: 用户认证
- **任务列表**: 查看和创建任务，一键打卡
- **积分钱包**: 查看积分余额和流水记录

## 开发计划

- [ ] 用户认证 (JWT)
- [ ] 任务提醒邮件/Webhook
- [ ] 积分兑换商品
- [ ] 任务统计报表