package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func expandRange(compressedRange string) []int {
	var expandedRange = []int{}
	delimiters := strings.Split(compressedRange, "-")
	start, _ := strconv.Atoi(delimiters[0])
	end, _ := strconv.Atoi(delimiters[1])
	for i := start; i <= end; i++ {
		expandedRange = append(expandedRange, i)
	}
	return expandedRange
}

func isFullOverlap(a, b string, full bool) bool {
	aDelimiters := strings.Split(a, "-")
	aMin, _ := strconv.Atoi(aDelimiters[0])
	aMax, _ := strconv.Atoi(aDelimiters[1])
	bDelimiters := strings.Split(b, "-")
	bMin, _ := strconv.Atoi(bDelimiters[0])
	bMax, _ := strconv.Atoi(bDelimiters[1])

	if full {
		if ((aMin <= bMin) && (aMax >= bMax)) ||
			((bMin <= aMin) && (bMax >= aMax)) {
			return true
		}
		return false
	}
	if (aMin >= bMin) && (aMin <= bMax) ||
		(aMax <= bMax) && (aMax >= bMin) ||
		(bMin >= aMin) && (bMin <= aMax) ||
		(bMax <= aMax) && (bMax >= aMin) {
		return true
	}
	return false
}

func countOverlaps(scanner *bufio.Scanner, full bool) int {
	overlapCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		ranges := strings.Split(line, ",")
		if isFullOverlap(ranges[0], ranges[1], full) {
			overlapCount++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return overlapCount
}

func main() {
	fmt.Println("Day 04")
	f, err := os.Open("04.input")
	if err != nil {
		log.Fatal()
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	fullOverlapCount := countOverlaps(scanner, true)
	fmt.Println(fullOverlapCount)

	_, err = f.Seek(0, io.SeekStart)
	scanner = bufio.NewScanner(f)
	if err != nil {
		log.Fatal(err)
	}
	OverlapCount := countOverlaps(scanner, false)
	fmt.Println(OverlapCount)
}
