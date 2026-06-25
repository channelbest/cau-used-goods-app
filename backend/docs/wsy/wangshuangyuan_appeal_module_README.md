# 王双媛申诉模块开发说明

## 一、模块目标

申诉模块用于处理用户对平台处置结果的复核请求。

典型场景：

- 商品被管理员下架后，卖家认为误判并提交申诉。
- 用户账号被禁用后，用户提交账号申诉。
- 订单被异常关闭后，买家或卖家提交申诉。
- 举报处理结果存在争议，用户提交申诉。

本模块第一版和第二版能力一起实现：

- 用户提交申诉。
- 用户查看自己的申诉。
- 管理员查看全部申诉。
- 管理员处理申诉。
- 申诉通过后记录处理结论；商品、账号或订单状态恢复由对应后台接口单独执行。
- 处理后写入站内消息通知申诉人。
- 处理后写入管理员操作日志。

## 二、代码位置

```text
backend/internal/appeal/
backend/cmd/server/main.go
backend/scripts/sql/schema.sql
```

新增模块文件：

```text
backend/internal/appeal/model.go
backend/internal/appeal/repository.go
backend/internal/appeal/service.go
backend/internal/appeal/handler.go
backend/internal/appeal/router.go
```

同时补充：

```text
backend/internal/admin/model.go
backend/internal/message/model.go
```

用于增加申诉相关日志类型和消息关联类型。

## 三、关联数据表

### 3.1 appeals

`appeals` 是申诉主表，保存申诉人、申诉对象、申诉理由和管理员处理结果。

核心字段：

```text
id             申诉 ID
appellant_id   申诉人用户 ID
target_type    申诉对象类型：PRODUCT / USER / ORDER / REPORT
target_id      申诉对象 ID
reason         申诉理由
status         申诉状态
handle_result  管理员处理说明
handler_id     处理管理员 ID
handle_time    处理时间
create_time    创建时间
update_time    更新时间
```

状态：

```text
PENDING      待处理
PROCESSING   处理中
APPROVED     申诉通过
REJECTED     申诉驳回
CLOSED       申诉关闭
```

### 3.2 appeal_images

`appeal_images` 是申诉凭证图片表，用于保存申诉关联的图片路径。

核心字段：

```text
id          图片 ID
appeal_id   申诉 ID
image_url   图片路径
sort_order  排序值
create_time 上传时间
```

## 四、接口列表

### 4.1 用户提交申诉

```http
POST /appeals
```

权限：

```text
登录用户
```

说明：用户端申诉不强制学生认证，因为账号被禁用或认证状态异常时也需要允许用户申诉。

请求体：

```json
{
  "targetType": "PRODUCT",
  "targetId": 12,
  "reason": "商品被误判违规，实际是普通教材",
  "evidenceUrls": [
    "/uploads/appeal/demo-1.jpg"
  ]
}
```

业务规则：

- `targetType` 只能是 `PRODUCT`、`USER`、`ORDER`、`REPORT`。
- `targetId` 必须存在。
- 只能申诉和自己有关的对象：自己的商品、自己的账号、自己参与的订单、自己提交的举报。
- `reason` 不能为空，最多 500 字。
- `evidenceUrls` 最多 9 张。
- 同一用户对同一对象存在 `PENDING` 或 `PROCESSING` 申诉时，不能重复提交。

返回示例：

```json
{
  "id": 1,
  "appellantId": 8,
  "targetType": "PRODUCT",
  "targetId": 12,
  "reason": "商品被误判违规，实际是普通教材",
  "status": "PENDING",
  "createTime": "2026-06-03 16:40:00",
  "updateTime": "2026-06-03 16:40:00",
  "images": [
    "/uploads/appeal/demo-1.jpg"
  ]
}
```

### 4.2 用户查看自己的申诉列表

```http
GET /appeals/my?page=1&pageSize=20
```

支持筛选：

```text
targetType
targetId
status
page
pageSize
```

示例：

```http
GET /appeals/my?status=PENDING&page=1&pageSize=20
```

### 4.3 用户查看申诉详情

```http
GET /appeals/:id
```

权限规则：

```text
只能查看自己的申诉
```

非本人访问返回：

```text
403 permission denied
```

### 4.4 管理员查看申诉列表

