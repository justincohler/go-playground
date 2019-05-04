package main

import (
	"bufio"
	feeder "feed"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	cADD      = "ADD"
	cREMOVE   = "REMOVE"
	cCONTAINS = "CONTAINS"
	cSTRING   = "STRING"
)

func parseLine(line string) (string, int64, string, int64, bool) {
	line = strings.Replace(line, "{", "", -1)
	line = strings.Replace(line, "}", "", -1)
	args := strings.Split(line, ",")
	commandName := args[0]

	if len(args) < 2 {
		return "", int64(0), "", int64(0), true
	}
	reqID, _ := strconv.ParseInt(args[1], 10, 64)

	var timestamp int64
	var body string
	if len(args) == 3 {
		timestamp, _ = strconv.ParseInt(args[2], 10, 64)
	} else if len(args) == 4 {
		body = args[2]
		timestamp, _ = strconv.ParseInt(args[3], 10, 64)
	}

	return commandName, reqID, body, timestamp, false
}

func executeLines(wg *sync.WaitGroup, feed feeder.Feed, lines []string) {
	defer wg.Done()

	for _, line := range lines {
		commandName, reqID, body, timestamp, err := parseLine(line)

		if err {
			break
		}

		var status string

		switch commandName {
		case cADD:
			fmt.Println("Adding...")
			feed.Add(body, timestamp)
			fmt.Println("{{", reqID, "}, {SUCCESS}}")
		case cREMOVE:
			fmt.Println("Removing...")
			if feed.Remove(timestamp) {
				status = "SUCCESS"
			} else {
				status = "FAILED"
			}
			fmt.Println("{{", reqID, "}, {", status, "}}")

		case cCONTAINS:
			fmt.Println("Containsing...")
			if feed.Contains(timestamp) {
				status = "YES"
			} else {
				status = "NO"
			}
			fmt.Println("{{", reqID, "}, {", status, "}}")

		case cSTRING:
			fmt.Println("Stringing...")
			fmt.Println("{{", reqID, "}, {", feed.String(), "}}")
		}
	}
}

func main() {
	var wg sync.WaitGroup

	nThreads, _ := strconv.Atoi(os.Args[1])
	blockSize, _ := strconv.Atoi(os.Args[2])

	runtime.GOMAXPROCS(nThreads)

	scanner := bufio.NewScanner(os.Stdin)

	var lines []string

	var res bool
	var f feeder.Feed
	f = feeder.NewFeed()
	for {
		lines = make([]string, blockSize)
		for i := 0; i < blockSize; i++ {
			res = scanner.Scan()
			if !res {
				break
			}
			lines[i] = scanner.Text()
		}
		if !res {
			break
		}
		wg.Add(1)
		go executeLines(&wg, f, lines)
	}

	// Add leftover records (N % blockSize)
	wg.Add(1)
	go executeLines(&wg, f, lines)

	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
