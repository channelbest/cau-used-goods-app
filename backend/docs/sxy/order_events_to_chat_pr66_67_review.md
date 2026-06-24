# PR 66/67 订单消息迁移到聊天会话审查记录

## 背景

- 当前迁移目标：将订单类消息从系统消息中心迁移到商品相关聊天会话中。
- 当前本地分支：`feature/sxy-test`。
- 备份/修复来源分支：`feature/backend-sxy-test`。
- 涉及 PR：
  - PR #66 `fix(backend): sync order events to chat flow`
  - PR #67 `fix(frontend): update chat product and order detail pages`
- 数据库已执行结构迁移：`chat_messages.sender_id/receiver_id` 改为可空，新增 `actor_type/order_id/event_type` 和 `idx_chat_messages_order/fk_chat_messages_order`。

## 当前本地状态

- 本地 `feature/sxy-test` 当前 HEAD 为 PR #66 的 merge commit：`18c0286a`。
- 远端 `origin/dev` 已前进到 PR #67 的 merge commit：`fef218e5`。
- 如果要让当前工作区包含 PR #67 前端改动，还需要再合入 `origin/dev`。

## PR #66 后端可靠性

结论：方向正确，可以作为迁移基础，但不是完全闭环。

已完成：

- 订单创建、确认、取消、完成、超时取消、异常关闭、管理员改状态会调用 `chat.RecordOrderEvent`。
- `RecordOrderEvent` 会通过 `EnsureConversation` 自动创建/复用会话，解决“直接预约没有聊天会话”的核心问题。
- `chat_messages` 支持 `ORDER_EVENT`，并记录 `actor_type/order_id/event_type`。
- `go test ./...` 在 backend 通过。

风险：

- 订单事件写入聊天不在订单主事务内，失败只打日志，不回滚订单操作。可能出现订单状态已变更但聊天事件缺失。
- `migrate_order_events_to_chat.sql` 实际只做表结构迁移，没有把历史 `messages` 表中的订单消息迁移到 `chat_messages`。
- 系统消息中心仍保留订单类消息展示逻辑。新订单事件不会再进入系统消息，但旧订单消息仍会显示，迁移期会出现新旧入口并存。
- 系统事件 `actor_type=SYSTEM` 设置 `receiver_id=NULL` 且 `read_status=READ`，同时增加双方会话未读数。消息本身不可被 `MarkRead` 按 `receiver_id` 标记已读，只能依赖会话未读数清零。当前能工作，但语义不一致，后续统计单条未读时可能踩坑。

建议修复：

- 增加补偿脚本或后台任务，扫描订单状态与 `chat_messages.order_id/event_type`，补齐缺失事件。
- 如果要迁移历史消息，单独写历史数据迁移脚本：按 `messages.related_type='ORDER'` 和订单买卖双方创建/复用会话，再插入 `ORDER_EVENT`。
- 明确是否保留旧订单系统消息。若完全迁移，应从系统消息筛选中移除订单类 `ORDER_CREATED/ORDER_CONFIRMED/ORDER_CANCELED/ORDER_TIMEOUT/ORDER_EXCEPTION_CLOSED`。
- 考虑给系统事件插入两条定向消息，或调整 `MarkRead` 支持 `actor_type=SYSTEM AND receiver_id IS NULL`，避免未读语义不一致。

## PR #67 前端可靠性

结论：能展示订单事件，但仍有产品流程问题。

已完成：

- 聊天页识别 `messageType=ORDER_EVENT`，展示订单事件卡片。
- 点击订单事件可跳转 `/pages/order/detail?id=...`。
- 订单详情新增“联系对方”，买家可 `createOrGetConversation`，卖家通过列表查找已有会话。
- 商品详情在 `LOCKED` 且当前用户是预约买家时允许聊天。

风险：

- 商品详情“提交预约”仍直接跳转预约页，不是先进入聊天页。后端会自动创建会话，但交互上还不是闲鱼式“先聊再交易”。
- 卖家在订单详情中查找会话依赖分页遍历 `listConversations`，最多 20 页，每页 100 条。数据量大时仍可能找不到。
- 聊天页订单事件跳订单详情时没有带 `fromMessage=1` 或 `conversationId`，返回路径不一定回到聊天会话。
- 系统消息页仍显示订单类系统消息，和聊天订单事件迁移目标不完全一致。

建议修复：

- 如果目标是贴近闲鱼，商品详情主按钮应优先进入聊天，预约动作放到聊天页或会话内订单卡片中触发。
- 后端增加按 `order_id` 查询会话接口，例如 `GET /chat/orders/:orderId/conversation`，替代前端分页查找。
- 聊天页打开订单详情时附带来源参数，例如 `fromChat=1&conversationId=...`，订单详情提供“返回会话”。
- 从系统消息中心移除订单类消息入口，或在历史迁移完成前标注为历史订单通知。

## 出问题时的排查与修复

### 1. 直接预约后没有聊天会话

排查 SQL：

```sql
SELECT id, product_id, buyer_id, seller_id
FROM chat_conversations
WHERE product_id = ? AND buyer_id = ? AND seller_id = ?;
```

修复：

- 确认后端已部署 PR #66。
- 确认订单服务初始化传入的是 `chatService`。
- 若会话缺失，补调用 `chat.RecordOrderEvent` 或执行补偿脚本创建会话并插入 `ORDER_CREATED`。

### 2. 会话存在但没有订单事件

排查 SQL：

```sql
SELECT id, conversation_id, order_id, event_type, actor_type, content, create_time
FROM chat_messages
WHERE order_id = ?
ORDER BY create_time;
```

修复：

- 查看后端日志中的 `record user order event failed` 或 `record system order event failed`。
- 检查 `chat_messages.order_id` 外键是否存在、订单 ID 是否正确。
- 补插对应 `ORDER_EVENT` 消息，并更新 `chat_conversations.last_message_id/last_message_content/last_message_time`。

### 3. 订单系统消息和聊天订单事件重复

原因：

- 新逻辑写聊天事件，旧数据仍在 `messages` 表。
- 系统消息页仍包含订单类 message type。

修复：

- 确认是否迁移历史数据。
- 迁移后可将旧订单系统消息标记已读/删除，或从前端系统消息类型中过滤掉订单类消息。

### 4. 卖家点“联系对方”找不到会话

原因：

- 前端通过分页查 `listConversations` 找会话，可能翻页上限内没查到。

修复：

- 短期：提高分页上限或在订单详情传入 `conversationId`。
- 长期：增加后端按订单 ID 查会话接口。

### 5. 系统事件未读数异常

原因：

- 系统事件 `receiver_id=NULL`，但会话双方未读数都增加。

修复：

- 若只依赖会话未读数，当前 `MarkRead` 清零可接受。
- 若要统计单条未读消息，应改为双定向系统消息，或扩展 `MarkRead` 处理 `receiver_id IS NULL` 的系统事件。
