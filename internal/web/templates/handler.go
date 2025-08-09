package templates

import (
	"log"
	"net/http"

	"app/internal/invoice"
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
	ShowSizeSelector bool
	HtmxFragmentName string
	FirstPage        int
	PrevPage         int
	Page             int
	NextPage         int
	LastPage         int
	Size             int
	SortBy           string
	SorDir           string
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
	TotalPaid    string
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
	mux.HandleFunc("/invoices/fragment", h.getInvoicesFragment)
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

func invoicesToInvoiceViews(invoices []invoice.Invoice) []InvoiceView {
	var invoiceViews []InvoiceView
	for _, invoice := range invoices {
		invoiceViews = append(invoiceViews, InvoiceView{
			ID:           invoice.ID,
			CustomerName: invoice.CustomerName,
			AmountDue:    util.Float64ToEuros(invoice.AmountDue),
			PaymentMean:  invoice.PaymentMean,
			InvoiceDate:  invoice.InvoiceDate.Format("2006-01-02"),
			DueDate:      invoice.DueDate.Format("2006-01-02"),
			TotalPaid:    util.Float64ToEuros(invoice.TotalPaid),
		})
	}
	return invoiceViews
}
