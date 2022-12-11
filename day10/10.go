package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Day 10")
	f, err := os.Open("10.input")
	//f, err := os.Open("10-example-01.input")
	//f, err := os.Open("10-example-02.input")
	if err != nil {
		log.Fatal()
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	ops := parseOps(scanner)
	cpu := cpu{register: 1}
	history := cpu.run(ops)
	fmt.Println("Sum of signal strengths:", sumSignalStrength(history))
	emulateCRT(history)
}

type cpu struct {
	cycle    int
	register int
}

func (c *cpu) run(ops []operation) []cpu {
	history := []cpu{}
	for _, op := range ops {
		switch op.name {
		case "noop":
			history = append(history, c.noop()...)
		case "addx":
			history = append(history, c.addx(op.value)...)
		}
	}
	// Fake record for state after we completed the last op
	history = append(history, cpu{cycle: c.cycle + 1, register: c.register})

	return history
}

func (c *cpu) noop() []cpu {
	c.cycle++
	return []cpu{*c}
}

func (c *cpu) addx(v int) []cpu {
	history := []cpu{}
	for i := 0; i < 2; i++ {
		c.cycle++
		history = append(history, *c)
	}
	c.register += v
	return history
}

type operation struct {
	name  string
	value int
}

func parseOps(scanner *bufio.Scanner) []operation {
	ops := []operation{}
	for scanner.Scan() {
		l := scanner.Text()
		o := strings.Split(l, " ")
		op := operation{name: o[0]}
		if op.name == "addx" {
			op.value, _ = strconv.Atoi(o[1])
		}
		ops = append(ops, op)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return ops
}

func sumSignalStrength(hist []cpu) int {
	sum := 0
	for _, h := range hist {
		if h.cycle%40 == 20 {
			sum += (h.cycle * h.register)
		}
	}
	return sum
}

func emulateCRT(hist []cpu) {
	for r := 0; r < 6; r++ {
		for p := 0; p < 40; p++ {
			sprite := 0
			cycle := (r * 40) + p + 1
			for _, c := range hist {
				if c.cycle == cycle {
					sprite = c.register
					break
				}
			}
			if (p >= sprite-1) && (p <= sprite+1) {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
}
