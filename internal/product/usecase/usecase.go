package productusecase

import (
	"context"
	"time"

	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/internal/model"
	productrepo "github.com/arfan21/vocagame/internal/product/repository"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	WithTx(tx pgx.Tx) *productrepo.Repository
	GetList(ctx context.Context) (res []entity.Product, err error)
	GetByID(ctx context.Context, id string, isForUpdate bool) (res entity.Product, err error)
	ReduceStock(ctx context.Context, id string, qty int) error
}

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc UseCase) Begin(ctx context.Context) (pgx.Tx, error) {
	return uc.repo.Begin(ctx)
}

func (uc UseCase) WithTx(tx pgx.Tx) *UseCase {
	return &UseCase{repo: uc.repo.WithTx(tx)}
}

func (uc UseCase) GetList(ctx context.Context) (res []model.ProductResponse, err error) {
	products, err := uc.repo.GetList(ctx)
	if err != nil {
		return nil, err
	}

	res = make([]model.ProductResponse, len(products))

	for i, product := range products {
		res[i] = model.ProductResponse{
			ID:        product.ID,
			Name:      product.Name,
			Price:     product.Price,
			Stock:     product.Stock,
			CreatedAt: product.CreatedAt.Format(time.DateTime),
			UpdatedAt: product.UpdatedAt.Format(time.DateTime),
		}
	}

	return res, nil
}

func (uc UseCase) GetByID(ctx context.Context, id string, isForUpdate bool) (res model.ProductResponse, err error) {
	product, err := uc.repo.GetByID(ctx, id, isForUpdate)
	if err != nil {
		return model.ProductResponse{}, err
	}

	res = model.ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt.Format(time.DateTime),
		UpdatedAt: product.UpdatedAt.Format(time.DateTime),
	}

	return res, nil
}

func (uc UseCase) ReduceStock(ctx context.Context, id string, qty int) error {
	return uc.repo.ReduceStock(ctx, id, qty)
}
