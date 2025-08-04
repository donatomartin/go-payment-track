package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"app/internal/invoice"
	invoiceRepository "app/internal/invoice/repository"
	"app/internal/platform/util"

	"github.com/xuri/excelize/v2"
)

func ImportFacturas(f *excelize.File) {

	var invoices []invoice.Invoice

	months := []string{
		"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
		"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
	}

	for _, month := range months {

		// Read all rows
		rows, err := f.GetRows(month)
		if err != nil {
			fmt.Println(err)
		}

		for j, row := range rows {

			if j < 2 {
				continue
			}

			if row[1] == "" {
				break
			}

			date, _ := time.Parse("01-02-06", row[3])
			due_date, _ := time.Parse("01-02-06", row[6])

			amount := row[4]
			amount = strings.TrimRight(strings.ReplaceAll(amount, ",", ""), " €")
			amount_due, _ := strconv.ParseFloat(amount, 64)

			id := row[2]
			id = strings.TrimSpace(id)

			invoice := &invoice.Invoice{
				ID:           id,
				CustomerName: row[1],
				AmountDue:    amount_due,
				InvoiceDate:  date,
				DueDate:      due_date,
				PaymentMean:  row[5],
			}

			invoices = append(invoices, *invoice)

		}

	}

	repo := invoiceRepository.NewInvoiceRepository(util.GetDB())

	if err := repo.CreateBatch(context.Background(), invoices); err != nil {
		fmt.Printf("Error creating invoices: %v\n", err)
		panic(err)
	}

	fmt.Printf("Successfully imported %d invoices\n", len(invoices))
}
