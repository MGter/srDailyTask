# 用户资料和头像上传设计

## 背景
用户希望可以修改自己的基础信息，并新增头像上传功能。当前用户信息只包含用户名、邮箱、积分，密码以明文存储，本次按用户确认暂不改密码存储方式。

## 入口和 UI
采用“首页独立资料卡片”方案。

资料卡片放在 `web/src/views/Tasks.vue` 左侧栏最上方，位于“今日打卡任务”上面。卡片展示：
- 圆形头像
- 用户名
- 邮箱
- 当前总积分
- 编辑资料按钮

没有头像时，前端显示用户名首字符作为默认头像。

点击“编辑资料”打开弹窗。弹窗字段包括：
- 用户名
- 邮箱
- 头像文件选择和预览
- 旧密码
- 新密码
- 确认新密码

不填写新密码时，只更新基础资料，不修改密码。填写新密码时必须填写旧密码，并校验确认密码一致。

## 后端数据模型
`users` 表新增字段：

```sql
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255);
```

`model.User` 新增 `AvatarURL string json:"avatar_url"`，并在用户查询、登录、注册返回中带回该字段。

## 后端接口
新增或扩展用户接口：

### 更新资料
`PUT /api/user/{id}`

请求体：
```json
{
  "username": "MGter",
  "email": "mgter@example.com",
  "old_password": "old-password",
  "new_password": "new-password"
}
```

规则：
- `username` 不能为空。
- 修改用户名时不能和其他用户重复。
- `email` 允许为空；非空时做基础邮箱格式校验。
- `new_password` 为空时不修改密码。
- `new_password` 非空时必须提供 `old_password`，旧密码正确才更新密码。
- 本次不改变现有明文密码存储方式。

### 上传头像
`POST /api/user/{id}/avatar`

请求类型：`multipart/form-data`，字段名：`avatar`。

规则：
- 只允许 `.jpg`、`.jpeg`、`.png`、`.webp`、`.gif`。
- 文件大小不超过 2MB。
- 保存到项目根目录下的 `uploads/avatars/`。
- 文件名使用用户 ID、时间戳和随机值生成，避免覆盖和路径注入。
- 上传成功后更新 `users.avatar_url`，返回最新用户信息。

## 静态文件访问
后端在 `SetupServer` 中新增 `/uploads/` 静态目录映射，浏览器通过 `avatar_url` 直接访问头像，例如：

```text
/uploads/avatars/1_1777450000_ab12cd34.webp
```

上传目录属于运行时数据，不提交到 git。

## 前端数据流
`web/src/api/index.js` 增加：
- `updateUser(id, data)` 调用 `PUT /user/{id}`。
- `uploadAvatar(id, formData)` 调用 `POST /user/{id}/avatar`，请求头使用 `multipart/form-data`。

`Tasks.vue`：
1. `loadUser()` 获取用户数据，包含 `avatar_url`。
2. 资料卡片展示用户信息。
3. 打开弹窗时用当前用户信息初始化表单。
4. 保存时先做前端校验：确认密码一致、头像类型和大小合法。
5. 调用更新资料接口。
6. 如果选择了头像，再调用上传头像接口。
7. 成功后刷新 `user`，同步写入 `localStorage.user`。

## 错误处理
- 后端返回明确错误：用户名重复、旧密码错误、头像过大、头像格式不支持。
- 前端保存失败时显示后端错误。
- 保存成功后关闭弹窗，并提示资料已更新。

## 验证
- 资料卡片显示头像、用户名、邮箱、总积分。
- 修改用户名和邮箱后刷新页面仍生效。
- 不填新密码时不会修改密码。
- 旧密码错误时无法修改密码。
- 新密码和确认密码不一致时前端阻止提交。
- 上传合法头像后，页面立即显示新头像。
- 上传非图片或超过 2MB 的文件会提示错误。
- `go test ./...` 成功。
- `cd web && npm run build` 成功。
- `systemctl --user restart srdailytask` 后 `curl -s http://localhost:18888/health` 返回 `{"status":"ok"}`。
