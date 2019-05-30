package main

import (
	"bufio"
	"fmt"
	"imgutil"
	"os"
	"queue"
	"strconv"
	"strings"
	"sync"
)

// Context allows for further parallelization
type Context struct {
	wg             sync.WaitGroup
	scanFilterLock sync.Mutex
	scanCond       *sync.Cond
	filterInCond   *sync.Cond
	qFilter        queue.Stack
	qSave          queue.Stack
	nThreads       int
	readComplete   bool
}

// NewContext returns a new app context
func NewContext(nThreads int) *Context {
	ctx := Context{}
	ctx.scanCond = sync.NewCond(&ctx.scanFilterLock)
	ctx.filterInCond = sync.NewCond(&ctx.scanFilterLock)
	ctx.qFilter = queue.NewUnbounded()
	ctx.qSave = queue.NewUnbounded()
	ctx.nThreads = nThreads
	return &ctx
}

// QueueRequest is utility to send over queue s.t. our I/O is constrained to one thread.
type QueueRequest struct {
	Image   *imgutil.PNGImage
	Filters []string
	OutFile string
}

func spawnImageProcessor(ctx *Context) {
	defer ctx.wg.Done()
	for !(ctx.readComplete && ctx.qFilter.Empty()) {
		ctx.filterInCond.L.Lock()
		ctx.filterInCond.Wait()
		if ctx.readComplete {
			ctx.filterInCond.L.Unlock()
			break
		}
		value := ctx.qFilter.Pop()
		ctx.filterInCond.L.Unlock()

		fmt.Println("Thread FILTER received message")
		request := value.(QueueRequest)
		request.Image = request.Image.ApplyFilters(request.Filters)

		// Tell the saver there's a file to save
		ctx.qSave.Push(request)
		fmt.Println("Thread FILTER sent SAVE message")

		// Tell the csv reader to send more work
		ctx.scanCond.L.Lock()
		fmt.Println("Thread FILTER ready for work!")
		ctx.scanCond.Signal()
		ctx.scanCond.L.Unlock()
	}
}

func parseLine(line string) (string, string, []string) {
	lineArgs := strings.Split(strings.Replace(line, " ", "", -1), ",")
	return lineArgs[0], lineArgs[1], lineArgs[2:]
}

func processSerial() {
	filePath := os.Args[1]
	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		inFile, outFile, filters := parseLine(scanner.Text())
		img, _ := imgutil.Load(inFile)
		img.Threads = 1

		fmt.Println("Serially applying filters")
		filteredImg := img.ApplyFilters(filters)
		fmt.Println("Serially saving out file")
		filteredImg.Save(outFile)
	}
}

func processParallel() {
	sThreads, filePath := os.Args[1], os.Args[2]
	nThreads, _ := strconv.Atoi(sThreads)
	ctx := NewContext(nThreads)

	ctx.wg.Add(1)
	go spawnImageProcessor(ctx)

	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		inFile, outFile, filters := parseLine(scanner.Text())

		img, _ := imgutil.Load(inFile)
		img.Threads = ctx.nThreads
		request := QueueRequest{Image: img, Filters: filters, OutFile: outFile}

		ctx.qFilter.Push(request)
		ctx.filterInCond.Signal()
		fmt.Println("Thread 0* sent FILTER signal")

		ctx.scanCond.L.Lock()
		ctx.scanCond.Wait()
		ctx.scanCond.L.Unlock()
	}
	ctx.filterInCond.L.Lock()
	ctx.readComplete = true
	ctx.filterInCond.Broadcast()
	ctx.filterInCond.L.Unlock()

	ctx.wg.Wait()

	for !ctx.qSave.Empty() {
		value := ctx.qSave.Pop()
		request := value.(QueueRequest)
		fmt.Println("Thread SAVE received message")

		request.Image.Save(request.OutFile)
	}

	fmt.Println("Finished Writing All Images.")
}

func main() {
	if len(os.Args) < 3 {
		processSerial()
	} else {
		processParallel()
	}
}
