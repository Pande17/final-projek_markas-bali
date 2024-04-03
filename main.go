package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/MasterDimmy/go-cls"
)


func main(){
	cls.CLS()
	fmt.Println("testing project")


	fmt.Println("\n======================================")
	fmt.Println("Tekan 'Enter' untuk melanjutkan...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}