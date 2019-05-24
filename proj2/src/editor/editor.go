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
	wg                sync.WaitGroup
	scanAndFilterLock sync.Mutex
	filterCond        *sync.Cond
	scanCond          *sync.Cond
	saveLock          sync.Mutex
	saveCond          *sync.Cond
	qFilter           queue.Stack
	qSave             queue.Stack
	nThreads          int
	readComplete      bool
}

// NewContext returns a new app context
func NewContext(nThreads int) *Context {
	ctx := Context{}
	ctx.filterCond = sync.NewCond(&ctx.scanAndFilterLock)
	ctx.scanCond = sync.NewCond(&ctx.scanAndFilterLock)
	ctx.saveCond = sync.NewCond(&ctx.saveLock)
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
	for {
		ctx.filterCond.L.Lock()
		ctx.filterCond.Wait()
		if ctx.readComplete {
			ctx.filterCond.L.Unlock()
			break
		}
		value := ctx.qFilter.Pop()
		request := value.(QueueRequest)
		ctx.filterCond.L.Unlock()

		request.Image = request.Image.ApplyFilters(request.Filters)

		// Tell the saver there's a file to save
		ctx.saveCond.L.Lock()
		ctx.qSave.Push(request)
		ctx.saveCond.Signal()
		ctx.saveCond.L.Unlock()

		// Tell the csv reader to send more work
		ctx.scanCond.Signal()
	}
}

func spawnImageWriter(ctx *Context, fileCount *int) {
	savedCount := 0
	for !ctx.readComplete || savedCount < *fileCount {
		ctx.saveCond.L.Lock()
		ctx.saveCond.Wait()
		fmt.Println("Received SAVE signal.")
		value := ctx.qSave.Pop()
		ctx.saveCond.L.Unlock()
		request := value.(QueueRequest)

		request.Image.Save(request.OutFile)
		savedCount++
	}
}

func parseLine(line string) (string, string, []string) {
	lineArgs := strings.Split(strings.Replace(line, " ", "", -1), ",")
	return lineArgs[0], lineArgs[1], lineArgs[2:]
}

func processSequential() {
	filePath := os.Args[1]
	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		inFile, outFile, filters := parseLine(scanner.Text())
		img, _ := imgutil.Load(inFile)
		img.Threads = 1

		filteredImg := img.ApplyFilters(filters)
		filteredImg.Save(outFile)
	}
}

func processParallel() {
	sThreads, filePath := os.Args[1], os.Args[2]
	nThreads, _ := strconv.Atoi(sThreads)
	ctx := NewContext(nThreads)
	fileCount := 0

	for i := 0; i < nThreads; i++ {
		ctx.wg.Add(1)
		go spawnImageProcessor(ctx)
	}

	go spawnImageWriter(ctx, &fileCount)

	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		inFile, outFile, filters := parseLine(scanner.Text())

		img, _ := imgutil.Load(inFile)
		img.Threads = ctx.nThreads
		request := QueueRequest{img, filters, outFile}

		ctx.filterCond.L.Lock()
		ctx.qFilter.Push(request)
		ctx.filterCond.Signal()
		ctx.filterCond.L.Unlock()

		ctx.scanCond.L.Lock()
		ctx.scanCond.Wait()
		ctx.scanCond.L.Unlock()
		fileCount++
	}
	ctx.filterCond.L.Lock()
	ctx.readComplete = true
	ctx.filterCond.Broadcast()
	ctx.filterCond.L.Unlock()

	ctx.wg.Wait()
}

func main() {
	if len(os.Args) < 3 {
		processSequential()
	} else {
		processParallel()
	}
}
