package main

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/ShilpThapak/mr/internal/handlers"
	"github.com/ShilpThapak/mr/internal/models"
	"github.com/ShilpThapak/mr/internal/mrRpc"
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
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

func requestTasksRPC() (task models.TaskRequestReply) {
	var args models.TaskRequestArgs
	var reply models.TaskRequestReply
	ok := call("Cordinator.HandleTaskRequests", &args, &reply)
	if !ok {
		panic(fmt.Errorf("[WORKER] RPC call to Coordinator.HandleTaskRequests failed"))
	}
	return reply
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
	for {
		reply := requestTasksRPC()
		wait := reply.Wait
		if wait {
			fmt.Println("[WORKER] Waiting for task..")
			time.Sleep(5 * time.Second)
			continue
		}

		task := reply.Task
		fmt.Println("[WORKER] Recieved task with id: ", task.Id)
		switch task.Type {
			case models.MapTask:
				handlers.HandleMapTasks(task)

			case models.ReduceTask:
				handlers.HandleReduceTasks(task)
		}
		fmt.Println("[WORKER] Updating task completion. ", reply.Task.Id)
		markTaskCompletionRPC(reply.Task)
		time.Sleep(5 * time.Second)
	}
}
