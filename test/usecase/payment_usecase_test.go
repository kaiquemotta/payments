package usecase

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"payments/domain"
	usecase2 "payments/usecase"
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
	useCase := usecase2.NewPaymentUseCase(mockRepo)

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

func TestGetPaymentByID(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	useCase := usecase2.NewPaymentUseCase(mockRepo)

	id := "507f191e810c19729de860ea"
	objectID, _ := primitive.ObjectIDFromHex(id)

	mockPayment := domain.Payment{
		ID:     objectID,
		Amount: 150.0, // Valor fictício para teste
	}

	mockRepo.On("GetByID", id).Return(mockPayment, nil)

	payment, err := useCase.GetPaymentByID(id)

	assert.Nil(t, err)
	assert.Equal(t, mockPayment.ID, payment.ID)
	assert.Equal(t, mockPayment.Amount, payment.Amount)

	mockRepo.AssertExpectations(t)
}

func TestCreatePayment(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	useCase := usecase2.NewPaymentUseCase(mockRepo)

	payment := &domain.Payment{
		ID:          primitive.NewObjectID(),
		Amount:      150.0,
		PaymentType: domain.Pix, // Supondo que seja um enum válido
	}

	mockRepo.On("Create", payment).Return("generated-id", nil)

	// Chama o método de criação de pagamento
	id, err := useCase.CreatePayment(payment)

	// Verificações
	assert.Nil(t, err)
	assert.Equal(t, "generated-id", id)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePayment(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	useCase := usecase2.NewPaymentUseCase(mockRepo)

	id := "507f191e810c19729de860ea"
	existingPayment := domain.Payment{
		ID:     primitive.NewObjectID(),
		Amount: 200.0,
	}

	updatedPayment := &domain.Payment{
		ID:     existingPayment.ID,
		Amount: 250.0,
	}

	// Mocka a busca pelo pagamento existente
	mockRepo.On("GetByID", id).Return(existingPayment, nil)

	// Mocka a atualização do pagamento
	mockRepo.On("Update", id, updatedPayment).Return(nil)

	// Executa o método
	err := useCase.UpdatePayment(id, updatedPayment)

	// Verificações
	assert.Nil(t, err)
	mockRepo.AssertExpectations(t)
}
func TestDeletePayment(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	useCase := usecase2.NewPaymentUseCase(mockRepo)

	paymentID := "507f191e810c19729de860ea"
	existingPayment := domain.Payment{
		ID:     primitive.NewObjectID(),
		Amount: 100.0,
	}

	// Cenário de sucesso: pagamento encontrado e deletado com sucesso
	mockRepo.On("GetByID", paymentID).Return(existingPayment, nil)
	mockRepo.On("Delete", paymentID).Return(nil)

	err := useCase.DeletePayment(paymentID)
	assert.Nil(t, err)

	// Cenário de erro: pagamento não encontrado
	mockRepo.On("GetByID", "invalid_id").Return(domain.Payment{}, fmt.Errorf("payment not found"))

	err = useCase.DeletePayment("invalid_id")
	assert.NotNil(t, err)
	assert.Equal(t, "payment not found: payment not found", err.Error())

	mockRepo.AssertExpectations(t)
}
