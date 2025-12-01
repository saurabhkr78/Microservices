package order

import (
	"context"
	"database/sql"

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
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, nil

}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) PutOrder(ctx context.Context, o Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Commit OR rollback
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Insert into orders table
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES ($1,$2,$3,$4)",
		o.ID, o.CreatedAt, o.AccountID, o.TotalPrice,
	)
	if err != nil {
		return err
	}

	// Prepare COPY statement
	stmt, err := tx.PrepareContext(
		ctx,
		pq.CopyIn("ordered_products", "order_id", "product_id", "quantity"),
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert products
	for _, p := range o.Products {
		_, err = stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity)
		if err != nil {
			return err
		}
	}

	// Finish COPY
	_, err = stmt.ExecContext(ctx)
	return err
}

func (r *postgresRepository) GetOrdersFromAccount(ctx context.Context, accountID string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
		o.id
		o.created_at,
		o.account_id,
		o.total_price::money::numeric::float8
		op.product_id,
		op.qunatity,
		FROM orders o JOIN order_products op ON (o.id=op.order_id)
		WHERE o.acount_id=$1,
		ORDER BY o.id`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := []order{}
	lastOrder := &order{}
	orderedProduct := &OrderedProduct{}
	products := []product{}

	for rows.Next() {
		if err = rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.AccountID,
			&order.TotalPrice,
			&ororderedProduct.ID,
			&ororderedProduct.Name,
			&ororderedProduct.Description,
			&ororderedProduct.Price,
			&ororderedProduct.Quantity,
		); err != nil {
			return nil, err
		}
		if lastOrder.ID != "" && lastOrder.ID != Order.ID {
			newOrder := Order{
				ID:         lastOrder.ID,
				AccountID:  lastOrder.AccountID,
				CreatedAt:  lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Products:   lastOrder.Products,
			}
			orders.append(orders, newOrder)
			products = []OrderedProduct{}

		}
		products = append(products, OrderedProduct{
			ID:       orderedProduct.ID,
			Quantity: orderedProduct.Quantity,
		})
		*lastOrder = *order
	}
	if lastOrder != nil {
		newOrder := Order{
			ID:         lastOrder.ID,
			AccountID:  lastOrder.AccountID,
			CreatedAt:  lastOrder.CreatedAt,
			TotalPrice: lastOrder.TotalPrice,
			Products:   lastOrder.Products,
		}
		orders.append(orders, newOrder)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil

}
