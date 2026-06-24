DROP PROCEDURE IF EXISTS `add_review_anonymous_column`;

DELIMITER //

CREATE PROCEDURE `add_review_anonymous_column`()
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'reviews'
      AND COLUMN_NAME = 'is_anonymous'
  ) THEN
    ALTER TABLE `reviews`
      ADD COLUMN `is_anonymous` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否匿名评价'
      AFTER `content`;
  END IF;
END//

DELIMITER ;

CALL `add_review_anonymous_column`();

DROP PROCEDURE IF EXISTS `add_review_anonymous_column`;
