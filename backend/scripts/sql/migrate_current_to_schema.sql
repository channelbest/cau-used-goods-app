USE `cau_used_goods`;

DROP PROCEDURE IF EXISTS `add_column_if_missing`;
DROP PROCEDURE IF EXISTS `add_index_if_missing`;

DELIMITER //

CREATE PROCEDURE `add_column_if_missing`(
  IN p_table_name VARCHAR(64),
  IN p_column_name VARCHAR(64),
  IN p_column_definition TEXT,
  IN p_after_column VARCHAR(64)
)
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = p_table_name
      AND COLUMN_NAME = p_column_name
  ) THEN
    SET @ddl = CONCAT(
      'ALTER TABLE `', p_table_name, '` ADD COLUMN `', p_column_name, '` ',
      p_column_definition,
      IF(p_after_column IS NULL OR p_after_column = '', '', CONCAT(' AFTER `', p_after_column, '`'))
    );
    PREPARE stmt FROM @ddl;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
  END IF;
END//

CREATE PROCEDURE `add_index_if_missing`(
  IN p_table_name VARCHAR(64),
  IN p_index_name VARCHAR(64),
  IN p_index_definition TEXT
)
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = p_table_name
      AND INDEX_NAME = p_index_name
  ) THEN
    SET @ddl = CONCAT('ALTER TABLE `', p_table_name, '` ADD ', p_index_definition);
    PREPARE stmt FROM @ddl;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
  END IF;
END//

DELIMITER ;

CALL `add_column_if_missing`(
  'users',
  'token_version',
  'INT NOT NULL DEFAULT 0 COMMENT ''访问凭证版本，递增后旧JWT失效''',
  'last_login_time'
);

ALTER TABLE `users`
  MODIFY COLUMN `role` VARCHAR(20) NOT NULL DEFAULT 'USER' COMMENT 'USER / ADMIN / SUPER_ADMIN',
  MODIFY COLUMN `account_status` VARCHAR(20) NOT NULL DEFAULT 'NORMAL' COMMENT 'NORMAL / DISABLED / BANNED / CANCELED';

CALL `add_column_if_missing`(
  'admin_logs',
  'related_type',
  'VARCHAR(30) NULL COMMENT ''操作依据类型：REPORT / APPEAL''',
  'ip_address'
);

CALL `add_column_if_missing`(
  'admin_logs',
  'related_id',
  'BIGINT UNSIGNED NULL COMMENT ''操作依据ID''',
  'related_type'
);

ALTER TABLE `admin_logs`
  MODIFY COLUMN `target_type` VARCHAR(30) NOT NULL COMMENT 'USER / PRODUCT / ORDER / NOTICE / WORD / CATEGORY / REPORT / APPEAL';

CALL `add_index_if_missing`(
  'admin_logs',
  'idx_admin_logs_related',
  'KEY `idx_admin_logs_related` (`related_type`, `related_id`)'
);

CALL `add_column_if_missing`(
  'chat_conversations',
  'buyer_hidden_at',
  'DATETIME NULL COMMENT ''买家隐藏会话时间，NULL表示买家聊天列表可见''',
  'seller_unread_count'
);

CALL `add_column_if_missing`(
  'chat_conversations',
  'seller_hidden_at',
  'DATETIME NULL COMMENT ''卖家隐藏会话时间，NULL表示卖家聊天列表可见''',
  'buyer_hidden_at'
);

ALTER TABLE `reports`
  MODIFY COLUMN `status` VARCHAR(30) NOT NULL DEFAULT 'PENDING' COMMENT 'PENDING / PROCESSING / APPROVED / REJECTED / CLOSED';

CALL `add_column_if_missing`(
  'reviews',
  'is_anonymous',
  'TINYINT(1) NOT NULL DEFAULT 0 COMMENT ''是否匿名评价''',
  'content'
);

DROP PROCEDURE IF EXISTS `add_index_if_missing`;
DROP PROCEDURE IF EXISTS `add_column_if_missing`;
