
-- +migrate Up
ALTER TABLE `user_favorites` DROP FOREIGN KEY `user_favorite_user_id`;
ALTER TABLE `user_favorites` MODIFY COLUMN `user_id` CHAR(64) NOT NULL;
ALTER TABLE `user_favorites` ADD CONSTRAINT `user_favorite_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
-- +migrate Down
