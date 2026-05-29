# 王双媛敏感词管理接口测试说明

## 一、测试范围

本说明用于测试管理员敏感词管理接口：

```http
GET /admin/sensitive-words
POST /admin/sensitive-words
PUT /admin/sensitive-words/:id
DELETE /admin/sensitive-words/:id
```

说明：

- 敏感词管理接口仅管理员可访问。
- 新增、更新、删除敏感词会写入 `admin_logs`。
- `DELETE` 当前按禁用处理，即设置 `status = DISABLED`。

## 二、测试前置条件

本地服务默认地址：

```text
http://localhost:8080
```

数据库需要已执行：

```text
backend/scripts/sql/schema.sql
```

依赖表：

```text
users
sensitive_words
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
  -Body '{"openid":"dev_admin_word_001","role":"ADMIN"}'

$adminToken = $adminLogin.data.token
$adminHeaders = @{ Authorization = "Bearer $adminToken" }
```

### 3.2 新增敏感词

```powershell
$createResult = Invoke-RestMethod `
  -Method Post `
  -Uri "$baseUrl/admin/sensitive-words" `
  -Headers $adminHeaders `
  -ContentType "application/json" `
  -Body '{"word":"测试敏感词001","wordType":"FORBIDDEN","status":"ENABLED"}'

$wordId = $createResult.data.id
$wordId
```

预期：

- 返回 `code = 0`。
- 返回敏感词 ID。
- `admin_logs` 中新增 `CREATE_WORD` 日志。

### 3.3 查询敏感词列表

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/admin/sensitive-words?page=1&pageSize=20" `
  -Headers $adminHeaders
```

预期：

- 返回敏感词列表。
- 包含刚创建的敏感词。

### 3.4 按状态筛选

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/admin/sensitive-words?status=ENABLED&page=1&pageSize=20" `
  -Headers $adminHeaders
```

预期：

- 只返回 `status = ENABLED` 的敏感词。

### 3.5 按类型筛选

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/admin/sensitive-words?wordType=FORBIDDEN&page=1&pageSize=20" `
  -Headers $adminHeaders
```

预期：

- 只返回 `wordType = FORBIDDEN` 的敏感词。

### 3.6 按关键字搜索

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/admin/sensitive-words?keyword=测试&page=1&pageSize=20" `
  -Headers $adminHeaders
```

预期：

- 返回 `word` 中包含 `测试` 的敏感词。

### 3.7 更新敏感词

```powershell
Invoke-RestMethod `
  -Method Put `
  -Uri "$baseUrl/admin/sensitive-words/$wordId" `
  -Headers $adminHeaders `
  -ContentType "application/json" `
  -Body '{"word":"测试敏感词001-已更新","wordType":"RISK","status":"ENABLED"}'
```

预期：

- 返回 `updated = true`。
- `admin_logs` 中新增 `UPDATE_WORD` 日志。

### 3.8 删除/禁用敏感词

```powershell
Invoke-RestMethod `
  -Method Delete `
  -Uri "$baseUrl/admin/sensitive-words/$wordId" `
  -Headers $adminHeaders
```

预期：

- 返回 `deleted = true`。
- 数据库中该敏感词 `status = DISABLED`。
- `admin_logs` 中新增 `DELETE_WORD` 日志。

## 四、异常场景测试

### 4.1 普通用户不能访问

```powershell
$userLogin = Invoke-RestMethod `
  -Method Post `
  -Uri "$baseUrl/auth/dev-login" `
  -ContentType "application/json" `
  -Body '{"openid":"dev_user_word_001","role":"USER"}'

$userHeaders = @{ Authorization = "Bearer $userLogin.data.token" }

try {
  Invoke-RestMethod `
    -Method Get `
    -Uri "$baseUrl/admin/sensitive-words" `
    -Headers $userHeaders
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
403
```

### 4.2 新增敏感词缺少 word

```powershell
try {
  Invoke-RestMethod `
    -Method Post `
    -Uri "$baseUrl/admin/sensitive-words" `
    -Headers $adminHeaders `
    -ContentType "application/json" `
    -Body '{"wordType":"FORBIDDEN","status":"ENABLED"}'
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
400
```

### 4.3 非法敏感词类型

```powershell
try {
  Invoke-RestMethod `
    -Method Post `
    -Uri "$baseUrl/admin/sensitive-words" `
    -Headers $adminHeaders `
    -ContentType "application/json" `
    -Body '{"word":"非法类型测试","wordType":"INVALID","status":"ENABLED"}'
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
400
```

### 4.4 不存在的敏感词

```powershell
try {
  Invoke-RestMethod `
    -Method Delete `
    -Uri "$baseUrl/admin/sensitive-words/999999" `
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
测试模块：管理员敏感词管理模块

测试接口：
- GET /admin/sensitive-words
- POST /admin/sensitive-words
- PUT /admin/sensitive-words/:id
- DELETE /admin/sensitive-words/:id

测试结果：
- 管理员鉴权正常
- 普通用户访问被拦截
- 敏感词新增、查询、更新、禁用正常
- 管理员操作日志写入正常
- 非法参数和不存在敏感词处理正常

结论：管理员敏感词管理模块接口测试通过
```

## 六、提交建议

```bash
git status
git add backend/internal/sensitive backend/internal/admin backend/cmd/server/main.go backend/docs/wangshuangyuan_module_README.md backend/docs/wangshuangyuan_sensitive_api_test.md
git commit -m "feat(sensitive): 完成敏感词管理接口"
```
