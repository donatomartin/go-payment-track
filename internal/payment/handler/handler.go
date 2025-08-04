package handler

import (
	"app/internal/payment/repository"
	"log"
	"net/http"
)

type ApiPaymentHandler struct {
	repo   repository.PaymentRepository
	logger *log.Logger
	mux    *http.ServeMux
}

func NewApiPaymentHandler(repo repository.PaymentRepository, logger *log.Logger, mux *http.ServeMux) *ApiPaymentHandler {
	return &ApiPaymentHandler{
		repo:   repo,
		logger: logger,
		mux:    mux,
	}
}

func (h *ApiPaymentHandler) RegisterRoutes() {
	h.mux.HandleFunc("/api/v1/payments/paged", h.getPagedPayments)
	h.mux.HandleFunc("/api/v1/payments", h.getPayments)
}
