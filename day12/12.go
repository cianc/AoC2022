package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

type vertex struct {
	x          int
	y          int
	height     int
	distance   int
	neighbours []neighbour
	visited    bool
	previous   neighbour
}

func (v vertex) String() string {
	return fmt.Sprintf("(%d:%d) -- h:%c -- d:%d -- neighbours: %s -- visited: %s -- previous: %s", v.x, v.y, v.height, v.distance, v.neighbours, v.visited, v.previous)
}

type neighbour struct {
	x          int
	y          int
	edgeLength int
}

func (n neighbour) String() string {
	return fmt.Sprintf("(%d:%d)::%d", n.x, n.y, n.edgeLength)
}

type graph [][]*vertex

type terrain [][]byte

func parseMap(lines []string) (terrain, neighbour, neighbour) {
	heightMap := make(terrain, len(lines))
	start := neighbour{}
	finish := neighbour{}

	for i := 0; i < len(lines); i++ {
		heightMap[i] = make([]byte, len(lines[i]))
	}

	for y, line := range lines {
		for x, b := range line {
			height := b
			if b == 'S' {
				height = 'a'
				start.x = x
				start.y = y
			} else if b == 'E' {
				height = 'z'
				finish.x = x
				finish.y = y
			}
			heightMap[y][x] = byte(height)
		}
	}
	return heightMap, start, finish
}

func printMap(m terrain) {
	for _, row := range m {
		for _, v := range row {
			fmt.Printf("%c", v)
		}
		fmt.Printf("\n")
	}
	hLine := strings.Repeat("=", len(m[0]))
	fmt.Printf("%s\n", hLine)
}

func createGraph(t terrain) graph {
	g := graph{}
	for i, r := range t {
		row := []*vertex{}
		for j, b := range r {
			v := &vertex{x: j, y: i, height: int(b), distance: math.MaxInt, previous: neighbour{x: -1, y: -1}}
			//Up
			if i > 0 {
				n := neighbour{}
				n.edgeLength = 1
				n.x = j
				n.y = i - 1
				v.neighbours = append(v.neighbours, n)
			}
			// Down
			if i < len(t)-1 {
				n := neighbour{}
				n.edgeLength = 1
				n.x = j
				n.y = i + 1
				v.neighbours = append(v.neighbours, n)
			}
			// Left
			if j > 0 {
				n := neighbour{}
				n.edgeLength = 1
				n.x = j - 1
				n.y = i
				v.neighbours = append(v.neighbours, n)
			}
			// Right
			if j < len(t[0])-1 {
				n := neighbour{}
				n.edgeLength = 1
				n.x = j + 1
				n.y = i
				v.neighbours = append(v.neighbours, n)
			}
			row = append(row, v)
		}
		g = append(g, row)
	}
	return g
}

// Start by recording distance of start point to itself: 0
// Then for current vertex, examine unvisited neighbours
// Calculate distance from starting point. If it's < known distance, update the distance and
// previous vertex values in the table.
// Add starting point to list of visited vertices. Won't be visiting here again.
// Move on to the unvisited vertex with the shortest distance from start
func (g *graph) calcShortestPaths(s neighbour) {
	start := (*g)[s.y][s.x]
	start.distance = 0

	unvisited := []*vertex{}
	for _, row := range *g {
		for _, v := range row {
			unvisited = append(unvisited, v)
		}
	}

	for {
		sort.Slice(unvisited, func(i, j int) bool { return unvisited[i].distance < unvisited[j].distance })

		if len(unvisited) == 0 {
			break
		}

		vertex := unvisited[0]
		vertex.visited = true
		unvisited = unvisited[1:]

		for _, n := range vertex.neighbours {
			nvertex := (*g)[n.y][n.x]
			if nvertex.visited {
				continue
			}

			// Artifically boost edgeLength when height gain > 1
			edgeLength := n.edgeLength
			if nvertex.height-vertex.height > 1 {
				edgeLength = int(math.Floor((math.MaxInt / 1e6)))
			}

			startDistance := vertex.distance + edgeLength
			if startDistance < nvertex.distance {
				nvertex.distance = startDistance
				nvertex.previous = neighbour{x: vertex.x, y: vertex.y, edgeLength: edgeLength}
			}
		}
	}
}

func printSolution(g graph, s neighbour, f neighbour) {
	fmt.Printf("============\n")

	goal := g[f.y][f.x]

	for {
		if goal.x == s.x && goal.y == s.y {
			break
		}
		if goal.previous.y > goal.y {
			g[goal.previous.y][goal.previous.x].height = '^'
		}
		if goal.previous.y < goal.y {
			g[goal.previous.y][goal.previous.x].height = 'v'
		}
		if goal.previous.x < goal.x {
			g[goal.previous.y][goal.previous.x].height = '>'
		}
		if goal.previous.x > goal.x {
			g[goal.previous.y][goal.previous.x].height = '<'
		}
		goal = g[goal.previous.y][goal.previous.x]
	}

	for _, row := range g {
		for _, v := range row {
			if v.x == f.x && v.y == f.y {
				fmt.Printf("E")
			} else if v.height != '^' && v.height != '<' && v.height != '>' && v.height != 'v' {
				fmt.Printf(".")
			} else {
				fmt.Printf("%c", v.height)
			}
		}
		fmt.Printf("\n")
	}
}

func findAllStarts(m terrain) []neighbour {
	starts := []neighbour{}

	for y, row := range m {
		for x, v := range row {
			if v == 'a' {
				starts = append(starts, neighbour{x: x, y: y})
			}
		}
	}
	return starts
}

func main() {
	input := flag.String("input", "example", "Run on the 'example' or 'real' input")
	part := flag.Int("part", 1, "Part 1 or 2 of the exercise.")
	printMaps := flag.Bool("printMaps", false, "Print maps")

	flag.Parse()
	fmt.Printf("Day 12\n========\n")

	filename := "12-example-01.input"
	if *input == "real" {
		filename = "12.input"
	}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal()
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	heightMap, start, finish := parseMap(lines)
	if *printMaps {
		printMap(heightMap)
	}
	graph := createGraph(heightMap)

	if *part == 1 {
		graph.calcShortestPaths(start)
		if *printMaps {
			printSolution(graph, start, finish)
		} else {
			fmt.Printf("Shortest path: %d\n", graph[finish.y][finish.x].distance)
		}

	} else {
		starts := findAllStarts(heightMap)
		graph.calcShortestPaths(starts[0])
		shortestPath := graph[finish.y][finish.x].distance
		if *printMaps {
			printSolution(graph, starts[0], finish)
		}
		for i := 1; i < len(starts); i++ {
			graph := createGraph(heightMap)
			graph.calcShortestPaths(starts[i])
			if *printMaps {
				printSolution(graph, starts[i], finish)
			} else {
				fmt.Printf(".")
			}
			sp := graph[finish.y][finish.x].distance

			if sp < shortestPath {
				shortestPath = sp
			}
		}
		fmt.Printf("Shortest path: %d\n", shortestPath)
	}
}
