CREATE TABLE `translations`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `text`        varchar(255) NOT NULL DEFAULT '',
    `translation` varchar(255) NOT NULL DEFAULT '',
    `lang_from`   varchar(2)   NOT NULL,
    `lang_to`     varchar(2)   NOT NULL,
    `manual`      tinyint(1)   NOT NULL,
    `created_at`  timestamp    NOT NULL,
    `updated_at`  timestamp    NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unq_text_lang_translation` (`text`, `lang_from`, `lang_to`, `translation`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;