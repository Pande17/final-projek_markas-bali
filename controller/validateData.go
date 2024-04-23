package controller

import (
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
)

func ValidateData(headers []string, rows [][]string) error {
	if len(headers) == 0 {
		return fmt.Errorf("CSV file must have headers")
	}

	for i, row := range rows {
		for j, value := range row {
			header := strings.ToLower(headers[j])
			switch header {
			case "email":
				if !govalidator.IsEmail(value) {
					return fmt.Errorf("Invalid email format at Row %d, column %d", i+2, j+1)
				}
			case "phone", "no", "telp", "hp":
				if !govalidator.IsNumeric(value) {
					return fmt.Errorf("Invalid phone format at Row %d, column %d", i+2, j+1)
				}
			}
		}
	}

	return nil
}
