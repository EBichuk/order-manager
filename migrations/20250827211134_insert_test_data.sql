-- +goose Up
-- +goose StatementBegin
INSERT INTO orders (
    order_uid, track_number, entry, locate, internal_signature,
    customer_id, delivery_service, shardkey, sm_id, date_created, off_shard
) VALUES (
    'b563feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', 'en', '',
    'test_customer', 'meest', '9', 99, '2021-11-26 06:22:19', '1'
);

INSERT INTO deliveries (
    order_uid, name, phone, zip, city, address, region, email
) VALUES (
    'b563feb7b2b84b6test', 'Test Testov', '+9720000000', '2639809', 
    'Kiryat Mozkin', 'Ploshad Mira 15', 'Kraiot', 'test@gmail.com'
);

INSERT INTO payments (
    order_uid, transaction, request_id, currency, provider, 
    amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
) VALUES (
    'b563feb7b2b84b6test', 'b563feb7b2b84b6test', '', 'USD', 
    'wbpay', 1817, 1637907727, 'alpha', 1500, 317, 0
);

INSERT INTO items (
    order_uid, chrt_id, track_number, price, rid, name_item, 
    sale, size, total_price, nm_id, brand, status
) VALUES 
(
    'b563feb7b2b84b6test', 9934930, 'WBILMTESTTRACK', 453, 
    'ab4219087a764ae0btest', 'Mascaras', 30, '0', 317, 
    2389212, 'Vivienne Sabo', 202
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM items;
DELETE FROM deliveries;
DELETE FROM payments;
DELETE FROM orders;
-- +goose StatementEnd
