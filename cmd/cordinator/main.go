package main

import (
	"fmt"
	"net/rpc"
	"os"
	"sync"
	"time"
	"net"
	"net/http"
	"github.com/ShilpThapak/mr/internal/models"
	"github.com/ShilpThapak/mr/internal/mrRpc"
	"github.com/ShilpThapak/mr/internal/utils"

)

const ReducerCnt int = 10

type Cordinator struct {
	Phase models.Phase
	Mu sync.Mutex
	Tasks []models.Task
	NReduce int
}

func (c *Cordinator) HandleTaskRequests(args *models.TaskRequestArgs, reply *models.TaskRequestReply) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()


	fmt.Printf("[CORDINATOR][%s] Handling task request.. \n", c.Phase)

	allMapTasksDone := true
	for idx, task := range c.Tasks {
		switch task.Status {
			case models.Pending:
				reply.Task = task
				c.Tasks[idx].Status = models.Inprogress
				return nil
			case models.Inprogress:
				allMapTasksDone = false
		}
	}

	// if map tasks are in-progress, wait for them to complete
	if !allMapTasksDone {
		reply.Wait = true
		return nil
	}

	reduceTaskCnt := 0
	for _, task := range c.Tasks {
		if task.Type == models.ReduceTask {
			reduceTaskCnt += 1
		}
	}

	// if no pending / in-progress map tasks, switch the phase and start again
	if reduceTaskCnt == 0 {
		fmt.Printf("[CORDINATOR][%s] All Map Tasks Done. Switching to reduce phase \n", c.Phase)
		c.Phase = models.ReducePhase
		for i := range c.NReduce {
			var reduceTask models.Task
			reduceTask.Id = i
			reduceTask.ReducerCnt = ReducerCnt
			reduceTask.Status = models.Pending
			reduceTask.Type = models.ReduceTask
			c.Tasks = append(c.Tasks, reduceTask)
		}
	}
	
	allReduceTasksDone := true
	for idx, task := range c.Tasks {
		switch task.Status {
			case models.Pending:
				reply.Task = task
				c.Tasks[idx].Status = models.Inprogress
				return nil
			case models.Inprogress:
				allReduceTasksDone = false
		}
	}

	if allReduceTasksDone {
		reply.Wait = true
		return nil
	}

	return nil
}

func (c *Cordinator) HandleTaskCompletion(args *models.TaskCompletionArgs, reply *models.TaskCompletionReply) error {
	taskId := args.Task.Id
	for idx, task := range c.Tasks {
		if task.Id == taskId {
			c.Tasks[idx].Status = models.Done
		}
	}
	fmt.Println("Updated the task with completion", args.Task, c.Tasks)
	return nil
}

func (c *Cordinator) Done(allTasksDone bool) bool {
c.Mu.Lock()
defer c.Mu.Unlock()
	if c.Phase == models.MapPhase {
		return false
	}

	allDone := true
	for _, task := range c.Tasks {
		if task.Status != models.Done {
			allDone = false
		}
	}
	return allDone
}

func (c *Cordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()

	sockname := mrRpc.CoordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	utils.Check(e)
	go http.Serve(l, nil)

	fmt.Printf("[CORDINATOR][%s] Started RPC Server.. \n", c.Phase)
}

func makeCordinator(filenames []string) *Cordinator {
	var c Cordinator
	c.Phase = models.MapPhase
	c.NReduce = ReducerCnt

	for idx, filename := range filenames {
		var mapTask models.Task
		mapTask.FileName = filename
		mapTask.Id = idx
		mapTask.Type = models.MapTask
		mapTask.Status = models.Pending
		mapTask.ReducerCnt = c.NReduce
		c.Tasks = append(c.Tasks, mapTask)
	}
	fmt.Printf("[CORDINATOR][%s phase] Added %d tasks \n", c.Phase, len(c.Tasks))
	c.server()
	return &c
}

func main() {
	fmt.Println("[CORDINATOR] Starting up..")
	var filenames []string
	if (len(os.Args) > 1) {
		filenames = os.Args[1:]
		fmt.Println("[CORDINATOR] Recieved the files to process: ", filenames)
	} else {
		fmt.Println(("No filenames passed. \nUsage: go run cmd/cordinator.go filename.txt"))
		return
	}

	c := makeCordinator(filenames)
	allTasksDone := false
	for c.Done(allTasksDone) == false {
		time.Sleep(1 * time.Second)
	}
	fmt.Println("[CORDINATOR] All tasks done. Shutting down..")
}
