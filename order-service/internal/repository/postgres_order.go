package repository

import (
	"context"
	"database/sql"
	"order-service/internal/domain"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Save(ctx context.Context, ord *domain.Order) error {
	query := `INSERT INTO orders (id, customer_id, item_name, amount, status, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, ord.ID, ord.CustomerID, ord.ItemName, ord.Amount, ord.Status, ord.CreatedAt)
	return err
}

func (r *PostgresOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	query := `SELECT id, customer_id, item_name, amount, status, created_at FROM orders WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var ord domain.Order
	err := row.Scan(&ord.ID, &ord.CustomerID, &ord.ItemName, &ord.Amount, &ord.Status, &ord.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrOrderNotFound
	}
	return &ord, err
}

func (r *PostgresOrderRepository) ListByAmountRange(ctx context.Context, minAmount, maxAmount int64) ([]domain.Order, error) {
	query := `SELECT id, customer_id, item_name, amount, status, created_at
		FROM orders
		WHERE amount >= $1 AND amount <= $2
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, minAmount, maxAmount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var ord domain.Order
		if err := rows.Scan(&ord.ID, &ord.CustomerID, &ord.ItemName, &ord.Amount, &ord.Status, &ord.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *PostgresOrderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
