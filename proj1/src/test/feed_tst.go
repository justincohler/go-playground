package main

import (
	"feed"
	"fmt"
	"time"
)

func main() {
	jackDorseyFeed := feed.NewFeed()
	rightNow := time.Now().Unix()
	jackDorseyFeed.Add("just setting up my twttr", rightNow)
	jackDorseyFeed.Add("here's another", rightNow+10)

	fmt.Println(jackDorseyFeed.String())
	jackDorseyFeed.Add("should be in the middle", rightNow+5)
	fmt.Println(jackDorseyFeed.String())

	fmt.Println(jackDorseyFeed.Contains(rightNow))      // true
	fmt.Println(jackDorseyFeed.Contains(rightNow + 10)) // true
	fmt.Println(jackDorseyFeed.Contains(rightNow - 1))  // false
	fmt.Println(jackDorseyFeed.Contains(rightNow + 20)) // false
	fmt.Println(jackDorseyFeed.String())
	fmt.Println(jackDorseyFeed.Remove(rightNow + 5))    // true
	fmt.Println(jackDorseyFeed.Contains(rightNow))      // true
	fmt.Println(jackDorseyFeed.Contains(rightNow + 10)) // true
	fmt.Println(jackDorseyFeed.String())
}
