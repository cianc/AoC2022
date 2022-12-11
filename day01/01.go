package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func parseCalories(scanner *bufio.Scanner) []int {
	var (
		calories []int
		count    int = -1
		line     string
	)
	for scanner.Scan() {
		line = scanner.Text()
		if line == "" {
			if count >= 0 {
				calories = append(calories, count)
				count = -1
			}
		} else {
			calory, err := strconv.Atoi(line)
			if err != nil {
				log.Fatal(err)
			}
			if count == -1 {
				count = calory
			} else {
				count += calory
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	if count >= 0 {
		calories = append(calories, count)
		count = -1
	}
	return calories
}

func topn(calories []int, n int) []int {
	sort.Slice(calories, func(i, j int) bool {
		return calories[i] > calories[j]
	})
	return calories[:3]

}

func main() {
	fmt.Println("Day 01")
	f, err := os.Open("01.input")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	cals := parseCalories(scanner)
	maxCals := topn(cals, 1)
	fmt.Println(maxCals[0])
	top3Cals := topn(cals, 3)
	fmt.Println(top3Cals[0] + top3Cals[1] + top3Cals[2])
}
