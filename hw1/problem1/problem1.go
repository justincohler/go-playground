package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func intersect(slice1 []string, slice2 []string) []string {
	// Map to contain unique slice1 values
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

	if scanner.Scan() {
		slice1 = strings.Split(strings.Replace(scanner.Text(), " ", "", -1), ",")
		fmt.Println("Enter list 2:")

		if scanner.Scan() {
			slice2 = strings.Split(strings.Replace(scanner.Text(), " ", "", -1), ",")
		} else {
			fmt.Println("Failed to read list 2.")
		}

	} else {
		fmt.Println("Failed to read list 1.")
	}

	intersection := intersect(slice1, slice2)
	fmt.Println(intersection)
}