```http
GET /admin/appeals?page=1&pageSize=20
```

权限：

```text
登录 + role=ADMIN
```

支持筛选：

```text
targetType
targetId
status
page
pageSize
```

### 4.5 管理员查看申诉详情

```http
GET /admin/appeals/:id
```

权限：

```text
登录 + role=ADMIN
```

### 4.6 管理员处理申诉

```http
POST /admin/appeals/:id/handle
```

请求体：

```json
{
  "status": "APPROVED",
  "handleResult": "申诉通过，商品恢复上架"
}
```

可处理状态：

```text
APPROVED
REJECTED
```

业务规则：

- 只有 `PENDING` 或 `PROCESSING` 的申诉可以处理。
- 处理后写入 `handle_result`、`handler_id`、`handle_time`。
- 处理后写入 `admin_logs`。
- 处理后通过 `messages` 通知申诉人。
- `APPROVED` 只表示申诉成立，不自动恢复商品、账号或订单状态。

## 五、申诉通过后的后续处置

### 5.1 商品申诉通过

```text
target_type = PRODUCT
status = APPROVED
```

后续处置：

```http
PUT /admin/products/:id/status
```

说明：

- 管理员确认申诉成立后，可通过商品状态接口恢复上架。
- 请求体可传 `relatedType=APPEAL`、`relatedId=申诉ID`，让商品状态日志关联该申诉。
- 已删除商品是否恢复仍由商品状态接口规则决定。

商品申诉入口补充：

- 卖家本人打开已下架商品详情时，如果该商品不能自行恢复上架，会显示“申诉下架”入口。
- 不能自行恢复的典型来源包括 `off_shelf_by = ADMIN` 和普通 `SYSTEM`。
- `off_shelf_by = USER` 或 `ACCOUNT_STATUS` 的商品不走该入口，用户可在“我的商品”中单个上架或一键上架。
- 从商品详情进入申诉页时，前端固定传入：

```text
targetType = PRODUCT
targetId = 当前商品 ID
lockTarget = 1
```

- 申诉页收到 `lockTarget=1` 后，申诉对象和对象 ID 不允许修改，避免用户将入口申诉改成其他目标。

### 5.2 账号申诉通过

```text
target_type = USER
status = APPROVED
```

后续处置：

```http
PUT /admin/users/:id/status
```

说明：

- 管理员确认申诉成立后，可通过用户状态接口恢复账号。
- 撤销永久封禁仍必须由超级管理员执行。
- 请求体可传 `relatedType=APPEAL`、`relatedId=申诉ID`，让用户状态日志关联该申诉。

### 5.3 订单申诉通过

```text
target_type = ORDER
status = APPROVED
```

当前只记录处理结果，不自动恢复订单。

原因：

- 订单状态通常涉及买卖双方和商品锁定状态。
- 自动恢复订单容易造成交易状态不一致。
- 如需调整订单，应通过订单后台状态接口或异常关闭接口处理，并传 `relatedType=APPEAL`、`relatedId=申诉ID`。

### 5.4 举报申诉通过

```text
target_type = REPORT
status = APPROVED
```

当前只记录处理结果，不自动修改原举报记录。

后续如需要，可以扩展为补充举报复核结果字段。

## 六、管理员日志和站内消息

### 6.1 管理员日志

申诉处理后写入：

```text
admin_logs
```

使用：

```text
operation_type = MARK_APPEAL_PROCESSING / APPROVE_APPEAL / REJECT_APPEAL
target_type = APPEAL
target_id = appeal.id
```

### 6.2 站内消息

申诉处理后写入：

```text
messages
```

使用：

```text
message_type = SYSTEM_NOTICE
related_type = APPEAL
related_id = appeal.id
```

申诉人可以在消息通知模块看到处理结果。

## 七、与 report 模块的区别

```text
report = 用户举报别人或对象
appeal = 用户对平台处理结果提出复核
```

两者方向不同，不建议合并。

## 八、后续扩展建议

后续可以扩展：

- 管理员将申诉状态改为 `PROCESSING`。
- 申诉通过后自动修改举报处理结果。
- 订单申诉通过后创建补偿通知或人工处理单。
- 管理员补充多次处理记录。
- 申诉撤回。
- 申诉超时提醒。
