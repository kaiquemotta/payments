package usecase

import (
	"payments/domain"
	"payments/repository"
)

type PaymentUseCase interface {
	GetAllPayments() ([]domain.Payment, error)
	GetPaymentByID(id string) (domain.Payment, error)
	CreatePayment(payment *domain.Payment) error
	UpdatePayment(id string, payment *domain.Payment) error
	DeletePayment(id string) error
}

type paymentUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewPaymentUseCase(repo repository.PaymentRepository) PaymentUseCase {
	return &paymentUseCase{repo}
}

func (uc *paymentUseCase) GetAllPayments() ([]domain.Payment, error) {
	return uc.paymentRepo.GetAll()
}

func (uc *paymentUseCase) GetPaymentByID(id string) (domain.Payment, error) {
	return uc.paymentRepo.GetByID(id)
}

func (uc *paymentUseCase) CreatePayment(payment *domain.Payment) error {
	return uc.paymentRepo.Create(payment)
}

func (uc *paymentUseCase) UpdatePayment(id string, payment *domain.Payment) error {
	return uc.paymentRepo.Update(id, payment)
}

func (uc *paymentUseCase) DeletePayment(id string) error {
	return uc.paymentRepo.Delete(id)
}
