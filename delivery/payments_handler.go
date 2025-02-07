package delivery

import (
	"github.com/gofiber/fiber/v2"
	"payments/domain"
	"payments/usecase"
)

// @Description Get All Payments
// @ID get-all-payments
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Payment
// @Router /payments [get]
func GetAllPayments(useCase usecase.PaymentUseCase) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		payments, err := useCase.GetAllPayments()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(payments)
	}
}

// @Description Get Payment by ID
// @ID get-payment-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Payment ID"
// @Success 200 {object} domain.Payment
// @Failure 404 {object} map[string]string
// @Router /payments/{id} [get]
func GetPaymentByID(useCase usecase.PaymentUseCase) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		payment, err := useCase.GetPaymentByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "Payment not found"})
		}
		return c.Status(fiber.StatusOK).JSON(payment)
	}
}

// @Description Create a Payment
// @ID create-payment
// @Accept  json
// @Produce  json
// @Param payment body domain.Payment true "Create Payment"
// @Success 201 {object} domain.Payment
// @Failure 400 {object} map[string]string
// @Router /payments [post]
func CreatePayment(useCase usecase.PaymentUseCase) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var payment domain.Payment
		if err := c.BodyParser(&payment); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		if err := useCase.CreatePayment(&payment); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.Status(fiber.StatusCreated).JSON(payment)
	}
}

// @Description Update a Payment
// @ID update-payment
// @Accept  json
// @Produce  json
// @Param id path string true "Payment ID"
// @Param payment body domain.Payment true "Update Payment"
// @Success 200 {object} domain.Payment
// @Failure 404 {object} map[string]string
// @Router /payments/{id} [put]
func UpdatePayment(useCase usecase.PaymentUseCase) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var payment domain.Payment
		if err := c.BodyParser(&payment); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		if err := useCase.UpdatePayment(id, &payment); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "Payment not found"})
		}
		return c.Status(fiber.StatusOK).JSON(payment)
	}
}

// @Description Delete a Payment
// @ID delete-payment
// @Accept  json
// @Produce  json
// @Param id path string true "Payment ID"
// @Success 200 {string} string "Payment deleted"
// @Failure 404 {object} map[string]string
// @Router /payments/{id} [delete]
func DeletePayment(useCase usecase.PaymentUseCase) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := useCase.DeletePayment(id); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "Payment not found"})
		}
		return c.Status(fiber.StatusOK).SendString("Payment deleted")
	}
}
