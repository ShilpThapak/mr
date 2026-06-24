package handlers

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"fmt"
	"github.com/ShilpThapak/mr/internal/models"
	"github.com/ShilpThapak/mr/internal/utils"
)

func HandleReduceTasks(task models.Task, reducef func(string, []string) string) {
	infilename := fmt.Sprintf("intermediate/mr-*-%d.txt", task.Id)
	filenames, err := filepath.Glob(infilename)
	utils.Check(err)

	fileKVMap := make(map[string][]string)

	for _, filename := range filenames {
		file, err := os.Open(filename)
		utils.Check(err)
		defer file.Close()
		
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			var kv models.KeyValue
			err3 := json.Unmarshal([]byte(line), &kv)
    	utils.Check(err3)

			fileKVMap[kv.Key] = append(fileKVMap[kv.Key], fmt.Sprint(kv.Value))
		}
		utils.Check(scanner.Err())
		
	}

	outFilename := fmt.Sprintf("outputs/mr-out-%d.txt", task.Id)
	outFile, err := os.Create(outFilename)
	utils.Check(err)
	defer outFile.Close()

	for key, value := range fileKVMap {
		reduceResult := reducef(key, value)
		_, err := fmt.Fprintf(outFile, "%s %s\n", key, reduceResult)
		utils.Check(err)
	}
}