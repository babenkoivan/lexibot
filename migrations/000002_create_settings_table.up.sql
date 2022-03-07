CREATE TABLE `settings`
(
    `user_id`            int unsigned                                                  NOT NULL,
    `lang_ui`            varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'en',
    `lang_dict`          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci          DEFAULT '',
    `auto_translate`     tinyint(1)                                                    NOT NULL DEFAULT '1',
    `words_per_training` int                                                           NOT NULL DEFAULT '10',
    `created_at`         timestamp                                                     NOT NULL,
    `updated_at`         timestamp                                                     NOT NULL,
    PRIMARY KEY (`user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;