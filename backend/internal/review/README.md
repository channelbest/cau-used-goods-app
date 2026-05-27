# 评价模块 (Review)

## 模块职责

评价模块负责处理交易完成后的评价功能，包括提交评价、查看评价列表、计算卖家评分等。

## 关联数据表

- `reviews` — 评价表
- `orders` — 订单表（校验订单状态）
- `products` — 商品表（关联查询商品信息）

## 核心规则

1. **只有买家可以评价**
2. **只有已完成的订单才能评价**
3. **一个订单只能评价一次**
4. 评分范围 1-5 星
5. 评价内容可选

## 主要接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/reviews` | 提交评价 | 登录+认证 |
| GET | `/reviews/:id` | 评价详情 | 公开 |
| GET | `/products/:productId/reviews` | 商品评价列表 | 公开 |
| GET | `/sellers/:sellerId/reviews` | 卖家评价列表 | 公开 |

## 权限校验

- 提交评价需要 JWT 登录 + 学生认证
- 查看评价列表和评价详情是公开接口，无需登录

## 异常情况

| 场景 | 错误信息 |
|------|---------|
| 非买家评价 | "only buyer can review" |
| 订单未完成 | "can only review completed orders" |
| 重复评价 | "you have already reviewed this order" |
| 评分超出范围 | "rating must be between 1 and 5" |

## 与概要设计一致的地方

- 评价权限控制一致
- 评分范围一致
