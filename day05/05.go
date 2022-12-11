package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Move struct {
	count int
	from  int
	to    int
}

type Stack []rune

func (s Stack) String() string {
	return fmt.Sprintf("%s", string(s))
}

func (s *Stack) Pop() (rune, error) {
	l := len(*s)
	if l == 0 {
		return 0, errors.New("Empty stack")
	}
	r := (*s)[l-1]
	*s = (*s)[:l-1]
	return r, nil
}

func (s *Stack) Push(r rune) {
	*s = append(*s, r)
}

func loadStartingStacks() []Stack {
	var stacks = []string{
		"STHFWR",
		"SGDQW",
		"BTW",
		"DRWTNQZJ",
		"FBHGLVTZ",
		"LPTCVBSG",
		"ZBRTWGP",
		"NGMTCJR",
		"LGBW",
	}
	// Move indexing starts from 1, so insert empty stack at index 0
	var startingStacks = []Stack{[]rune{}}
	for _, stack := range stacks {
		var s Stack
		for _, r := range stack {
			s.Push(rune(r))
		}
		startingStacks = append(startingStacks, s)
	}
	return startingStacks
}

func loadMoves(scanner *bufio.Scanner) []Move {
	moveRe := `^move\ ([0-9]+)\ from\ ([0-9])+\ to\ ([0-9]+)$`
	re := regexp.MustCompile(moveRe)

	var moves = []Move{}
	for scanner.Scan() {
		substrings := re.FindStringSubmatch(scanner.Text())
		if substrings == nil {
			continue
		}
		stack, _ := strconv.Atoi(substrings[1])
		from, _ := strconv.Atoi(substrings[2])
		to, _ := strconv.Atoi(substrings[3])
		moves = append(moves, Move{stack, from, to})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return moves
}

func moveStacks(stacks []Stack, moves []Move, mode string) []Stack {
	var tempItem rune
	var tempStack Stack
	for _, m := range moves {
		if mode == "9000" {
			for i := 0; i < m.count; i++ {
				tempItem, _ = stacks[m.from].Pop()
				stacks[m.to].Push(tempItem)
			}
		} else {
			for i := 0; i < m.count; i++ {
				tempItem, _ = stacks[m.from].Pop()
				tempStack.Push(tempItem)
			}
			for range tempStack {
				tempItem, _ = tempStack.Pop()
				stacks[m.to].Push(tempItem)
			}
		}
	}
	return stacks
}

func topCrates(stacks []Stack) Stack {
	var crates Stack
	for _, s := range stacks {
		tempItem, err := s.Pop()
		if err != nil {
			continue
		}
		crates = append(crates, tempItem)
	}
	return crates
}

func main() {

	fmt.Println("Day 05")

	moveType := flag.String("move_type", "9000", "Use the CrateMover 9000 or 9001")
	flag.Parse()

	startingStacks := loadStartingStacks()
	for _, s := range startingStacks {
		fmt.Println(s)
	}

	f, err := os.Open("05.input")
	if err != nil {
		log.Fatal()
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	moves := loadMoves(scanner)

	finalStacks := moveStacks(startingStacks, moves, *moveType)
	for _, s := range finalStacks {
		fmt.Println(s)
	}

	fmt.Println("\ntopCrates:", topCrates(finalStacks))
}
