# 王双媛消息通知模块接口测试说明

## 一、测试范围

本说明用于测试消息通知模块接口：

```http
GET /messages
GET /messages?readStatus=UNREAD
GET /messages?readStatus=READ
GET /messages/unread-count
GET /messages/:id
PUT /messages/:id/read
PUT /messages/read-all
```

模块代码位置：

```text
backend/internal/message/
backend/cmd/server/main.go
```

## 二、测试前置条件

本地服务默认地址：

```text
http://localhost:8080
```

确认 `backend/config/config.yaml` 中：

```yaml
server:
  env: dev
  port: 8080
```

`env=dev` 时，可以使用开发登录接口获取 JWT token。

数据库需要已执行：

```text
backend/scripts/sql/schema.sql
```

消息模块依赖表：

```text
users
messages
```

## 三、PowerShell 测试命令

### 3.1 获取测试 Token

```powershell
$baseUrl = "http://localhost:8080"

$login = Invoke-RestMethod `
  -Method Post `
  -Uri "$baseUrl/auth/dev-login" `
  -ContentType "application/json" `
  -Body '{"openid":"dev_message_user_001","role":"USER"}'

$token = $login.data.token
$userId = $login.data.user.id
$headers = @{ Authorization = "Bearer $token" }

$userId
```

### 3.2 插入测试消息

进入 MySQL：

```powershell
mysql -h 127.0.0.1 -P 3306 -u root -p cau_used_goods
```

将下面 SQL 中的 `1` 替换成上一步输出的 `$userId`：

```sql
INSERT INTO messages (
  receiver_id,
  sender_id,
  message_type,
  title,
  content,
  related_type,
  related_id,
  read_status
) VALUES
(1, NULL, 'SYSTEM_NOTICE', '测试系统通知', '这是一条未读测试消息', 'NOTICE', NULL, 'UNREAD'),
(1, NULL, 'ORDER_CONFIRMED', '订单已确认', '卖家已确认你的预约', 'ORDER', 1001, 'UNREAD'),
(1, NULL, 'REPORT_HANDLED', '举报处理结果', '你的举报已处理', 'REPORT', 2001, 'READ');
```

退出 MySQL：

```sql
exit;
```

### 3.3 查询消息列表

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/messages?page=1&pageSize=20" `
  -Headers $headers
```

预期：

- 返回当前用户的消息。
- `total` 大于等于 3。
- 按 `create_time DESC` 排序。

### 3.4 按未读状态筛选

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/messages?readStatus=UNREAD&page=1&pageSize=20" `
  -Headers $headers
```

预期：

- 只返回 `readStatus = UNREAD` 的消息。

### 3.5 按已读状态筛选

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/messages?readStatus=READ&page=1&pageSize=20" `
  -Headers $headers
```

预期：

- 只返回 `readStatus = READ` 的消息。

### 3.6 非法已读状态筛选

```powershell
try {
  Invoke-RestMethod `
    -Method Get `
    -Uri "$baseUrl/messages?readStatus=INVALID" `
    -Headers $headers
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期输出：

```text
400
```

### 3.7 查询未读数量

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/messages/unread-count" `
  -Headers $headers
```

预期：

- `data.count` 等于当前用户未读消息数量。

### 3.8 查询消息详情

先从消息列表中选择一条消息 ID，例如 `1`：

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/messages/1" `
  -Headers $headers
```

预期：

- 返回该消息完整内容。
- 返回的 `receiverId` 是当前登录用户 ID。

### 3.9 标记单条消息已读

```powershell
Invoke-RestMethod `
  -Method Put `
  -Uri "$baseUrl/messages/1/read" `
  -Headers $headers
```

预期：

```json
{
  "read": true
}
```

数据库检查：

```sql
SELECT id, receiver_id, read_status, read_time
FROM messages
WHERE id = 1;
```

预期：

```text
read_status = READ
read_time 不为空
```

### 3.10 全部标记已读

```powershell
Invoke-RestMethod `
  -Method Put `
  -Uri "$baseUrl/messages/read-all" `
  -Headers $headers
```

预期：

- `data.read = true`
- `data.count` 表示本次被更新为已读的消息数量。

再查询未读数量：

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/messages/unread-count" `
  -Headers $headers
```

预期：

```text
data.count = 0
```

## 四、异常场景测试

### 4.1 未登录访问

```powershell
try {
  Invoke-RestMethod `
    -Method Get `
    -Uri "$baseUrl/messages"
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
401
```

### 4.2 非法消息 ID

```powershell
try {
  Invoke-RestMethod `
    -Method Put `
    -Uri "$baseUrl/messages/abc/read" `
    -Headers $headers
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
400
```

### 4.3 不存在消息

```powershell
try {
  Invoke-RestMethod `
    -Method Get `
    -Uri "$baseUrl/messages/999999" `
    -Headers $headers
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
404
```

### 4.4 不能查看或操作他人消息

登录第二个用户：

```powershell
$login2 = Invoke-RestMethod `
  -Method Post `
  -Uri "$baseUrl/auth/dev-login" `
  -ContentType "application/json" `
  -Body '{"openid":"dev_message_user_002","role":"USER"}'

$token2 = $login2.data.token
$headers2 = @{ Authorization = "Bearer $token2" }
```

用用户 2 查询用户 1 的消息，例如消息 ID 为 `1`：

```powershell
try {
  Invoke-RestMethod `
    -Method Get `
    -Uri "$baseUrl/messages/1" `
    -Headers $headers2
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
404
```

用用户 2 标记用户 1 的消息已读：

```powershell
try {
  Invoke-RestMethod `
    -Method Put `
    -Uri "$baseUrl/messages/1/read" `
    -Headers $headers2
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
404
```

## 五、Apifox 测试建议

建议新建分组：

```text
王双媛-消息通知模块
```

环境变量：

```text
baseUrl = http://localhost:8080
token = dev-login 返回的 token
```

公共请求头：

```http
Authorization: Bearer {{token}}
Content-Type: application/json
```

接口列表：

```text
POST {{baseUrl}}/auth/dev-login
GET  {{baseUrl}}/messages?page=1&pageSize=20
GET  {{baseUrl}}/messages?readStatus=UNREAD&page=1&pageSize=20
GET  {{baseUrl}}/messages?readStatus=READ&page=1&pageSize=20
GET  {{baseUrl}}/messages/unread-count
GET  {{baseUrl}}/messages/:id
PUT  {{baseUrl}}/messages/:id/read
PUT  {{baseUrl}}/messages/read-all
```

## 六、测试结论模板

```text
测试模块：消息通知模块

测试接口：
- GET /messages
- GET /messages?readStatus=UNREAD
- GET /messages?readStatus=READ
- GET /messages/unread-count
- GET /messages/:id
- PUT /messages/:id/read
- PUT /messages/read-all

测试结果：
- 登录鉴权正常
- 消息列表查询正常
- 已读/未读筛选正常
- 消息详情查询正常
- 未读数量统计正常
- 单条消息标记已读正常
- 全部消息标记已读正常
- 非本人消息不可查看和操作
- 未登录访问被拦截
- 非法参数和不存在消息处理正常

结论：消息通知模块接口测试通过
```

## 七、提交建议

```bash
git status
git add backend/internal/message backend/docs/wangshuangyuan_message_api_test.md backend/docs/wangshuangyuan_module_README.md
git commit -m "feat(message): 补充消息详情和批量已读接口"
```
