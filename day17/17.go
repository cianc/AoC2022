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

type coordinate struct {
	x int
	y int
}

type rock struct {
	structure                            []coordinate
	lowest, highest, leftmost, rightmost coordinate
}

func (r rock) String() string {
	s := ""
	maxX, maxY := 0, 0
	minX, minY := math.MaxInt, math.MaxInt
	for _, point := range r.structure {
		maxX = int(math.Max(float64(point.x), float64(maxX)))
		minX = int(math.Min(float64(point.x), float64(minX)))
		maxY = int(math.Max(float64(point.y), float64(maxY)))
		minY = int(math.Min(float64(point.y), float64(minY)))
	}
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			foundRock := false
			for _, point := range r.structure {
				if point.x == x && point.y == y {
					foundRock = true
					break
				}
			}
			if foundRock {
				s = fmt.Sprintf("%s# ", s)
			} else {
				s = fmt.Sprintf("%s. ", s)
			}
		}
		s = fmt.Sprintf("%s\n", s)
	}
	return s
}

func (r *rock) move(x int, y int) {
	for i := range r.structure {
		r.structure[i].x += x
		r.structure[i].y += y
	}
	r.highest.y += y
	r.lowest.y += y
	r.leftmost.x += x
	r.rightmost.x += x
}

func (r *rock) shove(direction string, chamber chamber) bool {
	switch direction {
	case ">":
		r.move(1, 0)
		if chamber.checkForIntersection(*r) {
			r.move(-1, 0)
			return false
		}
	case "<":
		r.move(-1, 0)
		if chamber.checkForIntersection(*r) {
			r.move(1, 0)
			return false
		}
	case "down":
		r.move(0, -1)
		if chamber.checkForIntersection(*r) {
			r.move(0, 1)
			return false
		}
	default:
		log.Fatal("??mystery direction??", direction)
	}
	return true
}

func nextRock() func(int) rock {
	rocks := [5]rock{}

	rocks[0].structure = []coordinate{}
	rocks[0].structure = append(rocks[0].structure, coordinate{0, 0})
	rocks[0].structure = append(rocks[0].structure, coordinate{1, 0})
	rocks[0].structure = append(rocks[0].structure, coordinate{2, 0})
	rocks[0].structure = append(rocks[0].structure, coordinate{3, 0})
	rocks[0].highest = coordinate{0, 0}
	rocks[0].lowest = coordinate{0, 0}
	rocks[0].leftmost = coordinate{0, 0}
	rocks[0].rightmost = coordinate{3, 0}

	rocks[1].structure = []coordinate{}
	rocks[1].structure = append(rocks[1].structure, coordinate{1, 0})
	rocks[1].structure = append(rocks[1].structure, coordinate{0, 1})
	rocks[1].structure = append(rocks[1].structure, coordinate{1, 1})
	rocks[1].structure = append(rocks[1].structure, coordinate{2, 1})
	rocks[1].structure = append(rocks[1].structure, coordinate{1, 2})
	rocks[1].highest = coordinate{1, 2}
	rocks[1].lowest = coordinate{1, 0}
	rocks[1].leftmost = coordinate{0, 1}
	rocks[1].rightmost = coordinate{2, 1}

	rocks[2].structure = []coordinate{}
	rocks[2].structure = append(rocks[2].structure, coordinate{0, 0})
	rocks[2].structure = append(rocks[2].structure, coordinate{1, 0})
	rocks[2].structure = append(rocks[2].structure, coordinate{2, 0})
	rocks[2].structure = append(rocks[2].structure, coordinate{2, 1})
	rocks[2].structure = append(rocks[2].structure, coordinate{2, 2})
	rocks[2].highest = coordinate{2, 2}
	rocks[2].lowest = coordinate{0, 0}
	rocks[2].leftmost = coordinate{0, 0}
	rocks[2].rightmost = coordinate{2, 0}

	rocks[3].structure = []coordinate{}
	rocks[3].structure = append(rocks[3].structure, coordinate{0, 0})
	rocks[3].structure = append(rocks[3].structure, coordinate{0, 1})
	rocks[3].structure = append(rocks[3].structure, coordinate{0, 2})
	rocks[3].structure = append(rocks[3].structure, coordinate{0, 3})
	rocks[3].highest = coordinate{0, 3}
	rocks[3].lowest = coordinate{0, 0}
	rocks[3].leftmost = coordinate{0, 0}
	rocks[3].rightmost = coordinate{0, 0}

	rocks[4].structure = []coordinate{}
	rocks[4].structure = append(rocks[4].structure, coordinate{0, 0})
	rocks[4].structure = append(rocks[4].structure, coordinate{1, 0})
	rocks[4].structure = append(rocks[4].structure, coordinate{0, 1})
	rocks[4].structure = append(rocks[4].structure, coordinate{1, 1})
	rocks[4].highest = coordinate{0, 1}
	rocks[4].lowest = coordinate{0, 0}
	rocks[4].leftmost = coordinate{0, 0}
	rocks[4].rightmost = coordinate{1, 0}

	rockIndex := 0
	return func(yOffset int) rock {
		r := rock{}
		for _, point := range rocks[rockIndex].structure {
			r.structure = append(r.structure, coordinate{point.x, point.y})
		}
		r.highest = coordinate{rocks[rockIndex].highest.x, rocks[rockIndex].highest.y}
		r.lowest = coordinate{rocks[rockIndex].lowest.x, rocks[rockIndex].lowest.y}
		r.leftmost = coordinate{rocks[rockIndex].leftmost.x, rocks[rockIndex].leftmost.y}
		r.rightmost = coordinate{rocks[rockIndex].rightmost.x, rocks[rockIndex].rightmost.y}
		r.move(2, yOffset)
		rockIndex = (rockIndex + 1) % len(rocks)
		return r
	}
}

