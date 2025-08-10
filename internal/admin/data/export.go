package data

import (
	"context"
	"fmt"

	"github.com/xuri/excelize/v2"

	invoiceRepo "app/internal/invoice/repository"
	paymentRepo "app/internal/payment/repository"
)

var months = []string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}

func ExportData(ctx context.Context, pRepo paymentRepo.PaymentRepository, iRepo invoiceRepo.InvoiceRepository, year int) (*excelize.File, error) {
	f := excelize.NewFile()
	f.DeleteSheet("Sheet1")

	// Payments sheet
	payments, err := pRepo.GetAll(ctx, "id", "asc", 0, 100000)
	if err != nil {
		return nil, err
	}

	f.NewSheet("PAGOS")
	f.SetSheetRow("PAGOS", "A1", &[]interface{}{})
	f.SetSheetRow("PAGOS", "A2", &[]interface{}{"Num Pago", "Num Factura", "Fecha", "Cantidad"})
	rowIdx := 3
	for _, p := range payments {
		if p.Date.Year() != year {
			continue
		}
		row := []interface{}{p.ID, p.InvoiceID, p.Date.Format("01-02-06"), fmt.Sprintf("%.2f €", p.Amount)}
		axis, _ := excelize.JoinCellName("A", rowIdx)
		f.SetSheetRow("PAGOS", axis, &row)
		rowIdx++
	}

	// Invoices sheets by month
	invoices, err := iRepo.GetAll(ctx, "invoice_date", "asc", 0, 100000)
	if err != nil {
		return nil, err
	}

	for _, m := range months {
		f.NewSheet(m)
		f.SetSheetRow(m, "A1", &[]interface{}{})
		f.SetSheetRow(m, "A2", &[]interface{}{"", "Cliente", "ID", "Fecha", "Importe", "Forma de pago", "Vencimiento"})
	}

	for _, inv := range invoices {
		if inv.InvoiceDate.Year() != year {
			continue
		}
		sheet := months[int(inv.InvoiceDate.Month())-1]
		row := []interface{}{"", inv.CustomerName, inv.ID, inv.InvoiceDate.Format("01-02-06"), fmt.Sprintf("%.2f €", inv.AmountDue), inv.PaymentMean, inv.DueDate.Format("01-02-06")}
		rows, err := f.GetRows(sheet)
		if err != nil {
			return nil, err
		}
		axis, _ := excelize.JoinCellName("A", len(rows)+1)
		f.SetSheetRow(sheet, axis, &row)
	}

	if idx, err := f.GetSheetIndex("PAGOS"); err == nil {
		f.SetActiveSheet(idx)
	}
	return f, nil
}
