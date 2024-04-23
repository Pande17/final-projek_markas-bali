package controller

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MasterDimmy/go-cls"
)

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
