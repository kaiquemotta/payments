package delivery

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"payments/domain"
	"payments/usecase"
)

type PaymentHandler struct {
	useCase usecase.PaymentUseCase
}

func NewPaymentHandler(useCase usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{useCase: useCase}
}

// GetAllPayments retorna todos os pagamentos
func (h *PaymentHandler) GetAllPayments(c *fiber.Ctx) error {
	payments, err := h.useCase.GetAllPayments()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(payments)
}

// GetPaymentByID retorna um pagamento específico por ID
func (h *PaymentHandler) GetPaymentByID(c *fiber.Ctx) error {
	id := c.Params("id")
	payment, err := h.useCase.GetPaymentByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "Payment not found"})
	}
	return c.Status(fiber.StatusOK).JSON(payment)
}

// CreatePayment cria um novo pagamento
func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	var payment domain.Payment
	// Tenta fazer o parsing do corpo da requisição
	if err := c.BodyParser(&payment); err != nil {
		// Retorna erro 400 (Bad Request) se a requisição estiver malformada
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input: " + err.Error())
	}

	// Chama o caso de uso para criar o pagamento e obter o UUID
	uid, err := h.useCase.CreatePayment(&payment)
	if err != nil {
		// Verifica se o erro é de tipo inválido (InvalidPaymentTypeError)
		var invalidErr *domain.InvalidPaymentTypeError
		if errors.As(err, &invalidErr) {
			// Retorna erro 422 (Unprocessable Entity) se o tipo de pagamento for inválido
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": invalidErr.Error(),
			})
		}
		// Caso contrário, retorna erro 500 (Internal Server Error)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error: " + err.Error())
	}

	// Retorna a resposta 201 (Created) com o UUID gerado
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"uuid": uid,
	})
}

// UpdatePayment atualiza um pagamento existente
func (h *PaymentHandler) UpdatePayment(c *fiber.Ctx) error {
	id := c.Params("id")
	var payment domain.Payment
	if err := c.BodyParser(&payment); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.useCase.UpdatePayment(id, &payment); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "Payment not found"})
	}
	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *PaymentHandler) DeletePayment(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.useCase.DeletePayment(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "Payment not found"})
	}
	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (ph *PaymentHandler) Callback(c *fiber.Ctx) error {
	var callbackData domain.PaymentCallback

	if err := c.BodyParser(&callbackData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := ph.useCase.ProcessPaymentCallback(&callbackData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error processing payment callback",
		})
	}

	// Retornando um "OK" ao serviço de pagamento
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment processed successfully",
	})
}
