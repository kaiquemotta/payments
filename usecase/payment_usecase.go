package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"payments/domain"
	"payments/repository"
)

type PaymentUseCase interface {
	GetAllPayments() ([]domain.Payment, error)
	GetPaymentByID(id string) (domain.Payment, error)
	CreatePayment(payment *domain.Payment) (string, error) // Atualizado para retornar UUID (string) e erro
	UpdatePayment(id string, payment *domain.Payment) error
	DeletePayment(id string) error
	ProcessPaymentCallback(paymentCallback *domain.PaymentCallback) error
}

type paymentUseCase struct {
	paymentRepo repository.PaymentRepository
}

// NewPaymentUseCase cria uma nova instância do PaymentUseCase
func NewPaymentUseCase(repo repository.PaymentRepository) PaymentUseCase {
	return &paymentUseCase{repo}
}

func (uc *paymentUseCase) GetAllPayments() ([]domain.Payment, error) {
	return uc.paymentRepo.GetAll()
}

func (uc *paymentUseCase) GetPaymentByID(id string) (domain.Payment, error) {
	return uc.paymentRepo.GetByID(id)
}

func (uc *paymentUseCase) CreatePayment(payment *domain.Payment) (string, error) {
	// A validação de tipo de pagamento será feita no domínio
	if err := payment.PaymentType.IsValid(); err != nil {
		return "", err // Erro retornado pela camada de domínio
	}
	// Chama o repositório para criar o pagamento e obter o ID gerado
	return uc.paymentRepo.Create(payment)
}

func (uc *paymentUseCase) UpdatePayment(id string, payment *domain.Payment) error {
	// Chama o use case para verificar se o pagamento existe
	_, err := uc.GetPaymentByID(id)
	if err != nil {
		// Retorna erro se o pagamento não for encontrado
		return fmt.Errorf("payment not found: %v", err)
	}
	// Se o pagamento existe, chama o método de atualização no repositório
	return uc.paymentRepo.Update(id, payment)
}

func (uc *paymentUseCase) DeletePayment(id string) error {
	// Chama o use case para verificar se o pagamento existe
	_, err := uc.GetPaymentByID(id)
	if err != nil {
		// Retorna erro se o pagamento não for encontrado
		return fmt.Errorf("payment not found: %v", err)
	}
	// Se o pagamento existe, chama o método de deleção
	return uc.paymentRepo.Delete(id)
}

func (uc *paymentUseCase) ProcessPaymentCallback(callbackData *domain.PaymentCallback) error {
	// Primeiro, processe o pagamento como já estava fazendo
	payment, err := uc.paymentRepo.GetByID(callbackData.PaymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %v", err)
	}
	payment.Status = callbackData.Status

	fmt.Println("Pagamento processado:", payment)

	// Após processar o pagamento, faça um POST para o microserviço de pedidos para atualizar o status do pedido
	err = uc.updateOrderStatus(payment.OrderId, callbackData.Status)
	if err != nil {
		return fmt.Errorf("error updating order status: %v", err)
	}

	return nil
}

func (uc *paymentUseCase) updateOrderStatus(orderID string, status string) error {
	orderUpdate := map[string]interface{}{
		"status": "Finalizado",
	}
	jsonData, err := json.Marshal(orderUpdate)
	if err != nil {
		return fmt.Errorf("error marshalling order data: %v", err)
	}

	url := "https://api-ms-order-6ec42f917adf.herokuapp.com/orders/" + orderID
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to update order status: received status code %d", resp.StatusCode)
	}

	return nil
}
