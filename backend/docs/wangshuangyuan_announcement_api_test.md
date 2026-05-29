# 王双媛公告管理接口测试说明

## 一、测试范围

本说明用于测试管理员公告管理接口：

```http
GET /admin/announcements
POST /admin/announcements
PUT /admin/announcements/:id
PUT /admin/announcements/:id/status
DELETE /admin/announcements/:id
```

模块代码位置：

```text
backend/internal/admin/
backend/cmd/server/main.go
```

说明：

- 公告管理接口仅管理员可访问。
- 每次新增、编辑、改状态、删除公告都会写入 `admin_logs`。
- 由于 `announcements` 表没有 `is_deleted` 字段，`DELETE` 当前按下线处理，即设置 `status = OFFLINE`。

## 二、测试前置条件

本地服务默认地址：

```text
http://localhost:8080
```

数据库需要已执行：

```text
backend/scripts/sql/schema.sql
```

公告模块依赖表：

```text
users
announcements
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
  -Body '{"openid":"dev_admin_notice_001","role":"ADMIN"}'

$adminToken = $adminLogin.data.token
$adminHeaders = @{ Authorization = "Bearer $adminToken" }
```

### 3.2 创建公告

```powershell
$createResult = Invoke-RestMethod `
  -Method Post `
  -Uri "$baseUrl/admin/announcements" `
  -Headers $adminHeaders `
  -ContentType "application/json" `
  -Body '{"title":"测试公告","content":"这是一条测试公告","coverUrl":"/uploads/notices/test.png","status":"DRAFT"}'

$noticeId = $createResult.data.id
$noticeId
```

预期：

- 返回 `code = 0`。
- 返回公告 ID。
- `admin_logs` 中新增 `CREATE_NOTICE` 日志。

### 3.3 查询公告列表

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/admin/announcements?page=1&pageSize=20" `
  -Headers $adminHeaders
```

预期：

- 返回公告列表。
- 包含刚创建的公告。

### 3.4 按状态筛选公告

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/admin/announcements?status=DRAFT&page=1&pageSize=20" `
  -Headers $adminHeaders
```

预期：

- 只返回 `status = DRAFT` 的公告。

### 3.5 按关键字搜索公告

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/admin/announcements?keyword=测试&page=1&pageSize=20" `
  -Headers $adminHeaders
```

预期：

- 返回标题或内容中包含 `测试` 的公告。

### 3.6 编辑公告

```powershell
Invoke-RestMethod `
  -Method Put `
  -Uri "$baseUrl/admin/announcements/$noticeId" `
  -Headers $adminHeaders `
  -ContentType "application/json" `
  -Body '{"title":"测试公告-已编辑","content":"公告内容已更新","coverUrl":"/uploads/notices/updated.png"}'
```

预期：

- 返回 `updated = true`。
- `admin_logs` 中新增 `UPDATE_NOTICE` 日志。

### 3.7 发布公告

```powershell
Invoke-RestMethod `
  -Method Put `
  -Uri "$baseUrl/admin/announcements/$noticeId/status" `
  -Headers $adminHeaders `
  -ContentType "application/json" `
  -Body '{"status":"PUBLISHED"}'
```

预期：

- 返回 `updated = true`。
- 公告状态变为 `PUBLISHED`。
- `publish_time` 不为空。
- `admin_logs` 中新增 `STATUS_NOTICE` 日志。

### 3.8 下线公告

```powershell
Invoke-RestMethod `
  -Method Put `
  -Uri "$baseUrl/admin/announcements/$noticeId/status" `
  -Headers $adminHeaders `
  -ContentType "application/json" `
  -Body '{"status":"OFFLINE"}'
```

预期：

- 返回 `updated = true`。
- 公告状态变为 `OFFLINE`。

### 3.9 删除公告

```powershell
Invoke-RestMethod `
  -Method Delete `
  -Uri "$baseUrl/admin/announcements/$noticeId" `
  -Headers $adminHeaders
```

预期：

- 返回 `deleted = true`。
- 公告状态为 `OFFLINE`。
- `admin_logs` 中新增 `DELETE_NOTICE` 日志。

## 四、异常场景测试

### 4.1 普通用户不能访问公告管理

```powershell
$userLogin = Invoke-RestMethod `
  -Method Post `
  -Uri "$baseUrl/auth/dev-login" `
  -ContentType "application/json" `
  -Body '{"openid":"dev_user_notice_001","role":"USER"}'

$userHeaders = @{ Authorization = "Bearer $userLogin.data.token" }

try {
  Invoke-RestMethod `
    -Method Get `
    -Uri "$baseUrl/admin/announcements" `
    -Headers $userHeaders
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
403
```

### 4.2 创建公告缺少标题

```powershell
try {
  Invoke-RestMethod `
    -Method Post `
    -Uri "$baseUrl/admin/announcements" `
    -Headers $adminHeaders `
    -ContentType "application/json" `
    -Body '{"content":"缺少标题"}'
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
400
```

### 4.3 非法公告状态

```powershell
try {
  Invoke-RestMethod `
    -Method Put `
    -Uri "$baseUrl/admin/announcements/$noticeId/status" `
    -Headers $adminHeaders `
    -ContentType "application/json" `
    -Body '{"status":"INVALID"}'
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
400
```

### 4.4 不存在公告

```powershell
try {
  Invoke-RestMethod `
    -Method Delete `
    -Uri "$baseUrl/admin/announcements/999999" `
    -Headers $adminHeaders
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
404
```

## 五、测试结论模板

```text
测试模块：管理员公告管理模块

测试接口：
- GET /admin/announcements
- POST /admin/announcements
- PUT /admin/announcements/:id
- PUT /admin/announcements/:id/status
- DELETE /admin/announcements/:id

测试结果：
- 管理员鉴权正常
- 普通用户访问被拦截
- 公告创建、查询、编辑、发布、下线正常
- 删除接口按下线处理正常
- 管理员操作日志写入正常
- 非法参数和不存在公告处理正常

结论：管理员公告管理模块接口测试通过
```

## 六、提交建议

```bash
git status
git add backend/internal/admin backend/docs/wangshuangyuan_module_README.md backend/docs/wangshuangyuan_announcement_api_test.md
git commit -m "feat(admin): 完成公告管理接口"
```
