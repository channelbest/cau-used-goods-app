# 王双媛 举报申诉状态流转与管理员后续处理开发说明

## 一、改动目标

本次调整围绕举报、申诉的处理边界，以及管理员后续处置接口展开。

核心目标：

- `CLOSED` 表示普通用户在管理员接手前主动关闭或撤回。
- 管理员处理举报或申诉时，只能给出 `APPROVED` 或 `REJECTED`。
- 举报或申诉在管理员处理完成并通知用户后结束。
- 后续商品、订单、账号状态变化，由管理员另行调用对应管理接口完成。

## 二、举报状态流转

### 2.1 状态定义

```text
PENDING      待处理
PROCESSING   管理员处理中
APPROVED     举报成立
REJECTED     举报驳回
CLOSED       用户主动关闭或撤回
```

### 2.2 流转规则

```text
用户提交举报:
PENDING

用户主动关闭:
PENDING -> CLOSED

管理员接手:
PENDING -> PROCESSING

管理员处理:
PENDING / PROCESSING -> APPROVED / REJECTED
```

说明：

- 用户只能在 `PENDING` 状态关闭举报。
- 一旦举报进入 `PROCESSING`，说明管理员已经接手，用户不能再关闭。
- 管理员不能把举报处理为 `CLOSED`。
- 管理员处理为 `APPROVED` 或 `REJECTED` 后，写入处理结果、处理人、处理时间，并通知举报人。
- 举报处理完成后，不自动修改商品、订单或账号状态。

## 三、申诉状态流转

### 3.1 状态定义

```text
PENDING      待处理
PROCESSING   管理员处理中
APPROVED     申诉通过
REJECTED     申诉驳回
CLOSED       用户主动关闭或撤回
```

### 3.2 流转规则

```text
用户提交申诉:
PENDING

用户主动关闭:
PENDING -> CLOSED

管理员接手:
PENDING -> PROCESSING

管理员处理:
PENDING / PROCESSING -> APPROVED / REJECTED
```

说明：

- 用户只能在 `PENDING` 状态关闭申诉。
- 一旦申诉进入 `PROCESSING`，用户不能再关闭。
- 管理员不能把申诉处理为 `CLOSED`。
- 管理员处理为 `APPROVED` 或 `REJECTED` 后，写入处理结果、处理人、处理时间，并通知申诉人。
- 申诉通过后，不再自动恢复商品或账号状态。

## 四、处理结束边界

举报或申诉的处理结束点为：

```text
状态更新为 APPROVED / REJECTED
写入处理结果
写入 handler_id
写入 handle_time
写入管理员操作日志
发送站内消息给举报人或申诉人
```

到这里，举报或申诉流程结束。

后续处置由管理员独立调用管理接口完成，例如：

- 商品下架、恢复、删除。
- 订单取消、异常关闭、完成。
- 账号禁用、恢复、注销、删除。

## 五、新增或调整的用户接口

### 5.1 关闭举报

```http
POST /reports/:id/close
```

权限：

```text
登录 + 学生认证
```

业务规则：

- 只能关闭自己的举报。
- 只能关闭 `PENDING` 状态的举报。
- 成功后状态变为 `CLOSED`。
- `handler_id` 清空。

请求体：

```json
{
  "closeReason": "不再需要举报"
}
```

### 5.2 关闭申诉

```http
POST /appeals/:id/close
```

权限：

```text
登录
```

业务规则：

- 只能关闭自己的申诉。
- 只能关闭 `PENDING` 状态的申诉。
- 成功后状态变为 `CLOSED`。
- `handler_id` 清空。

请求体：

```json
{
  "closeReason": "不再需要申诉"
}
```

## 六、调整后的管理员处理接口

### 6.1 处理举报

```http
POST /admin/reports/:id/handle
```

允许状态：

```text
APPROVED
REJECTED
```

不再允许：

```text
CLOSED
RESOLVED
```

请求体：

```json
{
  "status": "APPROVED",
  "handleResult": "举报成立，后续将由管理员下架商品"
}
```

说明：该接口只处理举报案件状态，不接受 `accountStatus`。如需封禁用户、下架商品或关闭订单，应调用对应后台状态接口，并传入 `relatedType=REPORT`、`relatedId=举报ID`。

### 6.2 处理申诉

```http
POST /admin/appeals/:id/handle
```

允许状态：

```text
APPROVED
REJECTED
```

不再允许：

```text
CLOSED
```

请求体：

```json
{
  "status": "APPROVED",
  "handleResult": "申诉通过，后续将由管理员恢复商品状态"
}
```

说明：该接口只处理申诉案件状态，不接受 `accountStatus`。如需恢复用户或商品状态，应调用对应后台状态接口，并传入 `relatedType=APPEAL`、`relatedId=申诉ID`。

## 七、新增管理员后续处理接口

### 7.1 修改账号状态

```http
PUT /admin/users/:id/status
```

权限：

```text
登录 + 管理员
```

支持状态：

```text
NORMAL
DISABLED
BANNED
```

请求体：

```json
{
  "accountStatus": "DISABLED",
  "reason": "举报成立，禁用账号",
  "relatedType": "REPORT",
  "relatedId": 123
}
```

说明：

