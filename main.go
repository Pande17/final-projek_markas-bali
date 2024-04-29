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
	"testing"

	"FinalProject/Kelompok10/mockstruct" // Sesuaikan dengan struktur direktori Anda

	"github.com/MasterDimmy/go-cls"
	"github.com/asaskevich/govalidator"
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
	if len(records) < 1 {
		log.Fatal("File CSV terlalu sedikit untuk diproses")
	}
	
	// array slice index 0 dari records
	headers := records[0]

	// utk nampung data yg di convert jadi json
	var dataConvert map[string]any

	// nampung error valdasi
	var validationErrors []string

	// 
	wg := sync.WaitGroup{}

	// Membuat progress bar
	bar := progressbar.Default(int64(len(records)-1), "Memproses data csv")

	// Membuat channel untuk komunikasi antar goroutine
	chanProgress := make(chan int)
	chanRecords := make(chan mockstruct.CsvRecord, len(records)-1)
	chanConvertedData := make(chan map[string]any)
	chanErrors := make(chan error)

	// Menentukan jumlah goroutine yang akan dijalankan
	numRoutines := 1

	// 2 go rutin 
	go func() {
        for progress := range chanProgress {
            bar.Add(progress)
        }
    }()

    go func() {
        for errors := range chanErrors {
            validationErrors = append(validationErrors, errors.Error())
        }
    }()

	// looping part 3
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go validateRecords(chanProgress, chanRecords, chanConvertedData, chanErrors, wg, headers)

		for recordLooping := range records {
			if recordLooping == 0 {
				chanProgress <- 1
				continue
			}
			
			// Kirim data ke channel outputData
			chanRecords <- recordLooping
		}
	}
	close(chanRecords)

	go func() {
		// Menunggu selesai dari WaitGroup
		wg.Wait()
		
		// Menyelesaikan progress bar
		bar.Finish()
	
		// Menutup ketiga channel
		close(chanProgress)
		close(chanRecords)
		close(chanErrors)
	}()

	for loopingChanelDua := range chanRecords {

	}
	



// Proses semua record, kecuali header
for i := 1; i < len(records); i++ {
    record := mockstruct.CsvRecord{
        Index: i + 1, // Index dimulai dari 2 karena header dianggap indeks 1
        Data:  records[i],
    }
    chanRecords <- record
}
close(chanRecords) // Menutup channel setelah semua record dikirim

// Menunggu semua goroutine selesai
wg.Wait()

// Setelah semua goroutine selesai, tutup channel progress dan errors
close(chanProgress)
close(chanErrors)

// Cek apakah ada error validasi
if len(validationErrors) > 0 {
    for _, err := range validationErrors {
        log.Println(err)
    }
    log.Fatal("Ada error validasi, konversi dibatalkan")
}

	// Membuat channel untuk memberi tahu goroutine bahwa semua record telah diproses
	done := make(chan struct{})

	// Menjalankan goroutine untuk memvalidasi dan mengkonversi data
	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			for record := range chanRecords {
				outputJson := make(map[string]interface{})
				for j, value := range record.Data {
					// Proses validasi
					if headers[j] == "email" {
						if !isValidEmail(value) {
							chanErrors <- fmt.Errorf("Email di baris %d tidak valid: %s", record.Index, value)
							continue
						}
					}
					if headers[j] == "phone" {
						if !isValidPhoneNumber(value) {
							chanErrors <- fmt.Errorf("Nomor telepon di baris %d tidak valid: %s", record.Index, value)
							continue
						}
					}
					// Konversi data ke JSON
					outputJson[headers[j]] = value
				}

				// Mengirim data hanya jika channel belum ditutup
				select {
				case chanConvertedData <- outputJson:
				case <-done:
					return
				}
				chanProgress <- 1
			}
		}()
	}

	// Proses semua record, kecuali header
	for i := 1; i < len(records); i++ {
		record := mockstruct.CsvRecord{
			Index: i + 1, // Index dimulai dari 2 karena header dianggap indeks 1
			Data:  records[i],
		}
		chanRecords <- record
	}
	close(chanRecords)

	// Menunggu semua goroutine selesai
	wg.Wait()

	// Setelah semua goroutine selesai, tutup channel progress dan errors
	close(chanProgress)
	close(chanErrors)

	// Memberi sinyal bahwa semua record telah diproses
	close(done)

	// Cek apakah ada error validasi
	if len(validationErrors) > 0 {
		for _, err := range validationErrors {
			log.Println(err)
		}
		log.Fatal("Ada error validasi, konversi dibatalkan")
	}

	// Menggabungkan hasil konversi
	var convertedData []map[string]interface{}
	for i := 0; i < len(records)-1; i++ {
		convertedData = append(convertedData, <-chanConvertedData)
	}

	// Tulis data JSON ke file
	err = writeJSONToFile(convertedData, outputFile)
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

func isValidEmail(email string) bool {
	return govalidator.IsEmail(email)
}

func isValidPhoneNumber(phoneNumber string) bool {
	return govalidator.IsNumeric(phoneNumber)
}

func validateRecords(records <-chan mockstruct.CsvRecord, outputData chan map[string]any, errors chan error, progress chan int, wg *sync.WaitGroup, headers []string) {
    
    // Loop through records received from the channel
    for record := range records {
		data := map[string]any{}
        for index, value := range record.Data {
            // Proses validasi
            if headers[index] == "email" {
                if !isValidEmail(value) {
                    errors <- fmt.Errorf("Email di baris %d tidak valid: %s", record.Index, value)
                    continue
                }
            }
            if headers[index] == "phone" {
                if !isValidPhoneNumber(value) {
                    errors <- fmt.Errorf("Nomor telepon di baris %d tidak valid: %s", record.Index, value)
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

func Test_EmailValidation(t *testing.T) {
    t.Run("Success", func(t *testing.T) {

        mailString := "valid@email.com"
        header := "email"

        ch := make(chan mockstruct.CsvRecord)
        chOutput := make(chan map[string]any, 1)
        chErrors := make(chan error, 1)
        wg := new(sync.WaitGroup)

        mockData := mockstruct.CsvRecord{
            Index: 1,
            Data:  []string{mailString},
        }
        wg.Add(1)
        go utils.ValidateCsv(ch, chOutput, chErrors, make(chan int, 1), wg, []string{headers})

        ch <- mockData
        close(ch)

        wg.Wait()

        select {
        case err := <-chErrors:
            t.Errorf("Test failed, got error message :%s", err.Error())
        case <-chOutput:
            fmt.Println("Email Valid")
        }
    })
    t.Run("Failed", func(t *testing.T) {
        // Fill code
    })
	
}


