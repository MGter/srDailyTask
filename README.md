# srDailyTask - 日常打卡积分系统

srDailyTask 是一个日常任务、积分钱包和长期物品成本记录工具。后端使用 Go 标准库提供 API，前端使用 Vue 3 构建移动端优先的 Web 界面。

## 功能特性

- **用户账户**: 注册、登录、资料修改、密码修改、头像上传。
- **每日打卡**: 创建周期任务，支持打卡、取消打卡、跳过任务。
- **周期模式**: 单次、每周一、工作日、周末、每天。
- **积分钱包**: 手动添加收入/支出，删除记录，搜索钱包记录。
- **积分统计**: 月活跃度热力图、累计积分图、当日累计/消耗/结余/总计。
- **长期主义**: 记录购买物品、价格和购买日期，按使用天数计算日均成本。
- **报废管理**: 支持报废长期物品，报废后日均成本冻结并进入报废栏。
- **搜索排序**: 长期主义“使用中”支持全字段搜索、时间/价格/日均价格排序。
- **移动端适配**: 手机端按钮竖排、积分卡片紧凑布局、滚轮式日期选择。
- **系统服务**: 支持用户级 `systemd` 管理。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + Vite + Chart.js |
| 后端 | Go 1.21 + 标准库 `net/http` |
| 数据库 | MySQL |
| 定时任务 | robfig/cron |
| 部署 | 单二进制 + `web/dist` 静态资源 |

## 项目结构

```text
srDailyTask/
├── cmd/
│   └── main.go                         # 程序入口
├── config/
│   └── config.yaml                     # 服务和数据库配置
├── internal/
│   ├── config/                         # 配置加载
│   ├── handler/                        # HTTP 处理器和路由
│   │   ├── router.go
│   │   ├── user_handler.go
│   │   ├── task_handler.go
│   │   ├── point_handler.go
│   │   ├── wallet_handler.go
│   │   └── long_term_item_handler.go
│   ├── model/                          # 数据模型
│   │   ├── user.go
│   │   ├── task.go
│   │   ├── checkin.go
│   │   ├── wallet.go
│   │   └── long_term_item.go
│   ├── repository/                     # 数据访问层
│   ├── service/                        # 业务逻辑层
│   ├── scheduler/                      # 定时提醒
│   └── logger/                         # 日志工具
├── migrations/
│   ├── 001_init.sql
│   ├── 002_add_checkin_id.sql
│   ├── 003_add_user_avatar.sql
│   └── 004_add_long_term_items.sql
├── web/
│   ├── src/
│   │   ├── views/Tasks.vue             # 主界面
│   │   ├── api/index.js                # API 封装
│   │   ├── router/index.js
│   │   ├── App.vue
│   │   ├── main.js
│   │   └── style.css
│   ├── dist/                           # 前端构建产物
│   └── package.json
├── uploads/                            # 用户上传头像
├── pics/                               # 背景图源文件
├── go.mod
├── go.sum
└── Makefile
```

## 数据库迁移

初始化数据库：

```bash
mysql -u root -e "CREATE DATABASE daily_task CHARACTER SET utf8mb4;"
mysql -u root daily_task < migrations/001_init.sql
mysql -u root daily_task < migrations/002_add_checkin_id.sql
mysql -u root daily_task < migrations/003_add_user_avatar.sql
mysql -u root daily_task < migrations/004_add_long_term_items.sql
```

### 主要表

- `users`: 用户、积分余额、头像地址等资料。
- `tasks`: 打卡任务和周期配置。
- `checkins`: 打卡记录，跳过任务也会记录，但积分为 `0`。
- `wallet`: 积分钱包记录，打卡积分和手动收入/支出都会进入这里。
- `long_term_items`: 长期主义物品，记录价格、购买日期、报废日期和冻结日均成本。

## API 接口

### 用户相关

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/user/register` | 注册用户 |
| POST | `/api/user/login` | 用户登录 |
| GET | `/api/user/{id}` | 获取用户信息 |
| PUT | `/api/user/{id}` | 修改用户资料/密码 |
| POST | `/api/user/{id}/avatar` | 上传头像 |
| GET | `/api/users` | 获取用户列表 |

### 任务相关

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/task` | 创建任务 |
| GET | `/api/task/{id}` | 获取任务详情 |
| GET | `/api/task/user/{user_id}` | 获取用户任务 |
| GET | `/api/task/today/{user_id}` | 获取今日任务 |
| PUT | `/api/task/{id}` | 修改任务 |
| DELETE | `/api/task/{id}` | 删除任务 |

