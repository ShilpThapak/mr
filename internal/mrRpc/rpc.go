package mrRpc

import (
	"strconv"
	"os"
)

func CoordinatorSock() string {
	s := "/var/tmp/ShilpThapak-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}