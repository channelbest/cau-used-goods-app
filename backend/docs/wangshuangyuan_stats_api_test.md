# 王双媛统计接口测试说明

## 一、测试范围

本说明用于测试后台统计接口：

```http
GET /stats/products/overview
GET /stats/products/category-distribution
GET /stats/products/status-distribution
GET /stats/products/trend?days=7
GET /stats/orders/overview
GET /stats/users/overview
GET /stats/reports/overview
```

说明：

- 统计接口用于后台管理视图。
- 当前已统一增加管理员权限，必须登录且 `role=ADMIN`。

## 二、PowerShell 测试命令

### 2.1 获取管理员 Token

```powershell
$baseUrl = "http://localhost:8080"

$adminLogin = Invoke-RestMethod `
  -Method Post `
  -Uri "$baseUrl/auth/dev-login" `
  -ContentType "application/json" `
  -Body '{"openid":"dev_admin_stats_001","role":"ADMIN"}'

$adminHeaders = @{ Authorization = "Bearer $adminLogin.data.token" }
```

### 2.2 商品总览

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/stats/products/overview" `
  -Headers $adminHeaders
```

预期字段：

```text
totalProducts
onSaleProducts
offShelfProducts
lockedProducts
soldProducts
deletedProducts
totalViews
totalFavorites
averagePrice
```

### 2.3 商品分类分布

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/stats/products/category-distribution" `
  -Headers $adminHeaders
```

预期：

- 返回分类维度的商品数量、在售数量、平均价格、浏览量、收藏量。

### 2.4 商品状态分布

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/stats/products/status-distribution" `
  -Headers $adminHeaders
```

预期：

- 返回各商品状态对应数量。

### 2.5 商品发布趋势

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/stats/products/trend?days=7" `
  -Headers $adminHeaders
```

预期：

- 返回最近 7 天商品发布数量。
- `days` 最大按 90 处理。

### 2.6 订单总览

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/stats/orders/overview" `
  -Headers $adminHeaders
```

预期字段：

```text
totalOrders
pendingConfirmOrders
waitMeetOrders
completedOrders
canceledOrders
exceptionClosedOrders
totalCompletedAmount
averageCompletedPrice
```

### 2.7 用户总览

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/stats/users/overview" `
  -Headers $adminHeaders
```

预期字段：

```text
totalUsers
normalUsers
disabledUsers
verifiedUsers
pendingUsers
unverifiedUsers
adminUsers
```

### 2.8 举报总览

```powershell
Invoke-RestMethod `
  -Method Get `
  -Uri "$baseUrl/stats/reports/overview" `
  -Headers $adminHeaders
```

预期字段：

```text
totalReports
pendingReports
processingReports
resolvedReports
rejectedReports
closedReports
productReports
userReports
orderReports
```

## 三、权限测试

### 3.1 普通用户不能访问统计接口

```powershell
$userLogin = Invoke-RestMethod `
  -Method Post `
  -Uri "$baseUrl/auth/dev-login" `
  -ContentType "application/json" `
  -Body '{"openid":"dev_user_stats_001","role":"USER"}'

$userHeaders = @{ Authorization = "Bearer $userLogin.data.token" }

try {
  Invoke-RestMethod `
    -Method Get `
    -Uri "$baseUrl/stats/products/overview" `
    -Headers $userHeaders
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
403
```

### 3.2 未登录不能访问统计接口

```powershell
try {
  Invoke-RestMethod `
    -Method Get `
    -Uri "$baseUrl/stats/products/overview"
} catch {
  $_.Exception.Response.StatusCode.value__
}
```

预期：

```text
401
```

## 四、测试结论模板

```text
测试模块：后台统计接口

测试接口：
- GET /stats/products/overview
- GET /stats/products/category-distribution
- GET /stats/products/status-distribution
- GET /stats/products/trend
- GET /stats/orders/overview
- GET /stats/users/overview
- GET /stats/reports/overview

测试结果：
- 管理员鉴权正常
- 普通用户访问被拦截
- 未登录访问被拦截
- 商品统计正常
- 订单统计正常
- 用户统计正常
- 举报统计正常

结论：后台统计接口测试通过
```

## 五、提交建议

```bash
git status
git add backend/internal/stats backend/cmd/server/main.go backend/docs/wangshuangyuan_module_README.md backend/docs/wangshuangyuan_stats_api_test.md
git commit -m "feat(stats): 补充后台统计接口"
```
