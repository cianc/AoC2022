package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

type coord struct {
	x int
	y int
}

type coveredRange struct {
	from coord
	to   coord
}

type sensorBeaconPair struct {
	sensor coord
	beacon coord
}

func (sbp sensorBeaconPair) MinMaxX(y int) (int, int) {
	sbSeparation := distance(sbp.sensor, sbp.beacon)
	yOffset := int(math.Abs(float64(sbp.sensor.y - y)))
	maxX := sbp.sensor.x + (sbSeparation - yOffset)
	minX := sbp.sensor.x - (sbSeparation + yOffset)
	return minX, maxX
}

func main() {
	input := flag.String("input", "example", "Run on the 'example' or 'real' input")
	part := flag.Int("part", 1, "Part 1 or 2 of the exercise.")
	flag.Parse()

	fmt.Printf("Day 15\n========\n")

	filename := "15-example-part1.input"
	if *input == "real" {
		filename = "15-part1.input"
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
	sensorsAndBeacons := parsePositions(lines)
	yRow := 10
	maxCoord := 20
	if *input == "real" {
		yRow = 2000000
		maxCoord = 4000000
	}

	result := 0
	if *part == 1 {
		result = countPositionsWithNoBeacon(sensorsAndBeacons, yRow)
		fmt.Println("\nRESULT:", result)
	} else {
		result := findMissingBeacons(sensorsAndBeacons, maxCoord)
		fmt.Println("RESULT:", (result[0].x*4000000)+result[0].y)
	}
}

func parsePositions(lines []string) map[sensorBeaconPair]bool {
	locationRE := regexp.MustCompile(
		`^Sensor\ at\ x=(\-?[0-9]+),\ y=(\-?[0-9]+):\ closest\ beacon\ is\ at\ x=(\-?[0-9]+),\ y=(\-?[0-9]+)`)
	result := map[sensorBeaconPair]bool{}

	for _, line := range lines {
		lm := locationRE.FindStringSubmatch(line)
		if lm != nil {
			sensorX, err := strconv.Atoi(lm[1])
			if err != nil {
				log.Fatal("Failed to convert string ", lm[1], " to integer")
			}
			sensorY, err := strconv.Atoi(lm[2])
			if err != nil {
				log.Fatal("Failed to convert string ", lm[2], " to integer")
			}
			beaconX, err := strconv.Atoi(lm[3])
			if err != nil {
				log.Fatal("Failed to convert string ", lm[3], " to integer")
			}
			beaconY, err := strconv.Atoi(lm[4])
			if err != nil {
				log.Fatal("Failed to convert string ", lm[4], " to integer")
			}
			result[sensorBeaconPair{coord{sensorX, sensorY}, coord{beaconX, beaconY}}] = true
		}
	}
	return result
}

func distance(point1, point2 coord) int {
	return int(math.Abs(float64(point1.x-point2.x))) + int(math.Abs(float64(point1.y-point2.y)))
}

func beaconPossible(sbp sensorBeaconPair, point coord, existing bool) bool {
	sbSeparation := distance(sbp.sensor, sbp.beacon)

	if existing && sbp.beacon.x == point.x && sbp.beacon.y == point.y {
		return true
	}
	if distance(sbp.sensor, point) <= sbSeparation {
		return false
	}
	return true
}

func countPositionsWithNoBeacon(sensorsAndBeacons map[sensorBeaconPair]bool, yRow int) int {
	results := map[int]bool{}
	for sbp := range sensorsAndBeacons {
		minX, maxX := sbp.MinMaxX(yRow)
		for x := minX; x <= maxX; x++ {
			if !beaconPossible(sbp, coord{x, yRow}, true) {
				results[x] = true
			}
		}
	}
	return len(results)
}

func findMissingBeacons(sensorsAndBeacons map[sensorBeaconPair]bool, maxCoord int) []coord {
	result := []coord{}
	for y := 0; y <= maxCoord; y++ {
		if y%10000 == 0 {
			fmt.Printf(".")
		}
	OUTER:
		for x := 0; x <= maxCoord; x++ {
			possible := true
			for sbp := range sensorsAndBeacons {
				possible = beaconPossible(sbp, coord{x, y}, false)
				if !possible {
					_, maxX := sbp.MinMaxX(y)
					x = maxX
					continue OUTER
				}
			}
			if possible {
				result = append(result, coord{x, y})
			}
		}
	}
	return result
}
