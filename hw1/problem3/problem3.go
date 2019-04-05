package main

import (
	"bufio"
	"fmt"
	"os"
)

// Main simply reads in a file path and target word,
// then adds to a word count each time the target word
// is found in a file. Finally main() prints the word
// count.
//
// Assumptions: One word per line
func main() {

	filePath, targetWord := os.Args[1], os.Args[2]

	wordCount := 0

	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if word := scanner.Text(); word == targetWord {
			wordCount++
		}
	}

	fmt.Println(wordCount)
}
