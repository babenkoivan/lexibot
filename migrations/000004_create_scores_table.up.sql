CREATE TABLE `scores`
(
    `user_id`        int unsigned NOT NULL,
    `translation_id` int unsigned NOT NULL,
    `score`          int          NOT NULL DEFAULT '0',
    `created_at`     timestamp    NOT NULL,
    `updated_at`     timestamp    NOT NULL,
    PRIMARY KEY (`user_id`, `translation_id`),
    KEY `idx_updated_at` (`updated_at`),
    KEY `idx_translation_id` (`translation_id`),
    CONSTRAINT `fk_scores_translations` FOREIGN KEY (`translation_id`) REFERENCES `translations` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;