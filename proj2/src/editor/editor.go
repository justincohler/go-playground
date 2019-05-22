package main

import (
	"bufio"
	"imgutil"
	"os"
	"queue"
	"strconv"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	sThreads, filePath := os.Args[1], os.Args[2]
	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	inQueue := queue.NewUnbounded()
	// outQueue := queue.NewUnbounded()

	for scanner.Scan() {
		inQueue.Push(scanner.Text())
	}

	nThreads, _ := strconv.Atoi(sThreads)
	for i := 0; i < nThreads; i++ {
		wg.Add(1)
		go func(q queue.StringStack, nThreads int) {
			defer wg.Done()
			for !q.Empty() {
				line := q.Pop()
				lineArgs := strings.Split(strings.Replace(line, " ", "", -1), ",")

				fromFileName := lineArgs[0]
				toFileName := lineArgs[1]
				filters := lineArgs[2:]

				img, _ := imgutil.Load(fromFileName)
				img.Threads = nThreads
				var curr *imgutil.PNGImage
				curr = img
				for _, filter := range filters {
					curr = curr.ApplyFilter(filter)
				}
				curr.Save(toFileName)
			}
		}(inQueue, nThreads)
	}
	wg.Wait()
}
