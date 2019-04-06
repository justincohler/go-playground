package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Set is a generic interface for Hashset facilities.
type Set interface {
	Union(other Set) Set
	Intersect(other Set) Set
	Diff(other Set) Set
}

// IntSet is an integer set implementation of Set.
type IntSet struct {
	set map[int]bool
}

// Initialize adds a slice of string numbers to the set.
func (recv *IntSet) Initialize(line []string) {
	recv.set = make(map[int]bool)
	for _, str := range line {
		num, _ := strconv.Atoi(str)
		recv.set[num] = true
	}
}

// Print formats to set formatting
func (recv *IntSet) Print() string {
	res := "{"
	for k := range recv.set {
		res += strconv.Itoa(k) + ","
	}

	res = res[:len(res)-1]
	res += "}"
	return res
}

// Union returns the containing struct's set unioned with
// the "other" Intset provided.
func (recv *IntSet) Union(other IntSet) IntSet {
	union := IntSet{set: make(map[int]bool)}

	for k := range recv.set {
		union.set[k] = true
	}

	for k := range other.set {
		union.set[k] = true
	}
	return union
}

// Intersect returns elements shared between the
// containing struct's set and the "other" IntSet.
func (recv *IntSet) Intersect(other IntSet) IntSet {
	intersection := IntSet{set: make(map[int]bool)}

	for k := range other.set {
		if _, exists := recv.set[k]; exists {
			intersection.set[k] = true
		}
	}
	return intersection
}

// Diff returns differing elements between the
// containing struct's set and the "other" IntSet.
func (recv *IntSet) Diff(other IntSet) IntSet {
	diff := IntSet{set: make(map[int]bool)}

	union := recv.Union(other)
	intersection := recv.Intersect(other)

	for k := range union.set {
		if _, exists := intersection.set[k]; !exists {
			diff.set[k] = true
		}
	}

	return diff
}

func main() {

	action, filePath := os.Args[1], os.Args[2]

	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	line1 := strings.Split(scanner.Text(), " ")
	scanner.Scan()
	line2 := strings.Split(scanner.Text(), " ")

	var set1, set2 IntSet
	set1.Initialize(line1)
	set2.Initialize(line2)

	switch action {
	case "-u":
		union := set1.Union(set2)
		fmt.Println(union.Print())
	case "-i":
		intersection := set1.Intersect(set2)
		fmt.Println(intersection.Print())
	case "-d":
		diff := set1.Diff(set2)
		fmt.Println(diff.Print())
	}
}
