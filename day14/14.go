package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	regex "regexp"
	"strconv"
	"strings"
)

type coord struct {
	x int
	y int
}

type rockmap map[coord]string

func (rm rockmap) minmaxXY(part ...int) (int, int, int, int) {
	p := 1
	if len(part) > 0 {
		p = part[0]
	}

	maxX, maxY := math.MinInt, math.MinInt
	minX, minY := 500, 0
	for k := range rm {
		if k.x > maxX {
			maxX = k.x
		} else if k.x < minX {
			minX = k.x
		}
		if k.y > maxY {
			maxY = k.y
		} else if k.y < minY {
			minY = k.y
		}
	}
	if p == 2 {
		maxY += 2
	}
	return minX, maxX, minY, maxY
}

// Lower case string, not String because we don't want a real stringer because we want to be
// able to pass parameters.
func (rm rockmap) string(part ...int) string {
	p := 1
	if len(part) > 0 {
		p = part[0]
	}
	result := ""
	minX, maxX, minY, maxY := rm.minmaxXY(p)
	for y := minY; y <= maxY; y++ {
		result = fmt.Sprintf("%s%d\t", result, y)
		for x := minX; x <= maxX; x++ {
			if x == 500 && y == 0 {
				result = fmt.Sprintf("%s+", result)
				continue
			}
			if y == maxY {
				result = fmt.Sprintf("%s#", result)
				continue
			}
			v, ok := rm[coord{x, y}]
			switch {
			case !ok:
				result = fmt.Sprintf("%s.", result)
			case v == "#":
				result = fmt.Sprintf("%s#", result)
			case v == "o":
				result = fmt.Sprintf("%so", result)
			}
		}
		result = fmt.Sprintf("%s\n", result)
	}
	return result
}

func main() {
	input := flag.String("input", "example", "Run on the 'example' or 'real' input")
	part := flag.Int("part", 1, "Part 1 or 2 of the exercise.")
	flag.Parse()

	fmt.Printf("Day 13\n========\n")

	var filename string
	if *input == "real" {
		filename = "14-part1.input"
	} else {
		filename = "14-example-part1.input"
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
	rockMap := parseScans(lines)

	sum := 0
	_, _, _, maxY := rockMap.minmaxXY()
	for ; ; sum++ {
		result := simulateGrain(&rockMap, maxY+2, *part)
		fmt.Printf(".")

		if *part == 1 {
			if result == (coord{-1, -1}) {
				break
			}
		} else {
			if result == (coord{500, 0}) {
				break
			}
		}
	}
	fmt.Printf("\n%s\n", rockMap.string(*part))
	fmt.Printf("RESULT: %d\n", sum)
}

func parseScans(input []string) rockmap {
	coordRE := regex.MustCompile(`([0-9]+,[0-9]+)`)

	coords := map[coord]string{}
	var prevX, prevY int
	// Pull out the rock coordinates from the input
	for _, l := range input {
		for i, cm := range coordRE.FindAllStringSubmatch(l, -1) {
			if cm[1] != "" {
				stringCoords := strings.Split(cm[1], ",")
				if x, err := strconv.Atoi(stringCoords[0]); err != nil {
					log.Fatal(err)
				} else if y, err := strconv.Atoi(stringCoords[1]); err != nil {
					log.Fatal(err)
				} else {
					coords[coord{x, y}] = "#"
					if i > 0 {
						switch {
						case x-prevX > 0:
							for i := 1; i <= (x - prevX); i++ {
								coords[coord{x - i, y}] = "#"
							}
						case x-prevX < 0:
							for i := 1; i <= (prevX - x); i++ {
								coords[coord{x + i, y}] = "#"
							}
						case y-prevY > 0:
							for i := 1; i < (y - prevY); i++ {
								coords[coord{x, y - i}] = "#"
							}
						case y-prevY < 0:
							for i := 1; i < (prevY - y); i++ {
								coords[coord{x, y + i}] = "#"
							}
						}
					}
					prevX = x
					prevY = y
				}
			}
		}
	}
	return coords
}

func simulateGrain(rm *rockmap, maxY int, part int) coord {
	grain := coord{500, 0}

	blckd := func(below coord) bool {
		_, blocked := (*rm)[below]
		if part == 1 {
			return blocked
		}
		return blocked || below.y >= maxY
	}

	for {
		below := coord{grain.x, grain.y + 1}
		above := coord{}
		//_, blocked := (*rm)[below]
		blocked := blckd(below)
		if !blocked {
			above = grain
			grain = below
		} else {
			below := coord{grain.x - 1, grain.y + 1}
			//_, blocked = (*rm)[below]
			blocked := blckd(below)
			if !blocked {
				above = grain
				grain = below
			} else {
				below := coord{grain.x + 1, grain.y + 1}
				//_, blocked = (*rm)[below]
				blocked := blckd(below)
				if !blocked {
					above = grain
					grain = below
				} else {
					return coord{grain.x, grain.y}
				}
			}
		}

		if part == 1 {
			if grain.y > maxY {
				return coord{-1, -1}
			}
		}
		(*rm)[grain] = "o"
		delete((*rm), above)
	}
	return coord{-1, -1}
}
