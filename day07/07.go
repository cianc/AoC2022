package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Node struct {
	category rune
	name     string
	size     int
	parent   *Node
	children []*Node
}

func printFs(root *Node, depth int) {
	depth++
	padding := strings.Repeat(" ", depth)
	switch root.category {
	case 'd':
		fmt.Println(padding, "- ", root.name, " (dir=", root.size, ")")
	case 'f':
		fmt.Println(padding, "- ", root.name, " (file=", root.size, ")")
	}
	for _, child := range root.children {
		printFs(child, depth+2)
	}
}

func findTargetDir(rootDir *Node, neededSpace int) (bool, *Node) {
	var targetDir *Node
	success := false
	if (rootDir.category == 'd') && (rootDir.size >= neededSpace) {
		targetDir = rootDir
		success = true
	}
	for _, child := range rootDir.children {
		s, candidateDir := findTargetDir(child, neededSpace)
		if s && (candidateDir.size < targetDir.size) {
			targetDir = candidateDir
		}
	}
	return success, targetDir
}

func sumDirSizes(root *Node, max int) int {
	sum := 0
	if (root.category == 'd') && (root.size <= max) {
		sum += root.size
	}
	for _, child := range root.children {
		sum += sumDirSizes(child, max)
	}
	return sum
}

func setDirSizes(root *Node) int {
	for _, child := range root.children {
		if child.category == 'f' {
			root.size += child.size
		} else {
			s := setDirSizes(child)
			root.size += s
		}
	}
	return root.size
}

func constructFS(scanner *bufio.Scanner) *Node {
	cdRe := regexp.MustCompile(`^\$\ cd\ (.+)$`)
	lsRe := regexp.MustCompile(`^\$\ ls$`)
	dirRe := regexp.MustCompile(`^dir\ (.+)$`)
	fileRe := regexp.MustCompile(`^([0-9]+)\ (.+)$`)

	rootNode := Node{'d', "/", 0, nil, nil}
	currentNode := &rootNode

	for scanner.Scan() {
		line := scanner.Text()
		reMatch := []string{}
		if cdRe.MatchString(line) {
			reMatch = cdRe.FindStringSubmatch(line)
			target := reMatch[1]
			switch target {
			case "..":
				currentNode = currentNode.parent
			case "/":
				currentNode = &rootNode
			default:
				foundChild := false
				for _, c := range currentNode.children {
					if c.name == target {
						currentNode = c
						foundChild = true
						break
					}
				}
				if !foundChild {
					childNode := Node{'d', target, 0, currentNode, nil}
					currentNode.children = append(currentNode.children, &childNode)
					currentNode = &childNode

				}
			}
		} else if lsRe.MatchString(line) {
		} else if dirRe.MatchString(line) {
			reMatch = dirRe.FindStringSubmatch(line)
			dirName := reMatch[1]
			childNode := Node{'d', dirName, 0, currentNode, nil}
			currentNode.children = append(currentNode.children, &childNode)
		} else if fileRe.MatchString(line) {
			reMatch = fileRe.FindStringSubmatch(line)
			fSize, _ := strconv.Atoi(reMatch[1])
			fName := reMatch[2]
			childNode := Node{'f', fName, fSize, currentNode, nil}
			currentNode.children = append(currentNode.children, &childNode)
		} else {
			fmt.Println("MISS: ", line)
		}
	}
	return &rootNode
}

func main() {
	fmt.Println("Day 07")

	f, err := os.Open("07.input")
	//f, err := os.Open("07_test_input_01.input")
	if err != nil {
		log.Fatal()
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	fsRoot := constructFS(scanner)
	_ = setDirSizes(fsRoot)
	//printFs(fsRoot, 0)
	sum := sumDirSizes(fsRoot, 100000)
	fmt.Println("sum:", sum)
	neededSpace := 30000000 - (70000000 - fsRoot.size)
	fmt.Println("neededSpace:", neededSpace)
	_, targetDir := findTargetDir(fsRoot, neededSpace)
	fmt.Println("targetDir:", targetDir.name, ":", targetDir.size)
}
