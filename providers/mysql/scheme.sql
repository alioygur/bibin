SET foreign_key_checks = 0;
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `facebook_id` BIGINT unsigned NOT NULL,
  `first_name` varchar(255) NOT NULL,
  `last_name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `gender` TINYINT unsigned NOT NULL DEFAULT 0,
  `birthday` DATE NOT NULL,  
  `status` TINYINT unsigned NOT NULL DEFAULT 0,
  `is_admin` TINYINT unsigned NOT NULL DEFAULT 0,
  `account_type` TINYINT unsigned NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`email`),
  UNIQUE KEY (`facebook_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `reactions`;
CREATE TABLE `reactions` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `from_user_id` BIGINT unsigned NOT NULL,
  `to_user_id` BIGINT unsigned NOT NULL,
  `type` TINYINT unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`from_user_id`, `to_user_id`),
  CONSTRAINT `fk_reactions_from_user_id` FOREIGN KEY (`from_user_id`) references users(`id`),
  CONSTRAINT `fk_reactions_to_user_id` FOREIGN KEY (`to_user_id`) references users(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `credits`;
CREATE TABLE `credits` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT unsigned NOT NULL,
  `type` TINYINT unsigned NOT NULL,
  `amount` MEDIUMINT NOT NULL,
  `desc` TINYTEXT NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY (`user_id`),
  CONSTRAINT `fk_credits_user_id` FOREIGN KEY (`user_id`) references users(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- DROP TABLE IF EXISTS `friendships`;
-- CREATE TABLE `friendships` (
--   `user_id` BIGINT unsigned NOT NULL,
--   `friend_id` BIGINT unsigned NOT NULL,
--   PRIMARY KEY(`user_id`, `friend_id`),
--   CONSTRAINT `fk_friendships_user_id` FOREIGN KEY (`user_id`) references users(`id`),
--   CONSTRAINT `fk_friendships_friend_id` FOREIGN KEY (`friend_id`) references users(`id`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `matches`;
CREATE TABLE `matches` (
  `user_id` BIGINT unsigned NOT NULL,
  `friend_id` BIGINT unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(`user_id`, `friend_id`),
  CONSTRAINT `fk_matches_user_id` FOREIGN KEY (`user_id`) references users(`id`),
  CONSTRAINT `fk_matches_friend_id` FOREIGN KEY (`friend_id`) references users(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `abuses`;
CREATE TABLE `abuses` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT unsigned NOT NULL,
  `to_user_id` BIGINT unsigned NOT NULL,
  `reason` TINYTEXT NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(`id`),
  UNIQUE (`user_id`, `to_user_id`),
  CONSTRAINT `fk_abuses_user_id` FOREIGN KEY (`user_id`) references users(`id`),
  CONSTRAINT `fk_abuses_to_user_id` FOREIGN KEY (`to_user_id`) references users(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `images`;
CREATE TABLE `images` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `is_profile` TINYINT unsigned NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(`id`),
  UNIQUE (`name`),
  CONSTRAINT `fk_images_user_id` FOREIGN KEY (`user_id`) references users(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

SET foreign_key_checks = 1;