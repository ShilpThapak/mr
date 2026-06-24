package main

//
// same as crash.go but doesn't actually crash.
//
// go build -buildmode=plugin nocrash.go
//

import "github.com/ShilpThapak/mr/internal/models"
import crand "crypto/rand"
import "math/big"
import "strings"
import "os"
import "sort"
import "strconv"

func maybeCrash() {
	max := big.NewInt(1000)
	rr, _ := crand.Int(crand.Reader, max)
	if false && rr.Int64() < 500 {
		// crash!
		os.Exit(1)
	}
}

func Map(filename string, contents string) []models.KeyValue {
	maybeCrash()

	kva := []models.KeyValue{}
	kva = append(kva, models.KeyValue{"a", filename})
	kva = append(kva, models.KeyValue{"b", strconv.Itoa(len(filename))})
	kva = append(kva, models.KeyValue{"c", strconv.Itoa(len(contents))})
	kva = append(kva, models.KeyValue{"d", "xyzzy"})
	return kva
}

func Reduce(key string, values []string) string {
	maybeCrash()

	// sort values to ensure deterministic output.
	vv := make([]string, len(values))
	copy(vv, values)
	sort.Strings(vv)

	val := strings.Join(vv, " ")
	return val
}
