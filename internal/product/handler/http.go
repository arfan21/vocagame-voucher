package producthandler

import (
	"context"
	"net/http"

	"github.com/arfan21/vocagame/internal/model"
	"github.com/arfan21/vocagame/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
)

type UseCase interface {
	GetList(ctx context.Context) (res []model.ProductResponse, err error)
}

type HTTP struct {
	uc UseCase
}

func NewHTTP(uc UseCase) *HTTP {
	return &HTTP{uc}
}

// @Summary Get List Product
// @Description Get List Product
// @Tags Product
// @Accept json
// @Produce json
// @Success 200 {array} pkgutil.HTTPResponse{data=model.ProductResponse}
// @Router /api/v1/products [get]
func (h HTTP) GetList(c *fiber.Ctx) error {
	ctx := c.UserContext()
	res, err := h.uc.GetList(ctx)
	if err != nil {
		return err
	}

	return c.JSON(pkgutil.HTTPResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    res,
	})
}
