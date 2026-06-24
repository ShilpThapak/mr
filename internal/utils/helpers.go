package utils

import (
	"fmt"
	"plugin"

	"github.com/ShilpThapak/mr/internal/models"
)

func Check(e error){
	if e != nil {
		panic(e)
	}
}

// load the application Map and Reduce functions
// from a plugin file, e.g. ../mrapps/wc.so
func LoadPlugin(filename string) (func(string, string) []models.KeyValue, func(string, []string) string) {
	p, err := plugin.Open(filename)
	if err != nil {
		fmt.Printf("cannot load plugin %v", filename)
		panic(err)
	}
	xmapf, err := p.Lookup("Map")
	if err != nil {
		fmt.Printf("cannot find Map in %v", filename)
	}
	mapf := xmapf.(func(string, string) []models.KeyValue)
	xreducef, err := p.Lookup("Reduce")
	if err != nil {
		fmt.Printf("cannot find Reduce in %v", filename)
	}
	reducef := xreducef.(func(string, []string) string)

	return mapf, reducef
}