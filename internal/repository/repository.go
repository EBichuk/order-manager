package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"order-manager/internal/models"
	"order-manager/pkg/errorx"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetOrderByUID(orderUID string) (*models.Order, error) {
	query := `
		SELECT
			o.order_uid, o.track_number, o.entry, o.locate, o.internal_signature,
			o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.off_shard, 
			d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
			p.transaction, p.request_id, p.currency, p.provider, p.amount, 
			p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM 
			orders o
		JOIN 
			deliveries d ON d.order_uid = o.order_uid
		JOIN 
			payments p ON p.order_uid = o.order_uid
		WHERE 
			o.order_uid = $1
	`
	var order models.Order
	var delivery models.Delivery
	var payment models.Payment

	err := r.pool.QueryRow(context.Background(), query, orderUID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locate, &order.InternalSignature, &order.CustomerID,
		&order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OffShard,
		&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email,
		&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDt,
		&payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorx.ErrOrderNotFound
		}
		return nil, err
	}

	item := r.GetItemsByOrderUID(orderUID)

	order.Delivery = delivery
	order.Payment = payment
	order.Item = item
	return &order, nil
}

func (r *Repository) GetItemsByOrderUID(orderUID string) []models.Item {
	var items []models.Item
	query := `
		SELECT
			i.chrt_id, i.track_number, i.price, i.rid, i.name_item, 
			i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status
		FROM 
			items i
		WHERE 
			i.order_uid = $1
	`
	rows, err := r.pool.Query(context.Background(), query, orderUID)
	if err != nil {
		return nil
	}

	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid,
			&item.NameItem, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil
		}
		items = append(items, item)
	}
	return items
}

func (r *Repository) SaveOrder(order *models.Order) error {
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		log.Fatalf("Error starting transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `
		INSERT INTO orders (
    		order_uid, track_number, entry, locate, internal_signature,
    		customer_id, delivery_service, shardkey, sm_id, date_created, off_shard
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) 
		ON CONFLICT(order_uid) 
		DO UPDATE SET
			order_uid = $1, track_number = $2, entry = $3, locate = $4, internal_signature = $5,
    		customer_id = $6, delivery_service = $7, shardkey = $8, sm_id = $9, date_created = $10, off_shard = $11;`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locate, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OffShard)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), `
		INSERT INTO deliveries (
    		order_uid, name, phone, zip, city, address, region, email
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) 
		ON CONFLICT (order_uid) 
		DO UPDATE SET
			order_uid = $1, name = $2, phone = $3, zip = $4, city = $5, address = $6, region = $7, email = $8;`,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), `
		INSERT INTO payments (
			order_uid, transaction, request_id, currency, provider, 
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) 
		ON CONFLICT (transaction)
		DO UPDATE SET
			order_uid = $1, transaction = $2, request_id = $3, currency = $4, provider = $5, 
			amount = $6, payment_dt = $7, bank = $8, delivery_cost = $9, goods_total = $10, custom_fee = $11;`,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider,
		order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return err
	}

	queryItems := `
		INSERT INTO items (
    		order_uid, chrt_id, track_number, price, rid, name_item, 
    		sale, size, total_price, nm_id, brand, status)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		) 
		ON CONFLICT (rid) 
		DO UPDATE SET
			order_uid = $1, chrt_id = $2, track_number = $3, price = $4, rid = $5, name_item = $6, 
    		sale = $7, size = $8, total_price = $9, nm_id = $10, brand = $11, status = $12;`

	for _, item := range order.Item {
		_, err = tx.Exec(context.Background(), queryItems,
			order.OrderUID,
			item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.NameItem,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit(context.Background())
}

func (r *Repository) GetAllOrders(size int) ([]models.Order, error) {
	var orders []models.Order
	query := `
		SELECT 
			order_uid
		FROM 
			orders
		ORDER BY
			date_created
		LIMIT  
			$1
	`
	rows, err := r.pool.Query(context.Background(), query, size)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var orderUID string
		if err = rows.Scan(&orderUID); err != nil {
			continue
		}

		order, err := r.GetOrderByUID(orderUID)
		if err != nil {
			continue
		}
		orders = append(orders, *order)
	}
	return orders, nil
}
