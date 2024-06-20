package server

import (
	paymentmidtrans "github.com/arfan21/vocagame/client/payment/midtrans"
	paymentmethodhandler "github.com/arfan21/vocagame/internal/paymentmethod/handler"
	paymentmethodrepo "github.com/arfan21/vocagame/internal/paymentmethod/repository"
	paymentmethoduc "github.com/arfan21/vocagame/internal/paymentmethod/usecase"
	producthandler "github.com/arfan21/vocagame/internal/product/handler"
	productrepo "github.com/arfan21/vocagame/internal/product/repository"
	productusecase "github.com/arfan21/vocagame/internal/product/usecase"
	transactionhandler "github.com/arfan21/vocagame/internal/transaction/handler"
	transactionrepo "github.com/arfan21/vocagame/internal/transaction/repository"
	transactionuc "github.com/arfan21/vocagame/internal/transaction/usecase"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Routes() {

	api := s.app.Group("/api")
	api.Get("/health-check", func(c *fiber.Ctx) error {

		if err := s.db.Ping(c.UserContext()); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusOK)
	})

	midtransCoreAPI := paymentmidtrans.NewMidtransCoreAPI()

	productRepo := productrepo.New(s.db)
	productUsecase := productusecase.New(productRepo)
	productHandler := producthandler.NewHTTP(productUsecase)

	paymentMethodRepo := paymentmethodrepo.New(s.db)
	paymentMethodUsecase := paymentmethoduc.New(paymentMethodRepo)
	paymentMethodHandler := paymentmethodhandler.NewHTTP(paymentMethodUsecase)

	transactionRepo := transactionrepo.New(s.db)
	transactionUsecase := transactionuc.New(
		transactionRepo,
		productUsecase,
		paymentMethodUsecase,
		midtransCoreAPI,
	)
	transactionHandler := transactionhandler.NewHTTP(transactionUsecase)

	s.ProductRoutes(productHandler)
	s.PaymentMethodRoutes(paymentMethodHandler)
	s.TransactionRoutes(transactionHandler)
}

func (s Server) ProductRoutes(handler *producthandler.HTTP) {
	api := s.app.Group("/api/v1/products")
	api.Get("/", handler.GetList)
}

func (s Server) PaymentMethodRoutes(handler *paymentmethodhandler.HTTP) {
	api := s.app.Group("/api/v1/payment-methods")
	api.Get("/", handler.GetList)
}

func (s Server) TransactionRoutes(handler *transactionhandler.HTTP) {
	api := s.app.Group("/api/v1/transactions")
	api.Post("/", handler.Create)
	api.Get("/tracking", handler.Tracking)

	callback := api.Group("/callback")
	callback.Post("/midtrans", handler.MidtransNotification)
	callback.Get("/success", handler.SuccessCallback)
}
