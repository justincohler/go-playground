package main

import (
	"fmt"
	"mpcs52060/justincohler/proj1/src/feed"
	"time"
)

func main() {
	jackDorseyFeed := feed.NewFeed()
	rightNow := time.Now().Unix()
	jackDorseyFeed.Add("just setting up my twttr", rightNow)
	jackDorseyFeed.Add("here's another", rightNow+10)
	fmt.Println(jackDorseyFeed.Contains(rightNow))

	fmt.Println(jackDorseyFeed.Remove(rightNow + 5))
	fmt.Println(jackDorseyFeed)
}
