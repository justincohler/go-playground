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
	"sync/atomic"
)

// Context allows for further parallelization
type Context struct {
	wg             sync.WaitGroup
	scanFilterLock sync.Mutex
	scanCond       *sync.Cond
	filterCond     *sync.Cond
	saveLock       sync.Mutex
	saveCond       *sync.Cond
	qFilter        queue.Stack
	qSave          queue.Stack
	nThreads       int
	readComplete   bool
	fileCount      int32
	activeFilters  int32
}

// NewContext returns a new app context
func NewContext(nThreads int) *Context {
	ctx := Context{}
	ctx.scanCond = sync.NewCond(&ctx.scanFilterLock)
	ctx.filterCond = sync.NewCond(&ctx.scanFilterLock)
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
	goID    int
}

func spawnImageProcessor(ctx *Context, goID int) {
	defer ctx.wg.Done()
	for !(ctx.readComplete && ctx.qFilter.Empty()) {
		ctx.filterCond.L.Lock()
		ctx.filterCond.Wait()
		if ctx.readComplete {
			ctx.filterCond.L.Unlock()
			break
		}
		value := ctx.qFilter.Pop()
		ctx.filterCond.L.Unlock()

		fmt.Println("Thread", goID, "received FILTER message")
		request := value.(QueueRequest)
		request.goID = goID
		request.Image = request.Image.ApplyFilters(request.Filters)

		// Tell the saver there's a file to save
		ctx.saveCond.L.Lock()
		ctx.qSave.Push(request)
		// fmt.Println("Thread", goID, "sent SAVE message")
		ctx.saveCond.Signal()
		ctx.saveCond.L.Unlock()

		// Tell the csv reader to send more work
		ctx.scanCond.L.Lock()
		atomic.AddInt32(&ctx.activeFilters, -1)
		fmt.Println("Thread", goID, "ready for work!")
		ctx.scanCond.Signal()
		ctx.scanCond.L.Unlock()
	}
	fmt.Println("Thread", goID, "done filtering")
}

func spawnImageWriter(ctx *Context) {
	defer ctx.wg.Done()
	for !ctx.readComplete || atomic.LoadInt32(&ctx.fileCount) > 0 {

		ctx.saveCond.L.Lock()
		ctx.saveCond.Wait()
		value := ctx.qSave.Pop()
		ctx.saveCond.L.Unlock()
		request := value.(QueueRequest)
		// fmt.Println("Thread", request.goID, "SAVE signal received")

		request.Image.Save(request.OutFile)
		atomic.AddInt32(&ctx.fileCount, -1)
	}
	fmt.Println("Finished Writing All Images.")
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
	// nThreads = int(math.Sqrt(float64(nThreads))) // to allow for both data and functional parallelism
	ctx := NewContext(nThreads)

	for i := 0; i < nThreads; i++ {
		goID := i + 1
		ctx.wg.Add(1)
		go spawnImageProcessor(ctx, goID)
	}

	ctx.wg.Add(1)
	go spawnImageWriter(ctx)

	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		inFile, outFile, filters := parseLine(scanner.Text())

		img, _ := imgutil.Load(inFile)
		img.Threads = ctx.nThreads
		request := QueueRequest{Image: img, Filters: filters, OutFile: outFile}

		activeFilters := atomic.AddInt32(&ctx.activeFilters, 1)
		atomic.AddInt32(&ctx.fileCount, 1)

		ctx.qFilter.Push(request)
		ctx.filterCond.Signal()
		fmt.Println("Thread 0* sent FILTER signal")

		if activeFilters == 4 {
			ctx.scanCond.L.Lock()
			ctx.scanCond.Wait()
			ctx.scanCond.L.Unlock()
		}
	}
	ctx.filterCond.L.Lock()
	ctx.readComplete = true
	ctx.filterCond.Broadcast()
	ctx.filterCond.L.Unlock()

	ctx.wg.Wait()
}

func main() {
	if len(os.Args) < 3 {
		processSerial()
	} else {
		processParallel()
	}
}
