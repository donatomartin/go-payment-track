package payment

import (
	"context"
)

type PaymentService interface {
	GetAllPayments(ctx context.Context) ([]Payment, error)
	GetPagedPayments(ctx context.Context, sortBy, sortDir string, offset, limit int) ([]Payment, error)
	CreatePayment(ctx context.Context, payment *Payment) error
	CreatePayments(ctx context.Context, payments []Payment) error
	UpdatePayment(ctx context.Context, payment *Payment) error
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

func (s *paymentService) GetPagedPayments(ctx context.Context, sortBy, sortDir string, offset, limit int) ([]Payment, error) {
	return s.repo.GetPaged(ctx, sortBy, sortDir, offset, limit)
}

func (s *paymentService) CreatePayment(ctx context.Context, payment *Payment) error {
	return s.repo.AddPayment(ctx, payment)
}

func (s *paymentService) CreatePayments(ctx context.Context, payments []Payment) error {
	return s.repo.CreateBatch(ctx, payments)
}

func (s *paymentService) UpdatePayment(ctx context.Context, payment *Payment) error {
	return s.repo.UpdatePayment(ctx, payment)
}
