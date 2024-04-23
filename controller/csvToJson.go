package controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type BaseData struct {
	ID         string `json:"_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Postalcode string `json:"postal_code"`
}

func CsvToStruct(records [][]string) []BaseData {
	// Phone, err := strconv.Atoi()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	baseData := []BaseData{}
	for i, line := range records {
		if i > 0 {
			var rec BaseData
			for j, field := range line {
				if j == 0 {
					rec.ID = field
				} else if j == 1 {
					rec.Name = field
				} else if j == 2 {
					rec.Email = field
				} else if j == 3 {
					rec.Phone = field
				} else if j == 4 {
					rec.Address = field
				} else if j == 5 {
					rec.City = field
				} else if j == 6 {
					rec.Postalcode = field
				}
			}
			baseData = append(baseData, rec)
		}
	}
	return baseData
}

func convertToJson(data BaseData) []byte {
	encoded, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println(err)
	}

	return encoded
}

func saveJsonToFile(encoded []byte, name string) {
	file, err := os.Create(fmt.Sprintf("json/%s.json", name))
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err = file.Write(encoded)
	if err != nil {
		fmt.Println(err)
	}
}

func process(channel chan BaseData, wg *sync.WaitGroup) {
	for Data := range channel {
		encoded := convertToJson(Data)
		saveJsonToFile(encoded, Data.ID)
	}

	wg.Done()
}

func CsvConvert() {
	err := os.RemoveAll("json")
	if err != nil {
		fmt.Println(err)
	}

	os.Mkdir("json", 0777)

	fileCsv, err := os.Open("datas.csv")
	if err != nil {
		fmt.Println(err)
	}

	defer fileCsv.Close()

	reader := csv.NewReader(fileCsv)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	datas := CsvToStruct(records)

	startedAt := time.Now()

	wg := sync.WaitGroup{}

	var channel = make(chan BaseData)

	jml := 5
	fmt.Println("menjalankan", jml, "process goroutine")

	for i := 0; i < jml; i++ {
		wg.Add(1)
		go process(channel, &wg)
	}

	for _, data := range datas {
		channel <- data
	}

	close(channel)

	wg.Wait()

	fmt.Println("Success")
	fmt.Println(time.Since(startedAt))

}
