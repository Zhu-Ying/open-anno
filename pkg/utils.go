package pkg

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
)

func ToJSON[T any](data T) string {
	var buffer bytes.Buffer
	jsonEncoder := json.NewEncoder(&buffer)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(data)
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(buffer.String())
}

func FromJSON[T any](line string) T {
	var data T
	err := json.Unmarshal([]byte(line), &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func Min[T int | float64](a T, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T int | float64](a T, b T) T {
	if a < b {
		return b
	}
	return a
}

func Abs[T int | float64](a T) T {
	if a < 0 {
		return -a
	}
	return a
}
