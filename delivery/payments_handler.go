package delivery

import (
	"github.com/gofiber/fiber/v2"
	"payments/domain"
	"payments/usecase"
)

type PaymentHandler struct {
	paymentUseCase usecase.PaymentUseCase
}

func NewPaymentHandler(app *fiber.App, useCase usecase.PaymentUseCase) {
	handler := &PaymentHandler{useCase}

	app.Get("/payments", handler.GetAllPayments)
	app.Get("/payments/:id", handler.GetPaymentByID)
	app.Post("/payments", handler.CreatePayment)
	app.Put("/payments/:id", handler.UpdatePayment)
	app.Delete("/payments/:id", handler.DeletePayment)
}

func (h *PaymentHandler) GetAllPayments(c *fiber.Ctx) error {
	payments, err := h.paymentUseCase.GetAllPayments()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(payments)
}

func (h *PaymentHandler) GetPaymentByID(c *fiber.Ctx) error {
	id := c.Params("id")
	payment, err := h.paymentUseCase.GetPaymentByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Payment not found"})
	}
	return c.JSON(payment)
}

func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	var payment domain.Payment
	if err := c.BodyParser(&payment); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := h.paymentUseCase.CreatePayment(&payment)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(payment)
}

func (h *PaymentHandler) UpdatePayment(c *fiber.Ctx) error {
	id := c.Params("id")
	var payment domain.Payment
	if err := c.BodyParser(&payment); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := h.paymentUseCase.UpdatePayment(id, &payment)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(payment)
}

func (h *PaymentHandler) DeletePayment(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.paymentUseCase.DeletePayment(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