- 管理员不能通过该接口修改自己的账号状态。
- 管理员账号不能通过该接口被修改状态。
- 目标用户作为卖家或买家存在 `WAIT_MEET` 订单时，禁用/封禁失败；后端会按角色返回明细提示，前端应提示管理员先人工处理待面交订单，此时不会修改用户状态、不会异常关闭订单，也不会下架商品。
- 目标用户作为卖家或买家只存在 `PENDING_CONFIRM` 订单时，禁用/封禁继续执行；系统在同一事务内自动将这些待确认订单改为 `EXCEPTION_CLOSED`，并按角色处理商品状态。
- `relatedType/relatedId` 可选；管理员自行巡查发现问题时不传，基于举报/申诉处置时传 `REPORT/APPEAL + 案件ID`。
- 用户状态变更、待确认订单异常关闭、相关商品状态处理、自动下架用户在售商品和 `admin_logs` 写入同事务；禁用/封禁成功后，只有父操作传入 `relatedType/relatedId` 时，自动关闭订单日志和自动下架商品日志才会继承同一组关联来源。
- 事务提交后发送系统消息：受影响订单买卖双方收到订单异常关闭消息，目标用户收到账号状态变更消息；消息关联具体订单或用户，不继承举报/申诉作为用户侧消息关联对象。

### 7.2 修改商品状态

```http
PUT /admin/products/:id/status
```

权限：

```text
登录 + 管理员
```

支持状态：

```text
ON_SALE
OFF_SHELF
LOCKED
SOLD
DELETED
```

请求体：

```json
{
  "status": "OFF_SHELF",
  "reason": "举报成立，下架商品",
  "relatedType": "REPORT",
  "relatedId": 123
}
```

说明：

- `DELETED` 为逻辑删除。
- 设置 `DELETED` 时，会同时设置 `is_deleted = 1`。
- 设置 `OFF_SHELF` 时，会记录 `off_shelf_reason`。
- 前端管理员入口设置 `OFF_SHELF` 时必须填写具体下架原因，不再使用“管理员下架商品”作为泛化原因。
- `relatedType/relatedId` 可选；传入时只允许 `REPORT` 或 `APPEAL`，并校验案件目标为当前商品。
- 业务变更和 `admin_logs` 写入同事务。

商品下架来源说明：

| `off_shelf_by` | 场景 | 用户是否可自行上架 |
|---|---|---|
| `USER` | 用户主动下架 | 可以 |
| `ADMIN` | 管理员直接下架 | 不可以，应走申诉或管理员后续处置 |
| `SYSTEM` | 订单异常、系统规则等普通系统下架 | 不可以 |
| `ACCOUNT_STATUS` | 账号禁用/封禁时自动下架在售商品 | 账号恢复后可以由用户手动上架 |

账号恢复为 `NORMAL` 后，不自动恢复 `ACCOUNT_STATUS` 商品；商品保留在“我的商品/下架商品”中，由用户单个上架或一键上架确认恢复公开展示。

### 7.3 修改订单状态

```http
PUT /admin/orders/:id/status
```

权限：

```text
登录 + 管理员
```

支持状态：

```text
PENDING_CONFIRM
WAIT_MEET
COMPLETED
CANCELED
```

请求体：

```json
{
  "status": "CANCELED",
  "reason": "管理员取消订单",
  "relatedType": "REPORT",
  "relatedId": 123
}
```

说明：

- 修改为 `WAIT_MEET` 时，写入 `confirm_time`。
- 修改为 `COMPLETED` 时，写入 `finish_time`，并同步商品为 `SOLD`。
- 修改为 `CANCELED` 时，写入 `cancel_reason`、`cancel_by`、`close_time`，并同步商品为 `ON_SALE`。
- 修改为 `PENDING_CONFIRM` 或 `WAIT_MEET` 时，商品同步为 `LOCKED`。
- `relatedType/relatedId` 可选；传入时只允许 `REPORT` 或 `APPEAL`，并校验案件目标为当前订单。
- 业务变更、关联商品状态变化和 `admin_logs` 写入同事务。
- `EXCEPTION_CLOSED` 不允许通过该通用状态接口设置；管理员异常关闭订单必须使用 `POST /admin/orders/:id/exception-close`，并传入责任方 `responsibleParty=BUYER/SELLER`。

## 八、超时订单清理调整

原管理员接口已移除：

```http
POST /admin/orders/cleanup-expired
```

超时未确认订单现在由后端服务启动的定时任务自动清理，每 1 分钟执行一次。清理时使用数据库时间判断 `expire_time < NOW()`，取消订单的 SQL 带 `status='PENDING_CONFIRM'` 条件，并通过 `RowsAffected=1` 判断是否由当前执行者成功处理；只有成功处理的订单才会继续解锁商品和发送通知，避免并发或多实例部署时重复处理。

## 九、涉及代码位置

```text
backend/internal/report/
backend/internal/appeal/
backend/internal/user/
backend/internal/product/
backend/internal/order/
backend/internal/stats/
backend/cmd/server/main.go
backend/scripts/sql/schema.sql
```

## 十、数据库说明

### 10.1 逻辑删除

当前相关业务对象均不做物理删除：

- 用户主动注销使用 `account_status = CANCELED` 且 `is_deleted = 1`。
- 商品 `status = DELETED` 且 `is_deleted = 1`。
- 举报、申诉、订单只通过状态变化保留追溯记录。

### 10.2 举报状态迁移

如果历史数据里存在举报状态 `RESOLVED`，需要迁移为 `APPROVED`：

```sql
UPDATE reports
SET status = 'APPROVED'
WHERE status = 'RESOLVED';
```

## 十一、验证结果

后端编译测试已通过：

```powershell
go test ./...
```
