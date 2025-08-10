package handler

import (
	invoiceRepo "app/internal/invoice/repository"
	paymentRepo "app/internal/payment/repository"
	"log"
	"net/http"
)

type ApiAdminHandler struct {
	paymentRepo paymentRepo.PaymentRepository
	invoiceRepo invoiceRepo.InvoiceRepository
	logger      *log.Logger
	mux         *http.ServeMux
}

func NewApiAdminHandler(paymentRepo paymentRepo.PaymentRepository, invoiceRepo invoiceRepo.InvoiceRepository, logger *log.Logger, mux *http.ServeMux) *ApiAdminHandler {
	return &ApiAdminHandler{
		paymentRepo: paymentRepo,
		invoiceRepo: invoiceRepo,
		logger:      logger,
		mux:         mux,
	}
}

func (h *ApiAdminHandler) RegisterRoutes() {
	h.mux.HandleFunc("/api/v1/admin/clear", h.clearDatabase)
	h.mux.HandleFunc("/api/v1/admin/import", h.importData)
	h.mux.HandleFunc("/api/v1/admin/export", h.exportData)
}
