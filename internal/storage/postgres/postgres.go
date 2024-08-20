package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // driver for postgres
	"wb-order-service/internal/model"
	"wb-order-service/internal/storage"
)

type Database struct {
	db *sql.DB
}

func NewDB(storageConn string) (*Database, error) {
	db, err := sql.Open("postgres", storageConn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %s: %w", storage.ErrConnect, err)
	}

	return &Database{db: db}, nil
}

func (db *Database) InsertOrderToDB(order *model.Order) (int, error) {
	const op = "storage.postgres.InsertOrderToDB"
	var orderId int

	tx, err := db.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: error: %w", op, err)
	}
	defer tx.Rollback()

	err = tx.QueryRow(
		"INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		order.OrderUid,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.ShardKey,
		order.SmId,
		order.DateCreated,
		order.OofShard,
	).Scan(&orderId)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to create order: %w: %d", op, storage.ErrInternalServer, err)
	}

	_, err = tx.Exec(
		"INSERT INTO delivery (name, phone, zip, city, address, region, email, order_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
		orderId,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to create delivery: %w", op, err)
	}

	_, err = tx.Exec(
		"INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, order_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.Payment.Transaction,
		order.Payment.RequestId,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
		orderId,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to create payment: %w", op, err)
	}

	for _, item := range order.Items {
		_, err := tx.Exec(
			"INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			item.ChrtId,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmId,
			item.Brand,
			item.Status,
			orderId,
		)
		if err != nil {
			return 0, fmt.Errorf("%s: failed to create item: %w", op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to commit: %w", op, err)
	}

	return orderId, nil
}

func (db *Database) GetOrdersFromDB() ([]model.Order, error) {
	const op = "storage.postgres.GetOrdersFromDB"
	var orders []model.Order

	rows, err := db.db.Query(
		"SELECT * FROM orders",
	)
	if err != nil {
		return nil, fmt.Errorf("%s: order query failed: %w: %d", op, storage.ErrInternalServer, err)
	}
	for rows.Next() {
		var order model.Order
		var orderId int

		err := rows.Scan(
			&orderId,
			&order.OrderUid,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerId,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmId,
			&order.DateCreated,
			&order.OofShard,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan order failed: %w: %d", op, storage.ErrInternalServer, err)
		}

		err = db.db.QueryRow(
			"SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_id = $1",
			orderId).Scan(
			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to query order.delivery: %w", op, err)
		}

		err = db.db.QueryRow(
			"SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_id = $1",
			orderId).Scan(
			&order.Payment.Transaction,
			&order.Payment.RequestId,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDt,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to query order.payment: %w", op, err)
		}

		rows, err := db.db.Query(
			"SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_id = $1",
			orderId)
		if err != nil {
			return nil, fmt.Errorf("%s: item query failed: %w: %d", op, storage.ErrInternalServer, err)
		}
		for rows.Next() {
			var item model.Item
			err = rows.Scan(
				&item.ChrtId,
				&item.TrackNumber,
				&item.Price,
				&item.Rid,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.TotalPrice,
				&item.NmId,
				&item.Brand,
				&item.Status,
			)
			if err != nil {
				return nil, fmt.Errorf("%s: scan item failed: %w: %d", op, storage.ErrInternalServer, err)
			}
			order.Items = append(order.Items, item)
		}
		orders = append(orders, order)
	}
	return orders, nil
}
