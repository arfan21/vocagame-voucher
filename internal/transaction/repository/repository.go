package transactionrepo

import (
	"context"
	"fmt"

	"github.com/arfan21/vocagame/internal/entity"
	dbpostgres "github.com/arfan21/vocagame/pkg/db/postgres"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db dbpostgres.Queryer
}

func New(db dbpostgres.Queryer) *Repository {
	return &Repository{db}
}

func (r Repository) Begin(ctx context.Context) (tx pgx.Tx, err error) {
	return r.db.Begin(ctx)
}

func (r Repository) WithTx(tx pgx.Tx) *Repository {
	return &Repository{db: tx}
}

func (r Repository) Create(ctx context.Context, data entity.Transaction) (err error) {
	query := `INSERT INTO transactions (id, product_id, payment_method_id, phone_number, quantity, total_price) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = r.db.Exec(ctx, query, data.ID, data.ProductID, data.PaymentMethodID, data.PhoneNumber, data.Quantity, data.TotalPrice)
	if err != nil {
		err = fmt.Errorf("transaction.repo.Create: failed to create transaction: %w", err)
		return err
	}

	return nil
}

func (r Repository) UpdateStatus(ctx context.Context, id string, status entity.Status) (err error) {
	query := `UPDATE transactions SET status = $1 WHERE id = $2`
	_, err = r.db.Exec(ctx, query, status, id)
	if err != nil {
		err = fmt.Errorf("transaction.repo.UpdateStatus: failed to update status: %w", err)
		return err
	}

	return nil
}

func (r Repository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (res []entity.Transaction, err error) {
	query := `
		SELECT 
			transactions.id, 
			product_id, 
			payment_method_id, 
			phone_number, 
			quantity, 
			total_price, 
			status, 
			transactions.created_at, 
			transactions.updated_at,
			products.name as product_name,
			payment_methods.name as payment_method_name
		FROM transactions 
		JOIN payment_methods ON transactions.payment_method_id = payment_methods.id
		JOIN products ON transactions.product_id = products.id
		WHERE phone_number = $1`
	rows, err := r.db.Query(ctx, query, phoneNumber)
	if err != nil {
		err = fmt.Errorf("transaction.repo.GetByPhoneNumber: failed to get transaction by phone number: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction entity.Transaction
		err = rows.Scan(
			&transaction.ID,
			&transaction.ProductID,
			&transaction.PaymentMethodID,
			&transaction.PhoneNumber,
			&transaction.Quantity,
			&transaction.TotalPrice,
			&transaction.Status,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&transaction.Product.Name,
			&transaction.PaymentMethod.Name,
		)
		if err != nil {
			err = fmt.Errorf("transaction.repo.GetByPhoneNumber: failed to scan transaction: %w", err)
			return nil, err
		}

		res = append(res, transaction)
	}

	return res, nil
}
