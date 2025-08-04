package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"

	"app/cmd/scripts/commands"
)

func main() {
	f, err := excelize.OpenFile("data.xlsx")
	if err != nil {
		fmt.Println(err)
	}

	commands.ImportPagos(f)
	commands.ImportFacturas(f)

}
