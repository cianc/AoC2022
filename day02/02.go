package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func scoreStrategy(scanner *bufio.Scanner) int {
	shapeScores := map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
	}
	winningCombos := map[string]string{
		"A": "B",
		"B": "C",
		"C": "A",
	}
	losingCombos := map[string]string{
		"A": "C",
		"B": "A",
		"C": "B",
	}

	score := 0
	var round string
	var myShape string
	for scanner.Scan() {
		round = scanner.Text()
		s := strings.Split(round, " ")
		// Lose
		if s[1] == "X" {
			myShape = losingCombos[s[0]]
			// Draw
		} else if s[1] == "Y" {
			myShape = s[0]
			score += 3
			// Win
		} else {
			myShape = winningCombos[s[0]]
			score += 6
		}
		score += shapeScores[myShape]
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return score
}

func main() {
	fmt.Println("Day 02")
	f, err := os.Open("02.input")
	if err != nil {
		log.Fatal()
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	score := scoreStrategy(scanner)
	fmt.Println(score)
}
