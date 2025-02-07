package routes

import (
	"github.com/gofiber/fiber/v2"
	"payments/delivery"
	"payments/usecase"
)

// RegisterPaymentRoutes registra as rotas da API de pagamentos.
func RegisterPaymentRoutes(app *fiber.App, useCase usecase.PaymentUseCase) {
	// Registrando as rotas de pagamento
	app.Get("/payments", delivery.GetAllPayments(useCase))
	app.Get("/payments/:id", delivery.GetPaymentByID(useCase))
	app.Post("/payments", delivery.CreatePayment(useCase))
	app.Put("/payments/:id", delivery.UpdatePayment(useCase))
	app.Delete("/payments/:id", delivery.DeletePayment(useCase))
}
