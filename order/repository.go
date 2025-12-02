package order

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetOrdersFromAccount(ctx context.Context, accountID string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &postgresRepository{db: db}, nil
}

func (r *postgresRepository) Close() {
	_ = r.db.Close()
}

// ---------------------------------------------------------------
// INSERT ORDER + ORDERED PRODUCTS (COPY IN)
// ---------------------------------------------------------------
func (r *postgresRepository) PutOrder(ctx context.Context, o Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback() // safe rollback
	}()

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO orders(id, created_at, account_id, total_price)
		 VALUES ($1, $2, $3, $4)`,
		o.ID, o.CreatedAt, o.AccountID, o.TotalPrice,
	)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(
		ctx,
		pq.CopyIn("order_products", "order_id", "product_id", "quantity"),
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range o.Products {
		_, err = stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity)
		if err != nil {
			return err
		}
	}

	if _, err = stmt.ExecContext(ctx); err != nil {
		return err
	}

	return tx.Commit()
}

// ---------------------------------------------------------------
// GET ORDERS + PRODUCTS
// ---------------------------------------------------------------
func (r *postgresRepository) GetOrdersFromAccount(ctx context.Context, accountID string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
			o.id,
			o.created_at,
			o.account_id,
			(o.total_price::numeric)::float8,
			op.product_id,
			op.quantity
		FROM orders o
		JOIN order_products op ON o.id = op.order_id
		WHERE o.account_id = $1
		ORDER BY o.id, op.product_id`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	var current Order
	var currentID string

	for rows.Next() {
		var (
			orderID    string
			createdAt  time.Time
			accID      string
			totalPrice float64
			productID  string
			quantity   uint32
		)

		if err := rows.Scan(
			&orderID,
			&createdAt,
			&accID,
			&totalPrice,
			&productID,
			&quantity,
		); err != nil {
			return nil, err
		}

		if orderID != currentID {
			if currentID != "" {
				orders = append(orders, current)
			}
			current = Order{
				ID:         orderID,
				CreatedAt:  createdAt,
				AccountID:  accID,
				TotalPrice: totalPrice,
				Products:   []OrderedProduct{},
			}
			currentID = orderID
		}

		current.Products = append(current.Products, OrderedProduct{
			ID:       productID,
			Quantity: quantity,
		})
	}

	if currentID != "" {
		orders = append(orders, current)
	}

	return orders, rows.Err()
}
