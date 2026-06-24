package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ShilpThapak/mr/internal/models"
	"github.com/ShilpThapak/mr/internal/utils"
)

func BulkWriteToFile(kvMap map[string][]models.KeyValue) {
	for fileAddress, kva := range kvMap {
		dir := filepath.Dir(fileAddress)
		base := filepath.Base(fileAddress)
		ext := filepath.Ext(base)

		nameWithoutExt := base[:len(base)-len(ext)]
		pattern := nameWithoutExt + "-*" + ext

		tempFile, err := os.CreateTemp(dir, pattern)
		utils.Check(err)
		defer tempFile.Close()

		for _, kv := range kva {
			jsonData, e := json.Marshal(kv)
			utils.Check(e)
			_, e = fmt.Fprintln(tempFile, string(jsonData))
			utils.Check(e)
		}

		err = os.Rename(tempFile.Name(), fileAddress)
		utils.Check(err)
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
}


