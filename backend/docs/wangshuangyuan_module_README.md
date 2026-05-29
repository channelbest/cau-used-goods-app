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

本分支创建前，工作区已有以下后端改动，均为本人修改内容。

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

改动目的：

- 和现有路由中 `/:id` 的参数命名保持一致。
- 避免接口调用时因为参数名不匹配导致 `invalid product id` 或 `invalid seller id`。

## 三、当前已实现模块

### 3.1 消息通知模块

代码目录：

```text
backend/internal/message/
```

关联数据表：

```text
messages
```

已实现接口：

```http
GET /messages
GET /messages?readStatus=UNREAD
GET /messages?readStatus=READ
GET /messages/unread-count
GET /messages/:id
PUT /messages/:id/read
PUT /messages/read-all
```

已实现内部能力：

```go
message.Service.Create(ctx, input)
```

说明：

- 用户只能查看、读取自己的消息。
- 查询他人消息或标记他人消息已读时，统一返回 `message not found`。
- `Create` 不暴露为前端接口，供订单、举报、公告等后端模块创建站内消息。

### 3.2 管理员操作日志模块

代码目录：

```text
backend/internal/admin/
```

关联数据表：

```text
admin_logs
```

已实现接口：

```http
GET /admin/logs
```

支持筛选参数：

```text
adminId
operationType
targetType
targetId
startTime
endTime
page
pageSize
```

已实现内部能力：

```go
admin.Service.LogAction(ctx, input)
```

说明：

- 管理员日志接口必须登录且 `role=ADMIN`。
- `LogAction` 供公告管理、敏感词管理等模块记录管理员操作。

### 3.3 管理员公告管理模块

代码目录：

```text
backend/internal/admin/
```

关联数据表：

```text
announcements
admin_logs
```

已实现接口：

```http
GET /admin/announcements
POST /admin/announcements
PUT /admin/announcements/:id
PUT /admin/announcements/:id/status
DELETE /admin/announcements/:id
```

支持状态：

```text
DRAFT
PUBLISHED
OFFLINE
```

说明：

- 公告管理接口必须登录且 `role=ADMIN`。
- 支持按状态、关键字分页查询公告。
- 支持创建、编辑、发布、下线公告。
- 由于 `announcements` 表没有 `is_deleted` 字段，`DELETE` 当前按下线处理，即设置 `status = OFFLINE`。
- 新增、编辑、改状态、删除公告时写入 `admin_logs`。

下线和删除的前端理解：

- 下线接口：`PUT /admin/announcements/:id/status`，请求 `status = OFFLINE`。
- 删除接口：`DELETE /admin/announcements/:id`。
- 当前两者都会让公告变为 `OFFLINE`，即不再作为有效公告展示。
- 区别在业务语义和日志记录：下线记录 `STATUS_NOTICE`，删除记录 `DELETE_NOTICE`。
- 如果后续数据库增加 `is_deleted`、`deleted_time` 字段，删除接口可调整为真正逻辑删除。

### 3.4 管理员敏感词管理模块

代码目录：

```text
backend/internal/sensitive/
```

关联数据表：

```text
sensitive_words
admin_logs
```

已实现接口：

```http
GET /admin/sensitive-words
POST /admin/sensitive-words
PUT /admin/sensitive-words/:id
DELETE /admin/sensitive-words/:id
```

支持类型：

```text
FORBIDDEN
RISK
```

支持状态：

```text
ENABLED
DISABLED
```

说明：

- 敏感词管理接口必须登录且 `role=ADMIN`。
- 支持按状态、类型、关键字分页查询敏感词。
- 支持新增、更新、禁用敏感词。
- 由于 `sensitive_words` 表没有 `is_deleted` 字段，`DELETE` 当前按禁用处理，即设置 `status = DISABLED`。
- 新增、更新、禁用敏感词时写入 `admin_logs`。
- 原有 `CheckText` 检测能力保留，仍供商品、举报等业务模块使用。

## 四、后续负责模块

### 4.1 部分统计接口

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

## 五、建议开发顺序

建议按照以下顺序继续开发：

1. 完成消息模块、管理员日志模块、公告管理模块、敏感词管理模块接口测试并提交。
2. 检查并调整统计接口权限。
3. 补充接口测试和 README 说明。

## 六、需要注意的问题

- 不要使用 `git add .`，避免提交无关文件。
- 管理员接口必须使用 `middleware.Admin()`。
- 管理员操作必须写入 `admin_logs`。
- 消息接口必须校验当前用户身份，不能读取他人消息。
- 统计接口如果用于后台管理，需要和组内确认是否改为管理员权限。

## 七、阶段性交付物

本阶段建议交付：

- 王双媛模块说明文档。
- 消息模块接口实现。
- 消息模块接口测试说明。
- 管理员日志接口实现。
- 管理员日志接口测试说明。
- 管理员公告接口实现。
- 管理员公告接口测试说明。
- 敏感词管理接口实现。
- 敏感词管理接口测试说明。
- 统计接口权限说明或调整记录。
