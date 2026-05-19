# 后端协作 README

本目录用于 CAU 二手交易平台后端开发。当前阶段先完成后端模块拆分、数据库与业务流程复核、各模块设计文档整理，最后统一确认整体后端方案。

## 技术选型

以后端概要设计为准，推荐采用：

- 语言与框架：Go + Gin
- 数据库：MySQL
- ORM/数据访问：GORM 或 sqlx，组内统一后再使用
- 登录态：JWT/token
- 文件存储：初版使用本地文件存储，数据库只保存图片 URL 或文件路径
- 可选能力：Redis、Nginx、微信订阅消息、AI 服务
- 接口风格：RESTful API
- 接口测试：Apifox

## 当前任务

1. 按大模块分配后端任务
   - 登录认证与用户模块
   - 商品与分类模块
   - 图片上传与文件模块
   - 订单预约与状态流转模块
   - 收藏、评价、举报模块
   - 消息通知模块
   - 管理员、敏感词、公告、日志模块
   - AI 辅助发布与统计扩展模块

2. 共同阅读设计文档
   - `doc/CAU二手交易平台_软件需求分析.docx`
   - `doc/CAU二手交易平台概要设计报告V2.docx`
   - 重点核对数据库表、状态字段、权限规则、接口模块、核心业务流程。

3. 复核数据库与业务流程
   - 确认表结构是否能覆盖 P0/P1 功能。
   - 确认商品和订单状态流转是否完整。
   - 确认创建订单、取消订单、完成订单是否使用事务。
   - 发现问题时在各自模块文档中标注，不要直接改口径。

4. 整理各自模块设计文档
   - 模块职责
   - 关联数据表
   - 主要接口
   - 权限校验
   - 状态流转
   - 异常情况
   - 与概要设计不一致或有疑问的地方

5. 统一确认整体后端设计方案
   - 确认接口路径和返回格式。
   - 确认数据库字段和枚举值。
   - 确认错误码和权限策略。
   - 确认 P0 先做范围，P1/P2/扩展按进度排期。

## 模块分工建议

- 成员 A：登录认证、用户、学生认证、JWT、中间件。
- 成员 B：商品、分类、商品图片、文件上传、敏感词检测。
- 成员 C：订单预约、订单确认/取消/完成、状态流转、事务处理。
- 成员 D：收藏、评价、举报、消息通知。
- 成员 E：管理员、公告、敏感词管理、操作日志、统计扩展。

如果人数较少，可以合并为：

- 核心账户与权限组：登录、认证、用户、管理员权限。
- 核心交易组：商品、订单、评价、举报、消息。
- 基础设施组：数据库、文件上传、统一响应、错误码、接口文档。

## 后端核心规则

- 系统不做线上支付、物流、退款、资金托管。
- 用户必须完成微信登录，并通过学生认证后才能发布、预约、评价、举报。
- `auth_status=VERIFIED` 且 `account_status=NORMAL` 才允许核心交易操作。
- 商品状态至少包括：`ON_SALE`、`LOCKED`、`SOLD`、`OFF_SHELF`、`DELETED`。
- 订单状态至少包括：`PENDING_CONFIRM`、`WAIT_MEET`、`COMPLETED`、`CANCELED`、`EXCEPTION_CLOSED`。
- 同一商品同一时间只能有一个有效锁定订单。
- 创建订单与锁定商品必须在同一个 MySQL 事务中完成。
- 取消订单与恢复商品必须在同一个 MySQL 事务中完成。
- 完成订单与标记商品已售必须在同一个 MySQL 事务中完成。
- 商品、订单、举报、日志等追溯数据不要物理删除。
- 管理员操作必须写入 `admin_logs`。
- 商品详情页不返回明文联系方式；订单进入 `WAIT_MEET` 后仅交易双方可见。

## 推荐接口模块

后续可以按以下路由组组织，具体路径由后端组统一：

- `/auth`：微信登录、学生认证、认证状态查询。
- `/users`：个人信息、账号状态、隐私信息。
- `/categories`：分类列表、分类状态。
- `/products`：商品列表、搜索筛选、详情、发布、编辑、上下架。
- `/upload`：商品图片、举报凭证上传。
- `/orders`：提交预约、我的订单、订单详情、确认、取消、完成。
- `/favorites`：收藏、取消收藏、收藏列表。
- `/reviews`：提交评价、查看评价、评价状态处理。
- `/reports`：提交举报、我的举报、举报详情。
- `/messages`：消息列表、消息详情、标记已读。
- `/admin`：用户管理、商品管理、订单查看、举报处理、公告、敏感词、日志、统计。
- `/ai`：AI 标题优化、AI 描述生成。

