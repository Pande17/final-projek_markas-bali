package utils_test

import (
	"FinalProject/Kelompok10/mockstruct"
	"FinalProject/Kelompok10/utils"
	"fmt"
	"sync"
	"testing"
)

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
		go utils.ValidateRecords(ch, chOutput, chErrors, make(chan int, 1), wg, []string{header})

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
