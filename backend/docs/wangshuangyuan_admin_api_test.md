# 王双媛管理员日志模块接口测试说明

## 一、测试范围

本说明用于测试管理员操作日志模块：

```http
GET /admin/logs
```

内部能力：

```go
adminService.LogAction(ctx, input)
```

模块代码位置：

```text
backend/internal/admin/
backend/cmd/server/main.go
```

## 二、测试前置条件

本地服务默认地址：

```text
http://localhost:8080
```

数据库需要已执行：

```text
backend/scripts/sql/schema.sql
```

管理员日志模块依赖表：

```text
users
admin_logs
```

## 三、PowerShell 测试命令

### 3.1 获取管理员 Token

```powershell
$baseUrl = "http://localhost:8080"

$adminLogin = Invoke-RestMethod `
  -Method Post `
  -Uri "$baseUrl/auth/dev-login" `
  -ContentType "application/json" `
  -Body '{"openid":"dev_admin_user_001","role":"ADMIN"}'

$adminToken = $adminLogin.data.token
$adminId = $adminLogin.data.user.id
$adminHeaders = @{ Authorization = "Bearer $adminToken" }

$adminId
```

### 3.2 插入测试日志

进入 MySQL：

```powershell
mysql -h 127.0.0.1 -P 3306 -u root -p cau_used_goods
```

将下面 SQL 中的 `1` 替换成上一步输出的 `$adminId`：

```sql
INSERT INTO admin_logs (
  admin_id,
  operation_type,
  target_type,
  target_id,
  description,
  ip_address
) VALUES
(1, 'CREATE_NOTICE', 'NOTICE', 1001, '创建公告测试日志', '127.0.0.1'),
(1, 'UPDATE_WORD', 'WORD', 2001, '更新敏感词测试日志', '127.0.0.1');
```

退出 MySQL：

```sql
exit;
```

### 3.3 查询管理员日志列表

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/admin/logs?page=1&pageSize=20" `
  -Headers $adminHeaders
```

预期：

- 返回 `code = 0`。
- `data.items` 中包含测试日志。
- `data.total` 大于等于 2。

### 3.4 按操作类型筛选

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/admin/logs?operationType=CREATE_NOTICE&page=1&pageSize=20" `
  -Headers $adminHeaders
```

预期：

- 只返回 `operationType = CREATE_NOTICE` 的日志。

### 3.5 按目标类型筛选

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/admin/logs?targetType=WORD&page=1&pageSize=20" `
  -Headers $adminHeaders
```

预期：

- 只返回 `targetType = WORD` 的日志。

### 3.6 按管理员 ID 筛选

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/admin/logs?adminId=$adminId&page=1&pageSize=20" `
  -Headers $adminHeaders
```

预期：

- 只返回当前管理员的日志。

## 四、异常场景测试

### 4.1 未登录访问

```powershell
try {
  Invoke-RestMethod `
    -Method Get `
    -Uri "$baseUrl/admin/logs"
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
401
```

### 4.2 普通用户访问管理员接口

```powershell
$userLogin = Invoke-RestMethod `
  -Method Post `
  -Uri "$baseUrl/auth/dev-login" `
  -ContentType "application/json" `
  -Body '{"openid":"dev_normal_user_001","role":"USER"}'

$userToken = $userLogin.data.token
$userHeaders = @{ Authorization = "Bearer $userToken" }

try {
  Invoke-RestMethod `
    -Method Get `
    -Uri "$baseUrl/admin/logs" `
    -Headers $userHeaders
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
403
```

## 五、内部调用说明

公告管理、敏感词管理等管理员操作完成后，可调用：

```go
description := "创建公告"
ip := c.ClientIP()

logID, err := adminService.LogAction(ctx, admin.LogActionInput{
    AdminID:       adminID,
    OperationType: "CREATE_NOTICE",
    TargetType:    admin.TargetTypeNotice,
    TargetID:      noticeID,
    Description:   &description,
    IPAddress:     &ip,
})
```

说明：

- `LogAction` 不是前端接口。
- `admin_id`、`operation_type`、`target_type`、`target_id` 必填。
- `description` 最长 500。
- `ip_address` 最长 50。

## 六、测试结论模板

```text
测试模块：管理员操作日志模块

测试接口：
- GET /admin/logs

测试结果：
- 管理员鉴权正常
- 普通用户访问被拦截
- 日志列表查询正常
- 按管理员、操作类型、目标类型筛选正常
- 分页参数正常

结论：管理员操作日志模块接口测试通过
```

## 七、提交建议

```bash
git status
git add backend/internal/admin backend/cmd/server/main.go backend/docs/wangshuangyuan_module_README.md backend/docs/wangshuangyuan_admin_api_test.md
git commit -m "feat(admin): 完成管理员操作日志基础接口"
```
