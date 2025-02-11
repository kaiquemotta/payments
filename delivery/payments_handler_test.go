package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http/httptest"
	"payments/domain"
	"testing"
)

// MockPaymentUseCase é um mock para o PaymentUseCase
type MockPaymentUseCase struct {
	mock.Mock
}

// GetAllPayments é o método mockado para retornar uma lista de pagamentos ou erro
func (m *MockPaymentUseCase) GetAllPayments() ([]domain.Payment, error) {
	args := m.Called()
	return args.Get(0).([]domain.Payment), args.Error(1)
}

// GetPaymentByID é um método mockado para atender à interface PaymentUseCase
func (m *MockPaymentUseCase) GetPaymentByID(id string) (domain.Payment, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Payment), args.Error(1)
}

// CreatePayment é um método mockado para atender à interface PaymentUseCase
func (m *MockPaymentUseCase) CreatePayment(payment *domain.Payment) (string, error) {
	args := m.Called(payment)
	return args.String(0), args.Error(1)
}

// UpdatePayment é um método mockado para atender à interface PaymentUseCase
func (m *MockPaymentUseCase) UpdatePayment(id string, payment *domain.Payment) error {
	args := m.Called(id, payment)
	return args.Error(0)
}

// DeletePayment é um método mockado para atender à interface PaymentUseCase
func (m *MockPaymentUseCase) DeletePayment(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// ProcessPaymentCallback é o método que estava faltando, agora adicionado
func (m *MockPaymentUseCase) ProcessPaymentCallback(paymentCallback *domain.PaymentCallback) error {
	args := m.Called(paymentCallback)
	return args.Error(0)
}

func TestPaymentHandler_GetAllPayments_Success(t *testing.T) {
	mockUseCase := new(MockPaymentUseCase)
	handler := NewPaymentHandler(mockUseCase)
	mockPayments := []domain.Payment{
		{ID: primitive.NewObjectID(), Amount: 100},
		{ID: primitive.NewObjectID(), Amount: 200},
	}
	mockUseCase.On("GetAllPayments").Return(mockPayments, nil)
	app := fiber.New()
	app.Get("/payments", handler.GetAllPayments)
	req := httptest.NewRequest("GET", "/payments", nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	var payments []domain.Payment
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&payments))
	assert.Equal(t, len(payments), len(mockPayments))
}

func TestPaymentHandler_GetAllPayments_Error(t *testing.T) {
	mockUseCase := new(MockPaymentUseCase)
	handler := NewPaymentHandler(mockUseCase)
	mockUseCase.On("GetAllPayments").Return([]domain.Payment{}, errors.New("Internal server error"))
	app := fiber.New()
	app.Get("/payments", handler.GetAllPayments)
	req := httptest.NewRequest("GET", "/payments", nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestPaymentHandler_GetPaymentByID_Success(t *testing.T) {
	mockUseCase := new(MockPaymentUseCase)
	handler := NewPaymentHandler(mockUseCase)

	mockPayment := domain.Payment{
		ID:     primitive.NewObjectID(),
		Amount: 150,
	}
	mockUseCase.On("GetPaymentByID", mockPayment.ID.Hex()).Return(mockPayment, nil)
	app := fiber.New()
	app.Get("/payments/:id", handler.GetPaymentByID)
	req := httptest.NewRequest("GET", "/payments/"+mockPayment.ID.Hex(), nil)
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var payment domain.Payment
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&payment))
	assert.Equal(t, mockPayment.ID, payment.ID)
	assert.Equal(t, mockPayment.Amount, payment.Amount)
}

func TestPaymentHandler_GetPaymentByID_NotFound(t *testing.T) {
	mockUseCase := new(MockPaymentUseCase)
	handler := NewPaymentHandler(mockUseCase)

	mockUseCase.On("GetPaymentByID", "invalid-id").Return(domain.Payment{}, errors.New("Payment not found"))

	app := fiber.New()
	app.Get("/payments/:id", handler.GetPaymentByID)

	req := httptest.NewRequest("GET", "/payments/invalid-id", nil)
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestPaymentHandler_CreatePayment_Success(t *testing.T) {
	mockUseCase := new(MockPaymentUseCase)
	handler := NewPaymentHandler(mockUseCase)

	mockPayment := domain.Payment{
		Amount: 150,
	}
	mockUUID := "123e4567-e89b-12d3-a456-426614174000"
	mockUseCase.On("CreatePayment", &mockPayment).Return(mockUUID, nil)

	app := fiber.New()
	app.Post("/payments", handler.CreatePayment)

	body, _ := json.Marshal(mockPayment)
	req := httptest.NewRequest("POST", "/payments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response map[string]string
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&response))
	assert.Equal(t, mockUUID, response["uuid"])
}

