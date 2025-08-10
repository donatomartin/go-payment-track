package main

import (
	"app/internal/invoice/repository"
	"app/internal/platform/util"
	"context"
	"fmt"
)

func main() {

	db := util.GetDB()

	repo := repository.NewInvoiceRepository(db)

	invoices, err := repo.GetDelayedInvoices(context.Background(), 0, 10)
	if err != nil {
		panic(err)
	}

	for _, inv := range invoices {
		fmt.Println("Delayed Invoice ID:", inv.ID, "Customer:", inv.CustomerName, "Due Date:", inv.DueDate.String(), "Amount Due:", inv.AmountDue, "Total Paid:", inv.TotalPaid)
	}
	fmt.Println("Total Delayed Invoices:", len(invoices))

}
