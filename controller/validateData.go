package controller

import (
	"fmt"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
)

func ValidateData(headers []string, rows [][]string) error {
	var pilihanUser string
	isErrorFound := false // Menandakan apakah ada kesalahan yang ditemukan

	if len(headers) == 0 {
		return fmt.Errorf("CSV file harus memiliki header")
	}

	for i, row := range rows {
		for j, value := range row {
			header := strings.ToLower(headers[j])
			switch header {
			case "email":
				if !govalidator.IsEmail(value) {
					emailErr := fmt.Sprintf("Email di baris %d tidak valid (%s)", i+2, value)
					fmt.Println(emailErr)
					isErrorFound = true
				}
			case "phone", "no", "telp", "hp":
				if !isValidPhoneNumber(value) {
					phoneErr := fmt.Sprintf("No hp di baris %d tidak valid (%s)", i+2, value)
					fmt.Println(phoneErr)
					isErrorFound = true
				}
			}
		}
	}

	// Jika ada kesalahan, tampilkan pesan error dan lakukan validasi "Y/N"
	if isErrorFound {
		fmt.Printf("Ada data yang tidak benar dari data csv tersebut, apakah kamu yakin untuk menkoversi nya ? (Y/N) : ")
		_, err := fmt.Scanln(&pilihanUser)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
		}
		if pilihanUser == "Y" || pilihanUser == "y" {
			fmt.Println("Melanjutkan konversi...")
			return nil
		} else if pilihanUser == "N" || pilihanUser == "n" {
			os.Exit(0)
		} else {
			fmt.Println("Pilihan tidak ditemukan, mohon masukan jawaban Y/N")
		}
	}

	return nil
}

func isValidPhoneNumber(phoneNumber string) bool {
	if len(phoneNumber) > 0 && phoneNumber[0] == '+' {
		return govalidator.IsNumeric(phoneNumber[1:])
	}
	return govalidator.IsNumeric(phoneNumber)
}