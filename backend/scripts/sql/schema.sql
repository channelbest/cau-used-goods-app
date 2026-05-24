
CREATE DATABASE IF NOT EXISTS `cau_used_goods`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `cau_used_goods`;

CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户主键ID',
  `openid` VARCHAR(64) NOT NULL COMMENT '微信用户唯一标识',
  `nickname` VARCHAR(50) NULL COMMENT '微信昵称或用户昵称',
  `avatar_url` VARCHAR(255) NULL COMMENT '微信头像地址',
  `student_id` VARCHAR(30) NULL COMMENT '学号',
  `real_name` VARCHAR(30) NULL COMMENT '学生真实姓名',
  `college` VARCHAR(100) NULL COMMENT '学院信息',
  `phone` VARCHAR(20) NULL COMMENT '联系方式',
  `role` VARCHAR(20) NOT NULL DEFAULT 'USER' COMMENT 'USER / ADMIN',
  `auth_status` VARCHAR(20) NOT NULL DEFAULT 'UNVERIFIED' COMMENT 'UNVERIFIED / PENDING / VERIFIED / REJECTED',
  `account_status` VARCHAR(20) NOT NULL DEFAULT 'NORMAL' COMMENT 'NORMAL / DISABLED / CANCELED / DELETED',
  `last_login_time` DATETIME NULL COMMENT '最近登录时间',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_openid` (`openid`),
  UNIQUE KEY `uk_users_student_id` (`student_id`),
  KEY `idx_users_role` (`role`),
  KEY `idx_users_auth_account` (`auth_status`, `account_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `categories` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父分类ID',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序值',
  `status` VARCHAR(20) NOT NULL DEFAULT 'ENABLED' COMMENT 'ENABLED / DISABLED',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_categories_parent_name` (`parent_id`, `name`),
  KEY `idx_categories_status_sort` (`status`, `sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品分类表';

