package invoice

import (
	"context"
)

type InvoiceService struct {
	repo InvoiceRepository
}

func NewInvoiceService(repo InvoiceRepository) *InvoiceService {
	return &InvoiceService{repo: repo}
}

func (s *InvoiceService) GetAllInvoices(ctx context.Context) ([]Invoice, error) {
	return s.repo.GetAll(ctx)
}

func (s *InvoiceService) GetPagedInvoices(ctx context.Context, sortBy, sortDir string, offset, limit int) ([]Invoice, error) {
	return s.repo.GetPaged(ctx, sortBy, sortDir, offset, limit)
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, invoice *Invoice) error {
	return s.repo.AddInvoice(ctx, invoice)
}

func (s *InvoiceService) CreateInvoices(ctx context.Context, invoices []Invoice) error {
	return s.repo.CreateBatch(ctx, invoices)
}

func (s *InvoiceService) UpdateInvoice(ctx context.Context, invoice *Invoice) error {
	return s.repo.UpdateInvoice(ctx, invoice)
}

func (s *InvoiceService) DeleteInvoice(ctx context.Context, id int) error {
	return s.repo.DeleteInvoice(ctx, id)
}

func (s *InvoiceService) GetDelayedInvoicesCount(ctx context.Context) (int, error) {
	return s.repo.GetDelayedInvoicesCount(ctx)
}

func (s *InvoiceService) GetDelayedInvoicesAmount(ctx context.Context) (float64, error) {
	return s.repo.GetDelayedInvoicesAmount(ctx)
}

func (s *InvoiceService) GetDelayedInvoices(ctx context.Context) ([]Invoice, error) {
	return s.repo.GetDelayedInvoices(ctx)
}

func (s *InvoiceService) GetPendingInvoicesCount(ctx context.Context) (int, error) {
	return s.repo.GetPendingInvoicesCount(ctx)
}

func (s *InvoiceService) GetPendingInvoicesAmount(ctx context.Context) (float64, error) {
	return s.repo.GetPendingInvoicesAmount(ctx)
}

func (s *InvoiceService) GetPendingInvoices(ctx context.Context) ([]Invoice, error) {
	return s.repo.GetPendingInvoices(ctx)
}
