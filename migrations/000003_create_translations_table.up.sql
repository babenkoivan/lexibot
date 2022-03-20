CREATE TABLE `translations`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `user_id`     int unsigned NOT NULL,
    `text`        varchar(255) NOT NULL,
    `translation` varchar(255) NOT NULL,
    `lang_from`   varchar(2)   NOT NULL,
    `lang_to`     varchar(2)   NOT NULL,
    `manual`      tinyint(1)   NOT NULL,
    `score`       int          NOT NULL DEFAULT '0',
    `created_at`  timestamp    NOT NULL,
    `updated_at`  timestamp    NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unq_user_id_lang_from_text` (`user_id`, `lang_from`, `text`),
    KEY `idx_text_lang_from_lang_to` (`text`, `lang_from`, `lang_to`),
    KEY `idx_updated_at` (`updated_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;