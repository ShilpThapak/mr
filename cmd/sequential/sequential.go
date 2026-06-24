package main

import (
	"fmt"
	"os"

	"github.com/ShilpThapak/mr/internal/handlers"
	"github.com/ShilpThapak/mr/internal/models"
	"github.com/ShilpThapak/mr/internal/utils"
)

const ReducerCnt int = 10

func main(){
	var filenames []string

	if (len(os.Args) < 3) {
		fmt.Println("[SEQUENTIAL] No filenames passed. \nUsage: go run sample.go xxxx.so filename.txt ...")
		return
	}

	mapf, reducef := utils.LoadPlugin(os.Args[1])

	filenames = os.Args[2:]
	fmt.Println("[SEQUENTIAL] Recieved the files to process: ", filenames)

	for idx, filename := range filenames {
		var mapTask models.Task
		mapTask.Id = idx
		mapTask.FileName = filename
		mapTask.ReducerCnt = ReducerCnt
		mapTask.Type = "MapTask"
		mapTask.Status = "Pending"

		handlers.HandleMapTasks(mapTask, mapf)
	}

	for i := range ReducerCnt {
		var reduceTask models.Task
		reduceTask.Id = i
		reduceTask.ReducerCnt = ReducerCnt
		reduceTask.Status = "Pending"
		reduceTask.Type = "ReduceTask"

		handlers.HandleReduceTasks(reduceTask, reducef)
	}

}
