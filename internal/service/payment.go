package service

import (
	"context"
	"pagos-cesar/internal/repository"
)

type PaymentService interface {
	GetAllPayments(ctx context.Context) ([]repository.Payment, error)
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}

func (s *paymentService) GetAllPayments(ctx context.Context) ([]repository.Payment, error) {
	return s.repo.GetAll(ctx)
}
