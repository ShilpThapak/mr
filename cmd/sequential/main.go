package main

import (
	"fmt"
	"os"
	"github.com/ShilpThapak/mr/internal/models"
	"github.com/ShilpThapak/mr/internal/handlers"
)

const ReducerCnt int = 10

func main(){
	fmt.Println("Hello world")

	var filenames string
	// Get the filename
	if (len(os.Args) > 1) {
		filenames = os.Args[1]
		fmt.Println("Recieved the files to process: ", filenames)
	} else {
		// panic("No filename passes. Usage: go run sample.go filename.txt")
		fmt.Println(("No filenames passed. \nUsage: go run sample.go filename.txt"))
		return
	}

	for idx, filename := range os.Args[1:] {
		var mapTask models.Task
		mapTask.Id = idx
		mapTask.FileName = filename
		mapTask.ReducerCnt = ReducerCnt
		mapTask.Type = "MapTask"
		mapTask.Status = "Pending"

		handlers.HandleMapTasks(mapTask)
	}

	for i := range ReducerCnt {
		var reduceTask models.Task
		reduceTask.Id = i
		reduceTask.ReducerCnt = ReducerCnt
		reduceTask.Status = "Pending"
		reduceTask.Type = "ReduceTask"

		handlers.HandleReduceTasks(reduceTask)
	}

}
