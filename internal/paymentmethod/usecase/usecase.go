package paymentmethoduc

import (
	"context"

	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/internal/model"
)

type Repository interface {
	GetList(ctx context.Context) (res []entity.PaymentMethod, err error)
	GetByID(ctx context.Context, id int) (res entity.PaymentMethod, err error)
}

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc UseCase) GetList(ctx context.Context) (res []model.PaymentMethodResponse, err error) {
	paymentMethods, err := uc.repo.GetList(ctx)
	if err != nil {
		return nil, err
	}

	res = make([]model.PaymentMethodResponse, len(paymentMethods))

	for i, paymentMethod := range paymentMethods {
		res[i] = model.PaymentMethodResponse{
			ID:   paymentMethod.ID,
			Name: paymentMethod.Name,
		}
	}

	return res, nil
}

func (uc UseCase) GetByID(ctx context.Context, id int) (res model.PaymentMethodResponse, err error) {
	paymentMethod, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return model.PaymentMethodResponse{}, err
	}

	res = model.PaymentMethodResponse{
		ID:   paymentMethod.ID,
		Name: paymentMethod.Name,
	}

	return res, nil
}
