package util

import (
	"fmt"
	"strings"
)

func Float64ToEuros(amount float64) string {
	string := fmt.Sprintf("%.2f", amount)
	string = strings.ReplaceAll(string, ".", ",")

	// Put a dot every three digits from the right for the integer part
	for i := len(string) - 6; i > 0; i -= 3 {
		string = string[:i] + "." + string[i:]
	}

	string += " €"

	return string
}

