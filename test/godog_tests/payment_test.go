package godog_tests

import (
	"errors"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"payments/domain"
	"payments/usecase"
)

// Mock do repositório
type mockPaymentRepository struct {
	CreateFunc func(payment *domain.Payment) (string, error)
}

func (m *mockPaymentRepository) Create(payment *domain.Payment) (string, error) {
	return m.CreateFunc(payment)
}

// Outros métodos vazios para satisfazer a interface
func (m *mockPaymentRepository) GetAll() ([]domain.Payment, error) { return nil, nil }
func (m *mockPaymentRepository) GetByID(id string) (domain.Payment, error) {
	return domain.Payment{}, nil
}
func (m *mockPaymentRepository) Update(id string, payment *domain.Payment) error { return nil }
func (m *mockPaymentRepository) Delete(id string) error                          { return nil }

// Variáveis globais usadas nos testes
var (
	payment   *domain.Payment
	useCase   *usecase.PaymentUseCase
	lastError error
)

// Passos do Gherkin
func queTenhoUmPagamentoValido() error {
	// Criando um pagamento válido com tipo Pix
	payment = &domain.Payment{PaymentType: domain.Pix}
	return nil
}

func euTentoCriarOPagamento() error {
	// Mock do repositório com a lógica para simular sucesso
	mockRepo := &mockPaymentRepository{
		CreateFunc: func(payment *domain.Payment) (string, error) {
			if payment.PaymentType == domain.Pix { // Simulando sucesso para tipo Pix
				return "payment-123", nil // Pagamento criado com sucesso
			}
			return "", errors.New("invalid payment type")
		},
	}

	// Inicializando o caso de uso com o repositório simulado
	useCase := usecase.NewPaymentUseCase(mockRepo)
	_, lastError = useCase.CreatePayment(payment)
	return nil
}

func oPagamentoDeveSerCriadoComSucesso() error {
	if lastError != nil {
		return errors.New("esperava sucesso, mas ocorreu um erro")
	}
	return nil
}

// Inicializa os testes
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^que tenho um pagamento válido$`, queTenhoUmPagamentoValido)
	ctx.Step(`^eu tentar criar o pagamento$`, euTentoCriarOPagamento)
	ctx.Step(`^o pagamento deve ser criado com sucesso$`, oPagamentoDeveSerCriadoComSucesso)
}

// Rodando os testes
func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty", // Define o formato da saída
		Paths:  []string{"../../features/payment.feature"},
	}
	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()
	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
