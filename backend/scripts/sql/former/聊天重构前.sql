-- --------------------------------------------------------
-- 主机:                           127.0.0.1
-- 服务器版本:                        8.4.8 - MySQL Community Server - GPL
-- 服务器操作系统:                      Win64
-- HeidiSQL 版本:                  12.10.0.7000
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- 导出 cau_used_goods 的数据库结构
CREATE DATABASE IF NOT EXISTS `cau_used_goods` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `cau_used_goods`;

-- 导出  表 cau_used_goods.admin_logs 结构
CREATE TABLE IF NOT EXISTS `admin_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `admin_id` bigint unsigned NOT NULL COMMENT '管理员用户ID',
  `operation_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作类型',
  `target_type` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'USER / PRODUCT / ORDER / NOTICE / WORD / CATEGORY / REPORT / APPEAL',
  `target_id` bigint unsigned NOT NULL COMMENT '操作对象ID',
  `description` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作说明',
  `ip_address` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作IP',
  `related_type` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作依据类型：REPORT / APPEAL',
  `related_id` bigint unsigned DEFAULT NULL COMMENT '操作依据ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_admin_logs_admin_time` (`admin_id`,`create_time`),
  KEY `idx_admin_logs_target` (`target_type`,`target_id`),
  KEY `idx_admin_logs_operation` (`operation_type`),
  KEY `idx_admin_logs_related` (`related_type`,`related_id`),
  CONSTRAINT `fk_admin_logs_admin` FOREIGN KEY (`admin_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=102 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员操作日志表';

-- 正在导出表  cau_used_goods.admin_logs 的数据：~48 rows (大约)
INSERT INTO `admin_logs` (`id`, `admin_id`, `operation_type`, `target_type`, `target_id`, `description`, `ip_address`, `related_type`, `related_id`, `create_time`) VALUES
	(1, 2, 'STUDENT_VERIFY_APPROVE', 'USER', 1, '学生认证审核通过', NULL, NULL, NULL, '2026-06-19 18:05:35'),
	(2, 2, 'STUDENT_VERIFY_APPROVE', 'USER', 1, '学生认证审核通过', NULL, NULL, NULL, '2026-06-19 18:22:40'),
	(3, 2, 'STUDENT_VERIFY_APPROVE', 'USER', 3, '学生认证审核通过', NULL, NULL, NULL, '2026-06-19 19:28:36'),
	(4, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 1, 'off shelf product due account status DISABLED: 恶意发布多个商品', '127.0.0.1', NULL, NULL, '2026-06-19 20:10:56'),
	(5, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 2, 'off shelf product due account status DISABLED: 恶意发布多个商品', '127.0.0.1', NULL, NULL, '2026-06-19 20:10:56'),
	(6, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 3, 'off shelf product due account status DISABLED: 恶意发布多个商品', '127.0.0.1', NULL, NULL, '2026-06-19 20:10:56'),
	(7, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 4, 'off shelf product due account status DISABLED: 恶意发布多个商品', '127.0.0.1', NULL, NULL, '2026-06-19 20:10:56'),
	(8, 2, 'USER_DISABLE', 'USER', 1, '恶意发布多个商品', '127.0.0.1', NULL, NULL, '2026-06-19 20:10:56'),
	(9, 2, 'MARK_APPEAL_PROCESSING', 'APPEAL', 1, 'mark appeal #1 as PROCESSING', '127.0.0.1', NULL, NULL, '2026-06-19 20:29:03'),
	(10, 2, 'USER_ENABLE', 'USER', 1, '申诉理由合理', '127.0.0.1', NULL, NULL, '2026-06-19 20:31:02'),
	(11, 2, 'APPROVE_APPEAL', 'APPEAL', 1, 'handle appeal #1: APPROVED', '127.0.0.1', NULL, NULL, '2026-06-19 20:31:10'),
	(12, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 4, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-19 20:52:17'),
	(13, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 4, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-19 20:52:20'),
	(14, 2, 'USER_BAN', 'USER', 3, '测试', '127.0.0.1', NULL, NULL, '2026-06-19 23:06:59'),
	(15, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 1, 'off shelf product due account status DISABLED: hi', '127.0.0.1', NULL, NULL, '2026-06-19 23:52:02'),
	(16, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 2, 'off shelf product due account status DISABLED: hi', '127.0.0.1', NULL, NULL, '2026-06-19 23:52:02'),
	(17, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 3, 'off shelf product due account status DISABLED: hi', '127.0.0.1', NULL, NULL, '2026-06-19 23:52:02'),
	(18, 2, 'USER_DISABLE', 'USER', 1, 'hi', '127.0.0.1', NULL, NULL, '2026-06-19 23:52:02'),
	(19, 2, 'USER_ENABLE', 'USER', 1, '1', '127.0.0.1', NULL, NULL, '2026-06-20 00:11:47'),
	(20, 2, 'STUDENT_VERIFY_APPROVE', 'USER', 5, '学生认证审核通过', NULL, NULL, NULL, '2026-06-20 14:13:32'),
	(21, 2, 'MARK_APPEAL_PROCESSING', 'APPEAL', 3, 'mark appeal #3 as PROCESSING', '127.0.0.1', NULL, NULL, '2026-06-20 14:40:48'),
	(22, 2, 'APPROVE_APPEAL', 'APPEAL', 3, 'handle appeal #3: APPROVED', '127.0.0.1', NULL, NULL, '2026-06-20 14:40:57'),
	(23, 2, 'STUDENT_VERIFY_APPROVE', 'USER', 6, '学生认证审核通过', NULL, NULL, NULL, '2026-06-20 14:41:12'),
	(24, 2, 'STUDENT_VERIFY_APPROVE', 'USER', 7, '学生认证审核通过', NULL, NULL, NULL, '2026-06-20 14:42:52'),
	(25, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 6, 'off shelf product due account status DISABLED: 虚假宣传', '127.0.0.1', NULL, NULL, '2026-06-20 15:28:45'),
	(26, 2, 'USER_DISABLE', 'USER', 6, '虚假宣传', '127.0.0.1', NULL, NULL, '2026-06-20 15:28:45'),
	(27, 2, 'STUDENT_VERIFY_REJECT', 'USER', 4, '认证信息不符合要求', NULL, NULL, NULL, '2026-06-20 16:51:27'),
	(28, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 9, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 17:23:42'),
	(29, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 9, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 17:23:44'),
	(30, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 9, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 17:23:45'),
	(31, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 9, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 17:23:47'),
	(32, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 9, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 17:23:48'),
	(33, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 9, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 17:23:51'),
	(34, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 9, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 17:23:52'),
	(35, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 9, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 17:23:53'),
	(36, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 9, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 17:25:51'),
	(37, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 9, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 17:28:53'),
	(38, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 9, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 18:09:14'),
	(39, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 9, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 18:10:45'),
	(40, 2, 'MARK_REPORT_PROCESSING', 'REPORT', 2, 'mark report as PROCESSING', '127.0.0.1', NULL, NULL, '2026-06-20 18:10:53'),
	(41, 2, 'REJECT_REPORT', 'REPORT', 2, 'handle report #2: REJECTED; result=证据不足，无法认定违规', '127.0.0.1', NULL, NULL, '2026-06-20 18:11:00'),
	(42, 2, 'MARK_REPORT_PROCESSING', 'REPORT', 1, 'mark report as PROCESSING', '127.0.0.1', NULL, NULL, '2026-06-20 18:27:10'),
	(43, 2, 'ORDER_EXCEPTION_CLOSE', 'ORDER', 4, 'auto exception close by account status change; responsibleParty=BUYER; reason=可以', '127.0.0.1', 'REPORT', 1, '2026-06-20 18:28:35'),
	(44, 2, 'USER_DISABLE', 'USER', 1, '可以', '127.0.0.1', 'REPORT', 1, '2026-06-20 18:28:35'),
	(45, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 10, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 18:37:51'),
	(46, 2, 'MARK_APPEAL_PROCESSING', 'APPEAL', 2, 'mark appeal #2 as PROCESSING', '127.0.0.1', NULL, NULL, '2026-06-20 18:43:41'),
	(47, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 4, 'update product status to ON_SALE: 管理员上架商品', '127.0.0.1', 'APPEAL', 2, '2026-06-20 18:43:47'),
	(48, 2, 'APPROVE_APPEAL', 'APPEAL', 2, 'handle appeal #2: APPROVED', '127.0.0.1', NULL, NULL, '2026-06-20 18:43:54'),
	(49, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 4, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', 'APPEAL', 2, '2026-06-20 18:45:39'),
	(50, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 4, 'update product status to ON_SALE: 管理员上架商品', '127.0.0.1', 'APPEAL', 2, '2026-06-20 18:45:50'),
	(51, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 8, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-20 18:54:32'),
	(52, 2, 'USER_ENABLE', 'USER', 1, '你好', '127.0.0.1', 'REPORT', 1, '2026-06-20 19:12:34'),
	(53, 2, 'REJECT_REPORT', 'REPORT', 1, 'handle report #1: REJECTED; result=重复举报或恶意举报', '127.0.0.1', NULL, NULL, '2026-06-20 19:14:51'),
	(54, 2, 'MARK_REPORT_PROCESSING', 'REPORT', 3, 'mark report as PROCESSING', '127.0.0.1', NULL, NULL, '2026-06-20 19:16:57'),
	(55, 2, 'ORDER_EXCEPTION_CLOSE', 'ORDER', 5, 'responsibleParty=SELLER; reason=商品违规或信息异常', '127.0.0.1', 'REPORT', 3, '2026-06-20 19:17:12'),
	(56, 2, 'APPROVE_REPORT', 'REPORT', 3, 'handle report #3: APPROVED; result=举报已通过', '127.0.0.1', NULL, NULL, '2026-06-20 19:17:19'),
	(57, 2, 'USER_ENABLE', 'USER', 6, '正常用户', '127.0.0.1', NULL, NULL, '2026-06-20 19:53:32'),
	(58, 2, 'CREATE_NOTICE', 'NOTICE', 1, 'create announcement: 即将上线！', '127.0.0.1', NULL, NULL, '2026-06-20 20:40:07'),
	(59, 2, 'STATUS_NOTICE', 'NOTICE', 1, 'update announcement status: PUBLISHED', '127.0.0.1', NULL, NULL, '2026-06-20 20:40:13'),
	(60, 2, 'STATUS_NOTICE', 'NOTICE', 1, 'update announcement status: OFFLINE', '127.0.0.1', NULL, NULL, '2026-06-20 20:41:35'),
	(61, 2, 'STATUS_NOTICE', 'NOTICE', 1, 'update announcement status: PUBLISHED', '127.0.0.1', NULL, NULL, '2026-06-20 21:02:43'),
	(62, 2, 'STATUS_NOTICE', 'NOTICE', 1, 'update announcement status: OFFLINE', '127.0.0.1', NULL, NULL, '2026-06-20 21:03:16'),
	(63, 2, 'STATUS_NOTICE', 'NOTICE', 1, 'update announcement status: PUBLISHED', '127.0.0.1', NULL, NULL, '2026-06-20 21:05:29'),
	(64, 2, 'USER_DISABLE', 'USER', 6, '不喜欢', '127.0.0.1', NULL, NULL, '2026-06-20 21:17:39'),
	(65, 2, 'MARK_APPEAL_PROCESSING', 'APPEAL', 4, 'mark appeal #4 as PROCESSING', '127.0.0.1', NULL, NULL, '2026-06-20 21:30:21'),
	(66, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 10, 'update product status to ON_SALE: 管理员上架商品', '127.0.0.1', 'APPEAL', 4, '2026-06-20 21:31:05'),
	(67, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 4, 'off shelf product due account status DISABLED: 不喜欢', '127.0.0.1', NULL, NULL, '2026-06-20 21:31:16'),
	(68, 2, 'USER_DISABLE', 'USER', 1, '不喜欢', '127.0.0.1', NULL, NULL, '2026-06-20 21:31:16'),
	(69, 2, 'APPROVE_APPEAL', 'APPEAL', 4, 'handle appeal #4: APPROVED', '127.0.0.1', NULL, NULL, '2026-06-20 21:31:37'),
	(70, 2, 'USER_ENABLE', 'USER', 1, '错误', '127.0.0.1', NULL, NULL, '2026-06-20 21:32:07'),
	(71, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 10, 'off shelf product due account status DISABLED: 不喜欢', '127.0.0.1', NULL, NULL, '2026-06-20 21:32:18'),
	(72, 2, 'USER_DISABLE', 'USER', 7, '不喜欢', '127.0.0.1', NULL, NULL, '2026-06-20 21:32:18'),
	(73, 2, 'USER_ENABLE', 'USER', 7, '错误', '127.0.0.1', NULL, NULL, '2026-06-20 21:33:14'),
	(74, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 11, 'off shelf product due account status DISABLED: 不喜欢', '127.0.0.1', NULL, NULL, '2026-06-20 21:35:22'),
	(75, 2, 'USER_DISABLE', 'USER', 7, '不喜欢', '127.0.0.1', NULL, NULL, '2026-06-20 21:35:22'),
	(76, 2, 'USER_ENABLE', 'USER', 7, '错误', '127.0.0.1', NULL, NULL, '2026-06-20 21:38:01'),
	(77, 2, 'USER_ENABLE', 'USER', 6, '错误', '127.0.0.1', NULL, NULL, '2026-06-20 21:38:06'),
	(78, 2, 'MARK_REPORT_PROCESSING', 'REPORT', 4, 'mark report as PROCESSING', '127.0.0.1', NULL, NULL, '2026-06-20 21:39:11'),
	(79, 2, 'USER_DISABLE', 'USER', 6, '色情', '127.0.0.1', 'REPORT', 4, '2026-06-20 21:39:22'),
	(80, 2, 'APPROVE_REPORT', 'REPORT', 4, 'handle report #4: APPROVED; result=举报已通过', '127.0.0.1', NULL, NULL, '2026-06-20 21:39:26'),
	(81, 2, 'USER_ENABLE', 'USER', 6, '错误', '127.0.0.1', NULL, NULL, '2026-06-20 21:40:03'),
	(82, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 12, 'off shelf product due account status DISABLED: 不喜欢', '127.0.0.1', NULL, NULL, '2026-06-20 21:42:05'),
	(83, 2, 'USER_DISABLE', 'USER', 6, '不喜欢', '127.0.0.1', NULL, NULL, '2026-06-20 21:42:05'),
	(84, 2, 'USER_ENABLE', 'USER', 6, '错误', '127.0.0.1', NULL, NULL, '2026-06-20 21:44:47'),
	(85, 2, 'USER_DISABLE', 'USER', 6, '不喜欢', '127.0.0.1', NULL, NULL, '2026-06-20 21:44:53'),
	(86, 2, 'USER_ENABLE', 'USER', 6, '错误', '127.0.0.1', NULL, NULL, '2026-06-20 22:01:41'),
	(87, 11, 'USER_BAN', 'USER', 8, '违规', '127.0.0.1', NULL, NULL, '2026-06-21 12:58:01'),
	(88, 11, 'MARK_APPEAL_PROCESSING', 'APPEAL', 5, 'mark appeal #5 as PROCESSING', '127.0.0.1', NULL, NULL, '2026-06-21 13:02:58'),
	(89, 11, 'APPROVE_APPEAL', 'APPEAL', 5, 'handle appeal #5: APPROVED', '127.0.0.1', NULL, NULL, '2026-06-21 13:05:54'),
	(90, 11, 'USER_UNBAN', 'USER', 8, '该用户不违规', '127.0.0.1', 'APPEAL', 5, '2026-06-21 13:06:03'),
	(91, 2, 'APPROVE_APPEAL', 'APPEAL', 6, 'handle appeal #6: APPROVED', '127.0.0.1', NULL, NULL, '2026-06-21 13:25:45'),
	(92, 2, 'STUDENT_VERIFY_APPROVE', 'USER', 12, '学生认证审核通过', NULL, NULL, NULL, '2026-06-21 14:46:40'),
	(93, 2, 'UPDATE_PRODUCT_STATUS', 'PRODUCT', 14, 'update product status to OFF_SHELF: 管理员下架商品', '127.0.0.1', NULL, NULL, '2026-06-21 14:48:05'),
	(94, 2, 'MARK_APPEAL_PROCESSING', 'APPEAL', 8, 'mark appeal #8 as PROCESSING', '127.0.0.1', NULL, NULL, '2026-06-21 14:51:01'),
	(95, 2, 'APPROVE_APPEAL', 'APPEAL', 8, 'handle appeal #8: APPROVED', '127.0.0.1', NULL, NULL, '2026-06-21 14:51:06'),
	(96, 2, 'ORDER_EXCEPTION_CLOSE', 'ORDER', 6, 'responsibleParty=SELLER; reason=商品违规或信息异常', '127.0.0.1', NULL, NULL, '2026-06-21 15:17:05'),
	(97, 2, 'REJECT_REPORT', 'REPORT', 5, 'handle report #5: REJECTED; result=举报内容与对象不符', '127.0.0.1', NULL, NULL, '2026-06-21 15:18:54'),
	(98, 2, 'REJECT_APPEAL', 'APPEAL', 9, 'handle appeal #9: REJECTED', '127.0.0.1', NULL, NULL, '2026-06-21 15:22:10'),
	(99, 11, 'MARK_APPEAL_PROCESSING', 'APPEAL', 7, 'mark appeal #7 as PROCESSING', '127.0.0.1', NULL, NULL, '2026-06-21 20:39:18'),
	(100, 11, 'USER_UNBAN', 'USER', 3, '好吧', '127.0.0.1', 'APPEAL', 7, '2026-06-21 20:39:32'),
	(101, 11, 'APPROVE_APPEAL', 'APPEAL', 7, 'handle appeal #7: APPROVED', '127.0.0.1', NULL, NULL, '2026-06-21 20:39:43');

-- 导出  表 cau_used_goods.ai_generation_logs 结构
CREATE TABLE IF NOT EXISTS `ai_generation_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'AI生成记录ID',
  `user_id` bigint unsigned NOT NULL COMMENT '调用AI功能的用户ID',
  `product_id` bigint unsigned DEFAULT NULL COMMENT '关联商品ID',
  `generation_type` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'TITLE_OPTIMIZE / DESCRIPTION_GENERATE',
  `input_text` text COLLATE utf8mb4_unicode_ci COMMENT '用户输入内容',
  `output_text` text COLLATE utf8mb4_unicode_ci COMMENT 'AI生成内容',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'SUCCESS' COMMENT 'SUCCESS / FAILED',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '生成时间',
  PRIMARY KEY (`id`),
  KEY `idx_ai_generation_logs_user_time` (`user_id`,`create_time`),
  KEY `idx_ai_generation_logs_product` (`product_id`),
  KEY `idx_ai_generation_logs_type_status` (`generation_type`,`status`),
  CONSTRAINT `fk_ai_generation_logs_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_ai_generation_logs_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI生成记录表';

-- 正在导出表  cau_used_goods.ai_generation_logs 的数据：~0 rows (大约)

-- 导出  表 cau_used_goods.announcements 结构
CREATE TABLE IF NOT EXISTS `announcements` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '公告ID',
  `title` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '公告标题',
  `content` text COLLATE utf8mb4_unicode_ci COMMENT '公告内容',
  `cover_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '封面图或轮播图地址',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'DRAFT' COMMENT 'DRAFT / PUBLISHED / OFFLINE',
  `publish_time` datetime DEFAULT NULL COMMENT '发布时间',
  `create_by` bigint unsigned NOT NULL COMMENT '创建管理员ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_announcements_status_publish` (`status`,`publish_time`),
  KEY `idx_announcements_create_by` (`create_by`),
  CONSTRAINT `fk_announcements_create_by` FOREIGN KEY (`create_by`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公告表';

-- 正在导出表  cau_used_goods.announcements 的数据：~0 rows (大约)
INSERT INTO `announcements` (`id`, `title`, `content`, `cover_url`, `status`, `publish_time`, `create_by`, `create_time`, `update_time`) VALUES
	(1, '即将上线！', '我们的小程序将于6.21号上线！敬请期待！', NULL, 'PUBLISHED', '2026-06-20 21:05:29', 2, '2026-06-20 20:40:07', '2026-06-20 21:05:29');

-- 导出  表 cau_used_goods.appeals 结构
CREATE TABLE IF NOT EXISTS `appeals` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '申诉ID',
  `appellant_id` bigint unsigned NOT NULL COMMENT '申诉人用户ID',
  `target_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申诉对象类型：PRODUCT / USER / ORDER / REPORT',
  `target_id` bigint unsigned NOT NULL COMMENT '申诉对象ID',
  `reason` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申诉理由',
  `status` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'PENDING' COMMENT 'PENDING / PROCESSING / APPROVED / REJECTED / CLOSED',
  `handle_result` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '管理员处理结果',
  `handler_id` bigint unsigned DEFAULT NULL COMMENT '处理管理员ID',
  `handle_time` datetime DEFAULT NULL COMMENT '处理时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_appeals_appellant_status_time` (`appellant_id`,`status`,`create_time`),
  KEY `idx_appeals_target` (`target_type`,`target_id`),
  KEY `idx_appeals_status_time` (`status`,`create_time`),
  KEY `idx_appeals_handler` (`handler_id`),
  CONSTRAINT `fk_appeals_appellant` FOREIGN KEY (`appellant_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_appeals_handler` FOREIGN KEY (`handler_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='申诉表';

-- 正在导出表  cau_used_goods.appeals 的数据：~7 rows (大约)
INSERT INTO `appeals` (`id`, `appellant_id`, `target_type`, `target_id`, `reason`, `status`, `handle_result`, `handler_id`, `handle_time`, `create_time`, `update_time`) VALUES
	(1, 1, 'USER', 1, '不应该封禁我，我只是二手商品比较多', 'APPROVED', '申诉已通过', 2, '2026-06-19 20:31:10', '2026-06-19 20:26:59', '2026-06-19 20:31:10'),
	(2, 1, 'PRODUCT', 4, '111', 'APPROVED', '申诉已通过', 2, '2026-06-20 18:43:54', '2026-06-19 23:24:56', '2026-06-20 18:43:54'),
	(3, 5, 'ORDER', 2, '金钱纠纷', 'APPROVED', '申诉已通过', 2, '2026-06-20 14:40:57', '2026-06-20 14:36:59', '2026-06-20 14:40:57'),
	(4, 7, 'PRODUCT', 10, '为啥下架', 'APPROVED', '申诉已通过', 2, '2026-06-20 21:31:37', '2026-06-20 21:29:40', '2026-06-20 21:31:37'),
	(5, 8, 'USER', 8, '我没有违规，我什么都没有干，我甚至没有发布商品', 'APPROVED', '申诉已通过', 11, '2026-06-21 13:05:54', '2026-06-21 13:02:43', '2026-06-21 13:05:54'),
	(6, 3, 'USER', 3, '不应该封禁我', 'APPROVED', '申诉已通过', 2, '2026-06-21 13:25:45', '2026-06-21 13:25:16', '2026-06-21 13:25:45'),
	(7, 3, 'USER', 3, '还是没有解封禁我的用户', 'APPROVED', '申诉已通过', 11, '2026-06-21 20:39:43', '2026-06-21 13:29:43', '2026-06-21 20:39:43'),
	(8, 12, 'PRODUCT', 14, '为什么下架', 'APPROVED', '申诉已通过', 2, '2026-06-21 14:51:06', '2026-06-21 14:50:30', '2026-06-21 14:51:06'),
	(9, 6, 'ORDER', 6, '我的怎么了', 'REJECTED', '未提供有效证明', 2, '2026-06-21 15:22:10', '2026-06-21 15:18:17', '2026-06-21 15:22:10');

-- 导出  表 cau_used_goods.appeal_images 结构
CREATE TABLE IF NOT EXISTS `appeal_images` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '申诉凭证图片ID',
  `appeal_id` bigint unsigned NOT NULL COMMENT '申诉ID',
  `image_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '凭证图片访问路径',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '图片排序',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  PRIMARY KEY (`id`),
  KEY `idx_appeal_images_appeal_sort` (`appeal_id`,`sort_order`),
  CONSTRAINT `fk_appeal_images_appeal` FOREIGN KEY (`appeal_id`) REFERENCES `appeals` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='申诉凭证图片表';

-- 正在导出表  cau_used_goods.appeal_images 的数据：~0 rows (大约)

-- 导出  表 cau_used_goods.browse_history 结构
CREATE TABLE IF NOT EXISTS `browse_history` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '浏览记录ID',
  `user_id` bigint unsigned NOT NULL COMMENT '浏览用户ID',
  `product_id` bigint unsigned NOT NULL COMMENT '被浏览商品ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '浏览时间',
  PRIMARY KEY (`id`),
  KEY `idx_browse_history_user_time` (`user_id`,`create_time`),
  KEY `idx_browse_history_product` (`product_id`),
  CONSTRAINT `fk_browse_history_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_browse_history_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='浏览历史表';

-- 正在导出表  cau_used_goods.browse_history 的数据：~0 rows (大约)
INSERT INTO `browse_history` (`id`, `user_id`, `product_id`, `create_time`) VALUES
	(1, 3, 3, '2026-06-19 19:08:32'),
	(2, 3, 2, '2026-06-19 19:08:40'),
	(3, 3, 1, '2026-06-19 19:24:18'),
	(4, 4, 4, '2026-06-19 19:49:34'),
	(5, 3, 4, '2026-06-19 19:59:02'),
	(6, 4, 3, '2026-06-19 20:03:26'),
	(7, 4, 2, '2026-06-19 20:04:08'),
	(8, 4, 1, '2026-06-19 20:04:21'),
	(9, 3, 5, '2026-06-20 11:09:01'),
	(10, 5, 5, '2026-06-20 14:14:01'),
	(11, 7, 6, '2026-06-20 14:44:34'),
	(12, 1, 7, '2026-06-20 17:07:21'),
	(13, 1, 8, '2026-06-20 17:22:01'),
	(14, 1, 9, '2026-06-20 18:10:18'),
	(15, 3, 4, '2026-06-20 19:06:07'),
	(16, 5, 4, '2026-06-20 19:11:51'),
	(17, 3, 4, '2026-06-20 20:31:32'),
	(18, 7, 4, '2026-06-20 20:43:55'),
	(19, 6, 4, '2026-06-20 21:18:15'),
	(20, 1, 13, '2026-06-20 22:02:33'),
	(21, 7, 13, '2026-06-20 22:10:13'),
	(22, 8, 13, '2026-06-21 13:01:59'),
	(23, 3, 13, '2026-06-21 13:28:21'),
	(24, 1, 13, '2026-06-21 13:40:38'),
	(25, 12, 13, '2026-06-21 14:54:42'),
	(26, 1, 12, '2026-06-21 16:05:32'),
	(27, 6, 3, '2026-06-21 16:40:56'),
	(28, 1, 15, '2026-06-21 16:57:37'),
	(29, 12, 13, '2026-06-21 20:05:16');

-- 导出  表 cau_used_goods.categories 结构
CREATE TABLE IF NOT EXISTS `categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '父分类ID',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '排序值',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ENABLED' COMMENT 'ENABLED / DISABLED',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_categories_parent_name` (`parent_id`,`name`),
  KEY `idx_categories_status_sort` (`status`,`sort_order`)
) ENGINE=InnoDB AUTO_INCREMENT=73 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品分类表';

-- 正在导出表  cau_used_goods.categories 的数据：~66 rows (大约)
INSERT INTO `categories` (`id`, `name`, `parent_id`, `sort_order`, `status`, `create_time`, `update_time`) VALUES
	(1, '教材资料', 0, 10, 'ENABLED', '2026-06-16 17:15:11', '2026-06-19 14:35:44'),
	(2, '电子产品', 0, 20, 'ENABLED', '2026-06-16 17:15:11', '2026-06-19 14:35:44'),
	(3, '生活用品', 0, 30, 'ENABLED', '2026-06-16 17:15:11', '2026-06-19 14:35:44'),
	(4, '服饰鞋包', 0, 40, 'ENABLED', '2026-06-16 17:15:11', '2026-06-19 14:35:44'),
	(5, '运动户外', 0, 50, 'ENABLED', '2026-06-16 17:15:11', '2026-06-19 14:35:44'),
	(6, '其他', 0, 999, 'ENABLED', '2026-06-16 17:15:11', '2026-06-19 14:35:44'),
	(13, '公共课教材', 1, 101, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(14, '专业课教材', 1, 102, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(15, '考研资料', 1, 103, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(16, '考证资料', 1, 104, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(17, '外语学习', 1, 105, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(18, '课程笔记', 1, 106, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(19, '实验报告资料', 1, 107, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(20, '课外读物', 1, 108, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(21, '教辅习题', 1, 109, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(22, '打印复印资料', 1, 110, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(23, '手机通讯', 2, 201, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(24, '电脑笔记本', 2, 202, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(25, '平板设备', 2, 203, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(26, '耳机音响', 2, 204, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(27, '相机摄影', 2, 205, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(28, '键盘鼠标', 2, 206, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(29, '充电器线材', 2, 207, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(30, '存储设备', 2, 208, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(31, '智能穿戴', 2, 209, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(32, '数码配件', 2, 210, 'ENABLED', '2026-06-19 14:35:44', '2026-06-19 14:35:44'),
	(33, '宿舍用品', 3, 301, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(34, '学习文具', 3, 302, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(35, '收纳整理', 3, 303, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(36, '台灯照明', 3, 304, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(37, '小家电', 3, 305, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(38, '厨具餐具', 3, 306, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(39, '清洁用品', 3, 307, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(40, '床上用品', 3, 308, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(41, '美妆个护', 3, 309, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(42, '日用杂物', 3, 310, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(43, '男装', 4, 401, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(44, '女装', 4, 402, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(45, '鞋靴', 4, 403, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(46, '箱包', 4, 404, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(47, '帽子围巾', 4, 405, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(48, '手表饰品', 4, 406, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(49, '正装礼服', 4, 407, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(50, '运动服饰', 4, 408, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(51, '校服院服', 4, 409, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(52, '其他配饰', 4, 410, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(53, '球类用品', 5, 501, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(54, '健身器材', 5, 502, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(55, '户外装备', 5, 503, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(56, '骑行装备', 5, 504, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(57, '运动护具', 5, 505, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(58, '瑜伽舞蹈', 5, 506, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(59, '游泳用品', 5, 507, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(60, '校园代步', 5, 508, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(61, '运动鞋服', 5, 509, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(62, '露营旅行', 5, 510, 'ENABLED', '2026-06-19 14:35:45', '2026-06-19 14:35:45'),
	(63, '票券卡券', 6, 901, 'ENABLED', '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(64, '虚拟资料', 6, 902, 'ENABLED', '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(65, '乐器器材', 6, 903, 'ENABLED', '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(66, '宠物用品', 6, 904, 'ENABLED', '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(67, '手工艺品', 6, 905, 'ENABLED', '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(68, '模型玩具', 6, 906, 'ENABLED', '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(69, '活动周边', 6, 907, 'ENABLED', '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(70, '租借服务', 6, 908, 'ENABLED', '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(71, '闲置赠送', 6, 909, 'ENABLED', '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(72, '其他闲置', 6, 999, 'ENABLED', '2026-06-19 14:35:46', '2026-06-19 14:35:46');

-- 导出  表 cau_used_goods.chat_conversations 结构
CREATE TABLE IF NOT EXISTS `chat_conversations` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '聊天会话ID',
  `product_id` bigint unsigned NOT NULL COMMENT '关联商品ID',
  `buyer_id` bigint unsigned NOT NULL COMMENT '买家用户ID',
  `seller_id` bigint unsigned NOT NULL COMMENT '卖家用户ID',
  `last_message_id` bigint unsigned DEFAULT NULL COMMENT '最后一条消息ID',
  `last_message_content` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '最后一条消息内容',
  `last_message_time` datetime DEFAULT NULL COMMENT '最后一条消息时间',
  `buyer_unread_count` int NOT NULL DEFAULT '0' COMMENT '买家未读数',
  `seller_unread_count` int NOT NULL DEFAULT '0' COMMENT '卖家未读数',
  `buyer_hidden_at` datetime DEFAULT NULL COMMENT '买家隐藏会话时间，NULL表示买家聊天列表可见',
  `seller_hidden_at` datetime DEFAULT NULL COMMENT '卖家隐藏会话时间，NULL表示卖家聊天列表可见',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVE' COMMENT 'ACTIVE / CLOSED',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_chat_product_buyer_seller` (`product_id`,`buyer_id`,`seller_id`),
  KEY `idx_chat_buyer_time` (`buyer_id`,`last_message_time`),
  KEY `idx_chat_seller_time` (`seller_id`,`last_message_time`),
  CONSTRAINT `fk_chat_conversations_buyer` FOREIGN KEY (`buyer_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_chat_conversations_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_chat_conversations_seller` FOREIGN KEY (`seller_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天会话表';

-- 正在导出表  cau_used_goods.chat_conversations 的数据：~0 rows (大约)
INSERT INTO `chat_conversations` (`id`, `product_id`, `buyer_id`, `seller_id`, `last_message_id`, `last_message_content`, `last_message_time`, `buyer_unread_count`, `seller_unread_count`, `buyer_hidden_at`, `seller_hidden_at`, `status`, `create_time`, `update_time`) VALUES
	(1, 3, 3, 1, 8, 'jo', '2026-06-19 23:48:19', 0, 0, NULL, NULL, 'ACTIVE', '2026-06-19 19:29:20', '2026-06-21 20:25:14'),
	(2, 7, 2, 7, NULL, NULL, NULL, 0, 0, NULL, NULL, 'ACTIVE', '2026-06-20 14:56:36', '2026-06-20 16:10:58'),
	(3, 6, 2, 6, NULL, NULL, NULL, 0, 0, NULL, NULL, 'ACTIVE', '2026-06-20 14:58:20', '2026-06-21 14:42:41'),
	(4, 6, 7, 6, 17, '下班', '2026-06-20 22:13:30', 1, 0, NULL, NULL, 'ACTIVE', '2026-06-20 15:27:10', '2026-06-21 16:49:57'),
	(5, 4, 7, 1, 13, '怎么卖呀', '2026-06-20 20:44:09', 0, 0, NULL, NULL, 'ACTIVE', '2026-06-20 20:44:03', '2026-06-20 22:07:50'),
	(6, 4, 2, 1, 14, '你好，我是管理员', '2026-06-20 20:52:06', 0, 0, NULL, NULL, 'ACTIVE', '2026-06-20 20:51:43', '2026-06-20 20:52:26'),
	(7, 13, 1, 6, 18, 'nihao', '2026-06-21 16:49:47', 0, 0, NULL, NULL, 'ACTIVE', '2026-06-20 22:02:37', '2026-06-21 16:55:59'),
	(8, 15, 1, 12, NULL, NULL, NULL, 0, 0, NULL, NULL, 'ACTIVE', '2026-06-21 16:57:38', '2026-06-21 20:34:12');

-- 导出  表 cau_used_goods.chat_messages 结构
CREATE TABLE IF NOT EXISTS `chat_messages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '聊天消息ID',
  `conversation_id` bigint unsigned NOT NULL COMMENT '聊天会话ID',
  `sender_id` bigint unsigned NOT NULL COMMENT '发送人ID',
  `receiver_id` bigint unsigned NOT NULL COMMENT '接收人ID',
  `content` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '消息内容',
  `message_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'TEXT' COMMENT 'TEXT',
  `read_status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'UNREAD' COMMENT 'UNREAD / READ',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
  PRIMARY KEY (`id`),
  KEY `idx_chat_messages_conversation_time` (`conversation_id`,`create_time`),
  KEY `idx_chat_messages_receiver_read` (`receiver_id`,`read_status`),
  KEY `fk_chat_messages_sender` (`sender_id`),
  CONSTRAINT `fk_chat_messages_conversation` FOREIGN KEY (`conversation_id`) REFERENCES `chat_conversations` (`id`),
  CONSTRAINT `fk_chat_messages_receiver` FOREIGN KEY (`receiver_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_chat_messages_sender` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天消息表';

-- 正在导出表  cau_used_goods.chat_messages 的数据：~0 rows (大约)
INSERT INTO `chat_messages` (`id`, `conversation_id`, `sender_id`, `receiver_id`, `content`, `message_type`, `read_status`, `create_time`) VALUES
	(1, 1, 3, 1, '你好，便宜一点', 'TEXT', 'READ', '2026-06-19 19:29:27'),
	(2, 1, 3, 1, '快点', 'TEXT', 'READ', '2026-06-19 19:29:38'),
	(3, 1, 1, 3, '你好', 'TEXT', 'READ', '2026-06-19 19:30:50'),
	(4, 1, 3, 1, '订单', 'TEXT', 'READ', '2026-06-19 19:32:45'),
	(5, 1, 1, 3, 'hi', 'TEXT', 'READ', '2026-06-19 19:32:51'),
	(6, 1, 3, 1, '订单', 'TEXT', 'READ', '2026-06-19 19:33:06'),
	(7, 1, 1, 3, '还在吗', 'TEXT', 'READ', '2026-06-19 20:33:46'),
	(8, 1, 1, 3, 'jo', 'TEXT', 'READ', '2026-06-19 23:48:19'),
	(9, 4, 7, 6, 'hello九尾卖吗', 'TEXT', 'READ', '2026-06-20 15:27:17'),
	(10, 4, 6, 7, '你好，不卖', 'TEXT', 'READ', '2026-06-20 15:27:44'),
	(11, 4, 6, 7, '999是防拍价', 'TEXT', 'READ', '2026-06-20 15:27:58'),
	(12, 4, 7, 6, '不卖拉倒，那你卖宋亚轩吗', 'TEXT', 'READ', '2026-06-20 15:28:31'),
	(13, 5, 7, 1, '怎么卖呀', 'TEXT', 'READ', '2026-06-20 20:44:09'),
	(14, 6, 2, 1, '你好，我是管理员', 'TEXT', 'READ', '2026-06-20 20:52:06'),
	(15, 7, 1, 6, '你好', 'TEXT', 'READ', '2026-06-20 22:02:42'),
	(16, 4, 6, 7, '下班', 'TEXT', 'READ', '2026-06-20 22:12:58'),
	(17, 4, 6, 7, '下班', 'TEXT', 'UNREAD', '2026-06-20 22:13:30'),
	(18, 7, 6, 1, 'nihao', 'TEXT', 'READ', '2026-06-21 16:49:47');

-- 导出  表 cau_used_goods.favorites 结构
CREATE TABLE IF NOT EXISTS `favorites` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
  `user_id` bigint unsigned NOT NULL COMMENT '收藏用户ID',
  `product_id` bigint unsigned NOT NULL COMMENT '被收藏商品ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否取消收藏',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_favorites_user_product` (`user_id`,`product_id`),
  KEY `idx_favorites_product` (`product_id`),
  CONSTRAINT `fk_favorites_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_favorites_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收藏表';

-- 正在导出表  cau_used_goods.favorites 的数据：~0 rows (大约)
INSERT INTO `favorites` (`id`, `user_id`, `product_id`, `create_time`, `is_deleted`) VALUES
	(1, 3, 4, '2026-06-19 20:00:20', 0),
	(2, 2, 6, '2026-06-20 14:58:10', 1),
	(3, 7, 4, '2026-06-20 20:43:58', 1),
	(4, 1, 15, '2026-06-21 17:05:49', 0),
	(5, 1, 12, '2026-06-21 17:06:28', 1);

-- 导出  表 cau_used_goods.messages 结构
CREATE TABLE IF NOT EXISTS `messages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '消息ID',
  `receiver_id` bigint unsigned NOT NULL COMMENT '接收用户ID',
  `sender_id` bigint unsigned DEFAULT NULL COMMENT '发送者ID，系统消息可为空',
  `message_type` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'ORDER_CREATED / ORDER_CONFIRMED / ORDER_CANCELED / ORDER_TIMEOUT / REPORT_HANDLED / SYSTEM_NOTICE',
  `title` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '消息标题',
  `content` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '消息内容',
  `related_type` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'ORDER / PRODUCT / REPORT / NOTICE',
  `related_id` bigint unsigned DEFAULT NULL COMMENT '关联对象ID',
  `read_status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'UNREAD' COMMENT 'UNREAD / READ',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `read_time` datetime DEFAULT NULL COMMENT '读取时间',
  PRIMARY KEY (`id`),
  KEY `idx_messages_receiver_read_time` (`receiver_id`,`read_status`,`create_time`),
  KEY `idx_messages_sender` (`sender_id`),
  KEY `idx_messages_related` (`related_type`,`related_id`),
  CONSTRAINT `fk_messages_receiver` FOREIGN KEY (`receiver_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_messages_sender` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=95 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='站内消息表';

-- 正在导出表  cau_used_goods.messages 的数据：~87 rows (大约)
INSERT INTO `messages` (`id`, `receiver_id`, `sender_id`, `message_type`, `title`, `content`, `related_type`, `related_id`, `read_status`, `create_time`, `read_time`) VALUES
	(1, 1, 2, 'SYSTEM_NOTICE', '学生认证审核结果', '你的学生认证已通过，现在可以正常使用发布、收藏、预约和举报等功能。', 'USER', 1, 'READ', '2026-06-19 18:22:40', '2026-06-19 18:23:19'),
	(2, 3, 2, 'SYSTEM_NOTICE', '学生认证审核结果', '你的学生认证已通过，现在可以正常使用发布、收藏、预约和举报等功能。', 'USER', 3, 'READ', '2026-06-19 19:28:36', '2026-06-21 13:50:59'),
	(3, 1, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号已被临时禁用，原因：恶意发布多个商品', 'USER', 1, 'READ', '2026-06-19 20:10:56', '2026-06-20 16:24:15'),
	(4, 1, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号状态已恢复正常，原因：申诉理由合理', 'USER', 1, 'READ', '2026-06-19 20:31:02', '2026-06-19 20:33:26'),
	(5, 1, 2, 'SYSTEM_NOTICE', '申诉处理结果', '你的申诉已处理，结果：APPROVED。处理说明：申诉已通过', 'APPEAL', 1, 'READ', '2026-06-19 20:31:10', '2026-06-21 16:15:41'),
	(6, 3, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号已被永久封禁，原因：测试', 'USER', 3, 'READ', '2026-06-19 23:06:59', '2026-06-21 20:26:32'),
	(7, 1, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号已被临时禁用，原因：hi', 'USER', 1, 'READ', '2026-06-19 23:52:02', '2026-06-20 15:12:01'),
	(8, 1, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号状态已恢复正常，原因：1', 'USER', 1, 'READ', '2026-06-20 00:11:47', '2026-06-20 15:11:50'),
	(9, 5, 2, 'SYSTEM_NOTICE', '学生认证审核结果', '你的学生认证已通过，现在可以正常使用发布、收藏、预约和举报等功能。', 'USER', 5, 'READ', '2026-06-20 14:13:32', '2026-06-20 15:30:34'),
	(10, 1, NULL, 'ORDER_CREATED', '新订单提醒', '您的商品「6.20测试」有新的订单，买家已预约，请尽快确认。', 'ORDER', 1, 'READ', '2026-06-20 14:15:20', '2026-06-20 15:20:53'),
	(11, 5, NULL, 'ORDER_CANCELED', '订单已取消', '订单「6.20测试」已被取消，原因：不想卖了', 'ORDER', 1, 'READ', '2026-06-20 14:21:47', '2026-06-20 14:26:34'),
	(12, 1, NULL, 'ORDER_CREATED', '新订单提醒', '您的商品「6.20测试」有新的订单，买家已预约，请尽快确认。', 'ORDER', 2, 'READ', '2026-06-20 14:29:46', '2026-06-20 15:18:38'),
	(13, 5, NULL, 'ORDER_CONFIRMED', '订单已确认', '卖家已确认您的订单「6.20测试」，请按约定时间地点交易。', 'ORDER', 2, 'READ', '2026-06-20 14:30:48', '2026-06-20 15:26:05'),
	(14, 5, NULL, 'ORDER_CONFIRMED', '交易完成', '订单「6.20测试」已完成交易，欢迎评价。', 'ORDER', 2, 'READ', '2026-06-20 14:31:42', '2026-06-20 15:27:30'),
	(15, 5, 2, 'SYSTEM_NOTICE', '申诉处理结果', '你的申诉已处理，结果：APPROVED。处理说明：申诉已通过', 'APPEAL', 3, 'READ', '2026-06-20 14:40:57', '2026-06-20 15:52:37'),
	(16, 6, 2, 'SYSTEM_NOTICE', '学生认证审核结果', '你的学生认证已通过，现在可以正常使用发布、收藏、预约和举报等功能。', 'USER', 6, 'READ', '2026-06-20 14:41:12', '2026-06-21 15:35:36'),
	(17, 7, 2, 'SYSTEM_NOTICE', '学生认证审核结果', '你的学生认证已通过，现在可以正常使用发布、收藏、预约和举报等功能。', 'USER', 7, 'READ', '2026-06-20 14:42:52', '2026-06-20 21:22:59'),
	(18, 6, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号已被临时禁用，原因：虚假宣传', 'USER', 6, 'READ', '2026-06-20 15:28:45', '2026-06-20 16:49:21'),
	(19, 4, 2, 'SYSTEM_NOTICE', '学生认证审核结果', '你的学生认证未通过。认证信息不符合要求', 'USER', 4, 'READ', '2026-06-20 16:51:27', '2026-06-20 17:05:36'),
	(20, 7, NULL, 'ORDER_CREATED', '新订单提醒', '您的商品「小米台灯」有新的订单，买家已预约，请尽快确认。', 'ORDER', 3, 'READ', '2026-06-20 17:07:27', '2026-06-20 17:15:27'),
	(21, 1, NULL, 'ORDER_CONFIRMED', '订单已确认', '卖家已确认您的订单「小米台灯」，请按约定时间地点交易。', 'ORDER', 3, 'READ', '2026-06-20 17:13:10', '2026-06-20 19:29:36'),
	(22, 1, NULL, 'ORDER_CONFIRMED', '交易完成', '订单「小米台灯」已完成交易，欢迎评价。', 'ORDER', 3, 'READ', '2026-06-20 17:13:22', '2026-06-20 19:29:36'),
	(23, 7, NULL, 'ORDER_CREATED', '新订单提醒', '您的商品「菠萝手机」有新的订单，买家已预约，请尽快确认。', 'ORDER', 4, 'READ', '2026-06-20 17:22:06', '2026-06-20 17:22:13'),
	(24, 1, 2, 'REPORT_HANDLED', '举报处理结果', '你的举报已处理，结果：REJECTED。处理说明：证据不足，无法认定违规', 'REPORT', 2, 'READ', '2026-06-20 18:11:00', '2026-06-20 19:29:36'),
	(25, 1, 2, 'SYSTEM_NOTICE', '订单异常关闭', '订单「菠萝手机」因相关账号状态异常，已由管理员异常关闭。', 'ORDER', 4, 'READ', '2026-06-20 18:28:35', '2026-06-20 19:29:36'),
	(26, 7, 2, 'SYSTEM_NOTICE', '订单异常关闭', '订单「菠萝手机」因相关账号状态异常，已由管理员异常关闭。', 'ORDER', 4, 'READ', '2026-06-20 18:28:35', '2026-06-20 18:35:21'),
	(27, 1, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号已被临时禁用，原因：可以', 'USER', 1, 'READ', '2026-06-20 18:28:35', '2026-06-20 19:30:58'),
	(28, 1, 2, 'SYSTEM_NOTICE', '申诉处理结果', '你的申诉已处理，结果：APPROVED。处理说明：申诉已通过', 'APPEAL', 2, 'READ', '2026-06-20 18:43:54', '2026-06-21 16:14:23'),
	(29, 7, 2, 'SYSTEM_NOTICE', '商品下架通知', '你的商品「菠萝手机」已被管理员下架。原因：管理员下架商品', 'PRODUCT', 8, 'READ', '2026-06-20 18:54:32', '2026-06-20 21:47:16'),
	(30, 1, NULL, 'ORDER_CREATED', '新订单提醒', '您的商品「上好佳薯片」有新的订单，买家已预约，请尽快确认。', 'ORDER', 5, 'READ', '2026-06-20 19:11:59', '2026-06-20 19:13:11'),
	(31, 1, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号状态已恢复正常，原因：你好', 'USER', 1, 'READ', '2026-06-20 19:12:34', '2026-06-20 20:06:08'),
	(32, 5, NULL, 'ORDER_CONFIRMED', '订单已确认', '卖家已确认您的订单「上好佳薯片」，请按约定时间地点交易。', 'ORDER', 5, 'UNREAD', '2026-06-20 19:13:28', NULL),
	(33, 3, 2, 'REPORT_HANDLED', '举报处理结果', '你的举报已处理，结果：REJECTED。处理说明：重复举报或恶意举报', 'REPORT', 1, 'READ', '2026-06-20 19:14:51', '2026-06-21 13:50:55'),
	(34, 5, 2, 'SYSTEM_NOTICE', '订单异常关闭', '订单「上好佳薯片」已由管理员异常关闭，原因：商品违规或信息异常', 'ORDER', 5, 'UNREAD', '2026-06-20 19:17:12', NULL),
	(35, 1, 2, 'SYSTEM_NOTICE', '订单异常关闭', '订单「上好佳薯片」已由管理员异常关闭，原因：商品违规或信息异常', 'ORDER', 5, 'READ', '2026-06-20 19:17:12', '2026-06-20 19:29:36'),
	(36, 5, 2, 'REPORT_HANDLED', '举报处理结果', '你的举报已处理，结果：APPROVED。处理说明：举报已通过', 'REPORT', 3, 'UNREAD', '2026-06-20 19:17:19', NULL),
	(37, 6, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号状态已恢复正常，原因：正常用户', 'USER', 6, 'READ', '2026-06-20 19:53:32', '2026-06-21 15:36:46'),
	(38, 1, 2, 'SYSTEM_NOTICE', '平台公告：即将上线！', '我们的小程序将于6.21号上线！敬请期待！', 'NOTICE', 1, 'READ', '2026-06-20 21:02:43', '2026-06-20 21:21:11'),
	(39, 3, 2, 'SYSTEM_NOTICE', '平台公告：即将上线！', '我们的小程序将于6.21号上线！敬请期待！', 'NOTICE', 1, 'READ', '2026-06-20 21:02:43', '2026-06-21 13:51:08'),
	(40, 4, 2, 'SYSTEM_NOTICE', '平台公告：即将上线！', '我们的小程序将于6.21号上线！敬请期待！', 'NOTICE', 1, 'UNREAD', '2026-06-20 21:02:43', NULL),
	(41, 5, 2, 'SYSTEM_NOTICE', '平台公告：即将上线！', '我们的小程序将于6.21号上线！敬请期待！', 'NOTICE', 1, 'UNREAD', '2026-06-20 21:02:43', NULL),
	(42, 6, 2, 'SYSTEM_NOTICE', '平台公告：即将上线！', '我们的小程序将于6.21号上线！敬请期待！', 'NOTICE', 1, 'READ', '2026-06-20 21:02:43', '2026-06-21 15:36:41'),
	(43, 7, 2, 'SYSTEM_NOTICE', '平台公告：即将上线！', '我们的小程序将于6.21号上线！敬请期待！', 'NOTICE', 1, 'READ', '2026-06-20 21:02:43', '2026-06-20 21:26:38'),
	(44, 8, 2, 'SYSTEM_NOTICE', '平台公告：即将上线！', '我们的小程序将于6.21号上线！敬请期待！', 'NOTICE', 1, 'READ', '2026-06-20 21:02:43', '2026-06-21 13:03:40'),
	(45, 9, 2, 'SYSTEM_NOTICE', '平台公告：即将上线！', '我们的小程序将于6.21号上线！敬请期待！', 'NOTICE', 1, 'UNREAD', '2026-06-20 21:02:43', NULL),
	(53, 6, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号已被临时禁用，原因：不喜欢', 'USER', 6, 'READ', '2026-06-20 21:17:39', '2026-06-20 21:30:23'),
	(54, 7, 2, 'SYSTEM_NOTICE', '商品上架通知', '你的商品「宋亚轩」已被管理员上架。原因：管理员上架商品', 'PRODUCT', 10, 'READ', '2026-06-20 21:31:05', '2026-06-20 21:46:39'),
	(55, 1, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号已被临时禁用，原因：不喜欢', 'USER', 1, 'READ', '2026-06-20 21:31:16', '2026-06-20 21:31:29'),
	(56, 7, 2, 'SYSTEM_NOTICE', '申诉处理结果', '你的申诉已处理，结果：APPROVED。处理说明：申诉已通过', 'APPEAL', 4, 'READ', '2026-06-20 21:31:37', '2026-06-20 21:36:27'),
	(57, 1, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号状态已恢复正常，原因：错误', 'USER', 1, 'READ', '2026-06-20 21:32:07', '2026-06-21 13:27:22'),
	(58, 7, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号已被临时禁用，原因：不喜欢', 'USER', 7, 'READ', '2026-06-20 21:32:18', '2026-06-20 21:32:31'),
	(59, 7, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号状态已恢复正常，原因：错误', 'USER', 7, 'READ', '2026-06-20 21:33:14', '2026-06-20 21:35:51'),
	(60, 7, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号已被临时禁用，原因：不喜欢', 'USER', 7, 'READ', '2026-06-20 21:35:22', '2026-06-20 21:35:54'),
	(61, 7, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号状态已恢复正常，原因：错误', 'USER', 7, 'READ', '2026-06-20 21:38:01', '2026-06-20 21:53:19'),
	(62, 6, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号状态已恢复正常，原因：错误', 'USER', 6, 'READ', '2026-06-20 21:38:06', '2026-06-20 21:41:27'),
	(63, 6, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号已被临时禁用，原因：色情', 'USER', 6, 'READ', '2026-06-20 21:39:22', '2026-06-20 21:41:32'),
	(64, 7, 2, 'REPORT_HANDLED', '举报处理结果', '你的举报已处理，结果：APPROVED。处理说明：举报已通过', 'REPORT', 4, 'READ', '2026-06-20 21:39:26', '2026-06-20 21:39:38'),
	(65, 6, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号状态已恢复正常，原因：错误', 'USER', 6, 'READ', '2026-06-20 21:40:03', '2026-06-21 15:36:36'),
	(66, 6, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号已被临时禁用，原因：不喜欢', 'USER', 6, 'READ', '2026-06-20 21:42:05', '2026-06-21 15:35:25'),
	(67, 6, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号状态已恢复正常，原因：错误', 'USER', 6, 'READ', '2026-06-20 21:44:47', '2026-06-21 15:39:36'),
	(68, 6, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号已被临时禁用，原因：不喜欢', 'USER', 6, 'READ', '2026-06-20 21:44:53', '2026-06-21 15:36:32'),
	(69, 6, 2, 'SYSTEM_NOTICE', '账号状态变更', '你的账号状态已恢复正常，原因：错误', 'USER', 6, 'READ', '2026-06-20 22:01:41', '2026-06-21 15:32:59'),
	(70, 8, 11, 'SYSTEM_NOTICE', '账号状态变更', '你的账号已被永久封禁，原因：违规', 'USER', 8, 'READ', '2026-06-21 12:58:01', '2026-06-21 13:10:47'),
	(71, 8, 11, 'SYSTEM_NOTICE', '申诉处理结果', '你的申诉已处理，结果：APPROVED。处理说明：申诉已通过', 'APPEAL', 5, 'READ', '2026-06-21 13:05:54', '2026-06-21 13:10:05'),
	(72, 8, 11, 'SYSTEM_NOTICE', '账号状态变更', '你的账号状态已恢复正常，原因：该用户不违规', 'USER', 8, 'READ', '2026-06-21 13:06:03', '2026-06-21 13:10:29'),
	(73, 3, 2, 'SYSTEM_NOTICE', '申诉处理结果', '你的申诉已处理，结果：APPROVED。处理说明：申诉已通过', 'APPEAL', 6, 'READ', '2026-06-21 13:25:45', '2026-06-21 13:51:06'),
	(74, 12, 2, 'SYSTEM_NOTICE', '学生认证审核结果', '你的学生认证已通过，现在可以正常使用发布、收藏、预约和举报等功能。', 'USER', 12, 'READ', '2026-06-21 14:46:40', '2026-06-21 14:47:40'),
	(75, 12, 2, 'SYSTEM_NOTICE', '商品下架通知', '你的商品「小猫玩偶」已被管理员下架。原因：管理员下架商品', 'PRODUCT', 14, 'READ', '2026-06-21 14:48:05', '2026-06-21 14:55:16'),
	(76, 12, 2, 'SYSTEM_NOTICE', '申诉处理结果', '你的申诉已处理，结果：APPROVED。处理说明：申诉已通过', 'APPEAL', 8, 'READ', '2026-06-21 14:51:06', '2026-06-21 15:13:35'),
	(77, 6, NULL, 'ORDER_CREATED', '新订单提醒', '您的商品「九尾」有新的订单，买家已预约，请尽快确认。', 'ORDER', 6, 'READ', '2026-06-21 15:15:11', '2026-06-21 15:19:17'),
	(78, 12, NULL, 'ORDER_CONFIRMED', '订单已确认', '卖家已确认您的订单「九尾」，请按约定时间地点交易。', 'ORDER', 6, 'READ', '2026-06-21 15:16:28', '2026-06-21 16:57:19'),
	(79, 12, 2, 'SYSTEM_NOTICE', '订单异常关闭', '订单「九尾」已由管理员异常关闭，原因：商品违规或信息异常', 'ORDER', 6, 'READ', '2026-06-21 15:17:05', '2026-06-21 16:57:14'),
	(80, 6, 2, 'SYSTEM_NOTICE', '订单异常关闭', '订单「九尾」已由管理员异常关闭，原因：商品违规或信息异常', 'ORDER', 6, 'READ', '2026-06-21 15:17:05', '2026-06-21 15:17:27'),
	(81, 1, 2, 'REPORT_HANDLED', '举报处理结果', '你的举报已处理，结果：REJECTED。处理说明：举报内容与对象不符', 'REPORT', 5, 'READ', '2026-06-21 15:18:54', '2026-06-21 16:17:38'),
	(82, 6, 2, 'SYSTEM_NOTICE', '申诉处理结果', '你的申诉已处理，结果：REJECTED。处理说明：未提供有效证明', 'APPEAL', 9, 'READ', '2026-06-21 15:22:10', '2026-06-21 16:04:35'),
	(83, 6, NULL, 'ORDER_CREATED', '新订单提醒', '您的商品「九尾」有新的订单，买家已预约，请尽快确认。', 'ORDER', 7, 'READ', '2026-06-21 16:05:35', '2026-06-21 16:39:36'),
	(84, 6, NULL, 'ORDER_CANCELED', '订单已取消', '订单「九尾」已被取消，原因：用户取消', 'ORDER', 7, 'READ', '2026-06-21 16:16:12', '2026-06-21 16:39:31'),
	(85, 6, NULL, 'ORDER_CREATED', '新订单提醒', '您的商品「九尾」有新的订单，买家已预约，请尽快确认。', 'ORDER', 8, 'READ', '2026-06-21 16:16:31', '2026-06-21 16:39:26'),
	(86, 1, NULL, 'ORDER_CREATED', '新订单提醒', '您的商品「动物园门票」有新的订单，买家已预约，请尽快确认。', 'ORDER', 9, 'READ', '2026-06-21 16:41:00', '2026-06-21 18:30:24'),
	(87, 1, NULL, 'ORDER_CONFIRMED', '订单已确认', '卖家已确认您的订单「九尾」，请按约定时间地点交易。', 'ORDER', 8, 'READ', '2026-06-21 16:41:10', '2026-06-21 18:29:47'),
	(88, 12, NULL, 'ORDER_CREATED', '新订单提醒', '您的商品「电子」有新的订单，买家已预约，请尽快确认。', 'ORDER', 10, 'READ', '2026-06-21 17:00:10', '2026-06-21 17:24:47'),
	(89, 6, NULL, 'ORDER_CONFIRMED', '订单已确认', '卖家已确认您的订单「动物园门票」，请按约定时间地点交易。', 'ORDER', 9, 'READ', '2026-06-21 18:30:31', '2026-06-21 19:53:43'),
	(90, 6, NULL, 'ORDER_CREATED', '新订单提醒', '您的商品「九尾」有新的订单，买家已预约，请尽快确认。', 'ORDER', 11, 'READ', '2026-06-21 20:05:26', '2026-06-21 20:09:57'),
	(91, 6, NULL, 'ORDER_CANCELED', '订单已取消', '订单「九尾」已被取消，原因：用户取消', 'ORDER', 11, 'UNREAD', '2026-06-21 20:15:48', NULL),
	(92, 6, NULL, 'ORDER_CREATED', '新订单提醒', '您的商品「九尾」有新的订单，买家已预约，请尽快确认。', 'ORDER', 12, 'UNREAD', '2026-06-21 20:16:06', NULL),
	(93, 3, 11, 'SYSTEM_NOTICE', '账号状态变更', '你的账号状态已恢复正常，原因：好吧', 'USER', 3, 'UNREAD', '2026-06-21 20:39:33', NULL),
	(94, 3, 11, 'SYSTEM_NOTICE', '申诉处理结果', '你的申诉已处理，结果：APPROVED。处理说明：申诉已通过', 'APPEAL', 7, 'UNREAD', '2026-06-21 20:39:43', NULL);

-- 导出  表 cau_used_goods.orders 结构
CREATE TABLE IF NOT EXISTS `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '订单编号',
  `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
  `buyer_id` bigint unsigned NOT NULL COMMENT '买家用户ID',
  `seller_id` bigint unsigned NOT NULL COMMENT '卖家用户ID',
  `product_title_snapshot` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '下单时商品标题快照',
  `product_price_snapshot` decimal(10,2) NOT NULL COMMENT '下单时商品价格快照',
  `status` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'PENDING_CONFIRM' COMMENT 'PENDING_CONFIRM / WAIT_MEET / COMPLETED / CANCELED / EXCEPTION_CLOSED',
  `remark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '买家备注',
  `meet_time` datetime DEFAULT NULL COMMENT '约定面交时间',
  `meet_location` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '约定面交地点',
  `cancel_reason` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '取消原因',
  `cancel_by` bigint unsigned DEFAULT NULL COMMENT '取消操作人ID',
  `expire_time` datetime NOT NULL COMMENT '卖家确认截止时间',
  `confirm_time` datetime DEFAULT NULL COMMENT '卖家确认时间',
  `finish_time` datetime DEFAULT NULL COMMENT '交易完成时间',
  `close_time` datetime DEFAULT NULL COMMENT '取消或异常关闭时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '下单时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_orders_order_no` (`order_no`),
  KEY `idx_orders_product_status` (`product_id`,`status`),
  KEY `idx_orders_buyer_status_time` (`buyer_id`,`status`,`create_time`),
  KEY `idx_orders_seller_status_time` (`seller_id`,`status`,`create_time`),
  KEY `idx_orders_expire` (`status`,`expire_time`),
  KEY `idx_orders_cancel_by` (`cancel_by`),
  CONSTRAINT `fk_orders_buyer` FOREIGN KEY (`buyer_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_orders_cancel_by` FOREIGN KEY (`cancel_by`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_orders_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_orders_seller` FOREIGN KEY (`seller_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- 正在导出表  cau_used_goods.orders 的数据：~0 rows (大约)
INSERT INTO `orders` (`id`, `order_no`, `product_id`, `buyer_id`, `seller_id`, `product_title_snapshot`, `product_price_snapshot`, `status`, `remark`, `meet_time`, `meet_location`, `cancel_reason`, `cancel_by`, `expire_time`, `confirm_time`, `finish_time`, `close_time`, `create_time`, `update_time`) VALUES
	(1, 'O20260620141519371500', 5, 5, 1, '6.20测试', 10.00, 'CANCELED', '', '2026-06-20 17:16:00', '828', '不想卖了', 1, '2026-06-21 14:15:19', NULL, NULL, '2026-06-20 14:21:47', '2026-06-20 14:15:20', '2026-06-20 14:21:47'),
	(2, 'O20260620142946250800', 5, 5, 1, '6.20测试', 10.00, 'COMPLETED', '', '2026-06-20 15:31:00', '211', NULL, NULL, '2026-06-21 14:29:46', '2026-06-20 14:30:48', '2026-06-20 14:31:42', NULL, '2026-06-20 14:29:46', '2026-06-20 14:31:42'),
	(3, 'O20260620170727814100', 7, 1, 7, '小米台灯', 198.00, 'COMPLETED', '', '2026-06-20 18:07:00', '图书馆', NULL, NULL, '2026-06-21 17:07:27', '2026-06-20 17:13:10', '2026-06-20 17:13:22', NULL, '2026-06-20 17:07:27', '2026-06-20 17:13:22'),
	(4, 'O20260620172206912600', 8, 1, 7, '菠萝手机', 1999.00, 'EXCEPTION_CLOSED', '', '2026-06-20 18:22:00', '预约后协商', '可以', 2, '2026-06-21 17:22:06', NULL, NULL, '2026-06-20 18:28:35', '2026-06-20 17:22:06', '2026-06-20 18:28:35'),
	(5, 'O20260620191159134100', 4, 5, 1, '上好佳薯片', 250.00, 'EXCEPTION_CLOSED', '', '2026-06-20 20:11:00', '东区A座楼下', '商品违规或信息异常', 2, '2026-06-21 19:11:59', '2026-06-20 19:13:28', NULL, '2026-06-20 19:17:12', '2026-06-20 19:11:59', '2026-06-20 19:17:12'),
	(6, 'O20260621151511589600', 13, 12, 6, '九尾', 99.00, 'EXCEPTION_CLOSED', '', '2026-06-21 16:15:00', '东区图书馆门口', '商品违规或信息异常', 2, '2026-06-22 15:15:11', '2026-06-21 15:16:28', NULL, '2026-06-21 15:17:05', '2026-06-21 15:15:11', '2026-06-21 15:17:05'),
	(7, 'O20260621160535251000', 12, 1, 6, '九尾', 99.00, 'CANCELED', '', '2026-06-21 17:05:00', '东区图书馆门口', '用户取消', 1, '2026-06-22 16:05:35', NULL, NULL, '2026-06-21 16:16:12', '2026-06-21 16:05:35', '2026-06-21 16:16:12'),
	(8, 'O20260621161631127600', 12, 1, 6, '九尾', 99.00, 'WAIT_MEET', '', '2026-06-21 17:16:00', '东区图书馆门口', NULL, NULL, '2026-06-22 16:16:31', '2026-06-21 16:41:10', NULL, NULL, '2026-06-21 16:16:31', '2026-06-21 16:41:10'),
	(9, 'O20260621164100071300', 3, 6, 1, '动物园门票', 150.00, 'WAIT_MEET', '', '2026-06-21 17:40:00', '东区公一门口', NULL, NULL, '2026-06-22 16:41:00', '2026-06-21 18:30:31', NULL, NULL, '2026-06-21 16:41:00', '2026-06-21 18:30:31'),
	(10, 'O20260621170010147700', 15, 1, 12, '电子', 50.00, 'PENDING_CONFIRM', '', '2026-06-21 18:00:00', '预约后协商', NULL, NULL, '2026-06-22 17:00:10', NULL, NULL, NULL, '2026-06-21 17:00:10', '2026-06-21 17:00:10'),
	(11, 'O20260621200526658200', 13, 12, 6, '九尾', 99.00, 'CANCELED', '', '2026-06-21 21:05:00', '东区图书馆门口', '用户取消', 12, '2026-06-22 20:05:26', NULL, NULL, '2026-06-21 20:15:48', '2026-06-21 20:05:26', '2026-06-21 20:15:48'),
	(12, 'O20260621201606580400', 13, 12, 6, '九尾', 99.00, 'PENDING_CONFIRM', '', '2026-06-21 21:16:00', '东区图书馆门口', NULL, NULL, '2026-06-22 20:16:06', NULL, NULL, NULL, '2026-06-21 20:16:06', '2026-06-21 20:16:06');

-- 导出  表 cau_used_goods.products 结构
CREATE TABLE IF NOT EXISTS `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  `seller_id` bigint unsigned NOT NULL COMMENT '卖家用户ID',
  `category_id` bigint unsigned NOT NULL COMMENT '商品分类ID',
  `title` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '商品标题',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '商品描述',
  `original_price` decimal(10,2) DEFAULT NULL COMMENT '商品原价',
  `price` decimal(10,2) NOT NULL COMMENT '商品售价',
  `condition_level` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '商品成色',
  `meet_location` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '建议面交地点',
  `status` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ON_SALE' COMMENT 'ON_SALE / LOCKED / SOLD / OFF_SHELF / DELETED',
  `view_count` int NOT NULL DEFAULT '0' COMMENT '浏览量',
  `favorite_count` int NOT NULL DEFAULT '0' COMMENT '收藏数',
  `off_shelf_reason` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '下架原因',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '逻辑删除标记',
  PRIMARY KEY (`id`),
  KEY `idx_products_seller` (`seller_id`),
  KEY `idx_products_category_status` (`category_id`,`status`),
  KEY `idx_products_status_time` (`status`,`create_time`),
  KEY `idx_products_price` (`price`),
  KEY `idx_products_title` (`title`),
  CONSTRAINT `fk_products_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`),
  CONSTRAINT `fk_products_seller` FOREIGN KEY (`seller_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- 正在导出表  cau_used_goods.products 的数据：~9 rows (大约)
INSERT INTO `products` (`id`, `seller_id`, `category_id`, `title`, `description`, `original_price`, `price`, `condition_level`, `meet_location`, `status`, `view_count`, `favorite_count`, `off_shelf_reason`, `create_time`, `update_time`, `is_deleted`) VALUES
	(1, 1, 36, '小米台灯', '使用过一次，其它全新', 60.00, 30.00, '九成新', '东区A座楼下', 'OFF_SHELF', 2, 0, 'hi', '2026-06-19 18:59:03', '2026-06-19 23:52:02', 0),
	(2, 1, 15, '2027考研408王道书', '没有使用过，全新', 240.00, 230.00, '全新', '东区图书馆门口', 'OFF_SHELF', 2, 0, 'hi', '2026-06-19 19:01:10', '2026-06-19 23:52:02', 0),
	(3, 1, 63, '动物园门票', '买了票但有事去不了', 150.00, 150.00, '全新', '东区公一门口', 'LOCKED', 3, 0, '', '2026-06-19 19:03:31', '2026-06-21 16:41:00', 0),
	(4, 1, 42, '上好佳薯片', '非常好吃', 300.00, 250.00, '全新', '东区A座楼下', 'OFF_SHELF', 7, 1, '不喜欢', '2026-06-19 19:30:13', '2026-06-20 21:31:16', 0),
	(5, 1, 15, '6.20测试', '111', 10.00, 10.00, '全新', '211', 'SOLD', 2, 0, NULL, '2026-06-20 00:26:01', '2026-06-20 14:31:42', 0),
	(6, 6, 69, '全新九尾', '全新九尾，代打王者', 0.00, 999.00, '全新', '东区图书馆门口', 'ON_SALE', 1, 0, '', '2026-06-20 14:44:11', '2026-06-21 20:04:38', 0),
	(7, 7, 32, '小米台灯', '纯白智能调节光感', 398.00, 198.00, '九成新', '图书馆', 'SOLD', 1, 0, NULL, '2026-06-20 14:44:18', '2026-06-20 17:13:22', 0),
	(8, 7, 23, '菠萝手机', '完美无暇', 3999.00, 1999.00, '九成新', '', 'OFF_SHELF', 1, 0, '管理员下架商品', '2026-06-20 17:21:40', '2026-06-20 18:54:32', 0),
	(9, 7, 34, 'iii', '超级好用', 99.00, 10.00, '八成新', '', 'OFF_SHELF', 1, 0, '管理员下架商品', '2026-06-20 17:23:26', '2026-06-20 18:10:45', 0),
	(10, 7, 34, '宋亚轩', '超级好用', 999.00, 100.00, '九成新', '', 'OFF_SHELF', 0, 0, '不喜欢', '2026-06-20 18:36:37', '2026-06-20 21:32:18', 0),
	(11, 7, 40, '大宝贝', '超级好用', 1999.00, 19.90, '九成新', '', 'OFF_SHELF', 0, 0, '不喜欢', '2026-06-20 21:35:16', '2026-06-20 21:35:22', 0),
	(12, 6, 13, '九尾', '全新', 0.00, 99.00, '全新', '东区图书馆门口', 'LOCKED', 1, 0, '', '2026-06-20 21:41:20', '2026-06-21 17:06:29', 0),
	(13, 6, 13, '九尾', '全新', 0.00, 99.00, '全新', '东区图书馆门口', 'LOCKED', 7, 0, '', '2026-06-20 22:02:24', '2026-06-21 20:16:06', 0),
	(14, 12, 33, '小猫玩偶', '很可爱', 0.00, 15.00, '全新', '', 'OFF_SHELF', 0, 0, '管理员下架商品', '2026-06-21 14:47:31', '2026-06-21 14:48:05', 0),
	(15, 12, 23, '电子', '好用', 0.00, 50.00, '全新', '', 'LOCKED', 1, 1, NULL, '2026-06-21 16:57:03', '2026-06-21 17:05:49', 0);

-- 导出  表 cau_used_goods.product_images 结构
CREATE TABLE IF NOT EXISTS `product_images` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '图片ID',
  `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
  `image_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '图片访问地址',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '图片排序',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  PRIMARY KEY (`id`),
  KEY `idx_product_images_product_sort` (`product_id`,`sort_order`),
  CONSTRAINT `fk_product_images_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品图片表';

-- 正在导出表  cau_used_goods.product_images 的数据：~3 rows (大约)
INSERT INTO `product_images` (`id`, `product_id`, `image_url`, `sort_order`, `create_time`) VALUES
	(1, 1, '/uploads/products/1781866662439693900.jpg', 0, '2026-06-19 18:59:03'),
	(2, 1, '/uploads/products/1781866672143226900.jpg', 1, '2026-06-19 18:59:03'),
	(3, 1, '/uploads/products/1781866685043776500.jpg', 2, '2026-06-19 18:59:03'),
	(4, 2, '/uploads/products/1781866784657817400.jpg', 0, '2026-06-19 19:01:10'),
	(5, 2, '/uploads/products/1781866791971667900.jpg', 1, '2026-06-19 19:01:10'),
	(12, 4, '/uploads/products/1781868436488744200.jpg', 0, '2026-06-19 19:30:13'),
	(13, 5, '/uploads/products/1781886334065998600.jpg', 0, '2026-06-20 00:26:01'),
	(14, 6, '/uploads/products/1781937766852845900.jpg', 0, '2026-06-20 14:44:11'),
	(15, 6, '/uploads/products/1781937771161117500.jpg', 1, '2026-06-20 14:44:11'),
	(16, 6, '/uploads/products/1781937775944620600.jpg', 2, '2026-06-20 14:44:11'),
	(17, 6, '/uploads/products/1781937779638705800.jpg', 3, '2026-06-20 14:44:11'),
	(18, 7, '/uploads/products/1781937807389541300.jpg', 0, '2026-06-20 14:44:18'),
	(19, 8, '/uploads/products/1781947255628113600.jpg', 0, '2026-06-20 17:21:40'),
	(20, 9, '/uploads/products/1781947387711420200.jpg', 0, '2026-06-20 17:23:26'),
	(21, 10, '/uploads/products/1781951760948752400.jpg', 0, '2026-06-20 18:36:37'),
	(22, 11, '/uploads/products/1781962484804823700.jpg', 0, '2026-06-20 21:35:16'),
	(23, 12, '/uploads/products/1781962859150296100.jpg', 0, '2026-06-20 21:41:20'),
	(24, 13, '/uploads/products/1781964121932607300.jpg', 0, '2026-06-20 22:02:24'),
	(25, 14, '/uploads/products/1782024428345848600.jpg', 0, '2026-06-21 14:47:31'),
	(26, 3, 'http://62.234.163.176:7001/uploads/products/1781867153105830700.jpg', 0, '2026-06-21 16:29:27'),
	(27, 3, 'http://62.234.163.176:7001/uploads/products/1781867172056907800.jpg', 1, '2026-06-21 16:29:27'),
	(28, 3, '/uploads/products/1782030560884297500.jpg', 2, '2026-06-21 16:29:27'),
	(29, 15, '/uploads/products/1782032190956683700.jpg', 0, '2026-06-21 16:57:03');

-- 导出  表 cau_used_goods.reports 结构
CREATE TABLE IF NOT EXISTS `reports` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '举报ID',
  `reporter_id` bigint unsigned NOT NULL COMMENT '举报人ID',
  `target_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'PRODUCT / USER / ORDER',
  `target_id` bigint unsigned NOT NULL COMMENT '被举报对象ID',
  `reason_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '举报原因类型',
  `description` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '举报说明',
  `status` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'PENDING' COMMENT 'PENDING / PROCESSING / APPROVED / REJECTED / CLOSED',
  `handle_result` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '处理结果说明',
  `handler_id` bigint unsigned DEFAULT NULL COMMENT '处理管理员ID',
  `handle_time` datetime DEFAULT NULL COMMENT '处理时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '举报提交时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_reports_reporter` (`reporter_id`),
  KEY `idx_reports_target` (`target_type`,`target_id`),
  KEY `idx_reports_status_time` (`status`,`create_time`),
  KEY `idx_reports_handler` (`handler_id`),
  CONSTRAINT `fk_reports_handler` FOREIGN KEY (`handler_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_reports_reporter` FOREIGN KEY (`reporter_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='举报表';

-- 正在导出表  cau_used_goods.reports 的数据：~2 rows (大约)
INSERT INTO `reports` (`id`, `reporter_id`, `target_type`, `target_id`, `reason_type`, `description`, `status`, `handle_result`, `handler_id`, `handle_time`, `create_time`, `update_time`) VALUES
	(1, 3, 'USER', 1, 'TRADE_DISPUTE', '看你不顺眼', 'REJECTED', '重复举报或恶意举报', 2, '2026-06-20 19:14:51', '2026-06-19 20:37:26', '2026-06-20 19:14:51'),
	(2, 1, 'USER', 7, 'FAKE_PRODUCT', '有问题', 'REJECTED', '证据不足，无法认定违规', 2, '2026-06-20 18:11:00', '2026-06-20 18:10:28', '2026-06-20 18:11:00'),
	(3, 5, 'ORDER', 5, 'INAPPROPRIATE_CONTENT', '很不当啊', 'APPROVED', '举报已通过', 2, '2026-06-20 19:17:19', '2026-06-20 19:16:01', '2026-06-20 19:17:19'),
	(4, 7, 'USER', 6, 'INAPPROPRIATE_CONTENT', '提供色情服务', 'APPROVED', '举报已通过', 2, '2026-06-20 21:39:26', '2026-06-20 21:38:40', '2026-06-20 21:39:26'),
	(5, 1, 'PRODUCT', 13, 'FAKE_PRODUCT', '商品感觉有问题', 'REJECTED', '举报内容与对象不符', 2, '2026-06-21 15:18:54', '2026-06-21 13:40:59', '2026-06-21 15:18:54');

-- 导出  表 cau_used_goods.report_images 结构
CREATE TABLE IF NOT EXISTS `report_images` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '凭证图片ID',
  `report_id` bigint unsigned NOT NULL COMMENT '举报ID',
  `image_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '凭证图片访问地址',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '图片排序',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  PRIMARY KEY (`id`),
  KEY `idx_report_images_report_sort` (`report_id`,`sort_order`),
  CONSTRAINT `fk_report_images_report` FOREIGN KEY (`report_id`) REFERENCES `reports` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='举报凭证图片表';

-- 正在导出表  cau_used_goods.report_images 的数据：~0 rows (大约)
INSERT INTO `report_images` (`id`, `report_id`, `image_url`, `sort_order`, `create_time`) VALUES
	(1, 3, '/uploads/products/1781954159572359100.png', 0, '2026-06-20 19:16:01');

-- 导出  表 cau_used_goods.reviews 结构
CREATE TABLE IF NOT EXISTS `reviews` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '评价ID',
  `order_id` bigint unsigned NOT NULL COMMENT '订单ID',
  `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
  `reviewer_id` bigint unsigned NOT NULL COMMENT '评价人ID',
  `seller_id` bigint unsigned NOT NULL COMMENT '被评价卖家ID',
  `rating` int NOT NULL COMMENT '星级评分，取值1-5',
  `content` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文字评价',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'NORMAL' COMMENT 'NORMAL / HIDDEN / DELETED',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评价时间',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '逻辑删除标记',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reviews_order_reviewer` (`order_id`,`reviewer_id`),
  KEY `idx_reviews_product` (`product_id`),
  KEY `idx_reviews_seller_status` (`seller_id`,`status`),
  KEY `idx_reviews_reviewer` (`reviewer_id`),
  CONSTRAINT `fk_reviews_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
  CONSTRAINT `fk_reviews_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_reviews_reviewer` FOREIGN KEY (`reviewer_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_reviews_seller` FOREIGN KEY (`seller_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评价表';

-- 正在导出表  cau_used_goods.reviews 的数据：~0 rows (大约)

-- 导出  表 cau_used_goods.sensitive_words 结构
CREATE TABLE IF NOT EXISTS `sensitive_words` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '敏感词ID',
  `word` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '敏感词内容',
  `word_type` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FORBIDDEN / RISK',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ENABLED' COMMENT 'ENABLED / DISABLED',
  `create_by` bigint unsigned DEFAULT NULL COMMENT '创建管理员ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sensitive_words_word` (`word`),
  KEY `idx_sensitive_words_status` (`status`),
  KEY `idx_sensitive_words_create_by` (`create_by`),
  CONSTRAINT `fk_sensitive_words_create_by` FOREIGN KEY (`create_by`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=92 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='敏感词表';

-- 正在导出表  cau_used_goods.sensitive_words 的数据：~91 rows (大约)
INSERT INTO `sensitive_words` (`id`, `word`, `word_type`, `status`, `create_by`, `create_time`, `update_time`) VALUES
	(1, '毒品', 'FORBIDDEN', 'ENABLED', 5, '2026-06-16 22:10:11', '2026-06-16 22:10:11'),
	(2, '私下交易', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(3, '线下转账', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(4, '先款后货', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(5, '先付定金', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(6, '提前打款', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(7, '加微信交易', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(8, '加QQ交易', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(9, '绕过平台', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(10, '不走平台', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(11, '保证赚钱', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(12, '稳赚不赔', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(13, '内部渠道', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(14, '低价代购', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(15, '刷单', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(16, '刷信誉', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(17, '套现', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(18, '洗钱', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(19, '跑分', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(20, '薅羊毛项目', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(21, '身份证出售', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(22, '学生证出售', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(23, '校园卡出售', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(24, '银行卡出售', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(25, '电话卡出售', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(26, '手机号出售', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(27, '账号买卖', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(28, '游戏账号买卖', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(29, '代实名', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(30, '代认证', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(31, '个人信息出售', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(32, '隐私照片', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(33, '代写论文', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(34, '代写作业', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(35, '代做作业', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(36, '代考试', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(37, '替考', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(38, '论文代发', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(39, '毕业设计代做', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(40, '实验报告代写', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(41, '课程设计代做', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(42, '包过考试', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(43, '答案出售', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(44, '考试答案', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(45, '处方药', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(46, '违禁药品', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(47, '管制刀具', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(48, '仿真枪', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(49, '烟草', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(50, '香烟', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(51, '电子烟', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(52, '酒水转让', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(53, '白酒转让', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(54, '彩票', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(55, '博彩', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(56, '赌博', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(57, '成人用品', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(58, '危险化学品', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(59, '宠物活体', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(60, '盗版软件', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(61, '破解软件', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(62, '破解版', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(63, '盗版课程', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(64, '破解网课', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(65, '盗版电子书', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(66, '盗版教材', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(67, '外挂', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(68, '游戏外挂', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(69, '盗号工具', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(70, '校园贷', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(71, '高利贷', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(72, '裸贷', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(73, '借贷中介', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(74, '贷款套现', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(75, '信用卡套现', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(76, '包过', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(77, '全网最低', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(78, '假一赔十', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(79, '无理由退款保证', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(80, '绝对正品', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(81, '官方内部价', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(82, '低价秒杀', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(83, '限时返利', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(84, '傻逼', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(85, '垃圾人', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(86, '滚蛋', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(87, '骗子死全家', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(88, '脑残', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(89, '废物', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(90, '去死', 'FORBIDDEN', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46'),
	(91, '辱骂', 'RISK', 'ENABLED', NULL, '2026-06-19 14:35:46', '2026-06-19 14:35:46');

-- 导出  表 cau_used_goods.users 结构
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户主键ID',
  `openid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '微信用户唯一标识',
  `nickname` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '微信昵称或用户昵称',
  `avatar_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '微信头像地址',
  `student_id` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '学号',
  `real_name` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '学生真实姓名',
  `college` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '学院信息',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '联系方式',
  `role` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'USER' COMMENT 'USER / ADMIN / SUPER_ADMIN',
  `auth_status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'UNVERIFIED' COMMENT 'UNVERIFIED / PENDING / VERIFIED / REJECTED',
  `account_status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'NORMAL' COMMENT 'NORMAL / DISABLED / BANNED / CANCELED',
  `last_login_time` datetime DEFAULT NULL COMMENT '最近登录时间',
  `token_version` int NOT NULL DEFAULT '0' COMMENT '访问凭证版本，递增后旧JWT失效',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '逻辑删除标记',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_openid` (`openid`),
  UNIQUE KEY `uk_users_student_id` (`student_id`),
  KEY `idx_users_role` (`role`),
  KEY `idx_users_auth_account` (`auth_status`,`account_status`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 正在导出表  cau_used_goods.users 的数据：~12 rows (大约)
INSERT INTO `users` (`id`, `openid`, `nickname`, `avatar_url`, `student_id`, `real_name`, `college`, `phone`, `role`, `auth_status`, `account_status`, `last_login_time`, `token_version`, `create_time`, `update_time`, `is_deleted`) VALUES
	(1, 'frontend_a_dev_user', '王双媛', '/uploads/avatar/13e94b2aa9c9f3204cea52294109e37c.png', '2023308250232', '王双媛', '信息与电气工程学院', '13893530828', 'USER', 'VERIFIED', 'NORMAL', '2026-06-21 20:59:28', 0, '2026-06-19 17:31:26', '2026-06-21 20:59:28', 0),
	(2, 'admin_dev_user', '1号客服', '/uploads/avatar/e4da6967ca52535f104eae547fac792e.png', NULL, NULL, NULL, '13458932400', 'ADMIN', 'VERIFIED', 'NORMAL', '2026-06-21 20:58:32', 0, '2026-06-19 17:31:43', '2026-06-21 20:58:32', 0),
	(3, 'frontend_b_dev_user', '税馨乐', '/uploads/avatar/e90743f6d052eb23126bad4ab8e6ba2b.png', '2023308250229', '税馨乐', '信息与电气工程学院', '13568359833', 'USER', 'VERIFIED', 'NORMAL', '2026-06-21 20:46:43', 0, '2026-06-19 18:31:23', '2026-06-21 20:46:43', 0),
	(4, 'pending_dev_user', NULL, NULL, '2023308250200', '未认证', '信电', NULL, 'USER', 'REJECTED', 'NORMAL', '2026-06-20 17:05:29', 0, '2026-06-19 19:49:28', '2026-06-20 17:05:29', 0),
	(5, 'c-user', '愿欢', NULL, '2025523408232', '王双媛', '信电', '15039718898', 'USER', 'VERIFIED', 'NORMAL', '2026-06-20 19:11:46', 0, '2026-06-20 14:11:44', '2026-06-20 19:11:46', 0),
	(6, 'ysy', 'ysy', '/uploads/avatar/dd2fd4a11c16d694118bd604f150ebf2.jpg', '2023308250215', '尹双雨', '信电', '18255864872', 'USER', 'VERIFIED', 'NORMAL', '2026-06-21 20:09:47', 0, '2026-06-20 14:35:49', '2026-06-21 20:09:47', 0),
	(7, 'dx', 'queena', '/uploads/avatar/a163b3dccbed99c7253e3d6381056f8b.png', '2023308250208', '杜逊', '信电', '18743648176', 'USER', 'VERIFIED', 'NORMAL', '2026-06-20 22:07:42', 0, '2026-06-20 14:41:05', '2026-06-20 22:07:42', 0),
	(8, 'dev_openid_001', NULL, NULL, NULL, NULL, NULL, NULL, 'USER', 'UNVERIFIED', 'NORMAL', '2026-06-21 13:09:55', 0, '2026-06-20 15:15:40', '2026-06-21 13:09:55', 0),
	(9, 'ysy2', NULL, NULL, NULL, NULL, NULL, NULL, 'USER', 'UNVERIFIED', 'NORMAL', '2026-06-20 17:00:08', 0, '2026-06-20 17:00:08', '2026-06-20 17:00:08', 0),
	(11, 'super_admin_dev_user', '超级管理员', '/uploads/avatar/f5bcdf3d85e9ff86751e029f98e477f0.png', NULL, NULL, NULL, NULL, 'SUPER_ADMIN', 'UNVERIFIED', 'NORMAL', '2026-06-21 20:31:39', 0, '2026-06-21 12:55:58', '2026-06-21 20:31:39', 0),
	(12, 'wky', 'wky', '/uploads/avatar/c56024f84cf40b5ea5478a42fdeaf9af.jpg', '2023308250219', '王珂雅', '信电学院', NULL, 'USER', 'VERIFIED', 'NORMAL', '2026-06-21 20:14:58', 0, '2026-06-21 14:45:48', '2026-06-21 20:14:58', 0),
	(13, 'o-oxI3YG-EVs19t5iPnl4SfUVJtE', '求个好运气', '/uploads/avatar/cfc678162d7c9fd563b2f0285c9715b4.jpeg', NULL, NULL, NULL, '13568359833', 'USER', 'UNVERIFIED', 'NORMAL', '2026-06-21 20:09:02', 2, '2026-06-21 19:37:37', '2026-06-21 20:09:02', 0);

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
