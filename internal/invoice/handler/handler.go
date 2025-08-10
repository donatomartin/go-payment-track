package handler

import (
	"app/internal/invoice/repository"
	"log"
	"net/http"
)

type ApiInvoiceHandler struct {
	repo   repository.InvoiceRepository
	logger *log.Logger
	mux    *http.ServeMux
}

func NewApiInvoiceHandler(repo repository.InvoiceRepository, logger *log.Logger, mux *http.ServeMux) *ApiInvoiceHandler {
	return &ApiInvoiceHandler{
		repo:   repo,
		logger: logger,
		mux:    mux,
	}
}

func (h *ApiInvoiceHandler) RegisterRoutes() {
	h.mux.HandleFunc("/api/v1/invoices/add", h.addInvoice)
	h.mux.HandleFunc("/api/v1/invoices/update", h.updateInvoice)
}
