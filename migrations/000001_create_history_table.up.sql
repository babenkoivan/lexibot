CREATE TABLE `history`
(
    `user_id`    int unsigned NOT NULL,
    `type`       varchar(255) NOT NULL,
    `content`    json         NOT NULL,
    `created_at` timestamp    NOT NULL,
    `updated_at` timestamp    NOT NULL,
    PRIMARY KEY (`user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;