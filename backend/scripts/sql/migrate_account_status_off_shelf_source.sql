UPDATE products p
JOIN admin_logs l
  ON l.target_type = 'PRODUCT'
 AND l.target_id = p.id
 AND l.description LIKE 'off shelf product due account status%'
SET p.off_shelf_by = 'ACCOUNT_STATUS',
    p.update_time = CURRENT_TIMESTAMP
WHERE p.status = 'OFF_SHELF'
  AND p.off_shelf_by = 'SYSTEM'
  AND p.is_deleted = 0;
