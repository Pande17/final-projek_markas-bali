package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"FinalProject/Kelompok10/mockstruct" // Sesuaikan dengan struktur direktori Anda
	"FinalProject/Kelompok10/utils"

	"github.com/MasterDimmy/go-cls"
	"github.com/schollz/progressbar/v3"
)

func main() {
	cls.CLS()

	// Parsing flags
	inputFlag := flag.String("input", "", "set input file")
	outputFlag := flag.String("output", "", "set output file (optional)")
	flag.Parse()

	var inputPath string
	var outputAllFile string

	if *outputFlag != "" {
		outputAllFile = *outputFlag
	}

	// Meminta input path file jika tidak disediakan melalui flag
	if *inputFlag == "" {
		fmt.Println("==============================================")
		fmt.Print("Masukkan Path File CSV: ")
		fmt.Scanln(&inputPath)
		fmt.Print("==============================================\n")
	} else {
		inputPath = *inputFlag
	}

	// Mendapatkan nama file dan ekstensinya
	filename := filepath.Base(inputPath)
	extension := filepath.Ext(inputPath)

	// Menyusun nama file output
	if outputAllFile == "" {
		outputAllFile = strings.TrimSuffix(filename, filepath.Ext(filename))
	}
	outputAllFile = getOutputFileName(outputFlag, outputAllFile)
	outputFile := filepath.Join(outputAllFile)

	// Memeriksa apakah file output sudah ada
	if _, err := os.Stat(outputFile); err == nil {
		fmt.Printf("File %s sudah ada. Apakah Anda ingin mengkonversinya lagi? (y/n): ", outputFile)
		var convertLagi string
		fmt.Scanln(&convertLagi)
		if convertLagi != "y" && convertLagi != "Y" {
			fmt.Println("Konversi dibatalkan.")
			return
		}
	}

	// Memeriksa apakah file input adalah file CSV
	if extension != ".csv" {
		fmt.Printf("Input path file: %s is not a valid CSV file\n", inputPath)
		return
	}

	// Membuka file CSV
	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Println("Ups, terjadi sebuah error:", err)
		return
	}
	defer file.Close()

	// Membuat pembaca CSV
	reader := csv.NewReader(bufio.NewReader(file))

	// Membaca semua baris CSV
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading CSV:", err)
	}

	// cek panjang
	fmt.Println("Membaca Header...")
	if len(records) < 1 {
		log.Fatal("File CSV terlalu singkat untuk diproses")
	}

	// array slice index 0 dari records
	headers := records[0]

	// utk nampung data yg di convert jadi json
	dataConvert := []map[string]any{}

	// nampung error validasi
	var validationErrors []string

	wg := sync.WaitGroup{}

	// Membuat progress bar
	bar := progressbar.Default(int64(len(records)), "Memproses data csv")

	// Membuat channel untuk komunikasi antar goroutine
	chanProgress := make(chan int)                               // channel untuk persentase progress
	chanRecords := make(chan mockstruct.CsvRecord, len(records)) // channel untuk nampung data dari records
	chanConvertedData := make(chan map[string]any)               // channel untuk menampung data yang sudah diconvert
	chanErrors := make(chan error)                               // channel untuk menampung error

	// Menentukan jumlah goroutine yang akan dijalankan
	numRoutines := 5

	// go routine 1 
	go func() {
		for progress := range chanProgress {
			bar.Add(progress)
		}
	}()
	
	// go routine 2
	go func() {
		for errors := range chanErrors {
			validationErrors = append(validationErrors, errors.Error())
		}
	}()

	// looping part 3
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go utils.ValidateRecords(chanRecords, chanConvertedData, chanErrors, chanProgress, &wg, headers)
	}

	for index, record := range records {
		if index == 0 {
			chanProgress <- 1
			continue
		}

		// Kirim data ke channel outputData
		chanRecords <- mockstruct.CsvRecord{Index: index, Data: record}
	}

	close(chanRecords)

	// go routine 3
	go func() {
		// Menunggu selesai dari WaitGroup
		wg.Wait()

		// Menyelesaikan progress bar
		bar.Finish()

		// Menutup ketiga channel
		close(chanProgress)
		close(chanConvertedData)
		close(chanErrors)

	}()

	for loopingChanelDua := range chanConvertedData {
		dataConvert = append(dataConvert, loopingChanelDua)
	}

	// Cek apakah ada error validasi
	if len(validationErrors) > 0 {
		for _, err := range validationErrors {
			log.Println(err)
		}
		fmt.Print("Terjadi Error. Apakah Anda ingin tetap mengkonversinya atau tidak? (y/n): ")
		var mulaiUlang string 
		fmt.Scanln(&mulaiUlang)
		if mulaiUlang != "y" && mulaiUlang != "Y" {
			fmt.Println("Konversi dibatalkan.")
			return
		}
	}

	// Tulis data JSON ke file
	err = writeJSONToFile(dataConvert, outputFile)
	if err != nil {
		log.Fatal("Error writing JSON to file:", err)
	}

	fmt.Printf("Konversi dan Validasi File Berhasil, Data Tertulis ke file %s", outputFile)
}

func getOutputFileName(outputFlag *string, filename string) string {
	if *outputFlag != "" {
		return *outputFlag
	}
	outputFolder := "output_data_validasi"
	outputFile := "data.json"
	return filepath.Join(outputFolder, outputFile)
}

func writeJSONToFile(data []map[string]interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(data); err != nil {
		return err
	}
	return nil
}
