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
	"github.com/schollz/progressbar/v3"
)

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
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		fmt.Println("Error :", err)
		return
	}

	var rows [][]string
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error :", err)
			return
		}
		rows = append(rows, row)
	}

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
		return
	}
	defer output.Close()

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
}
