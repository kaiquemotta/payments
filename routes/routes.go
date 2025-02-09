// routes/routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"payments/delivery"
	"payments/usecase"
)

func RegisterPaymentRoutes(app *fiber.App, useCase usecase.PaymentUseCase) {
	handler := delivery.NewPaymentHandler(useCase)

	app.Get("/payments", handler.GetAllPayments)
	app.Get("/payments/:id", handler.GetPaymentByID)
	app.Post("/payments", handler.CreatePayment)
	app.Put("/payments/:id", handler.UpdatePayment)
	app.Delete("/payments/:id", handler.DeletePayment)
	app.Post("/payment/callback", handler.Callback)
}
