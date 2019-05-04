package main

import (
	"feed"
	"fmt"
	"sync"
	"time"
)

func add(wg *sync.WaitGroup, f feed.Feed, text string, timestamp int64) {
	defer wg.Done()
	f.Add(text, timestamp)
}
func remove(wg *sync.WaitGroup, f feed.Feed, timestamp int64) {
	defer wg.Done()
	fmt.Println(f.Remove(timestamp))
}

func main() {
	var wg sync.WaitGroup

	jackDorseyFeed := feed.NewFeed()
	rightNow := time.Now().Unix()
	wg.Add(1)
	go add(&wg, jackDorseyFeed, "just setting up my twttr", rightNow)
	wg.Add(1)
	go add(&wg, jackDorseyFeed, "here's another", rightNow+10)
	wg.Add(1)
	go add(&wg, jackDorseyFeed, "should be in the middle", rightNow+5)
	wg.Wait()
	fmt.Println(jackDorseyFeed.String())
	wg.Add(1)
	go remove(&wg, jackDorseyFeed, rightNow+5)
	wg.Wait()
	fmt.Println(jackDorseyFeed.String())
}
