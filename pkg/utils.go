package pkg

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
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

func FindArr[T string | int | float64](arr []T, k T) int {
	for i, v := range arr {
		if k == v {
			return i
		}
	}
	return -1
}

func FormatChrom(chrom string) string {
	chrom = strings.Replace(chrom, "chr", "", 1)
	if chrom == "M" {
		return "MT"
	}
	return chrom
}

func CopyFile(src string, dst string) {
	reader, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	writer, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(writer, reader); err != nil {
		log.Fatal(err)
	}

	defer writer.Close()
}
