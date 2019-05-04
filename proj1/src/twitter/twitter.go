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

// TaskQueue is a queue to hold a condition for reading.
type TaskQueue struct {
	sync.Mutex
	cond          *sync.Cond
	terminateFlag bool
}

// NewTaskQueue instantiates a condition for the queue.
func NewTaskQueue() *TaskQueue {
	q := TaskQueue{}
	q.cond = sync.NewCond(&q)
	return &q
}

func addQueueConsumer(wg *sync.WaitGroup, q *TaskQueue, f feed.Feed, lines *[]string) {
	defer wg.Done()
	for { // Continuously wait for work
		q.Lock()
		q.cond.Wait()
		executeLines(f, *lines)
		if q.terminateFlag {
			fmt.Println("Wrapping up")
			q.Unlock()
			return
		}
		q.Unlock()
	}
}

func main() {
	var wg sync.WaitGroup

	nThreads, _ := strconv.Atoi(os.Args[1])
	blockSize, _ := strconv.Atoi(os.Args[2])
	scanner := bufio.NewScanner(os.Stdin)

	q := NewTaskQueue()
	f := feed.NewFeed()

	var lines []string
	var res bool

	counter := 0

	for i := 0; i < nThreads; i++ {
		wg.Add(1)
		go addQueueConsumer(&wg, q, f, &lines)
	}

	for {
		q.Lock()
		lines = make([]string, blockSize)
		for i := 0; i < blockSize; i++ {
			res = scanner.Scan()
			if !res {
				break
			}
			lines[i] = scanner.Text()
			counter++
		}
		q.Unlock()
		q.cond.Signal()
		if !res {
			q.Lock()
			q.terminateFlag = true
			q.Unlock()
			q.cond.Broadcast()
			break
		}
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

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

func executeLines(f feed.Feed, lines []string) bool {
	if len(lines) == 0 {
		return false
	}
	for _, line := range lines {
		commandName, reqID, body, timestamp, err := parseLine(line)

		if err {
			break
		}

		var status string

		switch commandName {
		case cADD:
			f.Add(body, timestamp)
			fmt.Println("{{", reqID, "}, {SUCCESS}}")
		case cREMOVE:
			if f.Remove(timestamp) {
				status = "SUCCESS"
			} else {
				status = "FAILED"
			}
			fmt.Println("{{", reqID, "}, {", status, "}}")

		case cCONTAINS:
			if f.Contains(timestamp) {
				status = "YES"
			} else {
				status = "NO"
			}
			fmt.Println("{{", reqID, "}, {", status, "}}")

		case cSTRING:
			fmt.Println("{{", reqID, "}, {", f.String(), "}}")
		}
	}
	return true
}
