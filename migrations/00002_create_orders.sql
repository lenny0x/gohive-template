-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    order_no   VARCHAR(64)     NOT NULL,
    user_id    BIGINT UNSIGNED NOT NULL,
    amount     DECIMAL(10, 2)  NOT NULL DEFAULT 0.00,
    status     TINYINT         NOT NULL DEFAULT 0 COMMENT '0=pending, 1=paid, 2=completed, 3=cancelled',
    expired_at DATETIME        NOT NULL,
    created_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME        NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_order_no (order_no),
    KEY idx_user_id (user_id),
    KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS orders;
