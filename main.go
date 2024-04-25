package main

import (
	"FinalProject/Kelompok10/controller"
	"FinalProject/Kelompok10/model"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MasterDimmy/go-cls"
	"github.com/asaskevich/govalidator"
	"github.com/schollz/progressbar/v3"
)

<<<<<<< HEAD
const banyakData = 5 // banyak data untuk diproses secara pararel

func main() {
	cls.CLS()

	inputFlag := flag.String("input", "", "set input file")
	outputFlag := flag.String("output", "", "set output file (optional)")
	flag.Parse()

	
	var inputPath string
	var rows [][]string
	var outputAllFile string
    if *outputFlag != "" {
        outputAllFile = *outputFlag
    }

	if *inputFlag == "" {
		fmt.Println("==============================================")
		fmt.Print("Masukkan Path File CSV: ") 
		fmt.Scanln(&inputPath)
		fmt.Print("==============================================\n")
	} else {
		inputPath = *inputFlag
	}

	filename := filepath.Base(inputPath)
	extension := filepath.Ext(inputPath)

	if outputAllFile == "" {
        outputAllFile = strings.TrimSuffix(filename, filepath.Ext(filename)) // Menghapus ekstensi dari nama file
    }
	outputAllFile = getOutputFileName(outputFlag, outputAllFile)
	outputFile := filepath.Join(outputAllFile)

	if _, err := os.Stat(outputFile); err == nil {
        fmt.Printf("File %s sudah ada. Apakah Anda ingin mengkonversinya lagi? (y/n): ", outputFile)
        var convertLagi string
        fmt.Scanln(&convertLagi)
        if convertLagi != "y" && convertLagi != "Y" {
            fmt.Println("Konversi dibatalkan.")
            return
        }
    }

	if extension != ".csv" {
		fmt.Printf("Input path file: %s is not a valid CSV file\n", inputPath)
		return
	}

	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Println("Ups, terjadi sebuah error :", err)
=======
func main() {
	cls.CLS()
	inputUsr := flag.String("input" ,"", "set input file")
	outputUsr := flag.String("output", "", "set output file (optional)")
	flag.Parse()
	
	var inputFile string

	if *inputUsr == "" {
		fmt.Print("Masukkan Path File CSV: ")
		fmt.Scanln(&inputFile)
	} else {
		inputFile = *inputUsr
	}

	filename := filepath.Base(inputFile)
	extension := filepath.Ext(inputFile)

	if extension != ".csv" {
		fmt.Printf("Input file: %s is not a valid CSV file\n", inputFile)
		return
	}

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error :", err)
>>>>>>> 979380e1fc67e3937eea7c402c56c88683a145fe
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
<<<<<<< HEAD
		fmt.Println("Ups, terjadi sebuah error :", err)
		return
	}

=======
		fmt.Println("Error :", err)
		return
	}

	var rows [][]string
>>>>>>> 979380e1fc67e3937eea7c402c56c88683a145fe
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
<<<<<<< HEAD
			fmt.Println("Ups, terjadi sebuah error :", err)
=======
			fmt.Println("Error :", err)
>>>>>>> 979380e1fc67e3937eea7c402c56c88683a145fe
			return
		}
		rows = append(rows, row)
	}

<<<<<<< HEAD
	selesai := make(chan bool)
	defer close(selesai)

	var isError bool
	go func() {
		isError = validateAndConvert(headers, rows, outputFlag, filename)
		selesai <- true
	}()

	<-selesai
	if err := os.MkdirAll("output_data_validasi", 0755); err != nil {
        fmt.Println("Ups, terjadi sebuah error :", err)
        return
    }

	if isError {
		fmt.Println("Validasi gagal. Proses konversi dibatalkan.")
=======
	if err := controller.ValidateData(headers, rows); err != nil {
		fmt.Println("Error :", err)
		return
	} 

	jsonData := convertToJSON(headers, rows)

	var outputFile string
	if *outputUsr != "" {
		outputFile = *outputUsr
	} else {
		outputFile = strings.TrimSuffix(filename, extension) + ".json"
	}

	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error :", err)
>>>>>>> 979380e1fc67e3937eea7c402c56c88683a145fe
		return
	}
	defer output.Close()

<<<<<<< HEAD
	fmt.Printf("Konversi dan Validasi File Berhasil, Data Tertulis ke file %s", outputFile)
