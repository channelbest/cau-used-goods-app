# 王珂雅模块开发说明文档

## 一、模块负责人信息

负责人：王珂雅

负责内容：

- 商品模块（商品发布 / 浏览 / 编辑 / 删除 / 上下架）
- 商品分类模块
- 商品图片绑定模块
- 文件上传模块
- 敏感词检测模块
- AI 商品文案优化模块
- 商品数据统计模块
- 商品与订单联调状态接口说明

---

## 二、项目模块说明

本模块主要负责 CAU 二手交易平台中的商品相关业务逻辑开发。

目标是实现完整的商品生命周期管理，包括：

```text
商品发布
↓
图片上传
↓
AI 文案优化
↓
敏感词审核
↓
商品浏览
↓
商品搜索筛选
↓
商品上下架
↓
商品删除
↓
统计分析
```

同时为订单模块提供商品状态流转支持。

---

## 三、技术实现说明

后端技术栈：

- Go
- Gin Web Framework
- MySQL
- JWT 鉴权
- Resty HTTP Client
- 智谱 AI API

开发环境：

- Windows 10
- Go 1.24+
- MySQL 8.x
- 微信开发者工具（前端联调）

---

## 四、模块目录结构

本次新增 / 修改模块如下：

```text
backend/
├── cmd/server/main.go
├── config/config.yaml
├── internal/
│   ├── ai/
│   │   ├── handler.go
│   │   ├── repository.go
│   │   ├── router.go
│   │   └── service.go
│   │
│   ├── product/
│   │   ├── handler.go
│   │   ├── repository.go
│   │   ├── router.go
│   │   └── service.go
│   │
│   ├── sensitive/
│   │   ├── handler.go
│   │   ├── repository.go
│   │   ├── router.go
│   │   └── service.go
│   │
│   ├── stats/
│   │   ├── handler.go
│   │   ├── repository.go
│   │   ├── router.go
│   │   └── service.go
│   │
│   └── upload/
│       ├── handler.go
│       ├── router.go
│       └── service.go
│
└── docs/
    └── wangkeya_module_README.md
```

---

## 五、接口说明

### 5.1 分类模块

#### 获取商品分类

接口：

```http
GET /categories
```

说明：获取所有启用状态的商品分类。

返回示例：

```json
[
  {
    "id": 1,
    "name": "教材资料"
  }
]
```

---

### 5.2 商品模块

#### 1. 发布商品

接口：

```http
POST /products
```

说明：用户发布商品。

支持字段：

```json
{
  "categoryId": 1,
  "title": "高数教材",
  "description": "正版教材，九成新",
  "originalPrice": 59,
  "price": 20,
  "conditionLevel": "九成新",
  "meetLocation": "图书馆门口"
}
```

功能特点：

- 登录用户才能发布
- 自动进行敏感词检测
- 商品默认状态为 `ON_SALE`

---

#### 2. 商品浏览

接口：

```http
GET /products
```

支持以下查询方式。

##### 关键词搜索

```http
GET /products?keyword=教材
```

作用：搜索标题或描述。

##### 分类筛选

```http
GET /products?categoryId=1
```

作用：按分类查看商品。

##### 状态筛选

```http
GET /products?status=ON_SALE
```

支持状态：

```text
ON_SALE
OFF_SHELF
LOCKED
SOLD
DELETED
ALL
```

##### 价格区间筛选

```http
GET /products?minPrice=10&maxPrice=30
```

作用：筛选指定价格范围商品。

##### 排序

最新：

```http
GET /products?sort=newest
```

价格升序：

```http
GET /products?sort=price_asc
```

价格降序：

```http
GET /products?sort=price_desc
```

热度排序：

```http
GET /products?sort=popular
```

##### 分页

```http
GET /products?page=1&pageSize=10
```

返回：

```json
{
  "list": [],
  "page": 1,
  "pageSize": 10,
  "total": 20
}
```

---

#### 3. 商品详情

接口：

```http
GET /products/:id
```

说明：获取单个商品详细信息。

包含：

- 商品基本信息
- 图片列表
- 状态信息
- 浏览量
- 收藏量

---

#### 4. 我的商品

接口：

```http
GET /products/my
```

说明：获取当前登录用户发布的商品。

---

#### 5. 编辑商品

接口：

```http
PUT /products/:id
```

说明：修改商品信息。

限制：仅商品发布者可操作。

---

#### 6. 商品上下架

