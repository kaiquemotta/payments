package delivery

import (
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