统一响应结构：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": "YYYY-MM-DD HH:mm:ss",
  "requestId": "optional-request-id"
}
```

## 数据库复核重点

优先确认以下表：

- `users`
- `categories`
- `products`
- `product_images`
- `orders`
- `favorites`
- `reviews`
- `reports`
- `report_images`
- `messages`
- `sensitive_words`
- `admin_logs`
- `announcements`
- `browse_history`
- `ai_generation_logs`

重点检查：

- 主键、外键、唯一约束是否合理。
- `seller_id`、`buyer_id`、`reporter_id`、`handler_id` 等权限字段是否完整。
- 状态字段是否覆盖业务流程。
- 图片是否单独建表，数据库是否只保存访问路径。
- 订单并发预约是否能通过事务和条件更新保证一致性。

## Git 协作命令

本项目采用简化分支模型：

- `main`：稳定分支，由项目组长维护。
- `dev`：日常开发分支，前后端组长审核后合并到这里。
- `feature/*`：个人任务分支，组员在自己的功能分支上开发。

首次拉取：

```bash
git clone <仓库地址>
cd cau-used-goods-app
```

每天开始开发前：

```bash
git checkout dev
git pull origin dev
git status
```

新建个人分支：

```bash
git checkout -b feature/backend-姓名-模块名
```

提交修改：

```bash
git status
git add backend
git commit -m "feat(backend): 完成某某模块设计或接口"
```

推送分支：

```bash
git push origin feature/backend-姓名-模块名
```

发起 PR/MR：

- 源分支：`feature/backend-姓名-模块名`
- 目标分支：`dev`
- 审核人：后端组长

合并前同步最新 `dev`：

```bash
git checkout dev
git pull origin dev
git checkout feature/backend-姓名-模块名
git merge dev
```

冲突处理完成后：

```bash
git add .
git commit -m "chore(backend): resolve merge conflicts"
git push origin feature/backend-姓名-模块名
```

## 后端组长审核流程

后端组员完成模块设计或接口实现后，不直接合并到 `main`，也不直接推送到 `dev`。组员统一从个人 `feature/*` 分支发起 PR/MR，由后端组长审核后端相关内容，审核通过后合并到 `dev`。

1. 组员新建自己的模块分支：

```bash
git checkout dev
git pull origin dev
git checkout -b feature/backend-姓名-模块名
```

2. 完成修改后只提交后端相关文件：

```bash
git status
git add backend
git commit -m "feat(backend): 完成某某模块"
git push origin feature/backend-姓名-模块名
```

3. 在代码托管平台发起 Pull Request / Merge Request：

- 目标分支：`dev`。
- 审核人：后端组长。
- 标题建议：`feat(backend): 某某模块设计/接口实现`

4. PR/MR 描述需写清楚：

- 本次完成的模块职责。
- 新增或调整的接口。
- 关联的数据表和状态字段。
- 权限校验和异常处理方式。
- 与概要设计不一致或仍需讨论的问题。
- 是否影响前端联调。

5. 后端组长重点审核：

- 模块边界是否清楚，是否和其他成员职责冲突。
- 数据库字段、枚举值、状态流转是否符合概要设计。
- 创建订单、取消订单、完成订单等关键流程是否使用事务。
- 权限校验是否在服务端完成，管理员接口是否校验 `role=ADMIN`。
- 统一返回结构、错误码、日志记录是否一致。
- 是否误提交前端、个人配置或无关文件。

6. 审核结果处理：

- 审核通过：由后端组长合并到 `dev`。
- 需要修改：组员在原分支继续提交，PR/MR 会自动更新。
- 数据库或流程存在分歧：在 PR/MR 评论中标注，并同步到模块设计文档的“问题清单”。

7. 与主分支的关系：

- `dev` 用于前后端日常集成和接口联调。
- `main` 保持相对稳定，不接收普通组员功能分支直接合并。
- 阶段成果稳定后，由项目组长从 `dev` 合并到 `main`。

## 后端注意事项

- 所有权限校验必须在服务端完成，不能依赖前端隐藏按钮。
- 管理员接口需要中间件校验和业务层二次校验 `role=ADMIN`。
- 文件上传需要限制类型、大小和数量。
- 敏感词检测以后端为准，前端提示只能作为辅助。
- 订单状态变更接口必须检查当前状态，避免重复操作和越权操作。
- 消息写入失败不应破坏已经成功的核心交易事务，但需要记录异常并考虑补偿。
- AI 扩展失败不能影响普通商品发布流程。

## 本阶段交付物

- 后端模块分工表。
- 各模块设计文档。
- 数据库问题清单。
- 业务流程问题清单。
- 整体后端设计方案确认稿。
