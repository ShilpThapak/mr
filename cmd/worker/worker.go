package main

import (
	"fmt"
	"net/rpc"
	"time"
	"os"

	"github.com/ShilpThapak/mr/internal/handlers"
	"github.com/ShilpThapak/mr/internal/models"
	"github.com/ShilpThapak/mr/internal/mrRpc"
	"github.com/ShilpThapak/mr/internal/utils"
)

func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := mrRpc.CoordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		fmt.Println("[WORKER] dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err != nil {
			fmt.Println(err)
			return false
	}

	return true
}

func requestTasksRPC() (task models.TaskRequestReply, err error) {
	var args models.TaskRequestArgs
	var reply models.TaskRequestReply
	ok := call("Cordinator.HandleTaskRequests", &args, &reply)
	if !ok {
		// panic(fmt.Errorf("[WORKER] RPC call to Coordinator.HandleTaskRequests failed"))
		return reply, fmt.Errorf("[WORKER] RPC call to Coordinator.HandleTaskRequests failed")
	}
	return reply, nil
}

func markTaskCompletionRPC(task models.Task) {
	var args models.TaskCompletionArgs
	args.Task = task
	var reply models.TaskCompletionReply
	ok := call("Cordinator.HandleTaskCompletion", &args, &reply)
	if !ok{
		panic(fmt.Errorf("[WORKER] Task updation failed."))
	}
}

func main(){
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "[WORKER] Usage: go run cmd/worker.go xxx.so\n")
		os.Exit(1)
	}

	mapf, reducef := utils.LoadPlugin(os.Args[1])

	for {
		reply, err := requestTasksRPC()
		if err != nil {
			fmt.Println(err)
			fmt.Println("[WORKER] Can't get any tasks. Retrying...")
			time.Sleep(1 * time.Second)
			continue
		}

		wait := reply.Wait
		if wait {
			fmt.Println("[WORKER] Waiting for task..", reply.Task)
			time.Sleep(5 * time.Second)
			continue
		}

		allDone := reply.AllDone
		if allDone {
			os.Exit(0)
			return
		}

		task := reply.Task
		fmt.Println("[WORKER] Recieved task with id: ", task.Id, task.FileName)

		switch task.Type {
			case models.MapTask:
				handlers.HandleMapTasks(task, mapf)

			case models.ReduceTask:
				handlers.HandleReduceTasks(task, reducef)
		}
		fmt.Println("[WORKER] Updating task completion. ", reply.Task.Id)
		markTaskCompletionRPC(reply.Task)
	}
}
