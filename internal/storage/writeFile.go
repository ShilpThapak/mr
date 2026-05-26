package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ShilpThapak/mr/internal/models"
	"github.com/ShilpThapak/mr/internal/utils"
)

func BulkWriteToFile(kvMap map[string][]models.KeyValue) {
	for filepath, kva := range kvMap {
		file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		utils.Check(err)
		defer file.Close()

		for _, kv := range kva {
			jsonData, e2 := json.Marshal(kv)
			utils.Check(e2)
			_, e := fmt.Fprintln(file, string(jsonData))
			utils.Check(e)
			fmt.Println("written ", filepath, kv)
		}
	}
}

func WriteToFile(filepath string, kv models.KeyValue) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	utils.Check(err)
	defer file.Close()

	jsonData, e2 := json.Marshal(kv)
	utils.Check(e2)

	_, e := fmt.Fprintln(file, string(jsonData))
	utils.Check(e)
	fmt.Println("written ", filepath, kv)
}