CREATE TABLE IF NOT EXISTS `products` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  `seller_id` BIGINT UNSIGNED NOT NULL COMMENT '卖家用户ID',
  `category_id` BIGINT UNSIGNED NOT NULL COMMENT '商品分类ID',
  `title` VARCHAR(100) NOT NULL COMMENT '商品标题',
  `description` TEXT NULL COMMENT '商品描述',
  `original_price` DECIMAL(10,2) NULL COMMENT '商品原价',
  `price` DECIMAL(10,2) NOT NULL COMMENT '商品售价',
  `condition_level` VARCHAR(20) NULL COMMENT '商品成色',
  `meet_location` VARCHAR(100) NULL COMMENT '建议面交地点',
  `status` VARCHAR(30) NOT NULL DEFAULT 'ON_SALE' COMMENT 'ON_SALE / LOCKED / SOLD / OFF_SHELF / DELETED',
  `view_count` INT NOT NULL DEFAULT 0 COMMENT '浏览量',
  `favorite_count` INT NOT NULL DEFAULT 0 COMMENT '收藏数',
  `off_shelf_reason` VARCHAR(255) NULL COMMENT '下架原因',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`id`),
  KEY `idx_products_seller` (`seller_id`),
  KEY `idx_products_category_status` (`category_id`, `status`),
  KEY `idx_products_status_time` (`status`, `create_time`),
  KEY `idx_products_price` (`price`),
  KEY `idx_products_title` (`title`),
  CONSTRAINT `fk_products_seller` FOREIGN KEY (`seller_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_products_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

CREATE TABLE IF NOT EXISTS `product_images` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '图片ID',
  `product_id` BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
  `image_url` VARCHAR(255) NOT NULL COMMENT '图片访问地址',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '图片排序',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  PRIMARY KEY (`id`),
  KEY `idx_product_images_product_sort` (`product_id`, `sort_order`),
  CONSTRAINT `fk_product_images_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品图片表';

CREATE TABLE IF NOT EXISTS `orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no` VARCHAR(50) NOT NULL COMMENT '订单编号',
  `product_id` BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
  `buyer_id` BIGINT UNSIGNED NOT NULL COMMENT '买家用户ID',
  `seller_id` BIGINT UNSIGNED NOT NULL COMMENT '卖家用户ID',
  `product_title_snapshot` VARCHAR(100) NOT NULL COMMENT '下单时商品标题快照',
  `product_price_snapshot` DECIMAL(10,2) NOT NULL COMMENT '下单时商品价格快照',
  `status` VARCHAR(30) NOT NULL DEFAULT 'PENDING_CONFIRM' COMMENT 'PENDING_CONFIRM / WAIT_MEET / COMPLETED / CANCELED / EXCEPTION_CLOSED',
  `remark` VARCHAR(255) NULL COMMENT '买家备注',
  `meet_time` DATETIME NULL COMMENT '约定面交时间',
  `meet_location` VARCHAR(100) NULL COMMENT '约定面交地点',
  `cancel_reason` VARCHAR(255) NULL COMMENT '取消原因',
  `cancel_by` BIGINT UNSIGNED NULL COMMENT '取消操作人ID',
  `expire_time` DATETIME NOT NULL COMMENT '卖家确认截止时间',
  `confirm_time` DATETIME NULL COMMENT '卖家确认时间',
  `finish_time` DATETIME NULL COMMENT '交易完成时间',
  `close_time` DATETIME NULL COMMENT '取消或异常关闭时间',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '下单时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_orders_order_no` (`order_no`),
  KEY `idx_orders_product_status` (`product_id`, `status`),
  KEY `idx_orders_buyer_status_time` (`buyer_id`, `status`, `create_time`),
  KEY `idx_orders_seller_status_time` (`seller_id`, `status`, `create_time`),
  KEY `idx_orders_expire` (`status`, `expire_time`),
  KEY `idx_orders_cancel_by` (`cancel_by`),
  CONSTRAINT `fk_orders_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_orders_buyer` FOREIGN KEY (`buyer_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_orders_seller` FOREIGN KEY (`seller_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_orders_cancel_by` FOREIGN KEY (`cancel_by`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

CREATE TABLE IF NOT EXISTS `favorites` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '收藏用户ID',
  `product_id` BIGINT UNSIGNED NOT NULL COMMENT '被收藏商品ID',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否取消收藏',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_favorites_user_product` (`user_id`, `product_id`),
  KEY `idx_favorites_product` (`product_id`),
  CONSTRAINT `fk_favorites_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_favorites_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收藏表';

CREATE TABLE IF NOT EXISTS `reviews` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评价ID',
  `order_id` BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
  `product_id` BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
  `reviewer_id` BIGINT UNSIGNED NOT NULL COMMENT '评价人ID',
  `seller_id` BIGINT UNSIGNED NOT NULL COMMENT '被评价卖家ID',
  `rating` INT NOT NULL COMMENT '星级评分，取值1-5',
  `content` VARCHAR(500) NULL COMMENT '文字评价',
  `status` VARCHAR(20) NOT NULL DEFAULT 'NORMAL' COMMENT 'NORMAL / HIDDEN / DELETED',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评价时间',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reviews_order_reviewer` (`order_id`, `reviewer_id`),
  KEY `idx_reviews_product` (`product_id`),
  KEY `idx_reviews_seller_status` (`seller_id`, `status`),
  KEY `idx_reviews_reviewer` (`reviewer_id`),
  CONSTRAINT `fk_reviews_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
  CONSTRAINT `fk_reviews_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_reviews_reviewer` FOREIGN KEY (`reviewer_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_reviews_seller` FOREIGN KEY (`seller_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评价表';

CREATE TABLE IF NOT EXISTS `reports` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '举报ID',
  `reporter_id` BIGINT UNSIGNED NOT NULL COMMENT '举报人ID',
  `target_type` VARCHAR(20) NOT NULL COMMENT 'PRODUCT / USER / ORDER',
  `target_id` BIGINT UNSIGNED NOT NULL COMMENT '被举报对象ID',
  `reason_type` VARCHAR(50) NOT NULL COMMENT '举报原因类型',
  `description` VARCHAR(500) NULL COMMENT '举报说明',
  `status` VARCHAR(30) NOT NULL DEFAULT 'PENDING' COMMENT 'PENDING / PROCESSING / RESOLVED / REJECTED / CLOSED',
  `handle_result` VARCHAR(500) NULL COMMENT '处理结果说明',
  `handler_id` BIGINT UNSIGNED NULL COMMENT '处理管理员ID',
  `handle_time` DATETIME NULL COMMENT '处理时间',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '举报提交时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_reports_reporter` (`reporter_id`),
  KEY `idx_reports_target` (`target_type`, `target_id`),
  KEY `idx_reports_status_time` (`status`, `create_time`),
  KEY `idx_reports_handler` (`handler_id`),
  CONSTRAINT `fk_reports_reporter` FOREIGN KEY (`reporter_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_reports_handler` FOREIGN KEY (`handler_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='举报表';

CREATE TABLE IF NOT EXISTS `report_images` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '凭证图片ID',
  `report_id` BIGINT UNSIGNED NOT NULL COMMENT '举报ID',
  `image_url` VARCHAR(255) NOT NULL COMMENT '凭证图片访问地址',
  `sort_order` INT NOT NULL DEFAULT 0 COMMENT '图片排序',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  PRIMARY KEY (`id`),
  KEY `idx_report_images_report_sort` (`report_id`, `sort_order`),
  CONSTRAINT `fk_report_images_report` FOREIGN KEY (`report_id`) REFERENCES `reports` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='举报凭证图片表';

CREATE TABLE IF NOT EXISTS `messages` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '消息ID',
  `receiver_id` BIGINT UNSIGNED NOT NULL COMMENT '接收用户ID',
  `sender_id` BIGINT UNSIGNED NULL COMMENT '发送者ID，系统消息可为空',
  `message_type` VARCHAR(30) NOT NULL COMMENT 'ORDER_CREATED / ORDER_CONFIRMED / ORDER_CANCELED / ORDER_TIMEOUT / REPORT_HANDLED / SYSTEM_NOTICE',
  `title` VARCHAR(100) NOT NULL COMMENT '消息标题',
  `content` VARCHAR(500) NOT NULL COMMENT '消息内容',
  `related_type` VARCHAR(30) NULL COMMENT 'ORDER / PRODUCT / REPORT / NOTICE',
  `related_id` BIGINT UNSIGNED NULL COMMENT '关联对象ID',
  `read_status` VARCHAR(20) NOT NULL DEFAULT 'UNREAD' COMMENT 'UNREAD / READ',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `read_time` DATETIME NULL COMMENT '读取时间',
  PRIMARY KEY (`id`),
  KEY `idx_messages_receiver_read_time` (`receiver_id`, `read_status`, `create_time`),
  KEY `idx_messages_sender` (`sender_id`),
  KEY `idx_messages_related` (`related_type`, `related_id`),
  CONSTRAINT `fk_messages_receiver` FOREIGN KEY (`receiver_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_messages_sender` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='站内消息表';

CREATE TABLE IF NOT EXISTS `sensitive_words` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '敏感词ID',
  `word` VARCHAR(100) NOT NULL COMMENT '敏感词内容',
  `word_type` VARCHAR(30) NOT NULL COMMENT 'FORBIDDEN / RISK',
  `status` VARCHAR(20) NOT NULL DEFAULT 'ENABLED' COMMENT 'ENABLED / DISABLED',
  `create_by` BIGINT UNSIGNED NULL COMMENT '创建管理员ID',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sensitive_words_word` (`word`),
  KEY `idx_sensitive_words_status` (`status`),
  KEY `idx_sensitive_words_create_by` (`create_by`),
  CONSTRAINT `fk_sensitive_words_create_by` FOREIGN KEY (`create_by`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='敏感词表';

CREATE TABLE IF NOT EXISTS `admin_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `admin_id` BIGINT UNSIGNED NOT NULL COMMENT '管理员用户ID',
  `operation_type` VARCHAR(50) NOT NULL COMMENT '操作类型',
  `target_type` VARCHAR(30) NOT NULL COMMENT 'USER / PRODUCT / ORDER / REPORT / NOTICE / WORD',
  `target_id` BIGINT UNSIGNED NOT NULL COMMENT '操作对象ID',
  `description` VARCHAR(500) NULL COMMENT '操作说明',
  `ip_address` VARCHAR(50) NULL COMMENT '操作IP',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_admin_logs_admin_time` (`admin_id`, `create_time`),
  KEY `idx_admin_logs_target` (`target_type`, `target_id`),
  KEY `idx_admin_logs_operation` (`operation_type`),
  CONSTRAINT `fk_admin_logs_admin` FOREIGN KEY (`admin_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员操作日志表';

CREATE TABLE IF NOT EXISTS `announcements` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '公告ID',
  `title` VARCHAR(100) NOT NULL COMMENT '公告标题',
  `content` TEXT NULL COMMENT '公告内容',
  `cover_url` VARCHAR(255) NULL COMMENT '封面图或轮播图地址',
  `status` VARCHAR(20) NOT NULL DEFAULT 'DRAFT' COMMENT 'DRAFT / PUBLISHED / OFFLINE',
  `publish_time` DATETIME NULL COMMENT '发布时间',
  `create_by` BIGINT UNSIGNED NOT NULL COMMENT '创建管理员ID',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_announcements_status_publish` (`status`, `publish_time`),
  KEY `idx_announcements_create_by` (`create_by`),
  CONSTRAINT `fk_announcements_create_by` FOREIGN KEY (`create_by`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公告表';

CREATE TABLE IF NOT EXISTS `browse_history` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '浏览记录ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '浏览用户ID',
  `product_id` BIGINT UNSIGNED NOT NULL COMMENT '被浏览商品ID',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '浏览时间',
  PRIMARY KEY (`id`),
  KEY `idx_browse_history_user_time` (`user_id`, `create_time`),
  KEY `idx_browse_history_product` (`product_id`),
  CONSTRAINT `fk_browse_history_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_browse_history_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='浏览历史表';

CREATE TABLE IF NOT EXISTS `ai_generation_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'AI生成记录ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '调用AI功能的用户ID',
  `product_id` BIGINT UNSIGNED NULL COMMENT '关联商品ID',
  `generation_type` VARCHAR(30) NOT NULL COMMENT 'TITLE_OPTIMIZE / DESCRIPTION_GENERATE',
  `input_text` TEXT NULL COMMENT '用户输入内容',
  `output_text` TEXT NULL COMMENT 'AI生成内容',
  `status` VARCHAR(20) NOT NULL DEFAULT 'SUCCESS' COMMENT 'SUCCESS / FAILED',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '生成时间',
  PRIMARY KEY (`id`),
  KEY `idx_ai_generation_logs_user_time` (`user_id`, `create_time`),
  KEY `idx_ai_generation_logs_product` (`product_id`),
  KEY `idx_ai_generation_logs_type_status` (`generation_type`, `status`),
  CONSTRAINT `fk_ai_generation_logs_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_ai_generation_logs_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI生成记录表';

INSERT INTO `categories` (`name`, `parent_id`, `sort_order`, `status`)
VALUES
  ('教材资料', 0, 10, 'ENABLED'),
  ('电子产品', 0, 20, 'ENABLED'),
  ('生活用品', 0, 30, 'ENABLED'),
  ('服饰鞋包', 0, 40, 'ENABLED'),
  ('运动户外', 0, 50, 'ENABLED'),
  ('其他', 0, 999, 'ENABLED')
ON DUPLICATE KEY UPDATE
  `sort_order` = VALUES(`sort_order`),
  `status` = VALUES(`status`),
  `update_time` = CURRENT_TIMESTAMP;
