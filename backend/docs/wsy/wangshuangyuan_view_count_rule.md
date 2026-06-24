# 商品浏览量统计规则调整说明


原有逻辑是在调用公开商品详情接口 `GET /products/:id` 时直接执行 `view_count + 1`，没有区分访问者身份。

因此以下访问都会增加浏览量：

- 买家查看商品详情；
- 卖家查看自己的商品详情；
- 管理员在后台页面复用公开商品详情接口；
- 同一用户短时间内重复刷新详情页。

这会导致浏览量混入管理行为、卖家自查行为和重复刷新数据，不够接近真实商品曝光。

## 调整后的规则

商品浏览量只在满足以下条件时增加：

- 商品存在；
- 商品未删除；
- 商品状态为 `ON_SALE`；
- 当前访问者已登录；
- 当前访问者不是商品卖家本人；
- 当前访问者不是 `ADMIN` 或 `SUPER_ADMIN`；
- 当前访问者 1 小时内没有浏览过同一个商品。

未经过学生认证的普通用户仍然可以计入浏览量，因为这类用户在当前业务中相当于游客浏览行为。

## 不增加浏览量的情况

- 未登录请求；
- 卖家本人查看自己的商品；
- 管理员或超级管理员查看商品；
- 商品已下架、已锁定、已售出或已删除；
- 同一用户 1 小时内重复打开同一商品详情。

## 实现说明

后端新增了可选登录识别中间件：

- 没有 token：不阻断公开接口访问；
- token 无效：不阻断公开接口访问；
- token 有效：写入当前用户 ID 和角色，供商品详情接口判断浏览量是否应增加。

商品详情接口 `GET /products/:id` 接入可选登录识别后，服务层将访问者信息传入浏览量统计逻辑。

公开商品详情接口的可见性规则同步调整为：

- 未登录用户和普通买家只能查看 `ON_SALE` 在售商品；
- 商品卖家本人可通过 `GET /products/:id` 查看自己未删除的非在售商品，例如 `OFF_SHELF` 下架商品，用于消息通知、我的商品列表等自查场景；
- 订单买卖双方可通过 `GET /products/:id` 查看其订单关联的 `PENDING_CONFIRM`、`WAIT_MEET`、`COMPLETED` 商品详情，即使商品状态已变为 `LOCKED` 或 `SOLD`，用于订单详情、交易沟通和售后追溯；
- 卖家本人查看自己的非在售商品不增加 `view_count`；
- 订单买卖双方查看非在售订单商品不增加 `view_count`，浏览量递增仍只作用于 `ON_SALE` 商品；
- 管理员仍建议使用后台详情接口 `GET /admin/products/:id` 查看管理对象。

管理员后台不再依赖公开商品详情接口查看管理对象。后台商品状态页应使用：

```text
GET /admin/products/:id
```

该接口用于管理场景，可返回下架、锁定、已售出或已删除商品，不会触发浏览量递增。

浏览量更新使用 `browse_history` 表做 1 小时去重：

```sql
UPDATE products
SET view_count = view_count + 1,
    update_time = CURRENT_TIMESTAMP
WHERE id = ?
  AND is_deleted = 0
  AND status = 'ON_SALE'
  AND seller_id <> ?
  AND NOT EXISTS (
      SELECT 1
      FROM browse_history
      WHERE user_id = ?
        AND product_id = ?
        AND create_time >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
  )
```

如果本次成功增加浏览量，则同步写入一条 `browse_history` 记录。

## 前端配合

普通商品详情页仍然访问公开接口：

```text
GET /products/:id
```

但前端不再强制 `auth: false`。这样用户已登录时会自动携带 token，后端才能识别：

- 是否卖家本人；
- 是否管理员；
- 是否需要按 1 小时规则去重。

未登录时仍可请求商品详情，但不会计入浏览量。

管理员商品状态页访问后台接口：

```text
GET /admin/products/:id
```

后台接口不参与浏览量统计，也不受公开详情 `status='ON_SALE'` 可见性规则限制。

## 建议验证场景

本次按要求未运行代码，后续可按以下场景手动验证：

1. 买家第一次打开在售商品详情，`view_count + 1`。
2. 同一买家 1 小时内再次打开同一商品，`view_count` 不变。
3. 未认证普通用户打开在售商品详情，`view_count + 1`。
4. 卖家打开自己的商品详情，`view_count` 不变。
5. 卖家从商品下架通知或我的商品列表打开自己已下架商品详情，可返回商品详情，`view_count` 不变。
6. 管理员通过 `GET /admin/products/:id` 打开商品详情，`view_count` 不变。
7. 普通用户访问下架、锁定、已售出或删除商品，公开详情接口仍不可用，`view_count` 不变。
