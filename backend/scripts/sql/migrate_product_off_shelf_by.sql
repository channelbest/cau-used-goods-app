DROP PROCEDURE IF EXISTS `add_product_off_shelf_by_column`;

DELIMITER //

CREATE PROCEDURE `add_product_off_shelf_by_column`()
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'products'
      AND COLUMN_NAME = 'off_shelf_by'
  ) THEN
    ALTER TABLE `products`
      ADD COLUMN `off_shelf_by` VARCHAR(20) NULL COMMENT '下架来源：USER / ADMIN / SYSTEM'
      AFTER `off_shelf_reason`;
  END IF;
END//

DELIMITER ;

CALL `add_product_off_shelf_by_column`();

DROP PROCEDURE IF EXISTS `add_product_off_shelf_by_column`;
