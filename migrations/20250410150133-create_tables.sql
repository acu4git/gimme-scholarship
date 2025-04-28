-- +migrate Up
CREATE TABLE IF NOT EXISTS `education_levels` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(8) COMMENT '学部や院などの情報'
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `users` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `uuid` CHAR(36) NOT NULL UNIQUE,
  `email` VARCHAR(255) NOT NULL UNIQUE,
  `name` VARCHAR(255) NOT NULL COMMENT '表示名',
  `education_level_id` INT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT `users_education_level_id`
    FOREIGN KEY (`education_level_id`)
    REFERENCES `education_levels`(`id`)
    ON DELETE RESTRICT
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `user_auth` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `user_id` INT NOT NULL,
  `sub` VARCHAR(255) NOT NULL,
  `provider` VARCHAR(64) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT `fk_user_id`
    FOREIGN KEY (`user_id`)
    REFERENCES `users` (`id`),
  UNIQUE KEY `unique_sub_provider` (`sub`, `provider`)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `scholarships` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(128) NOT NULL,
  `address` TEXT NOT NULL,
  `target_detail` TEXT NOT NULL,
  `amount_detail` TEXT NOT NULL,
  `type_detail` TEXT NOT NULL,
  `capacity_detail` TEXT NOT NULL,
  `deadline` DATE NOT NULL,
  `deadline_detail` TEXT NOT NULL,
  `contact_point` TEXT NOT NULL,
  `remark` TEXT NOT NULL,
  `posting_date` DATE NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY (`name`)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `temporary_scholarships` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(128) NOT NULL,
  `address` TEXT NOT NULL,
  `target_detail` TEXT NOT NULL,
  `amount_detail` TEXT NOT NULL,
  `type_detail` TEXT NOT NULL,
  `capacity_detail` TEXT NOT NULL,
  `deadline` DATE NOT NULL,
  `deadline_detail` TEXT NOT NULL,
  `contact_point` TEXT NOT NULL,
  `remark` TEXT NOT NULL,
  `posting_date` DATE NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY (`name`)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `scholarship_target` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `scholarship_id` INT NOT NULL,
  `education_level_id` INT NOT NULL,
  UNIQUE KEY (`scholarship_id`, `education_level_id`),
  CONSTRAINT `fk_scholarship_id`
    FOREIGN KEY (`scholarship_id`)
    REFERENCES `scholarships` (`id`)
    ON DELETE CASCADE,
  CONSTRAINT `fk_education_level_id`
    FOREIGN KEY (`education_level_id`)
    REFERENCES `education_levels` (`id`)
    ON DELETE RESTRICT
) ENGINE = InnoDB;
-- +migrate Down
