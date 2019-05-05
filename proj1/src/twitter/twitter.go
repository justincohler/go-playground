package main

import (
	"bufio"
	"feed"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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
		ctx.readCond.Wait()
		if len(ctx.queue) < ctx.blockSize {
			executeLines(f, ctx.queue)
			ctx.queue = nil
		} else {
			executeLines(f, ctx.queue[:ctx.blockSize])
			ctx.queue = ctx.queue[ctx.blockSize:]
		}
		ctx.Unlock()
		if !ctx.terminateFlag {
			ctx.writeCond.Signal()
		} else {
			break
		}
	}
}

func addProducer(ctx *TaskContext) {
	// defer ctx.wg.Done()
	scanner := bufio.NewScanner(os.Stdin)
	var res bool
	for {
		time.Sleep(10 * time.Millisecond)
		ctx.Lock()
		// ctx.writeCond.Wait()
		for i := 0; i < ctx.blockSize; i++ {
			res = scanner.Scan()
			if !res {
				break
			}
			ctx.queue = append(ctx.queue, scanner.Text())
		}
		ctx.Unlock()
		ctx.readCond.Signal()
		if !res { // Wrap up the remaining threads
			ctx.Lock()
			ctx.terminateFlag = true
			ctx.Unlock()
			for len(ctx.queue) > 0 {
				ctx.readCond.Signal()
			}
			ctx.readCond.Broadcast()
			break
		}
	}
}

func main() {
	nThreads, _ := strconv.Atoi(os.Args[1])
	blockSize, _ := strconv.Atoi(os.Args[2])

	ctx := NewTaskContext(blockSize)
	f := feed.NewFeed()

	for i := 0; i < nThreads; i++ {
		ctx.wg.Add(1)
		go addConsumer(ctx, f)
	}

	// ctx.wg.Add(1)
	// go addProducer(ctx)
	// ctx.writeCond.Signal()
	addProducer(ctx)
	ctx.wg.Wait()
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
