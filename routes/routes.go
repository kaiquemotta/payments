package routes

import (
	"github.com/gofiber/fiber/v2"
	"payments/delivery"
	"payments/usecase"
)

// RegisterPaymentRoutes registra as rotas relacionadas aos pagamentos
func RegisterPaymentRoutes(app *fiber.App, useCase usecase.PaymentUseCase) {
	delivery.NewPaymentHandler(app, useCase)
}
