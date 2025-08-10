package data

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"app/internal/invoice"
	invoiceRepo "app/internal/invoice/repository"
	"app/internal/payment"
	paymentRepo "app/internal/payment/repository"

	"github.com/xuri/excelize/v2"
)

func ImportPagos(ctx context.Context, f *excelize.File, repo paymentRepo.PaymentRepository) error {
	rows, err := f.GetRows("PAGOS")
	if err != nil {
		return err
	}

	var payments []payment.Payment
	for i, row := range rows {
		if i < 2 {
			continue
		}
		if len(row) < 2 {
			break
		}

		id, _ := strconv.ParseInt(row[0], 10, 8)
		date, _ := time.Parse("01-02-06", row[2])

		invoiceID := strings.TrimSpace(row[1])
		amount := strings.TrimRight(strings.ReplaceAll(row[3], ",", ""), " €")
		floatAmount, _ := strconv.ParseFloat(amount, 64)

		payments = append(payments, payment.Payment{
			ID:        int(id),
			InvoiceID: invoiceID,
			Date:      date,
			Amount:    floatAmount,
		})
	}

	if err := repo.CreateBatch(ctx, payments); err != nil {
		return fmt.Errorf("creating payments: %w", err)
	}

	return nil
}

func ImportFacturas(ctx context.Context, f *excelize.File, repo invoiceRepo.InvoiceRepository) error {
	months := []string{
		"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
		"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
	}

	var invoices []invoice.Invoice
	for _, month := range months {
		rows, err := f.GetRows(month)
		if err != nil {
			return err
		}

		for j, row := range rows {
			if j < 2 {
				continue
			}
			if len(row) < 7 || row[1] == "" {
				break
			}

			date, _ := time.Parse("01-02-06", row[3])
			dueDate, _ := time.Parse("01-02-06", row[6])

			amount := strings.TrimRight(strings.ReplaceAll(row[4], ",", ""), " €")
			amountDue, _ := strconv.ParseFloat(amount, 64)

			id := strings.TrimSpace(row[2])

			invoices = append(invoices, invoice.Invoice{
				ID:           id,
				CustomerName: row[1],
				AmountDue:    amountDue,
				InvoiceDate:  date,
				DueDate:      dueDate,
				PaymentMean:  row[5],
			})
		}
	}

	if err := repo.CreateBatch(ctx, invoices); err != nil {
		return fmt.Errorf("creating invoices: %w", err)
	}

	return nil
}
