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
    SELECT 1 FROM information_schema.COLUMNS
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
    SELECT 1 FROM information_schema.STATISTICS
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

ALTER TABLE `chat_messages`
  MODIFY COLUMN `sender_id` BIGINT UNSIGNED NULL COMMENT '发送人ID，系统事件为空',
  MODIFY COLUMN `receiver_id` BIGINT UNSIGNED NULL COMMENT '接收人ID，双方可见的系统事件为空',
  MODIFY COLUMN `message_type` VARCHAR(30) NOT NULL DEFAULT 'TEXT' COMMENT 'TEXT / ORDER_EVENT';

CALL `add_column_if_missing`(
  'chat_messages',
  'actor_type',
  'VARCHAR(20) NOT NULL DEFAULT ''USER'' COMMENT ''USER / SYSTEM''',
  'message_type'
);

CALL `add_column_if_missing`(
  'chat_messages',
  'order_id',
  'BIGINT UNSIGNED NULL COMMENT ''关联订单ID''',
  'actor_type'
);

CALL `add_column_if_missing`(
  'chat_messages',
  'event_type',
  'VARCHAR(40) NULL COMMENT ''ORDER_CREATED / ORDER_CONFIRMED / ORDER_CANCELED / ORDER_COMPLETED / ORDER_TIMEOUT / ORDER_EXCEPTION_CLOSED / ORDER_STATUS_UPDATED''',
  'order_id'
);

CALL `add_index_if_missing`(
  'chat_messages',
  'idx_chat_messages_order',
  'KEY `idx_chat_messages_order` (`order_id`)'
);

SET @fk_exists = (
  SELECT COUNT(*)
  FROM information_schema.REFERENTIAL_CONSTRAINTS
  WHERE CONSTRAINT_SCHEMA = DATABASE()
    AND TABLE_NAME = 'chat_messages'
    AND CONSTRAINT_NAME = 'fk_chat_messages_order'
);
SET @fk_sql = IF(
  @fk_exists = 0,
  'ALTER TABLE `chat_messages` ADD CONSTRAINT `fk_chat_messages_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)',
  'SELECT 1'
);
PREPARE stmt FROM @fk_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

DROP PROCEDURE IF EXISTS `add_index_if_missing`;
DROP PROCEDURE IF EXISTS `add_column_if_missing`;
