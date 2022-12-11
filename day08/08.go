package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Tree struct {
	height int
	seen   bool
}

func (t *Tree) calcSs() {
}

func mapTrees(scanner *bufio.Scanner) [][]*Tree {
	treeMap := [][]*Tree{}

	for scanner.Scan() {
		l := strings.Split(scanner.Text(), "")
		row := []*Tree{}
		for _, t := range l {
			h, _ := strconv.Atoi(t)
			row = append(row, &Tree{h, false})
		}
		treeMap = append(treeMap, row)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return treeMap
}

func printTreeMap(treeMap [][]*Tree) {
	colour := "\033[0m"
	for _, row := range treeMap {
		for _, t := range row {
			if t.seen {
				colour = "\033[32m"
			} else {
				colour = "\033[0m"
			}
			fmt.Printf("%s%d ", colour, t.height)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\033[0m--------------\n")
}

func markVisibilityFromLeft(row []*Tree) {
	tallest := -1
	for _, t := range row {
		if t.height > tallest {
			tallest = t.height
			t.seen = true
		}
	}
}

func markVisibilityFromRight(row []*Tree) {
	tallest := -1
	for i := len(row) - 1; i >= 0; i-- {
		t := row[i]
		if t.height > tallest {
			tallest = t.height
			t.seen = true
		}
	}
}

func markVisibilityFromTop(treeMap [][]*Tree) {
	for i := 0; i < len(treeMap[0]); i++ {
		tallest := -1
		for _, row := range treeMap {
			t := row[i]
			if t.height > tallest {
				tallest = t.height
				t.seen = true
			}
		}
	}
}

func markVisibilityFromBottom(treeMap [][]*Tree) {
	for i := 0; i < len(treeMap); i++ {
		tallest := -1
		for j := len(treeMap[0]) - 1; j >= 0; j-- {
			t := treeMap[j][i]
			if t.height > tallest {
				tallest = t.height
				t.seen = true
			}
		}
	}
}

func countVisibleTrees(treeMap [][]*Tree) int {
	var count int = 0
	for _, row := range treeMap {
		markVisibilityFromLeft(row)
		markVisibilityFromRight(row)
	}
	markVisibilityFromTop(treeMap)
	markVisibilityFromBottom(treeMap)
	for _, row := range treeMap {
		for _, t := range row {
			if t.seen {
				count++
			}
		}
	}
	return count
}

func calculateMaxSS(treeMap [][]*Tree) (int, [2]int) {
	maxSS := 0
	var maxSSCoords [2]int

	numRows := len(treeMap)
	numCols := len(treeMap[0])
	for ri, row := range treeMap {
		for ci, t := range row {
			SScore, upScore, downScore, leftScore, rightScore := 0, 0, 0, 0, 0
			// Up
			for i := ri - 1; i >= 0; i-- {
				upScore++
				if treeMap[i][ci].height >= t.height {
					break
				}
			}
			// Down
			for i := ri + 1; i < numRows; i++ {
				downScore++
				if treeMap[i][ci].height >= t.height {
					break
				}
			}
			// Left
			for i := ci - 1; i >= 0; i-- {
				leftScore++
				if treeMap[ri][i].height >= t.height {
					break
				}
			}
			// Right
			for i := ci + 1; i < numCols; i++ {
				rightScore++
				if treeMap[ri][i].height >= t.height {
					break
				}
			}
			SScore = upScore * downScore * leftScore * rightScore
			if SScore > maxSS {
				maxSS = SScore
				maxSSCoords[0] = ri
				maxSSCoords[1] = ci
			}
		}
	}
	return maxSS, maxSSCoords
}

func main() {
	fmt.Println("Day 08")
	var treeMap [][]*Tree

	f, err := os.Open("08.input")
	//f, err := os.Open("08_example.input")
	if err != nil {
		log.Fatal()
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	treeMap = mapTrees(scanner)
	count := countVisibleTrees(treeMap)
	printTreeMap(treeMap)
	fmt.Println("count: ", count)
	maxSS, maxSSCoords := calculateMaxSS(treeMap)
	fmt.Println("maxSs: ", maxSS, "maxSSCoords: ", maxSSCoords)

}
