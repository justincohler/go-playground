package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
)

// Violations structures all years' Violation data.
type Violations struct {
	years map[int]*Year
	sync.Mutex
}

// Inc increments a year-quarter pair by a given amount
func (v *Violations) Inc(y int, q int, amount uint64) {
	v.Lock()
	defer v.Unlock()
	_, exists := v.years[y]
	if !exists {
		year := Year{year: y}
		v.years[y] = &year
		year.Inc(q, amount)
	} else {
		year := v.years[y]
		year.Inc(q, amount)
	}
}

// Print outputs a pretty, sorted output of total violations each year-quarter.
func (v *Violations) Print() {
	yearKeys := make([]int, 0, len(v.years))
	for yearKey := range v.years {
		yearKeys = append(yearKeys, yearKey)
	}
	sort.Ints(yearKeys)

	for _, yearKey := range yearKeys {
		fmt.Println(yearKey)
		year := v.years[yearKey]

		for q := 1; q <= 4; q++ {
			fmt.Printf("\tQ%v=%v\n", q, year.quarters[q])
		}
	}
}

// Year structures quarterly counts.
type Year struct {
	year     int
	quarters map[int]uint64
	sync.Mutex
}

// Inc increments a quarter's count.
func (y *Year) Inc(q int, amount uint64) {
	y.Lock()
	defer y.Unlock()
	if y.quarters == nil {
		y.quarters = make(map[int]uint64)
	}
	y.quarters[q] += amount
}

// ProcessLines loops through lines and adds
// citations to the appropriate year/quarter.
func ProcessLines(wg *sync.WaitGroup, violations *Violations, lines [][]string, thread int) {
	defer wg.Done()
	const violationDateIdx int = 3
	const violationsIdx int = 4
	for _, line := range lines {
		// Reference Go Layout Time: Mon Jan 2 15:04:05 MST 2006
		dateFormat := "01/02/2006"
		t, err := time.Parse(dateFormat, line[violationDateIdx])

		// e.g. Ignore header line, erroneous data
		if err != nil {
			continue
		}

		q := int((t.Month()-1)/3) + 1
		if q > 4 {
			fmt.Println(t.Month())
		}
		violationCount, _ := strconv.ParseUint(line[violationsIdx], 10, 64)
		violations.Inc(t.Year(), q, violationCount)
	}
}

func main() {
	THREADS := runtime.NumCPU()
	var wg sync.WaitGroup

	filePath := os.Args[1]
	file, _ := os.Open(filePath)
	defer file.Close()

	reader := csv.NewReader(file)
	lines, _ := reader.ReadAll()

	N := len(lines)

	violations := Violations{years: make(map[int]*Year)}

	for i := 0; i < THREADS; i++ {
		wg.Add(1)
		if i == THREADS-1 {
			go ProcessLines(&wg, &violations, lines[i*N/THREADS:], i)
		} else {
			go ProcessLines(&wg, &violations, lines[i*N/THREADS:(i+1)*N/THREADS], i)
		}

	}

	wg.Wait()
	violations.Print()
}
