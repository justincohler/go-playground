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

func spawnImageProcessor(ctx *Context) {
	defer ctx.wg.Done()
	for {
		ctx.filterCond.L.Lock()
		ctx.filterCond.Wait()
		if ctx.readComplete {
			ctx.filterCond.L.Unlock()
			break
		}
		line := ctx.qFilter.Pop()
		ctx.filterCond.L.Unlock()

		lineArgs := strings.Split(strings.Replace(line.(string), " ", "", -1), ",")

		fromFileName := lineArgs[0]
		toFileName := lineArgs[1]
		filters := lineArgs[2:]

		img, _ := imgutil.Load(fromFileName)
		img.Threads = ctx.nThreads
		var curr *imgutil.PNGImage
		curr = img
		for _, filter := range filters {
			curr = curr.ApplyFilter(filter)
		}

		element := make(map[string]*imgutil.PNGImage)
		element[toFileName] = curr

		// Tell the saver there's a file to save
		ctx.saveCond.L.Lock()
		ctx.qSave.Push(element)
		ctx.saveCond.Signal()
		ctx.saveCond.L.Unlock()

		// Tell the csv reader to send more work
		ctx.scanCond.Signal()
	}
}

func spawnImageWriter(ctx *Context, fileCount *int) {
	savedCount := 0
	for !ctx.readComplete || savedCount < *fileCount {
		ctx.saveLock.Lock()
		ctx.saveCond.Wait()
		fmt.Println("Received SAVE signal.")
		element := ctx.qSave.Pop()
		ctx.saveLock.Unlock()
		for fileName, image := range element.(map[string]*imgutil.PNGImage) {
			image.Save(fileName)
			break
		}
		savedCount++
	}
}

func main() {
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
		ctx.filterCond.L.Lock()
		ctx.qFilter.Push(scanner.Text())
		ctx.filterCond.Signal()
		ctx.filterCond.L.Unlock()

		ctx.scanCond.L.Lock()
		ctx.scanCond.Wait()
		fmt.Println("Received CONTINUE-SCAN signal.")
		ctx.scanCond.L.Unlock()
		fileCount++
	}
	ctx.filterCond.L.Lock()
	ctx.readComplete = true
	ctx.filterCond.Broadcast()
	ctx.filterCond.L.Unlock()

	ctx.wg.Wait()
}
