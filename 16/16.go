package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type valveState struct {
	flowRate   int
	neighbours []string
	open       bool
}

type valves map[string]valveState

func main() {
	input := flag.String("input", "example", "Run on the 'example' or 'real' input")
	part := flag.Int("part", 1, "Part 1 or 2 of the exercise.")
	flag.Parse()

	fmt.Printf("Day 16\n========\n")

	filename := "16-example-part1.input"
	if *input == "real" {
		filename = "16-part1.input"
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	valveMap := createValves(lines)
	distances, _, namesToIndices := shortestPaths(valveMap)
	// printDistances(distances, namesToIndices)
	// printPaths(next, distances, namesToIndices)
	var pressureReleased int
	if *part == 1 {
		pressureReleased = findBestRoute(valveMap, distances, namesToIndices, 30, 1)
	} else {
		pressureReleased = findBestRoute(valveMap, distances, namesToIndices, 26, 2)
	}
	fmt.Println("PRESSURERELEASED:", pressureReleased)
}

func createValves(lines []string) valves {
	valveMap := valves{}
	valveRe := regexp.MustCompile("^Valve ([A-Za-z]+) has flow rate=([0-9]+); tunnels? leads? to valves? ([A-Za-z, ]+)$")
	for _, line := range lines {
		vm := valveRe.FindStringSubmatch(line)
		if vm != nil {
			v := valveState{}
			name := vm[1]
			fr, err := strconv.Atoi(vm[2])
			if err != nil {
				log.Fatal(err)
			}
			v.flowRate = fr
			strings.Split(vm[3], ",")
			for _, n := range strings.Split(vm[3], ",") {
				n = strings.TrimSpace(n)
				v.neighbours = append(v.neighbours, n)
			}
			valveMap[name] = v

		}
	}
	return valveMap
}

func shortestPaths(valves valves) ([][]int, [][]int, map[string]int) {
	distances := make([][]int, len(valves))
	next := make([][]int, len(valves))
	for i := range distances {
		distances[i] = make([]int, len(valves))
		next[i] = make([]int, len(valves))
	}
	namesToIndices := map[string]int{}

	for i := 0; i < len(distances); i++ {
		for j := 0; j < len(distances[0]); j++ {
			distances[i][j] = math.MaxInt / int(1e6)
			next[i][j] = -1
		}
		distances[i][i] = 0
		next[i][i] = i
	}

	i := 0
	for name := range valves {
		namesToIndices[name] = i
		i++
	}

	for name, v := range valves {
		i := namesToIndices[name]
		for _, n := range v.neighbours {
			j := namesToIndices[n]
			distances[i][j] = 1
			distances[j][i] = 1
			next[i][j] = j
			next[j][i] = i
		}
	}

	for k := 0; k < len(distances); k++ {
		for i := 0; i < len(distances); i++ {
			for j := 0; j < len(distances); j++ {
				if distances[i][j] > distances[i][k]+distances[k][j] {
					distances[i][j] = distances[i][k] + distances[k][j]
					next[i][j] = next[i][k]
				}
			}
		}
	}
	return distances, next, namesToIndices
}

func printDistances(paths [][]int, namesToIndices map[string]int) {
	names := []string{}
	for name := range namesToIndices {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Printf("\t")
	for _, name := range names {
		fmt.Printf("%s\t", name)
	}
	fmt.Printf("\n")

	for _, n1 := range names {
		fmt.Printf("%s\t", n1)
		i1 := namesToIndices[n1]
		for _, n2 := range names {
			i2 := namesToIndices[n2]
			fmt.Printf("%d\t", paths[i1][i2])
		}
		fmt.Print("\n")
	}
}

func printPaths(next [][]int, distances [][]int, namesToIndices map[string]int) {
	indicesToNames := map[int]string{}
	names := []string{}
	for name, index := range namesToIndices {
		indicesToNames[index] = name
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name1 := range names {
		i := namesToIndices[name1]
		for _, name2 := range names {
			j := namesToIndices[name2]
			fmt.Printf("%s => %s (%d) ", name1, name2, distances[i][j])
			path := reconstructPath(i, j, next)
			namedPath := []string{}
			for _, p := range path {
				namedPath = append(namedPath, indicesToNames[p])
			}
			fmt.Printf("%s\n", strings.Join(namedPath, "->"))
		}
	}
}

func reconstructPath(from, to int, next [][]int) []int {
	path := []int{}
	for {
		path = append(path, from)
		if from == to {
			break
		}
		from = next[from][to]
	}
	return path
}

// Our queue structures for findBestRoute
type queueEntry struct {
	valveIndex       int
	timeRemaining    int
	valvesOpened     map[int]bool
	pressureReleased int
}
type queue []queueEntry

func (q *queue) enqueue(entry queueEntry) {
	*q = append(*q, entry)
}
func (q *queue) dequeue() (queueEntry, bool) {
	if len(*q) == 0 {
		return queueEntry{}, false
	}
	r := (*q)[0]
	*q = (*q)[1:]
	return r, true
}

// Modified breadth first search to walk the graph of valve opening sequences
func findBestRoute(valveMap valves, distances [][]int, namesToIndices map[string]int, totalTime int, numAgents int) int {
	// We keep track of the max pressure achieved for every unordered set of opened valves.
	// We can use this to calculate answers to part 1 and part 2.
	// The map key is a bitmap of indices, with index one represented as the first bet set, index
	// 2 the second bit etc.
	maxPressureRelasedPerSet := map[int]int{}

	indicesToNames := map[int]string{}
	for name, index := range namesToIndices {
		indicesToNames[index] = name
	}

	// Tee up the first item in the queue, "AA"
	queue := queue{
		queueEntry{
			valveIndex: namesToIndices["AA"], timeRemaining: totalTime, pressureReleased: 0}}

	// Time to walk the graph
	for {
		parent, ok := queue.dequeue()
		if !ok {
			break
		}

		// Figure out what valves are left to open
		valvesToOpen := map[int]bool{}
		for name, info := range valveMap {
			index := namesToIndices[name]
			if _, ok := parent.valvesOpened[index]; ok {
				continue
			}
			if info.flowRate > 0 {
				valvesToOpen[index] = true
			}
		}

		// Iterate over each of these valves left to open
		for childValve := range valvesToOpen {
			// 1 second to open the valve
			timeToOpen := distances[parent.valveIndex][childValve] + 1
			if timeToOpen <= parent.timeRemaining {
				timeRemaining := parent.timeRemaining - timeToOpen
				flowRate := valveMap[indicesToNames[childValve]].flowRate
				// Child inherits all of parents opened valves, plus itself
				childValvesOpened := map[int]bool{childValve: true}
				for vto := range parent.valvesOpened {
					childValvesOpened[vto] = true
				}
				pressureReleased := parent.pressureReleased + (timeRemaining * flowRate)
				// Record pressure for unordered valve set if it's greater than the current max.
				maxPressureKey := 0
				for cvo := range childValvesOpened {
					maxPressureKey |= (1 << cvo)
				}
				mp, ok := maxPressureRelasedPerSet[maxPressureKey]
				if ok {
					if pressureReleased > mp {
						maxPressureRelasedPerSet[maxPressureKey] = pressureReleased
					}
				} else {
					maxPressureRelasedPerSet[maxPressureKey] = pressureReleased
				}
				queue.enqueue(
					queueEntry{
						valveIndex: childValve, timeRemaining: timeRemaining, valvesOpened: childValvesOpened,
						pressureReleased: pressureReleased})
			}
		}
	}
	result := 0
	if numAgents == 1 {
		for _, v := range maxPressureRelasedPerSet {
			if v > result {
				result = v
			}
		}
	} else if numAgents == 2 {
		maskForAllValves := 0
		for index, name := range indicesToNames {
			if valveMap[name].flowRate != 0 {
				maskForAllValves |= (1 << index)
			}
		}
		// We'd like to iterate over all pairs of keys only once, but this is
		// easier than building another datastructure.
		for set1 := range maxPressureRelasedPerSet {
			for set2 := range maxPressureRelasedPerSet {
				// Only want disjoint sets, those are guaranteed to have higher scores
				if set1&set2 == 0 {
					if maxPressureRelasedPerSet[set1]+maxPressureRelasedPerSet[set2] > result {
						result = maxPressureRelasedPerSet[set1] + maxPressureRelasedPerSet[set2]
					}
				}
			}
		}
	} else {
		log.Fatal("Unconfigured number of agents: ", numAgents)
	}
	return result
}

// Heap's algorithm - turns out we don't need this. But it was a pain to debug an off by one error
// in my first try so I'm leaving it here for posterity.
func permute(input []int, size int, distances [][]int) [][]int {
	results := [][]int{}

	if size == 1 {
		results = append(results, input)
		return results
	}

	results = append(results, permute(input, size-1, distances)...)

	swap := func(one *int, two *int) {
		tmp := *one
		*one = *two
		*two = tmp
	}

	for i := 0; i < size-1; i++ {
		if size%2 == 1 {
			swap(&input[0], &input[size-1])
		} else {
			swap(&input[i], &input[size-1])

		}
		results = append(results, permute(input, size-1, distances)...)
	}
	return results
}