func TestPaymentHandler_CreatePayment_InvalidInput(t *testing.T) {
	mockUseCase := new(MockPaymentUseCase)
	handler := NewPaymentHandler(mockUseCase)

	app := fiber.New()
	app.Post("/payments", handler.CreatePayment)

	req := httptest.NewRequest("POST", "/payments", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestPaymentHandler_CreatePayment_UnprocessableEntity(t *testing.T) {
	mockUseCase := new(MockPaymentUseCase)
	handler := NewPaymentHandler(mockUseCase)

	mockPayment := domain.Payment{
		Amount: 150,
	}
	mockUseCase.On("CreatePayment", &mockPayment).Return("", &domain.InvalidPaymentTypeError{Type: "CreditCard"})

	app := fiber.New()
	app.Post("/payments", handler.CreatePayment)

	body, _ := json.Marshal(mockPayment)
	req := httptest.NewRequest("POST", "/payments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

func TestPaymentHandler_CreatePayment_InternalServerError(t *testing.T) {
	mockUseCase := new(MockPaymentUseCase)
	handler := NewPaymentHandler(mockUseCase)

	mockPayment := domain.Payment{
		Amount: 150,
	}
	mockUseCase.On("CreatePayment", &mockPayment).Return("", errors.New("unexpected error"))

	app := fiber.New()
	app.Post("/payments", handler.CreatePayment)

	body, _ := json.Marshal(mockPayment)
	req := httptest.NewRequest("POST", "/payments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestPaymentHandler_UpdatePayment_Success(t *testing.T) {
	mockUseCase := new(MockPaymentUseCase)
	handler := NewPaymentHandler(mockUseCase)

	paymentID := "60c72b2f9af1c88b8f8d3b4a"
	mockPayment := domain.Payment{Amount: 300}

	mockUseCase.On("UpdatePayment", paymentID, &mockPayment).Return(nil)

	app := fiber.New()
	app.Put("/payments/:id", handler.UpdatePayment)

	body, _ := json.Marshal(mockPayment)
	req := httptest.NewRequest("PUT", "/payments/"+paymentID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func TestPaymentHandler_UpdatePayment_BadRequest(t *testing.T) {
	mockUseCase := new(MockPaymentUseCase)
	handler := NewPaymentHandler(mockUseCase)

	app := fiber.New()
	app.Put("/payments/:id", handler.UpdatePayment)

	invalidBody := "{invalid json}"
	req := httptest.NewRequest("PUT", "/payments/60c72b2f9af1c88b8f8d3b4a", bytes.NewReader([]byte(invalidBody)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestPaymentHandler_UpdatePayment_NotFound(t *testing.T) {
	mockUseCase := new(MockPaymentUseCase)
	handler := NewPaymentHandler(mockUseCase)

	paymentID := "60c72b2f9af1c88b8f8d3b4a"
	mockPayment := domain.Payment{Amount: 300}

	mockUseCase.On("UpdatePayment", paymentID, &mockPayment).Return(errors.New("Payment not found"))

	app := fiber.New()
	app.Put("/payments/:id", handler.UpdatePayment)

	body, _ := json.Marshal(mockPayment)
	req := httptest.NewRequest("PUT", "/payments/"+paymentID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestPaymentHandler_DeletePayment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(MockPaymentUseCase)
		handler := NewPaymentHandler(mockUseCase)

		mockUseCase.On("DeletePayment", "valid-id").Return(nil)

		app := fiber.New()
		app.Delete("/payments/:id", handler.DeletePayment)

		req := httptest.NewRequest("DELETE", "/payments/valid-id", nil)
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("Payment Not Found", func(t *testing.T) {
		mockUseCase := new(MockPaymentUseCase)
		handler := NewPaymentHandler(mockUseCase)

		mockUseCase.On("DeletePayment", "invalid-id").Return(errors.New("payment not found"))

		app := fiber.New()
		app.Delete("/payments/:id", handler.DeletePayment)

		req := httptest.NewRequest("DELETE", "/payments/invalid-id", nil)
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]string
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Payment not found", response["error"])
	})
}

func TestPaymentHandler_Callback(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(MockPaymentUseCase)
		handler := NewPaymentHandler(mockUseCase)

		callbackData := domain.PaymentCallback{
			PaymentID: "123456",
			Status:    "success",
			Amount:    100.50,
			Message:   "Payment received",
			OrderId:   "order789",
		}

		mockUseCase.On("ProcessPaymentCallback", &callbackData).Return(nil)

		app := fiber.New()
		app.Post("/callback", handler.Callback)

		body, _ := json.Marshal(callbackData)
		req := httptest.NewRequest("POST", "/callback", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]string
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Payment processed successfully", response["message"])
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		mockUseCase := new(MockPaymentUseCase)
		handler := NewPaymentHandler(mockUseCase)

		app := fiber.New()
		app.Post("/callback", handler.Callback)

		req := httptest.NewRequest("POST", "/callback", bytes.NewReader([]byte("{invalid json}")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]string
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Invalid request body", response["error"])
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockUseCase := new(MockPaymentUseCase)
		handler := NewPaymentHandler(mockUseCase)

		callbackData := domain.PaymentCallback{
			PaymentID: "123456",
			Status:    "failed",
			Amount:    100.50,
			Message:   "Payment error",
			OrderId:   "order789",
		}

		mockUseCase.On("ProcessPaymentCallback", &callbackData).Return(errors.New("processing error"))

		app := fiber.New()
		app.Post("/callback", handler.Callback)

		body, _ := json.Marshal(callbackData)
		req := httptest.NewRequest("POST", "/callback", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]string
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Error processing payment callback", response["error"])
	})
}
