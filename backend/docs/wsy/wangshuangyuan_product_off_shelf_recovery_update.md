# 商品下架来源、恢复上架与申诉入口补充说明

## 一、背景

为区分不同下架来源，商品表使用：

| 字段 | 说明 |
|---|---|
| `off_shelf_by` | 下架来源 |
| `off_shelf_reason` | 下架原因 |

本次补充 `ACCOUNT_STATUS` 来源，用于表示“因用户账号禁用或封禁而自动下架的在售商品”。

## 二、下架来源枚举

| `off_shelf_by` | 触发场景 | 用户是否可自行上架 |
|---|---|---|
| `USER` | 用户主动下架 | 可以 |
| `ADMIN` | 管理员直接下架商品 | 不可以 |
| `SYSTEM` | 订单异常、系统规则等普通系统下架 | 不可以 |
| `ACCOUNT_STATUS` | 账号被禁用/封禁时，系统自动下架该用户在售商品 | 账号恢复后可以 |

## 三、账号恢复后的商品处理

账号从 `DISABLED` 或 `BANNED` 恢复为 `NORMAL` 后，不自动恢复商品上架。

处理规则：

- `ACCOUNT_STATUS` 商品继续停留在 `OFF_SHELF`。
- 用户登录后可在“我的商品”中看到下架原因。
- 用户可单个上架，也可使用“一键上架”。
- 上架前前端弹窗提示：商品上架后将重新公开展示，请确认商品无违规、物品仍可正常交易。
- 用户上架成功后，后端清空 `off_shelf_by` 和 `off_shelf_reason`。

一键上架接口：

```http
POST /products/batch-on-sale
```

只处理当前用户自己的：

```text
status = OFF_SHELF
off_shelf_by IN (USER, ACCOUNT_STATUS)
is_deleted = 0
```

## 四、不能自行上架的商品申诉入口

卖家本人打开商品详情页时，如果商品满足：

```text
status = OFF_SHELF
off_shelf_by NOT IN (USER, ACCOUNT_STATUS)
```

前端显示“申诉下架”入口。

跳转申诉页时固定：

```text
targetType = PRODUCT
targetId = 当前商品 ID
lockTarget = 1
```

申诉页收到 `lockTarget=1` 后，申诉对象和对象 ID 不允许修改。

## 五、管理员下架原因

管理员下架商品时，前端必须填写具体下架原因。

影响范围：

- 后台商品状态页。
- 商品详情页管理员操作栏。

原因会作为 `reason` 提交到：

```http
PUT /admin/products/:id/status
```

并写入：

- `products.off_shelf_reason`
- `admin_logs.description`
- 商品下架系统通知内容

当前后端仍建议继续增加 `OFF_SHELF` 时 `reason` 必填校验，避免绕过前端直接调用接口。

## 六、历史数据修正

旧版本中，账号状态导致自动下架的商品可能写为：

```text
off_shelf_by = SYSTEM
```

已新增脚本：

```text
backend/scripts/sql/migrate_account_status_off_shelf_source.sql
```

该脚本根据 `admin_logs.description LIKE 'off shelf product due account status%'` 将可识别的历史数据修正为：

```text
off_shelf_by = ACCOUNT_STATUS
```
