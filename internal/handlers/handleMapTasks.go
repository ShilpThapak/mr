package handlers

import (
	"os"
	"github.com/ShilpThapak/mr/internal/models"
	"github.com/ShilpThapak/mr/internal/storage"
	"github.com/ShilpThapak/mr/internal/utils"
	"strconv"
	"hash/fnv"
)

const IntermediateFilesBasePath = "intermediate/"

func ihash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func WriteIntermediateFile(lines []models.KeyValue, ReducerCnt int, TaskId int) {
	kvMap := make(map[string][]models.KeyValue)
	for _, i := range lines {
		key := i.Key
		id := int(ihash(key)) % ReducerCnt
		idStr := strconv.Itoa(id)
		taskIdStr := strconv.Itoa(TaskId)
		filepath := IntermediateFilesBasePath + "mr-" + taskIdStr + "-" + idStr + ".txt"
		kvMap[filepath] = append(kvMap[filepath], i)
	}
	storage.BulkWriteToFile(kvMap)
}

func HandleMapTasks(task models.Task, mapf func(string, string) []models.KeyValue) {
	filename := task.FileName

	// Read the file
	content, err := os.ReadFile(filename)
	utils.Check(err)

	mrOutput := mapf(filename, string(content))
	WriteIntermediateFile(mrOutput, task.ReducerCnt, task.Id)
}