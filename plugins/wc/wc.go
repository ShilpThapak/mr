package main

import (
	"github.com/ShilpThapak/mr/internal/models"
	"unicode"
	"strings"
	"strconv"
)

func Map(filename string, contents string) []models.KeyValue {
	ff := func(r rune) bool { return !unicode.IsLetter(r) }

	words := strings.FieldsFunc(contents, ff)

	kva := []models.KeyValue{}
	for _, w := range words {
		kv := models.KeyValue{
			Key: w,
			Value: "1",
		}
		kva = append(kva, kv)
	}
	return kva
}

func Reduce(key string, val []string) string {
  return strconv.Itoa(len(val))
}
