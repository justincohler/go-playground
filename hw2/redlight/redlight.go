package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"sync"
)

// Year structures quarterly counts.
type Year struct {
	year     int
	quarters map[string]uint64
	mux      sync.Mutex
}

// Inc increments a quarter's count.
func (y *Year) Inc(q string) {
	y.mux.Lock()
	defer y.mux.Unlock()
	y.quarters[q]++
}

// Value returns a quarter's count.
func (y *Year) Value(q string) uint64 {
	y.mux.Lock()
	defer y.mux.Unlock()
	return y.quarters[q]
}

// ProcessLines loops through lines and adds
// citations to the appropriate year/quarter.
func ProcessLines(wg *sync.WaitGroup, yearsPtr *map[int]Year, lines [][]string) {
	defer wg.Done()
	const violationDateIdx int = 3
	const violationsIdx int = 4
	// years := *yearsPtr
	for i, line := range lines {
		if i == 0 {
			fmt.Println(line[violationDateIdx], line[violationsIdx])
		}
	}
}

func main() {
	THREADS := runtime.NumCPU()
	fmt.Println(THREADS)

	var wg sync.WaitGroup

	filePath := os.Args[1]
	file, _ := os.Open(filePath)
	defer file.Close()

	reader := csv.NewReader(file)
	lines, _ := reader.ReadAll()

	N := len(lines)

	years := make(map[int]Year)
	for i := 0; i < THREADS; i++ {
		fmt.Println("Starting thread", i)
		wg.Add(1)
		if i == THREADS-1 {
			go ProcessLines(&wg, &years, lines[i*N/THREADS:])
		} else {
			go ProcessLines(&wg, &years, lines[i*N/THREADS:(i+1)*N/THREADS])
		}

	}

	wg.Wait()
}