=======
	jsonEncoder := json.NewEncoder(output)
	jsonEncoder.SetIndent("", "  ")
	if err := jsonEncoder.Encode(jsonData); err != nil {
		fmt.Println("Error :", err)
		return
	}
	
	bar := progressbar.Default(100)

	for i := 0; i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
		bar.Add(1)
	}
	bar.Clear()

	fmt.Printf("Konversi dan Validasi File Berhasil, Data Tertulis ke file %s di Folder ", outputFile)
}

func convertToJSON(headers []string, rows [][]string) model.JSONData {
	var jsonData model.JSONData
	jsonData.Data = make([]map[string]string, len(rows))

	for i, row := range rows {
		jsonData.Data[i] = make(map[string]string)
		for j, value := range row {
			jsonData.Data[i][headers[j]] = value
		}
	}

	return jsonData
>>>>>>> 979380e1fc67e3937eea7c402c56c88683a145fe
}

func validateAndConvert(headers []string, rows [][]string, outputFlag *string, filename string) bool {
	if err := ValidateData(headers, rows); err != nil {
		fmt.Println("Ups, terjadi sebuah error :", err)
		return true
	}

	jsonData := convertToJSON(headers, rows)
	outputFile := getOutputFileName(outputFlag, filename)
	if err := writeJSONToFile(jsonData, outputFile); err != nil {
		fmt.Println("Ups, terjadi sebuah error :", err)
		return true
	}

	// Membuat progress bar
	bar := progressbar.Default(int64(len(rows)), "Memproses data csv")

	// Channel untuk mengirim sinyal setiap kali sudah memproses sejumlah banyakData data
	progressSignal := make(chan bool, banyakData)

	// luping data csv
	for i, row := range rows {
		_ = row
		if (i+1)%banyakData == 0 || i == len(rows)-1 {
			progressSignal <- true
		}

		// update bar
		bar.Add(banyakData)
		time.Sleep(10 * time.Millisecond) 

		if (i+1)%banyakData == 0 {
			<-progressSignal
		}
	}

	close(progressSignal)

	return false
}

func getOutputFileName(outputFlag *string, filename string) string {
	if *outputFlag != "" {
		return *outputFlag
	}
	outputFolder := "output_data_validasi"
	outputFile := "data.json"
	return filepath.Join(outputFolder, outputFile)
}

func writeJSONToFile(jsonData model.JSONData, outputFile string) error {
	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	jsonEncoder := json.NewEncoder(output)
	jsonEncoder.SetIndent("", "  ")
	if err := jsonEncoder.Encode(jsonData); err != nil {
		return err
	}
	return nil
}

func convertToJSON(headers []string, rows [][]string) model.JSONData {
	var jsonData model.JSONData
	jsonData.Data = make([]map[string]string, len(rows))

	for i, row := range rows {
		jsonData.Data[i] = make(map[string]string)
		for j, value := range row {
			jsonData.Data[i][headers[j]] = value
		}
	}

	return jsonData
}

func ValidateData(headers []string, rows [][]string) error {
	var pilihanUser string
	isError := false

	if len(headers) == 0 {
		return fmt.Errorf("CSV file harus memiliki sebuah header")
	}

	for i, row := range rows {
		for j, value := range row {
			header := strings.ToLower(headers[j])
			switch header {
			case "email":
				if !govalidator.IsEmail(value) {
					emailErr := fmt.Sprintf("Email di baris %d tidak valid (%s)", i+2, value)
					fmt.Println(emailErr)
					isError = true
				}
			case "phone", "no", "telp", "hp":
				if !isValidPhoneNumber(value) {
					phoneErr := fmt.Sprintf("No hp di baris %d tidak valid (%s)", i+2, value)
					fmt.Println(phoneErr)
					isError = true
				}
			}
		}
	}

	if isError {
		for {
			fmt.Printf("Ada data yang tidak benar dari data csv tersebut, apakah anda yakin untuk melanjutkan ke tahap konversi ? (Y/N) : ")
			_, err := fmt.Scanln(&pilihanUser)
			if err != nil {
				fmt.Println("Ups, terjadi sebuah error :", err)
			}
			pilihanUser = strings.TrimSpace(pilihanUser)
			pilihanUser = strings.ToUpper(pilihanUser)
			if pilihanUser == "Y" || pilihanUser == "y" {
				fmt.Println("Melanjutkan konversi...")
				return nil
			} else if pilihanUser == "N" || pilihanUser == "n" {
				return fmt.Errorf("konversi dibatalkan")
			} else {
				fmt.Println("Pilihan tidak ditemukan, mohon masukkan jawaban Y/N")
			}
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