### 打卡相关

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/checkin/{task_id}` | 打卡并获得积分 |
| DELETE | `/api/checkin/{task_id}` | 取消打卡 |
| POST | `/api/checkin/{task_id}/skip` | 跳过任务，不计积分 |
| GET | `/api/checkin/user/{user_id}` | 获取打卡历史 |
| GET | `/api/checkin/today/{user_id}` | 获取今日已处理任务 ID |

### 积分钱包相关

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/wallet/{user_id}` | 获取钱包记录 |
| GET | `/api/wallet/{user_id}/balance` | 获取积分余额 |
| POST | `/api/wallet/spend` | 消费积分 |
| POST | `/api/wallet/add` | 添加收入/支出 |
| DELETE | `/api/wallet/delete/{id}` | 删除钱包记录 |

### 积分统计相关

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/points/{user_id}` | 获取积分历史 |
| GET | `/api/points/daily/{user_id}?days=N` | 获取每日统计 |

### 长期主义相关

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/long-term-items/{user_id}` | 获取长期主义物品和汇总 |
| POST | `/api/long-term-items` | 新增物品 |
| PUT | `/api/long-term-items/{id}` | 修改名称、价格、购买日期 |
| POST | `/api/long-term-items/{id}/scrap` | 报废物品并冻结日均成本 |
| DELETE | `/api/long-term-items/{id}` | 删除物品 |

## 前端界面

主界面位于 `web/src/views/Tasks.vue`。

### 个人账户

- 显示用户名、邮箱、头像。
- 显示当日累计、当日消耗、当日结余、总计、日均消费。
- 支持编辑资料、修改密码、上传头像。

### 每日打卡

- 今日任务列表：未打卡任务靠前，已处理任务灰色显示。
- 支持创建、修改、删除任务。
- 支持打卡、取消打卡、跳过任务。
- 跳过任务只写入打卡记录，不写入钱包积分。

### 积分钱包

- 支持添加收入/支出。
- 支持滚轮式日期时间选择。
- 支持全字段搜索钱包记录。
- 支持重复记录和删除记录。

### 月活跃度和累计积分

- 月活跃度使用类似 GitHub 的热力图展示最近两个月每日结余。
- 绿色表示结余，红色表示消耗，颜色深度封顶。
- 累计积分图支持切换 7/14/30/180/360 天。

### 长期主义

- 记录“买了什么、价格、购买日期”。
- 当前日均成本 = 价格 / 已拥有天数。
- 使用中的物品参与个人账户“日均消费”汇总。
- 物品报废后进入报废栏，日均成本冻结。
- 使用中列表支持全字段搜索。
- 使用中列表支持时间顺序、时间倒序、价格顺序、价格倒序、日均价格倒序。
- 默认按日均价格倒序排列。

### 移动端适配

- 手机端优先使用紧凑卡片布局。
- 操作按钮统一为浅色胶囊风格。
- 任务操作按钮在手机端竖排显示，便于点击。
- 个人账户积分卡片在手机端第一行固定显示 4 项，日均消费单独占一行。
- 页面限制横向溢出，避免手机端需要手动缩放。

## 快速开始

### 1. 安装依赖

```bash
go mod download
cd web && npm install
```

### 2. 配置数据库

编辑 `config/config.yaml`：

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

### 3. 构建和运行

```bash
# 完整打包：构建前端和后端
make pack

# 或分步构建
cd web && npm run build
cd ..
go build -o daily_task ./cmd

# 运行
./daily_task
```

访问：`http://localhost:18888`

## 用户级 systemd 服务

项目可以作为用户级 systemd 服务运行，服务名为 `srdailytask`。

```bash
# 启动/停止/重启
systemctl --user start srdailytask
systemctl --user stop srdailytask
systemctl --user restart srdailytask

# 查看状态和日志
systemctl --user status srdailytask
journalctl --user -u srdailytask -f

# 开机自启
systemctl --user enable srdailytask
systemctl --user disable srdailytask
```

后端代码变更后需要重新构建二进制再重启：

```bash
go build -o daily_task ./cmd
systemctl --user restart srdailytask
```

前端代码变更后需要重新构建 `web/dist` 再重启或刷新服务资源：

```bash
cd web && npm run build
systemctl --user restart srdailytask
```

## 部署文件

```text
daily_task              # 后端二进制
config/config.yaml      # 配置文件
web/dist/               # 前端静态文件
uploads/                # 用户头像上传目录
```

## 常用命令

```bash
# Go 测试
go test ./...

# 前端构建
cd web && npm run build

# 健康检查
curl http://localhost:18888/health
```

## 更新历史

- 2026-04-18: 项目从 C++ 迁移到 Go。
- 2026-04-19: 添加任务级别、周期模式、积分钱包、统计图表和背景图。
- 2026-04-20: 添加跳过任务、分页显示和手机端按钮竖排。
- 2026-04-24: 添加累计积分图、打包命令和用户级 systemd 服务。
- 2026-04-29:
  - 添加用户资料编辑和头像上传。
  - 添加月活跃度热力图。
  - 添加长期主义模块。
  - 添加长期主义搜索、排序、修改、报废和日均消费汇总。
  - 积分钱包支持搜索和滚轮式日期时间选择。
  - 统一玻璃卡片 UI、浅色胶囊按钮和手机端布局。
