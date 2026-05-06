-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    username   VARCHAR(64)     NOT NULL,
    email      VARCHAR(128)    NOT NULL,
    password   VARCHAR(255)    NOT NULL,
    nickname   VARCHAR(64)     NOT NULL DEFAULT '',
    avatar     VARCHAR(256)    NOT NULL DEFAULT '',
    status     TINYINT         NOT NULL DEFAULT 1 COMMENT '1=active, 0=disabled',
    created_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME        NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_email (email),
    UNIQUE KEY uk_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS users;
