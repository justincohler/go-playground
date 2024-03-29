package main

import (
	"bufio"
	"feed"
	"flag"
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

// TaskContext is a queue to hold a condition for reading.
type TaskContext struct {
	sync.Mutex
	readCond      *sync.Cond
	writeCond     *sync.Cond
	blockSize     int
	queue         []string
	terminateFlag bool
	wg            sync.WaitGroup
}

// NewTaskContext instantiates a condition for the queue.
func NewTaskContext(blockSize int) *TaskContext {
	ctx := TaskContext{blockSize: blockSize}
	ctx.readCond = sync.NewCond(&ctx)
	ctx.writeCond = sync.NewCond(&ctx)
	return &ctx
}

func addConsumer(ctx *TaskContext, f feed.Feed) {
	defer ctx.wg.Done()
	for { // until records are finished writing
		ctx.Lock()
		ctx.writeCond.Signal()
		ctx.readCond.Wait()
		executeLines(f, ctx.queue)
		ctx.queue = nil
		if ctx.terminateFlag {
			ctx.Unlock()
			break
		}
		ctx.Unlock()
	}
}

func addProducer(ctx *TaskContext) {
	defer ctx.wg.Done()
	scanner := bufio.NewScanner(os.Stdin)
	var res bool
	for {
		ctx.Lock()
		for i := 0; i < ctx.blockSize; i++ {
			res = scanner.Scan()
			if !res {
				break
			}
			ctx.queue = append(ctx.queue, scanner.Text())
		}
		ctx.readCond.Signal()

		if !res { // Wrap up the remaining threads
			ctx.terminateFlag = true
			ctx.readCond.Broadcast()
			ctx.Unlock()
			break
		}
		ctx.writeCond.Wait()
		ctx.Unlock()
	}
}

func processSerially(f feed.Feed) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		executeLine(f, scanner.Text())
	}
}

func main() {
	serial := flag.Bool("s", false, "Process records serially")
	flag.Parse()

	if *serial {
		serialFeed := feed.NewSerialFeed()
		processSerially(serialFeed)
	} else {
		nThreads, _ := strconv.Atoi(os.Args[1])
		blockSize, _ := strconv.Atoi(os.Args[2])
		ctx := NewTaskContext(blockSize)
		parallelFeed := feed.NewParallelFeed()
		for i := 0; i < nThreads; i++ {
			ctx.wg.Add(1)
			go addConsumer(ctx, parallelFeed)
		}
		ctx.wg.Add(1)
		go addProducer(ctx)
		ctx.wg.Wait()
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
		executeLine(f, line)
	}
	return true
}

func executeLine(f feed.Feed, line string) {
	commandName, reqID, body, timestamp, _ := parseLine(line)

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
