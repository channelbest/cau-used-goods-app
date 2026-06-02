# 商品模块说明

本目录基于 `feature/frontend-A-login-home-detail-v2` 的 uni-app 源码继续开发。`unpackage/` 是编译产物，已由 `.gitignore` 排除，不应作为业务源码修改。

## 微信开发者工具运行

在上一级 `frontend/` 目录执行：

```bash
npm run build:mp-weixin:cau
```

然后在微信开发者工具中导入：

```text
frontend/cau-used-goods-uni/unpackage/dist/build/mp-weixin
```

不要直接导入 `frontend/cau-used-goods-uni` 源码目录。微信开发者工具需要编译后生成的 `app.json`。

## 已实现功能

- 首页：分类入口、最新商品、分页加载、下拉刷新。
- 分类与搜索：关键词、分类、价格区间、成色、排序、分页和空状态。
- 商品详情：图片轮播、商品状态、卖家脱敏信息、收藏、预约和举报入口。
- 发布商品：表单校验、最多 9 张图片、单图 5 MB 限制、上传、删除图片。
- AI 辅助：标题优化候选、描述生成。AI 接口失败不会阻止用户手动发布。
- 权限入口：发布、收藏、预约和举报操作会检查登录与学生认证状态。

## 商品相关接口

```text
GET    /categories
GET    /products
GET    /products/:id
POST   /products
POST   /upload/products
POST   /ai/optimize-title
POST   /ai/generate-description
POST   /favorites
DELETE /favorites/:productId
GET    /favorites/check?productId=:productId
POST   /orders
POST   /reports
```

列表查询参数：

```text
keyword categoryId minPrice maxPrice conditionLevel sort page pageSize
```

图片上传接口预期返回：

```json
{
  "code": 0,
  "data": {
    "imageUrl": "/uploads/products/example.jpg"
  }
}
```
