package feed

import (
	"strconv"
)

//Feed represents a user's twitter feed
// You will add to this interface the implementations as you complete them.
type Feed interface {
	CheckOrder([]int64) (string, bool) // Do not remove this function!!! Needed by testing!
	Add(body string, timestamp int64)
	Remove(timestamp int64) bool
	Contains(timestamp int64) bool
}

//feed is the internal representation of a user's twitter feed (hidden from outside packages)
// You CAN add to this structure but you cannot remove any of the original fields. You must use
// the original fields in your implementation. You can assume the feed will not have duplicate posts
type feed struct {
	start *post // a pointer to the beginning post

}

//post is the internal representation of a post on a user's twitter feed (hidden from outside packages)
// You CAN add to this structure but you cannot remove any of the original fields. You must use
// the original fields in your implementation.
type post struct {
	body      string // the text of the post
	timestamp int64  // Unix timestamp of the post
	next      *post  // the next post in the feed
}

//NewPost creates and returns a new post value given its body and timestamp
func newPost(body string, timestamp int64, next *post) *post {
	return &post{body, timestamp, next}
}

//NewFeed creates a empy user feed
func NewFeed() Feed {
	return &feed{start: nil}
}

// Add inserts a new post to the feed. The feed is always ordered by the timestamp where
// the most recent timestamp is at the beginning of the feed followed by the second most
// recent timestamp, etc. You may need to insert a new post somewhere in the feed because
// the given timestamp may not be the most recent.
func (f *feed) Add(body string, timestamp int64) {
	if f.start == nil {
		f.start = newPost(body, timestamp, nil)
		return
	}
	iter := f.start
	for iter.next != nil {
		iter = iter.next
	}
	iter.next = newPost(body, timestamp, nil)
}

// Remove deletes the post with the given timestamp. If the timestamp
// is not included in a post of the feed then the feed remains
// unchanged. Return true if the deletion was a success, otherwise return false
func (f *feed) Remove(timestamp int64) bool {

	parent := f.start
	if parent == nil {
		return false
	} else if parent.timestamp == timestamp {
		f.start = parent.next
		return true
	}

	for child := parent.next; parent.next != nil; parent = child {
		if child.timestamp == timestamp {
			parent.next = child.next
			return true
		}
	}
	return false
}

// Contains determines whether a post with the given timestamp is
// inside a feed. The function returns true if there is a post
// with the timestamp, otherwise, false.
func (f *feed) Contains(timestamp int64) bool {
	if f.start == nil {
		return false
	}
	iter := f.start
	for iter.next != nil {
		if post := *iter; post.timestamp == timestamp {
			return true
		}
		iter = iter.next
	}
	return false
}

// String converts a feed into a string representation so you can
// print it out. Right now this method is NOT thread safe.
func (f *feed) String() string {
	str := ""
	for iter := f.start; iter != nil; iter = iter.next {
		str += (*iter).body + "\n"
	}
	return str
}

// DO NOT MODIFY THIS FUNCTION OR REMOVE IT! This is needed for the testing module.
func (f *feed) CheckOrder(order []int64) (string, bool) {

	var comma = ","
	var result = "["
	currNode := f.start
	validOrder := true
	i := 0
	for currNode != nil {
		if currNode.next == nil {
			comma = ""
		}
		if currNode.timestamp != order[i] {
			validOrder = false
		}

		result += (strconv.FormatInt(currNode.timestamp, 10) + comma)
		currNode = currNode.next
		i++
	}
	result = "]"
	if i != len(order) {
		validOrder = false
	}
	return result, validOrder
}
