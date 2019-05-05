package feed

// sequentialFeed is the internal representation of a user's twitter feed (hidden from outside packages)
type sequentialFeed struct {
	start *post // a pointer to the beginning post
}

// NewSequentialFeed creates a empy user feed
func NewSequentialFeed() Feed {
	return &sequentialFeed{start: nil}
}

// Add inserts a new post to the feed. The feed is always ordered by the timestamp where
// the most recent timestamp is at the beginning of the feed followed by the second most
// recent timestamp, etc. You may need to insert a new post somewhere in the feed because
// the given timestamp may not be the most recent.
func (f *sequentialFeed) Add(body string, timestamp int64) {
	if f.start == nil {
		f.start = newPost(body, timestamp, nil)
		return
	} else if f.start.timestamp > timestamp {
		f.start = newPost(body, timestamp, f.start)
		return
	}

	parent := f.start
	for parent.next != nil {
		if parent.next.timestamp > timestamp {
			parent.next = newPost(body, timestamp, parent.next)
			return
		}
		parent = parent.next
	}
	parent.next = newPost(body, timestamp, nil)
}

// Remove deletes the post with the given timestamp. If the timestamp
// is not included in a post of the feed then the feed remains
// unchanged. Return true if the deletion was a success, otherwise return false
func (f *sequentialFeed) Remove(timestamp int64) bool {
	parent := f.start
	if parent == nil {
		return false
	} else if parent.timestamp == timestamp {
		f.start = parent.next
		return true
	}

	for parent.next != nil {
		if parent.next.timestamp == timestamp {
			parent.next = parent.next.next
			return true
		}
		parent = parent.next
	}
	return false
}

// Contains determines whether a post with the given timestamp is
// inside a feed. The function returns true if there is a post
// with the timestamp, otherwise, false.
func (f *sequentialFeed) Contains(timestamp int64) bool {
	if f.start == nil {
		return false
	}
	curr := f.start
	for curr != nil {
		if curr.timestamp == timestamp {
			return true
		}
		curr = curr.next
	}
	return false
}

// String converts a feed into a string representation so you can
// print it out. Right now this method is NOT thread safe.
func (f *sequentialFeed) String() string {
	var str string
	curr := f.start
	for curr != nil {
		unixTimeUTC := time.Unix(curr.timestamp, 0)
		unitTimeInRFC3339 := unixTimeUTC.Format(time.RFC3339)
		str += fmt.Sprintf("(%v,%v)->", curr.body, unitTimeInRFC3339)
		curr = curr.next
	}
	return str
}
