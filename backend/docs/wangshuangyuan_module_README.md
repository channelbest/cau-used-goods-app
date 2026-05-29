# 王双媛模块开发说明文档

## 一、负责人信息

负责人：王双媛

负责内容：

- 消息通知模块
- 管理员模块
- 公告管理
- 敏感词管理
- 管理员操作日志
- 部分统计接口

当前分支：

```text
feature/backend-王双媛-admin-message
```

## 二、当前已完成的改动说明

本次分支创建前，工作区已有以下后端改动，均为本人修改内容。

### 2.1 服务初始化依赖调整

涉及文件：

```text
backend/cmd/server/main.go
```

改动内容：

- 提前初始化 `productRepo`、`productService`、`productHandler`。
- 将 `productService` 注入订单模块：

```go
orderService := order.NewService(orderRepo, productService)
```

- 将数据库连接和敏感词服务注入举报模块：

```go
reportService := report.NewService(reportRepo, db.DB(), sensitiveService)
```

改动目的：

- 订单模块需要和商品模块联动，支持下单锁定商品、取消订单恢复商品、完成订单标记售出等状态流转。
- 举报模块需要使用数据库事务能力，并复用敏感词检测能力，保证举报内容审核逻辑和商品内容审核逻辑一致。

### 2.2 评价接口路由参数修正

涉及文件：

```text
backend/internal/review/handler.go
```

改动内容：

- `ListByProduct` 从路由参数 `id` 读取商品 ID。
- `ListBySeller` 从路由参数 `id` 读取卖家 ID。

改动前：

```go
c.Param("productId")
c.Param("sellerId")
```

改动后：

```go
c.Param("id")
```

改动目的：

- 和现有路由中 `/:id` 的参数命名保持一致。
- 避免接口调用时因为参数名不匹配导致 `invalid product id` 或 `invalid seller id`。

## 三、接下来负责模块说明

王双媛接下来主要围绕 `message`、`admin`、`stats` 三类后端能力展开。

### 3.1 消息通知模块

计划目录：

```text
backend/internal/message/
  handler.go
  service.go
  repository.go
  router.go
  model.go
```

关联数据表：

```text
messages
```

计划接口：

```http
GET /messages
GET /messages/:id
GET /messages/unread-count
PUT /messages/:id/read
PUT /messages/read-all
```

权限要求：

- 必须登录。
- 用户只能查看、读取自己的消息。

核心逻辑：

- 查询当前用户消息列表。
- 查看消息详情。
- 标记单条消息已读。
- 标记全部消息已读。
- 查询未读消息数量。

### 3.2 管理员公告模块

计划归属：

```text
backend/internal/admin/
```

关联数据表：

```text
announcements
admin_logs
```

计划接口：

```http
GET /admin/announcements
POST /admin/announcements
PUT /admin/announcements/:id
PUT /admin/announcements/:id/status
DELETE /admin/announcements/:id
```

权限要求：

- 必须登录。
- 必须为管理员角色：`role=ADMIN`。
- 路由层使用 `middleware.Admin()`。
- 服务层保留管理员权限二次校验。

核心逻辑：

- 管理员创建公告。
- 管理员编辑公告。
- 管理员发布、下线公告。
- 管理员逻辑删除公告。
- 每次管理员操作写入 `admin_logs`。

### 3.3 敏感词管理模块

现有目录：

```text
backend/internal/sensitive/
```

当前状态：

- 已有 `repository.go`。
- 已有 `service.go`。
- 暂无管理员维护接口。

计划补充：

```text
backend/internal/sensitive/handler.go
backend/internal/sensitive/router.go
```

关联数据表：

```text
sensitive_words
admin_logs
```

计划接口：

```http
GET /admin/sensitive-words
POST /admin/sensitive-words
PUT /admin/sensitive-words/:id
DELETE /admin/sensitive-words/:id
```

权限要求：

- 仅管理员可操作。
- 新增、修改、删除敏感词时记录管理员操作日志。

核心逻辑：

- 敏感词列表查询。
- 新增敏感词。
- 更新敏感词状态。
- 逻辑删除或禁用敏感词。

### 3.4 管理员操作日志模块

计划归属：

```text
backend/internal/admin/
```

关联数据表：

```text
admin_logs
```

计划接口：

```http
GET /admin/logs
```

权限要求：

- 仅管理员可查看。

核心逻辑：

- 查询管理员操作记录。
- 支持按管理员、操作类型、目标类型、时间范围分页筛选。
- 为公告管理、敏感词管理等管理员操作提供统一日志写入方法。

### 3.5 部分统计接口

现有目录：

```text
backend/internal/stats/
```

当前已有接口：

```http
GET /stats/products/overview
GET /stats/products/category-distribution
GET /stats/products/status-distribution
GET /stats/products/trend
```

后续计划：

- 确认统计接口是否全部需要管理员权限。
- 如作为后台统计，应改为管理员接口或增加 `middleware.Admin()`。
- 补充订单、用户、举报、消息相关统计时，继续复用 `stats` 模块。

可能新增接口：

```http
GET /stats/orders/overview
GET /stats/users/overview
GET /stats/reports/overview
```

## 四、建议开发顺序

建议按照以下顺序开发：

1. 完成 `message` 模块目录和基础接口。
2. 完成管理员操作日志写入能力。
3. 完成公告管理接口。
4. 完成敏感词管理接口。
5. 检查并调整统计接口权限。
6. 补充接口测试和 README 说明。

优先做消息模块的原因：

- 当前仓库还没有 `internal/message`。
- 数据库中已经有 `messages` 表。
- 消息模块依赖较少，适合作为王双媛负责部分的第一个落地点。

## 五、需要注意的问题

- 不要使用 `git add .`，避免提交无关文件。
- 提交前建议只添加本次相关文件：

```bash
git add backend/docs/wangshuangyuan_module_README.md
git add backend/internal/message
git add backend/internal/admin
git add backend/internal/sensitive
git add backend/internal/stats
```

- 管理员接口必须使用 `middleware.Admin()`。
- 管理员操作必须写入 `admin_logs`。
- 消息接口必须校验当前用户身份，不能读取他人消息。
- 统计接口如果用于后台管理，需要和组内确认是否改为管理员权限。

## 六、阶段性交付物

本阶段建议交付：

- 王双媛模块说明文档。
- 消息模块接口实现。
- 管理员公告接口实现。
- 敏感词管理接口实现。
- 管理员日志记录实现。
- 统计接口权限说明或调整记录。
- Apifox 或接口测试说明。
