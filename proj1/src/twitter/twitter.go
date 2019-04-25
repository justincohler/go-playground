package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func parseLine(line string) []string {
	line = strings.Replace(line, "{", "", -1)
	line = strings.Replace(line, "}", "", -1)
	args := strings.Split(line, ",")
	return args
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var lines [][]string
	for scanner.Scan() {
		lines = append(lines, parseLine(scanner.Text()))
	}

	fmt.Println(lines[0])

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
