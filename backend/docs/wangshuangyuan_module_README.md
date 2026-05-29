# 王双媛模块开发说明文档

## 一、负责信息

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

## 三、当前已实现模块

### 3.1 消息通知模块

代码目录：

```text
backend/internal/message/
  handler.go
  service.go
  repository.go
  router.go
  model.go
```

入口注册：

```text
backend/cmd/server/main.go
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

权限要求：

- 必须登录。
- 用户只能查看、读取自己的消息。
- 查询他人消息或标记他人消息已读时，统一返回 `message not found`。

核心逻辑：

- 查询当前用户消息列表。
- 支持按 `READ`、`UNREAD` 筛选消息。
- 查询当前用户未读消息数量。
- 查询当前用户自己的消息详情。
- 标记单条消息已读。
- 标记当前用户全部未读消息为已读。
- 提供 `message.Service.Create(ctx, input)` 内部方法，供订单、举报、公告等后端模块创建站内消息。

内部创建能力：

```go
messageID, err := messageService.Create(ctx, message.CreateMessageInput{
    ReceiverID:  buyerID,
    SenderID:    nil,
    MessageType: message.MessageTypeOrderConfirmed,
    Title:       "订单已确认",
    Content:     "卖家已确认你的预约",
    RelatedType: ptrString(message.RelatedTypeOrder),
    RelatedID:   &orderID,
})
```

说明：

- `Create` 不是前端接口，不在 `router.go` 中暴露。
- 创建出的消息默认 `read_status = UNREAD`。
- `receiver_id`、`message_type`、`title`、`content` 必填。
- 标题长度不超过 100，内容长度不超过 500，和数据库表结构保持一致。
- 后续订单、举报、公告模块只需要注入 `message.Service` 即可复用消息写入能力。

## 四、接下来负责模块说明

王双媛后续主要围绕 `admin`、`sensitive`、`stats` 三类后端能力展开。

### 4.1 管理员公告模块

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

### 4.2 敏感词管理模块

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

### 4.3 管理员操作日志模块

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

### 4.4 部分统计接口

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

1. 完成消息模块接口测试并提交。
2. 完成管理员操作日志写入能力。
3. 完成公告管理接口。
4. 完成敏感词管理接口。
5. 检查并调整统计接口权限。
6. 补充接口测试和 README 说明。

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
- 管理员公告接口实现。
- 敏感词管理接口实现。
- 管理员日志记录实现。
- 统计接口权限说明或调整记录。
