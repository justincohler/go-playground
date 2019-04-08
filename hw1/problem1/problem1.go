package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Intersect takes in two slices and returns a
// slice of intersecting elements from the inputs.
// These intersecting elements must contain a "2".
func Intersect(slice1 []string, slice2 []string) []string {

	slice1Map := make(map[string]bool)
	for _, num := range slice1 {
		if strings.Contains(num, "2") {
			slice1Map[num] = true
		}
	}

	// Map to contain unique final values present in both lists
	twos := make(map[string]bool)
	for _, num := range slice2 {
		// Credit: https://stackoverflow.com/questions/2050391/how-to-check-if-a-map-contains-a-key-in-go
		if _, exists := slice1Map[num]; exists && strings.Contains(num, "2") {
			twos[num] = true
		}
	}

	var intersection []string
	for key := range twos {
		intersection = append(intersection, key)
	}
	return intersection
}

func main() {

	var slice1, slice2 []string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter list 1:")
	scanner.Scan()
	slice1 = strings.Split(strings.Replace(scanner.Text(), " ", "", -1), ",")

	fmt.Println("Enter list 2:")
	scanner.Scan()
	slice2 = strings.Split(strings.Replace(scanner.Text(), " ", "", -1), ",")

	intersection := Intersect(slice1, slice2)
	fmt.Println(intersection)
}