type chamber struct {
	width       int
	rocks       []coordinate
	highestRock coordinate
}

func (c chamber) String() string {
	s := ""
	for y := c.highestRock.y; y >= 0; y-- {
		s = fmt.Sprintf("%s|", s)
		for x := 0; x < c.width; x++ {
			foundRock := false
			for _, rock := range c.rocks {
				if rock.x == x && rock.y == y {
					foundRock = true
					break
				}
			}
			if foundRock {
				s = fmt.Sprintf("%s# ", s)
			} else {
				s = fmt.Sprintf("%s. ", s)
			}
		}
		s = fmt.Sprintf("%s|\n", s)
	}
	s = fmt.Sprintf("%s+%s+\n", s, strings.Repeat("-", c.width*2))
	return s
}

func (c chamber) checkForIntersection(r rock) bool {
	// Hit the sides or the bottom
	if r.leftmost.x < 0 || r.lowest.y < 0 || r.rightmost.x >= c.width {
		return true
	}

	// Hit a resting rock
	if r.lowest.y > c.highestRock.y {
		return false
	}
	for _, point := range r.structure {
		for i := len(c.rocks) - 1; i >= 0; i-- {
			if point.x == c.rocks[i].x && point.y == c.rocks[i].y {
				return true
			}
		}
	}
	return false
}

func (c chamber) fingerprint() [100]coordinate {
	result := [100]coordinate{}
	// For the 100 highest rocks, record their vertical distance to the top highest rock.
	sort.Slice(c.rocks, func(i int, j int) bool { return c.rocks[i].y < c.rocks[j].y })
	maxY := c.rocks[len(c.rocks)-1].y
	for i := 0; i < 100; i++ {
		result[i] = c.rocks[len(c.rocks)-200+i]
		result[i].y = maxY - result[i].y
	}
	return result
}

func (c chamber) fingerprint2() [50]int {
	result := [50]int{}

	// Create one bitmap per highest 50 rows in the chamber.
	sort.Slice(c.rocks, func(i int, j int) bool { return c.rocks[i].y < c.rocks[j].y })
	for i := 0; i < 50; i++ {
		rowResult := 0
		y := c.highestRock.y - i
		for j := len(c.rocks) - 1; j >= 0; j-- {
			if c.rocks[j].y == y {
				rowResult |= (1 << c.rocks[j].x)
			} else if c.rocks[j].y > y {
				continue
			} else if c.rocks[j].y < y {
				break
			}
		}
		result[i] = rowResult
	}
	return result
}

