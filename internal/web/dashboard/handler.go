package dashboard

import (
	"log"
	"net/http"

	invoiceRepo "app/internal/invoice/repository"
	"app/internal/payment"
	paymentRepo "app/internal/payment/repository"
	"app/internal/platform/util"
)

type DashboardHandler struct {
	logger      *log.Logger
	paymentRepo paymentRepo.PaymentRepository
	invoiceRepo invoiceRepo.InvoiceRepository
}

type Pagination struct {
	ShowPagination bool
	FirstPage      int
	PrevPage       int
	Page           int
	NextPage       int
	LastPage       int
	Size           int
	SortBy         string
	SorDir         string
}

func (p Pagination) GetOffset() int {
	return (p.Page - 1) * p.Size
}

type PaymentView struct {
	ID         int
	InvoiceID  string
	Amount     string
	Date       string
	ClientName string
}

type InvoiceView struct {
	ID           string
	CustomerName string
	AmountDue    string
	PaymentMean  string
	InvoiceDate  string
	DueDate      string
}

func NewDashboardHandler(

	paymentRepository paymentRepo.PaymentRepository,
	invoiceRepository invoiceRepo.InvoiceRepository,
	logger *log.Logger,

) *DashboardHandler {
	return &DashboardHandler{
		logger:      logger,
		paymentRepo: paymentRepository,
		invoiceRepo: invoiceRepository,
	}
}

func (h *DashboardHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.getDashboard) // Redirect root to dashboard
	mux.HandleFunc("/dashboard", h.getDashboard)
	mux.HandleFunc("/invoices", h.getInvoicesTable)
	mux.HandleFunc("/payments", h.getPaymentsTable)
	mux.HandleFunc("/payments/fragment", h.getPaymentsFragment)
}

func paymentsToPaymentViews(payments []payment.Payment) []PaymentView {
	var paymentViews []PaymentView
	for _, payment := range payments {
		paymentViews = append(paymentViews, PaymentView{
			ID:         payment.ID,
			InvoiceID:  payment.InvoiceID,
			Amount:     util.Float64ToEuros(payment.Amount),
			Date:       payment.Date.Format("2006-01-02"),
			ClientName: payment.ClientName,
		})
	}
	return paymentViews
}
