package models

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

type TaskRequestArgs struct {

}

type TaskRequestReply struct {
	Task Task
	Wait bool
	AllDone bool
}

type TaskCompletionArgs struct {
	Task Task
}

type TaskCompletionReply struct {
	
}