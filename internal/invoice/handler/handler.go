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
}
