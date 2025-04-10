
-- +migrate Up
CREATE TABLE IF NOT EXISTS `scholarships` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `association_name` VARCHAR(64) NOT NULL,

)
-- +migrate Down
