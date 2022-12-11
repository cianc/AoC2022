package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"
)

func scoreRune(r rune) int {
	if unicode.IsUpper(r) {
		return int(r) - 38
	}
	return int(r) - 96
}

func sumMisplacedPriorities(scanner *bufio.Scanner) int {
	var (
		backpack, pocket1, pocket2 string
		prioritySum                = 0
		pocket1Set                 map[rune]bool
	)

	for scanner.Scan() {
		backpack = scanner.Text()
		pocket1 = backpack[:len(backpack)/2]
		pocket2 = backpack[len(backpack)/2:]

		pocket1Set = make(map[rune]bool)
		for _, v := range pocket1 {
			pocket1Set[v] = true
		}
		for _, v := range pocket2 {
			if pocket1Set[v] {
				prioritySum += scoreRune(v)
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return prioritySum
}

func backpacksIntersection(backpacks []string) []rune {
	var itemsCount = make(map[rune]int)
	var intersection []rune
	var numBackpacks = len(backpacks)
	for i, b := range backpacks {
		for _, r := range b {
			if itemsCount[r] == i+1 {
				continue
			} else if itemsCount[r] == i {
				itemsCount[r] = i + 1
				continue
			}
		}
	}
	for r, i := range itemsCount {
		if i == numBackpacks {
			intersection = append(intersection, r)
		}
	}
	return intersection
}

func sumBadgePriorites(scanner *bufio.Scanner) int {
	var backpacks []string
	var scanCounter = 0
	var intersection = []rune{}
	var prioritySum = 0

	for scanner.Scan() {
		scanCounter++
		backpacks = append(backpacks, scanner.Text())
		if scanCounter == 3 {
			scanCounter = 0
			intersection = backpacksIntersection(backpacks)
			backpacks = nil
			for _, v := range intersection {
				prioritySum += scoreRune(v)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return prioritySum
}

func main() {
	fmt.Println("Day 03")
	f, err := os.Open("03.input")
	if err != nil {
		log.Fatal()
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	fmt.Println("misplacedPrioritySum:", sumMisplacedPriorities(scanner))
	_, err = f.Seek(0, io.SeekStart)
	scanner = bufio.NewScanner(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("badgePrioritySum", sumBadgePriorites(scanner))
}
