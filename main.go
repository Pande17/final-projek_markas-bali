package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MasterDimmy/go-cls"
)

// public variable
var FilePath string

// private function
func importFileCsv() {
	// ini case CLI saat dijalankan tanpa package flag
	cls.CLS()
	var err error

	fmt.Print("Masukkan path lokasi file CSV: ")
	fmt.Scanln(&FilePath)

	FilePath, err = filepath.Abs(FilePath)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n======================================")
	fmt.Println("   Tekan 'Enter' untuk melanjutkan...  ")
	fmt.Println("======================================")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func main() {
	importFileCsv()
	// ngebukak file csv nya
	file, err := os.Open(FilePath)
	if err != nil {
		// case error nya
		fmt.Println("Error:", err)
		return
	}
	// close file csv
	defer file.Close()

	// variabel yg isinya baru ngebaca file csv
	reader := csv.NewReader(file)

	// ni ngebaca semuanya
	records, err := reader.ReadAll()
	if err != nil {
		// case error
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(records)

	fmt.Println("\n======================================")
	fmt.Println("Tekan 'Enter' untuk melanjutkan...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
