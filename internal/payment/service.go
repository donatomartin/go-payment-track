package payment

import (
	"context"
)

type PaymentService interface {
	GetAllPayments(ctx context.Context) ([]Payment, error)
}

type paymentService struct {
	repo PaymentRepository
}

func NewPaymentService(repo PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}

func (s *paymentService) GetAllPayments(ctx context.Context) ([]Payment, error) {
	return s.repo.GetAll(ctx)
}
