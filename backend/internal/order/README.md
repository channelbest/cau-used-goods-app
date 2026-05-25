# 订单模块 (Order)

## 模块职责

订单模块负责处理二手交易的核心交易流程，包括订单创建、卖家确认、订单取消、订单完成等操作，以及订单状态流转和并发控制。

## 关联数据表

- `orders` — 订单主表
- `products` — 商品表（订单操作会修改商品状态）

## 订单状态流转

```
PENDING_CONFIRM（待确认）
    ├── 卖家确认 → WAIT_MEET（待见面）
    ├── 买家/卖家取消 → CANCELED（已取消）
    └── 超时未确认 → CANCELED（系统自动取消）

WAIT_MEET（待见面）
    ├── 卖家完成 → COMPLETED（已完成）
    └── 买家/卖家取消 → CANCELED（已取消）
```

## 核心规则

1. **同一商品同一时间只能有一个有效锁定订单**
2. **创建订单与锁定商品必须在同一个 MySQL 事务中完成**
3. **取消订单与恢复商品必须在同一个 MySQL 事务中完成**
4. **完成订单与标记商品已售必须在同一个 MySQL 事务中完成**
5. **订单创建后 24 小时内卖家未确认，系统自动取消**

## 主要接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/orders` | 创建订单 | 登录+认证 |
| GET | `/orders` | 我的订单列表 | 登录+认证 |
| GET | `/orders/:id` | 订单详情 | 登录+认证（买卖双方）|
| POST | `/orders/:id/confirm` | 卖家确认订单 | 登录+认证（卖家）|
| POST | `/orders/:id/cancel` | 取消订单 | 登录+认证（买卖双方）|
| POST | `/orders/:id/complete` | 完成订单 | 登录+认证（卖家）|
| POST | `/admin/orders/cleanup-expired` | 清理超时订单 | 登录+认证 |

## 权限校验

- 所有接口需要 JWT 登录
- 所有接口需要学生认证（`auth_status=VERIFIED`）
- 订单详情、确认、取消、完成接口校验操作人身份

## 异常情况

| 场景 | 错误信息 |
|------|---------|
| 购买自己的商品 | "cannot buy your own product" |
| 重复预约同一商品 | "you already have an active order for this product" |
| 商品已被锁定/售出 | "product not available" |
| 非卖家确认订单 | "permission denied" |
| 订单状态不允许操作 | "order cannot be confirmed/cancelled/completed" |
| 订单不存在 | "order not found" |

## 并发控制

通过数据库条件更新实现：
```sql
UPDATE products SET status = 'LOCKED' WHERE id = ? AND status = 'ON_SALE'
```

如果 `RowsAffected = 0`，说明商品已被其他订单锁定，返回错误。

## 与概要设计一致的地方

- 订单状态定义一致
- 事务处理要求一致
- 并发控制方案一致

## 仍需讨论的问题

- 超时自动取消目前是手动触发接口，是否需要定时任务？
- 订单完成后是否需要自动发送消息通知？
