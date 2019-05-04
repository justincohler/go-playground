package main

import (
	"bufio"
	"feed"
	"fmt"
	"os"
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

func parseLine(line string) (string, int64, string, int64) {
	fmt.Println(line)
	line = strings.Replace(line, "{", "", -1)
	line = strings.Replace(line, "}", "", -1)
	args := strings.Split(line, ",")
	commandName := args[0]

	fmt.Println(args)
	reqID, _ := strconv.ParseInt(args[1], 10, 64)

	var timestamp int64
	var body string
	if len(args) == 3 {
		timestamp, _ = strconv.ParseInt(args[2], 10, 64)
	} else if len(args) == 4 {
		body = args[2]
		timestamp, _ = strconv.ParseInt(args[4], 10, 64)
	}

	return commandName, reqID, body, timestamp
}

func executeLines(wg *sync.WaitGroup, feed feed.Feed, lines []string) {
	defer wg.Done()

	for _, line := range lines {
		fmt.Println(line)
		commandName, reqID, body, timestamp := parseLine(line)

		fmt.Println(commandName, reqID, body, timestamp)
		var status string

		switch commandName {
		case cADD:
			fmt.Println("{{", reqID, "}, {SUCCESS}}")
			feed.Add(body, timestamp)

		case cREMOVE:
			if feed.Remove(timestamp) {
				status = "SUCCESS"
			} else {
				status = "FAILED"
			}
			fmt.Println("{{", reqID, "}, {", status, "}}")

		case cCONTAINS:
			if feed.Contains(timestamp) {
				status = "YES"
			} else {
				status = "NO"
			}
			fmt.Println("{{", reqID, "}, {", status, "}}")

		case cSTRING:
			fmt.Println("{{", reqID, "}, {", feed.String(), "}}")
		}
	}
}

func main() {
	var feed feed.Feed
	var wg sync.WaitGroup

	// nThreads, _ := strconv.Atoi(os.Args[1])
	blockSize, _ := strconv.Atoi(os.Args[2])

	scanner := bufio.NewScanner(os.Stdin)

	var lines []string

	var res bool
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
		go executeLines(&wg, feed, lines)
	}

	// Add leftover records (N % blockSize)
	wg.Add(1)
	go executeLines(&wg, feed, lines)

	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
