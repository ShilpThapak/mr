package models

type TaskStatus string
type TaskType string
type Phase string

const (
	Pending TaskStatus = "Pending"
	Inprogress TaskStatus = "Inprogress"
	Done TaskStatus = "Done"
)

const (
	MapPhase Phase = "MapPhase"
	ReducePhase Phase = "ReducePhase"
)

const (
  MapTask TaskType = "MapTask"
  ReduceTask TaskType = "ReduceTask"
)

type KeyValue struct {
	Key string
	Value int
}

type Task struct {
  Id int
	Type TaskType
	FileName string
	ReducerCnt int
	Status TaskStatus
}
