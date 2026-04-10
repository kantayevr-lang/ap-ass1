package repository

import (
	"context"
	"database/sql"
	"log"
	"payment-service/internal/domain"

	"github.com/google/uuid"
)

type PostgresPaymentRepository struct {
	db *sql.DB
}

func NewPostgresPaymentRepository(db *sql.DB) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{db: db}
}

func (r *PostgresPaymentRepository) Save(ctx context.Context, p *domain.Payment) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}

	if p.Status == "" {
		p.Status = "pending"
	}

	query := `INSERT INTO payments (id, order_id, transaction_id, amount, status) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.OrderID, p.TransactionID, p.Amount, p.Status)

	if err != nil {
		log.Println("DB ERROR:", err)
		return err
	}

	return err
}

func (r *PostgresPaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	query := `SELECT id, order_id, transaction_id, amount, status FROM payments WHERE order_id = $1`
	var p domain.Payment
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status)
	return &p, err
}
