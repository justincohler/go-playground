package feed

import (
	"lock"
)

// parallelFeed is the internal representation of a user's twitter feed (hidden from outside packages)
type parallelFeed struct {
	lock.RWLock
	sf *serialFeed
}

// NewParallelFeed creates a empy user feed
func NewParallelFeed() Feed {
	sf := &serialFeed{start: nil}
	return &parallelFeed{sf: sf}
}

// Add inserts a new post to the feed. The feed is always ordered by the timestamp where
// the most recent timestamp is at the beginning of the feed followed by the second most
// recent timestamp, etc. You may need to insert a new post somewhere in the feed because
// the given timestamp may not be the most recent.
func (f *parallelFeed) Add(body string, timestamp int64) {
	f.Lock()
	defer f.Unlock()

	f.sf.Add(body, timestamp)
}

// Remove deletes the post with the given timestamp. If the timestamp
// is not included in a post of the feed then the feed remains
// unchanged. Return true if the deletion was a success, otherwise return false
func (f *parallelFeed) Remove(timestamp int64) bool {
	f.Lock()
	defer f.Unlock()

	return f.sf.Remove(timestamp)
}

// Contains determines whether a post with the given timestamp is
// inside a feed. The function returns true if there is a post
// with the timestamp, otherwise, false.
func (f *parallelFeed) Contains(timestamp int64) bool {
	f.RLock()
	defer f.RUnlock()

	return f.sf.Contains(timestamp)
}

// String converts a feed into a string representation so you can
// print it out. Right now this method is NOT thread safe.
func (f *parallelFeed) String() string {
	f.RLock()
	defer f.RUnlock()

	return f.sf.String()
}

// DO NOT MODIFY THIS FUNCTION OR REMOVE IT! This is needed for the testing module.
func (f *parallelFeed) CheckOrder(order []int64) (string, bool) {
	return f.sf.CheckOrder(order)
}