接口：

```http
PUT /products/:id/status
```

请求：

```json
{
  "status": "OFF_SHELF",
  "reason": "暂时不卖"
}
```

支持：

```text
ON_SALE
OFF_SHELF
```

---

#### 7. 删除商品

接口：

```http
DELETE /products/:id
```

说明：逻辑删除。

实现：

```text
status = DELETED
is_deleted = 1
```

---

### 5.3 图片模块

#### 商品图片绑定

接口：

```http
POST /products/:id/images
```

请求：

```json
{
  "images": [
    "/uploads/products/test1.jpg"
  ]
}
```

说明：为商品绑定图片。

限制：

- 最多 9 张
- 仅商品所有者可操作

---

### 5.4 文件上传模块

#### 图片上传

接口：

```http
POST /upload/image
```

说明：上传商品图片。

返回：

```json
{
  "url": "/uploads/products/xxx.jpg"
}
```

---

### 5.5 敏感词模块

说明：商品发布与编辑时自动调用。

检测内容：

```text
违禁词
违规内容
敏感内容
```

若命中：

```json
{
  "message": "内容包含敏感词"
}
```

作用：防止违规商品发布。

---

### 5.6 AI 模块

#### 商品文案优化

接口：

```http
POST /ai/optimize-product
```

请求：

```json
{
  "title": "高数教材",
  "description": "正版教材，九成新，适合期末复习"
}
```

返回：

```json
{
  "optimizedTitle": "高数教材 九成新 期末复习必备",
  "optimizedDescription": "正版高等数学教材，保存良好，适合课程学习和期末复习使用。"
}
```

实现：调用智谱 AI API。

说明：AI Key 从 `config/config.yaml` 读取。

提交代码时：不得提交真实 API Key。

---

### 5.7 统计模块

#### 商品总览

接口：

```http
GET /stats/products/overview
```

统计：

- 商品总数
- 在售数量
- 下架数量
- 已售数量
- 删除数量
- 总浏览量
- 总收藏量
- 平均价格

---

#### 分类分布

接口：

```http
GET /stats/products/category-distribution
```

统计每个分类：

- 商品数量
- 在售数量
- 平均价格
- 浏览量
- 收藏量

---

#### 状态分布

接口：

```http
GET /stats/products/status-distribution
```

统计商品状态：

```text
ON_SALE
OFF_SHELF
LOCKED
SOLD
DELETED
```

---

#### 发布趋势

接口：

```http
GET /stats/products/trend?days=7
```

统计：最近 N 天商品发布趋势。

---

## 六、商品状态说明

### ON_SALE

在售。

允许：

```text
浏览
搜索
下单
```

### LOCKED

订单锁定。

说明：买家预约成功后进入该状态。

### SOLD

已售出。

说明：交易完成。

### OFF_SHELF

卖家下架。

说明：暂不出售。

### DELETED

逻辑删除。

说明：商品不可见。

---

## 七、与订单模块联调说明

订单模块与商品模块联调规则：

### 创建订单

前提：

```text
商品必须为 ON_SALE
```

成功后：

```text
ON_SALE → LOCKED
```

### 取消订单

状态恢复：

```text
LOCKED → ON_SALE
```

### 完成订单

状态变更：

```text
LOCKED → SOLD
```

### 不允许下单

以下状态禁止下单：

```text
OFF_SHELF
DELETED
SOLD
```

---

## 八、测试说明

已测试：

```text
JWT 登录
商品发布
商品详情
商品编辑
商品删除
商品上下架
商品搜索
关键词查询
价格筛选
分类筛选
排序
分页
图片上传
图片绑定
敏感词拦截
AI 文案优化
统计模块
```

测试结果：全部通过。

命令：

```bash
go test ./...
```

结果：全部 package 正常通过。

---

## 九、注意事项

不要提交：

```text
真实 AI API Key
uploads/
微信 private config
```

建议 Git 提交前执行：

```bash
git status
```

确认无敏感文件。

---

## 十、当前待联调事项

与订单模块确认：

```text
商品状态流转是否统一
LOCKED 状态更新方式
订单取消恢复逻辑
订单完成售出逻辑
并发下单冲突处理
```

---

## 十一、总结

本模块已完成商品业务完整闭环：

```text
商品发布
商品浏览
商品搜索
商品筛选
商品管理
图片上传
内容审核
AI 优化
数据统计
```

具备实际校园二手交易平台基础能力。
