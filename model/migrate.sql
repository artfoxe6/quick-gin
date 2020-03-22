CREATE TABLE IF NOT EXISTS `user` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_name` varchar(30) DEFAULT NULL,
    `age` tinyint(4) DEFAULT NULL,
    `last_login_at` timestamp NULL DEFAULT NULL,
    `password` char(100) DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL,
    `deleted_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `user_user_name_IDX` (`user_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE IF NOT EXISTS `article` (
     `id` int(11) NOT NULL AUTO_INCREMENT,
     `title` varchar(30) DEFAULT NULL,
     `content` varchar(100) DEFAULT NULL,
     `favorite` int(11) DEFAULT '0',
     `uid` int(11) NOT NULL,
     `created_at` timestamp NULL DEFAULT NULL,
     `updated_at` timestamp NULL DEFAULT NULL,
     `deleted_at` timestamp NULL DEFAULT NULL,
     PRIMARY KEY (`id`),
     KEY `article_title_IDX` (`title`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;