// routes/routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"payments/delivery" // Para importar o handler (delivery)
	"payments/usecase"
)

// RegisterPaymentRoutes registra todas as rotas de pagamentos
func RegisterPaymentRoutes(app *fiber.App, useCase usecase.PaymentUseCase) {
	// Cria uma instância do handler, que é responsável por interagir com a delivery
	handler := delivery.NewPaymentHandler(useCase)

	// Registra as rotas com seus respectivos métodos
	app.Get("/payments", handler.GetAllPayments)
	app.Get("/payments/:id", handler.GetPaymentByID)
	app.Post("/payments", handler.CreatePayment)
	app.Put("/payments/:id", handler.UpdatePayment)
	app.Delete("/payments/:id", handler.DeletePayment)
}
