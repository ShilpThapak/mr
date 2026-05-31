package models

import (
	"sync"
)

type Cordinator struct {
	Phase string
	Mu sync.Mutex
	Tasks []Task
	NReduce int
}