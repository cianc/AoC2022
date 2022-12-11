package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type coord struct {
	x int
	y int
}

type vector struct {
	direction string
	magnitude int
}

type shortRope struct {
	head coord
	tail coord
}

type longRope struct {
	ropes [9]*shortRope
}

func main() {
	fmt.Println("Day 09")

	//f, err := os.Open("09-example.input")
	//f, err := os.Open("09-example-02.input")
	f, err := os.Open("09.input")
	if err != nil {
		log.Fatal()
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	moves := parseMoves(scanner)
	v1 := pullTheShortRope(moves)
	fmt.Println("Exercise 1:", len(v1))
	v2 := pullTheLongRope(moves)
	fmt.Println("Exercise 2:", len(v2))
}

func parseMoves(scanner *bufio.Scanner) []vector {
	moves := []vector{}
	for scanner.Scan() {
		move := strings.Split(scanner.Text(), " ")
		mag, _ := strconv.Atoi(move[1])
		moves = append(moves, vector{move[0], mag})
	}
	return moves
}

func pullTheShortRope(moves []vector) map[coord]bool {
	visited := make(map[coord]bool)
	rope := shortRope{}
	for _, m := range moves {
		vs := rope.move(m)
		for k := range vs {
			visited[k] = true
		}
	}
	return visited
}

func (rope *shortRope) move(m vector) map[coord]bool {
	visited := make(map[coord]bool)
	visited[coord{rope.tail.x, rope.tail.y}] = true
	for i := 0; i < m.magnitude; i++ {
		switch m.direction {
		case "U":
			rope.head.y++
		case "D":
			rope.head.y--
		case "L":
			rope.head.x--
		case "R":
			rope.head.x++
		}
		rope.tailFollowHead()
		visited[coord{rope.tail.x, rope.tail.y}] = true
	}
	return visited
}

func (rope *shortRope) tailFollowHead() {
	xDiff := rope.head.x - rope.tail.x
	yDiff := rope.head.y - rope.tail.y
	diff := coord{xDiff, yDiff}
	switch diff {
	// Move up
	case coord{0, 2}:
		rope.tail.y++
		// Move down
	case coord{0, -2}:
		rope.tail.y--
		// Move right
	case coord{2, 0}:
		rope.tail.x++
		// Move left
	case coord{-2, 0}:
		rope.tail.x--
		// Move diagonal up&right
	case coord{1, 2}, coord{2, 1}, coord{2, 2}:
		rope.tail.x++
		rope.tail.y++
		// Move diagonal down&right
	case coord{1, -2}, coord{2, -1}, coord{2, -2}:
		rope.tail.x++
		rope.tail.y--
		// Move diagonal down&left
	case coord{-1, -2}, coord{-2, -1}, coord{-2, -2}:
		rope.tail.x--
		rope.tail.y--
		// Move diagonal up&left
	case coord{-1, 2}, coord{-2, 1}, coord{-2, 2}:
		rope.tail.x--
		rope.tail.y++
	}
}

func pullTheLongRope(moves []vector) map[coord]bool {
	visited := make(map[coord]bool)
	lrope := longRope{}
	for i := 0; i < len(lrope.ropes); i++ {
		lrope.ropes[i] = &shortRope{}
	}

	for _, m := range moves {
		fmt.Println("Move:", m)
		vs := lrope.move(m)
		for k := range vs {
			visited[k] = true
		}
	}
	return visited
}

func (lrope *longRope) move(m vector) map[coord]bool {
	visited := make(map[coord]bool)
	l := lrope.ropes[len(lrope.ropes)-1]
	r1 := lrope.ropes[0]
	visited[coord{l.tail.x, l.tail.y}] = true

	for i := 0; i < m.magnitude; i++ {
		switch m.direction {
		case "U":
			r1.head.y++
		case "D":
			r1.head.y--
		case "L":
			r1.head.x--
		case "R":
			r1.head.x++
		}
		fmt.Println("i: ", i)
		for j, r := range lrope.ropes {
			if j > 0 {
				r.head = lrope.ropes[j-1].tail
			}
			r.tailFollowHead()
			fmt.Println(j, "--------", coord{r.tail.x, r.tail.y})

			// correct end state on exercise 1, but incorrect tail position count
			if j == len(lrope.ropes)-1 {
				fmt.Println("tail visit: ", coord{r.tail.x, r.tail.y})
				visited[coord{r.tail.x, r.tail.y}] = true
				fmt.Println(j, len(lrope.ropes)-1, visited)
			}
		}
	}
	return visited
}