func runSimulation(findPattern bool, nextRock func(int) rock, chamber *chamber, jetInput []string, jetIndex int, rockCount int, debug bool) int {
	lcm := len(jetInput) * 5
	fingerprints := map[[50]int][2]int{}

	for i := 0; i < rockCount; i++ {
		// Trim the rock list periodically
		if i%500 == 0 && i > 0 {
			chamber.rocks = chamber.rocks[len(chamber.rocks)-2000:]
		}

		// Periodically check for repeating pattern
		if findPattern {
			if i%lcm == 0 && i > 0 {
				fingerprint := chamber.fingerprint2()
				lastSeenInfo, ok := fingerprints[fingerprint]
				if !ok {
					fingerprints[fingerprint] = [2]int{i, chamber.highestRock.y + 1}
				} else {
					prefixHeight := lastSeenInfo[1]
					currentHeight := chamber.highestRock.y + 1
					repeatingHeight := currentHeight - prefixHeight
					repeatingPeriod := i - lastSeenInfo[0]
					numberOfRepititions := int(math.Floor(float64(rockCount) / float64(repeatingPeriod)))
					remainderRuns := (rockCount - lastSeenInfo[0]) % repeatingPeriod
					runSimulation(false, nextRock, chamber, jetInput, jetIndex, remainderRuns, debug)
					remainderHeight := chamber.highestRock.y + 1 - currentHeight
					return prefixHeight + (repeatingHeight * numberOfRepititions) + remainderHeight
				}
			}
		}

		// Generate a new rock and let it move and fall until it gets stuck.
		rock := nextRock(chamber.highestRock.y + 4)
		if debug {
			printChamberInFlight(*chamber, rock)
		}
		for {
			rock.shove(jetInput[jetIndex], *chamber)
			jetIndex = (jetIndex + 1) % len(jetInput)
			if debug {
				printChamberInFlight(*chamber, rock)
			}
			ok := rock.shove("down", *chamber)
			if !ok {
				chamber.rocks = append(chamber.rocks, rock.structure...)
				if rock.highest.y > chamber.highestRock.y {
					chamber.highestRock = rock.highest
				}
				break
			}
			if debug {
				printChamberInFlight(*chamber, rock)
			}
		}
	}
	return chamber.highestRock.y + 1
}

func printChamberInFlight(chamber chamber, rock rock) {
	height := int(math.Max(float64(rock.highest.y), float64(chamber.highestRock.y)))
	fmt.Printf("\n")
	for y := height; y >= 0; y-- {
		fmt.Printf("|")
		for x := 0; x < chamber.width; x++ {
			foundRock := false
			for _, point := range rock.structure {
				if point.x == x && point.y == y {
					foundRock = true
					break
				}
			}
			for _, rock := range chamber.rocks {
				if rock.x == x && rock.y == y {
					foundRock = true
					break
				}
			}
			if foundRock {
				fmt.Printf("# ")
			} else {
				fmt.Printf(". ")
			}
		}
		fmt.Printf("|\n")
	}
	fmt.Printf("+%s+\n", strings.Repeat("-", chamber.width*2))
}

func readInput(filename string) []string {
	input := []string{}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}
	return input
}

func main() {
	inputType := flag.String("input", "example", "Run on the 'example' or 'real' input")
	part := flag.Int("part", 1, "Part 1 or 2 of the exercise.")
	debug := flag.Bool("debug", false, "turn on debug output")
	flag.Parse()

	fmt.Printf("Day 17\n========\n")

	filename := "17-example-part1.input"
	if *inputType == "real" {
		filename = "17-part1.input"
	}
	input := readInput(filename)
	jetInput := []string{}
	for _, i := range input {
		jetInput = append(jetInput, strings.Split(i, "")...)
	}

	c := chamber{width: 7, highestRock: coordinate{-1, -1}}
	result := -1
	if *part == 1 {
		result = runSimulation(false, nextRock(), &c, jetInput, 0, 2022, *debug)
	} else {
		result = runSimulation(true, nextRock(), &c, jetInput, 0, 1000000000000, *debug)
	}
	fmt.Println("RESULT", result)
}
