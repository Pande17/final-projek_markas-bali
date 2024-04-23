package main

import (
	"FinalProject/Kelompok10/model"
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/MasterDimmy/go-cls"
	"github.com/schollz/progressbar/v3"
)

// public variable
var FilePath string

// private function
func importFileCsv() {
	// ini case CLI saat dijalankan tanpa package flag
	cls.CLS()
	var err error

	fmt.Print("Masukkan input file : ")
	fmt.Scanln(&FilePath)

	FilePath, err = filepath.Abs(FilePath)
	if err != nil {
		panic(err)
	}

	fmt.Println("=========== PROSES COMPLETE ===========")
	fmt.Printf("File berhasil divalidasi dan konversi : %s", FilePath)
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// testing flag package ( masih belajar makek :v)
	inputFile := flag.String("input", "", "Set input file")
	fmt.Println(inputFile)
}

func convWithGoroutine(ch <-chan model.BaseData, wg *sync.WaitGroup, noUrut int) {
	for csvFile := range ch {
		dataJson, err := json.Marshal(csvFile)
		if err != nil {
			fmt.Println("Terjadi error:", err)
		}

		err = os.WriteFile(fmt.Sprintf("books/%s.json", csvFile.Name), dataJson, 0644)
		if err != nil {
			fmt.Println("Terjadi error:", err)
		}

		fmt.Printf("Antrian No %d Memproses Kode Buku : %s!\n", noUrut, csvFile.Name)
	}
	wg.Done()
}

func convertCsvToJson() {
	csvFile, err := os.Open(FilePath)
	if err != nil {
		fmt.Println("Terjadi error : ", err)
		return
	}
	defer csvFile.Close()

	// parse file csv
	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Terjadi Error : ", err)
		return
	}

	// persiapan untuk menyimpan file JSON
	var jsonData []map[string]string

	// mengonversi setiap baris CSV menjadi map
	for _, row := range records {
		record := make(map[string]string)
		for i, column := range row {
			record[fmt.Sprintf("colum%d")]
		}
	}
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

	// Convert CSV ke json
	var jsonData []map[string]string
	for _, row := range records {
		entry := make(map[string]string)
		for i, value := range row {
			entry[records[0][i]] = value
		}
		jsonData = append(jsonData, entry)
	}

	// Convert JSON to string
	// pakek marsal inden supaya file josn nya rapi kebawah.. ga nyambung terus kesamping
	jsonString, err := json.MarshalIndent(jsonData, "", " ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print JSON string
	fmt.Println(string(jsonString))

	// tambahan hari ini ( 23/04/2024)
	_ = os.Mkdir("csv_convert", 0777)
	ch := make(chan model.BaseData)
	wg := sync.WaitGroup{}
	jumlahProses := 5
	for i := 0; i < jumlahProses; i++ {
		wg.Add(1)
		go convWithGoroutine(ch, &wg, i)
	}
	for _, csv := range records {
		ch <- csv
	}

	close(ch)

	wg.Wait()

	// testing progress bar ( ini udah berhasil.. tinggal copas & benerin logicny sesuai dengan case yg dibutuhkan)
	csvData := records // contoh data ngambil semua isi csv

	// variabel  progress bar dari total valuye
	bar := progressbar.Default(int64(len(csvData)), "Memproses Data")

	// loop sesuai isi data
	for _, value := range csvData {
		// proses nampilin log
		fmt.Println("Processing value:", value)

		// itungan lambat bar
		time.Sleep(40 * time.Millisecond)
		// tambah cls untuk estetika :v
		cls.CLS()
		bar.Add(1)
	}
	// CsvConvert()
	// Hapus abis suud prosesny
	bar.Clear()

	fmt.Println("\n======================================")
	fmt.Println("Tekan 'Enter' untuk melanjutkan...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
