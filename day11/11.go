package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	part := flag.Int("part", 1, "Which part of the exercise")
	input := flag.String("input", "example", "Run on the 'example' or 'real' input")
	flag.Parse()
	fmt.Println("Day 11")

	filename := "11-example-01.input"
	if *input == "real" {
		filename = "11.input"
	}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal()
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	monkeys := parseMonkeys(scanner)
	roundCount := 20
	if *part == 2 {
		roundCount = 10000
	}
	for i := 1; i <= roundCount; i++ {
		runRound(monkeys)
		if (i == 1) || (i == 20) || (i%1000 == 0) {
			fmt.Println("== After round", i, "==	")
			for j, monkey := range monkeys {
				fmt.Println("Monkey ", j, ":", monkey.inspections)
			}
		}
	}
	fmt.Println("Monkey business: ", calcMonkeyBusiness(monkeys))
}

type monkey struct {
	// Per-item worry levels
	worry []int
	// [0] is the operator, and [1] and [2] are the operands (polish notation)
	operation [3]string
	// [0] is modulo value to satisfy, [1] is monkey index to pass to on true, [2] on false
	test        [3]int
	inspections int
}

func parseMonkeys(scanner *bufio.Scanner) []*monkey {
	monkeyRE := regexp.MustCompile(`^Monkey\ [0-9]+:$`)

	monkeys := []*monkey{}
	reMatch := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		reMatch = monkeyRE.FindStringSubmatch(line)
		if reMatch != nil {
			m := monkey{}
			scanner.Scan()
			line = scanner.Text()
			m.worry = parseStartingItems(line)
			scanner.Scan()
			line = scanner.Text()
			m.operation = parseOperations(line)
			lines := []string{}
			for i := 0; i < 3; i++ {
				scanner.Scan()
				lines = append(lines, scanner.Text())
			}
			m.test = parseTest(lines)
			monkeys = append(monkeys, &m)
		}
	}
	return monkeys
}

func parseStartingItems(input string) []int {
	startingRE := regexp.MustCompile(`\s+Starting\ items:\ (([0-9]+(, )?)+)$`)
	reMatch := startingRE.FindStringSubmatch(input)
	results := []int{}

	if reMatch == nil {
		log.Fatal("Failed to parse input, expecting 'Starting' line, got: ", input)
	}
	startingItems := strings.Split(reMatch[1], ", ")
	for _, startingItem := range startingItems {
		s, err := strconv.Atoi(startingItem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, s)
	}
	return results
}

func parseOperations(input string) [3]string {
	operationRE := regexp.MustCompile(`^\s+Operation:\ new\ =\ ((old|[0-9])+\ [\+\*]\ (old|[0-9])+)$`)

	reMatch := operationRE.FindStringSubmatch(input)
	if reMatch == nil {
		log.Fatal("Failed to parse input, expecting 'Operation' line, got: ", input)
	}

	o := strings.Split(reMatch[1], " ")
	return [3]string{o[1], o[0], o[2]}
}

func parseTest(inputs []string) [3]int {
	testRE := regexp.MustCompile(`^\s+Test:\ divisible\ by\ ([0-9]+)$`)
	outcomeRE := regexp.MustCompile(`\s+If\ (true|false):\ throw\ to\ monkey\ ([0-9]+)$`)
	result := [3]int{}

	reMatch := testRE.FindStringSubmatch(inputs[0])
	if reMatch == nil {
		log.Fatal("Failed to parse input, expecting 'Test' line, got: ", inputs[0])
	}
	x, _ := strconv.Atoi(reMatch[1])
	result[0] = x

	for i := 1; i <= 2; i++ {
		reMatch = outcomeRE.FindStringSubmatch(inputs[i])
		if reMatch == nil {
			log.Fatal("Failed to parse input, expecting test output line, got: ", inputs[i])
		}
		if reMatch[1] == "true" {
			x, _ := strconv.Atoi(reMatch[2])
			result[1] = x
		} else {
			x, _ := strconv.Atoi(reMatch[2])
			result[2] = x
		}
	}
	return result
}

func runRound(monkeys []*monkey) {
	modProduct := calcModProduct(monkeys)
	for _, monkey := range monkeys {
		monkey.inspect(modProduct)
		monkey.testAndThrow(monkeys)
	}
}

func (monkey *monkey) inspect(modProduct int) {
	operator := monkey.operation[0]
	operands := [2]int{}
	var newWorry int

	for i, oldWorry := range monkey.worry {
		for i, o := range monkey.operation[1:] {
			if o == "old" {
				operands[i] = oldWorry
			} else {
				operands[i], _ = strconv.Atoi(o)
			}
		}

		switch operator {
		case "*":
			newWorry = operands[0] * operands[1]
		case "+":
			newWorry = operands[0] + operands[1]
		default:
			log.Fatal("Unexpected operator '", operator, "'")
		}
		monkey.worry[i] = newWorry % modProduct
	}
}

func (monkey *monkey) testAndThrow(monkeys []*monkey) {
	destMonkey := 0
	for _, worry := range monkey.worry {
		monkey.inspections++
		if worry%monkey.test[0] == 0 {
			destMonkey = monkey.test[1]
		} else {
			destMonkey = monkey.test[2]
		}
		monkeys[destMonkey].worry = append(monkeys[destMonkey].worry, worry)
	}
	monkey.worry = []int{}
}

func calcModProduct(monkeys []*monkey) int {
	modProduct := 1
	for _, monkey := range monkeys {
		modProduct *= monkey.test[0]
	}
	return modProduct
}

func calcMonkeyBusiness(monkeys []*monkey) int {
	topInspections := [2]int{}
	for _, monkey := range monkeys {
		if monkey.inspections > topInspections[0] {
			topInspections[1] = topInspections[0]
			topInspections[0] = monkey.inspections
		} else if monkey.inspections > topInspections[1] {
			topInspections[1] = monkey.inspections
		}
	}
	return topInspections[0] * topInspections[1]
}
