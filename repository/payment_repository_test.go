package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"payments/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

var db *mongo.Database
var paymentRepo PaymentRepository

// Setup de testes, criando uma conexão com o MongoDB em memória
func TestMain(m *testing.M) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("Error creating Mongo client:", err)
		return
	}
	err = client.Connect(context.Background())
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	db = client.Database("testdb")
	paymentRepo = NewPaymentRepository(db)

	// Limpa a coleção antes de começar os testes
	_ = db.Collection("payments").Drop(context.Background())

	m.Run()
}

// Testando o método GetAll
func TestPaymentRepository_GetAll(t *testing.T) {
	// Alterando ID para ObjectID diretamente
	payment := domain.Payment{
		ID:     primitive.NewObjectID(), // Gerando um novo ObjectID
		Amount: 100.50,
		Status: "completed",
	}

	_, err := paymentRepo.Create(&payment)
	assert.Nil(t, err)

	payments, err := paymentRepo.GetAll()
	assert.Nil(t, err)
	assert.Len(t, payments, 1)

	// Verificando se o ID está correto
	assert.Equal(t, payment.ID.Hex(), payments[0].ID.Hex()) // Hex() para comparar os ObjectIDs
}

// Testando o método GetByID
func TestPaymentRepository_GetByID(t *testing.T) {
	payment := domain.Payment{
		ID:     primitive.NewObjectID(), // Gerando um novo ObjectID
		Amount: 50.75,
		Status: "pending",
	}

	createdID, err := paymentRepo.Create(&payment)
	assert.Nil(t, err)

	// Tentando recuperar com um ID válido
	payment, err = paymentRepo.GetByID(createdID) // Convertendo ObjectID para string com Hex()
	assert.Nil(t, err)
	assert.Equal(t, payment.ID.Hex(), createdID) // Comparando ObjectIDs com Hex()

	// Tentando recuperar com um ID inválido
	_, err = paymentRepo.GetByID("invalid-id")
	assert.NotNil(t, err)
}

// Testando o método Create
func TestPaymentRepository_Create(t *testing.T) {
	payment := &domain.Payment{
		ID:     primitive.NewObjectID(), // Gerando um novo ObjectID
		Amount: 200.99,
		Status: "completed",
	}

	id, err := paymentRepo.Create(payment)
	assert.Nil(t, err)
	assert.NotEmpty(t, id)

	// Verifica se o pagamento foi criado corretamente
	storedPayment, err := paymentRepo.GetByID(id) // Convertendo ObjectID para string com Hex()
	assert.Nil(t, err)
	assert.Equal(t, payment.Amount, storedPayment.Amount)
	assert.Equal(t, payment.Status, storedPayment.Status)
}

// Testando o método Update
func TestPaymentRepository_Update(t *testing.T) {
	payment := domain.Payment{
		ID:     primitive.NewObjectID(), // Gerando um novo ObjectID
		Amount: 300.00,
		Status: "completed",
	}

	createdID, err := paymentRepo.Create(&payment)
	assert.Nil(t, err)

	// Atualizando o pagamento
	payment.Amount = 350.00
	err = paymentRepo.Update(createdID, &payment) // Convertendo ObjectID para string com Hex()
	assert.Nil(t, err)

	// Verifica se a atualização ocorreu corretamente
	updatedPayment, err := paymentRepo.GetByID(createdID) // Convertendo ObjectID para string com Hex()
	assert.Nil(t, err)
	assert.Equal(t, payment.Amount, updatedPayment.Amount)
}

// Testando o método Delete
func TestPaymentRepository_Delete(t *testing.T) {
	payment := domain.Payment{
		ID:     primitive.NewObjectID(), // Gerando um novo ObjectID
		Amount: 150.00,
		Status: "completed",
	}

	createdID, err := paymentRepo.Create(&payment)
	assert.Nil(t, err)

	// Deletando o pagamento
	err = paymentRepo.Delete(createdID) // Convertendo ObjectID para string com Hex()
	assert.Nil(t, err)

	// Verificando se o pagamento foi deletado
	_, err = paymentRepo.GetByID(createdID) // Convertendo ObjectID para string com Hex()
	assert.NotNil(t, err)
}
