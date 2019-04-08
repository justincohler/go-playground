package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// AddRangesToSet takes a list of ranges
// and returns a set of distinct numbers
// in string form to the calling function.
func AddRangesToSet(ranges *[]string) map[int]bool {
	set := make(map[int]bool)
	for _, r := range *ranges {
		if strings.Contains(r, "-") {
			elements := strings.Split(r, "-")
			from, _ := strconv.Atoi(elements[0])
			to, _ := strconv.Atoi(elements[1])
			for i := from; i <= to; i++ {
				set[i] = true
			}
		} else {
			num, _ := strconv.Atoi(r)
			set[num] = true
		}
	}
	return set
}

// MergeRanges takes a string of comma-
// separated ranges (or single numbers)
// and returns a list of unique integers
// comprising the given ranges.
func MergeRanges(text string) []int {
	ranges := strings.Split(strings.Replace(text, " ", "", -1), ",")
	rangeSet := AddRangesToSet(&ranges)

	var merged []int
	for key := range rangeSet {
		merged = append(merged, key)
	}
	sort.Ints(merged)
	return merged
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter integers:")
	scanner.Scan()
	merged := MergeRanges(scanner.Text())
	fmt.Println(merged)
}
