package utils

import (
	"FinalProject/Kelompok10/mockstruct"
	"fmt"
	"strings"
	"sync"

	"github.com/asaskevich/govalidator"
)

func isValidPhoneNumber(phoneNumber string) bool {
	if len(phoneNumber) > 0 && phoneNumber[0] == '+' {
		return govalidator.IsNumeric(phoneNumber[1:])
	}
	return govalidator.IsNumeric(phoneNumber)
}

func ValidateRecords(records <-chan mockstruct.CsvRecord, outputData chan map[string]any, errors chan error, progress chan int, wg *sync.WaitGroup, headers []string) {

	// Loop through records received from the channel
	for record := range records {
		data := map[string]any{}
		for index, value := range record.Data {
			// Proses validasi
			switch strings.ToLower(headers[index]) {
			case "email":
				if !isValidEmail(value) {
					errors <- fmt.Errorf("email di baris %d tidak valid: %s", record.Index, value)
					continue
				}
			case "phone", "hp", "no_telp":
				if !isValidPhoneNumber(value) {
					errors <- fmt.Errorf("nomor telepon di baris %d tidak valid: %s", record.Index, value)
					continue
				}
			}
			// Buat kunci untuk map menggunakan sprintf
			data[fmt.Sprintf("%v", headers[index])] = value

			outputData <- data
		}
		// Kirim sinyal progress ke channel progress
		progress <- 1
	}
	wg.Done()
}

func isValidEmail(email string) bool {
	return govalidator.IsEmail(email)
}
