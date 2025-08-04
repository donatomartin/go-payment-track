package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"app/internal/payment"
	paymentRepository "app/internal/payment/repository"
	"app/internal/platform/util"

	"github.com/xuri/excelize/v2"
)

func ImportPagos(f *excelize.File) {

	// Read all rows
	rows, err := f.GetRows("PAGOS")
	if err != nil {
		fmt.Println(err)
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

		invoiceID := row[1]
		invoiceID = strings.TrimSpace(invoiceID)

		amount := row[3]
		amount = strings.TrimRight(strings.ReplaceAll(amount, ",", ""), " €")
		float_amount, _ := strconv.ParseFloat(amount, 64)

		payment := &payment.Payment{
			ID:        int(id),
			InvoiceID: invoiceID,
			Date:      date,
			Amount:    float_amount,
		}

		payments = append(payments, *payment)

	}

	service := paymentRepository.NewPaymentRepository(util.GetDB())
	if err := service.CreateBatch(context.Background(), payments); err != nil {
		fmt.Printf("Error creating payments: %v\n", err)
		panic(err)
	}

	fmt.Printf("Successfully imported %d payments\n", len(payments))
}
