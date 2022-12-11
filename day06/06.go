package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func isSopm(buf []byte) bool {
	l := len(buf)
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			if buf[i] == buf[j] {
				return false
			}
		}
	}
	return true
}

func findStartMarker(f *os.File, distinctCount int) (int, []byte, error) {
	buf := make([]byte, distinctCount)

	for i := 0; ; i++ {
		_, err := f.Seek(int64(i), 0)
		if err != nil {
			log.Fatal(err)
		}
		b, err := f.Read(buf)
		if b < distinctCount {
			break
		}
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		if isSopm(buf) {
			return i + distinctCount, buf, nil
		}
	}
	return 0, buf, errors.New("Could not find start-of-packet marker")
}

func main() {
	fmt.Println("Day 06")

	f, err := os.Open("06.input")
	if err != nil {
		log.Fatal()
	}
	defer f.Close()
	sopmIndex, sopm, err := findStartMarker(f, 4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sopmIndex, ":", string(sopm))
	somIndex, som, err := findStartMarker(f, 14)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(somIndex, ":", string(som))
}
