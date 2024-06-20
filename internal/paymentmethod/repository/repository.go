package paymentmethodrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/pkg/constant"
	dbpostgres "github.com/arfan21/vocagame/pkg/db/postgres"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db dbpostgres.Queryer
}

func New(db dbpostgres.Queryer) *Repository {
	return &Repository{db}
}

func (r Repository) GetList(ctx context.Context) (res []entity.PaymentMethod, err error) {
	query := `SELECT id, name FROM payment_methods`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("paymentmethod.repo.GetList: failed to get list payment method: %w", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var paymentMethod entity.PaymentMethod
		err = rows.Scan(&paymentMethod.ID, &paymentMethod.Name)
		if err != nil {
			err = fmt.Errorf("paymentmethod.repo.GetList: failed to scan payment method: %w", err)
			return nil, err
		}

		res = append(res, paymentMethod)
	}

	return res, nil
}

func (r Repository) GetByID(ctx context.Context, id int) (res entity.PaymentMethod, err error) {
	query := `SELECT id, name FROM payment_methods WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	err = row.Scan(&res.ID, &res.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrPaymentMethodNotFound
		}
		err = fmt.Errorf("paymentmethod.repo.GetByID: failed to get payment method by id: %w", err)
		return entity.PaymentMethod{}, err
	}

	return res, nil
}
