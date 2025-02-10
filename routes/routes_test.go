package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"payments/domain"
	"strings"
	"testing"
)

type MockPaymentUseCase struct {
	mock.Mock
}

func (m *MockPaymentUseCase) GetAllPayments() ([]domain.Payment, error) {
	args := m.Called()
	return args.Get(0).([]domain.Payment), args.Error(1)
}

func (m *MockPaymentUseCase) GetPaymentByID(id string) (domain.Payment, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Payment), args.Error(1)
}

func (m *MockPaymentUseCase) CreatePayment(payment *domain.Payment) (string, error) {
	args := m.Called(payment)
	return args.String(0), args.Error(1)
}

func (m *MockPaymentUseCase) UpdatePayment(id string, payment *domain.Payment) error {
	args := m.Called(id, payment)
	return args.Error(0)
}

func (m *MockPaymentUseCase) DeletePayment(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPaymentUseCase) ProcessPaymentCallback(callback *domain.PaymentCallback) error {
	args := m.Called(callback)
	return args.Error(0)
}

func TestRegisterPaymentRoutes(t *testing.T) {
	mockUseCase := new(MockPaymentUseCase)

	mockUseCase.On("GetAllPayments").Return([]domain.Payment{}, nil)
	mockUseCase.On("GetPaymentByID", "1").Return(domain.Payment{}, nil)
	mockUseCase.On("CreatePayment", mock.Anything).Return("123", nil)
	mockUseCase.On("UpdatePayment", "1", mock.Anything).Return(nil)
	mockUseCase.On("DeletePayment", "1").Return(nil)
	mockUseCase.On("ProcessPaymentCallback", mock.Anything).Return(nil)

	app := fiber.New()

	RegisterPaymentRoutes(app, mockUseCase)

	t.Run("Test GetAllPayments Route", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/payments", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Test GetPaymentByID Route", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/payments/1", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Test CreatePayment Route", func(t *testing.T) {
		reqBody := strings.NewReader(`{
			  "order_id": "123456",
			  "amount": 100.50,
			  "method": "credit_card",
			  "status": "pending",
			  "payment_type": "online"
			}
		`)
		req := httptest.NewRequest("POST", "/payments", reqBody)
		req.Header.Set("Content-Type", "application/json") // Definir cabeçalho correto
		resp, _ := app.Test(req)
		assert.Equal(t, 201, resp.StatusCode)
	})

	t.Run("Test UpdatePayment Route", func(t *testing.T) {
		paymentJSON := `{
		"order_id": "123456",
		"amount": 150.75,
		"method": "credit_card",
		"status": "approved",
		"payment_type": "online"
	}`

		req := httptest.NewRequest("PUT", "/payments/1", strings.NewReader(paymentJSON))
		req.Header.Set("Content-Type", "application/json") // Definir o Content-Type é essencial

		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Test DeletePayment Route", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/payments/1", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Test Callback Route", func(t *testing.T) {
		reqBody := strings.NewReader(`    {
			"id": "67a8ffa093a5fa72f000452b",
			"sale_id": "11111",
			"amount": 20,
			"method": "123",
			"status": "teste",
			"payment_type": "PIX"
		}
		`)
		req := httptest.NewRequest("POST", "/payment/callback", reqBody)
		req.Header.Set("Content-Type", "application/json") // Definir cabeçalho correto
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})
	mockUseCase.AssertExpectations(t)
}
