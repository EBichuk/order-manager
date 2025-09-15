-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(50) NOT NULL,
    locate VARCHAR(10) NOT NULL,
    internal_signature VARCHAR(255) NOT NULL,
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(100) NOT NULL,
    shardkey VARCHAR(10) NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMP NOT NULL,
    off_shard VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS deliveries (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    zip VARCHAR(15) NOT NULL,
    city VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payments (
    order_uid VARCHAR(255) NOT NULL,
    transaction VARCHAR(255) PRIMARY KEY,
    request_id VARCHAR(255),
    currency VARCHAR(100) NOT NULL,
    provider VARCHAR(100) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt INTEGER NOT NULL,
    bank VARCHAR(100) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL,
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) NOT NULL,
    chrt_id INTEGER NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR(255) UNIQUE NOT NULL,
    name_item VARCHAR(255) NOT NULL,
    sale INTEGER NOT NULL,
    size INTEGER NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id INTEGER NOT NULL,
    brand VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL,
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
