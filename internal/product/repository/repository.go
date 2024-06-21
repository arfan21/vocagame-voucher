package productrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/pkg/constant"
	dbpostgres "github.com/arfan21/vocagame/pkg/db/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db dbpostgres.Queryer
}

func New(db dbpostgres.Queryer) *Repository {
	return &Repository{db}
}

func (repo Repository) Begin(ctx context.Context) (pgx.Tx, error) {
	return repo.db.Begin(ctx)
}

func (repo Repository) WithTx(tx pgx.Tx) *Repository {
	return &Repository{tx}
}

func (repo Repository) GetList(ctx context.Context) (res []entity.Product, err error) {
	query := `SELECT id, name, price, stock, created_at, updated_at FROM products`
	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("product.repo.GetList: failed to get list product: %w", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product entity.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			err = fmt.Errorf("product.repo.GetList: failed to scan product: %w", err)
			return nil, err
		}

		res = append(res, product)
	}

	return res, nil
}

func (repo Repository) GetByID(ctx context.Context, id string, isForUpdate bool) (res entity.Product, err error) {
	query := `SELECT id, name, price, stock, created_at, updated_at FROM products WHERE id = $1`

	if isForUpdate {
		query += " FOR UPDATE"
	}

	row := repo.db.QueryRow(ctx, query, id)

	err = row.Scan(&res.ID, &res.Name, &res.Price, &res.Stock, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLInvalidUUID {
				err = constant.ErrProductNotFound
			}
		}

		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrProductNotFound
		}
		err = fmt.Errorf("product.repo.GetByID: failed to get product by id: %w", err)
		return entity.Product{}, err
	}

	return res, nil
}

func (repo Repository) ReduceStock(ctx context.Context, id string, qty int) (err error) {
	query := `UPDATE products SET stock = stock - $1 WHERE id = $2 AND (stock - $1) >= 0`
	cmd, err := repo.db.Exec(ctx, query, qty, id)
	if err != nil {
		err = fmt.Errorf("product.repo.ReduceStock: failed to reduce stock: %w", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		err = fmt.Errorf("product.repo.ReduceStock: failed to reduce stock: stock not enough, %w", constant.ErrOutOfStock)
		return err
	}

	return nil
}

func (repo Repository) IncreaseStock(ctx context.Context, id string, qty int) (err error) {
	query := `UPDATE products SET stock = stock + $1 WHERE id = $2`
	_, err = repo.db.Exec(ctx, query, qty, id)
	if err != nil {
		err = fmt.Errorf("product.repo.IncreaseStock: failed to increase stock: %w", err)
		return err
	}

	return nil
}
