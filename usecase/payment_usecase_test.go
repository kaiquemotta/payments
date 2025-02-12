package usecase

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"payments/domain"
	"testing"
)

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) GetAll() ([]domain.Payment, error) {
	args := m.Called()
	return args.Get(0).([]domain.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetByID(id string) (domain.Payment, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Create(payment *domain.Payment) (string, error) {
	args := m.Called(payment)
	return args.String(0), args.Error(1)
}

func (m *MockPaymentRepository) Update(id string, payment *domain.Payment) error {
	args := m.Called(id, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPaymentRepository) ProcessPaymentCallback(callbackData *domain.PaymentCallback) error {
	args := m.Called(callbackData)
	return args.Error(0)
}

func TestGetAllPayments(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	useCase := NewPaymentUseCase(mockRepo)

	id1, _ := primitive.ObjectIDFromHex("507f191e810c19729de860ea")
	id2, _ := primitive.ObjectIDFromHex("507f191e810c19729de860eb")

	mockPayments := []domain.Payment{
		{ID: id1, Amount: 100.0}, // Usando 100.0 como float64
		{ID: id2, Amount: 200.0}, // Usando 200.0 como float64
	}

	mockRepo.On("GetAll").Return(mockPayments, nil)

	payments, err := useCase.GetAllPayments()

	assert.Nil(t, err)
	assert.Equal(t, len(payments), 2)
	assert.Equal(t, payments[0].ID, id1)
	assert.Equal(t, payments[1].Amount, 200.0) // Esperando 200.0 como float64

	mockRepo.AssertExpectations(t)
}
