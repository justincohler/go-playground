package feed

//Feed represents a user's twitter feed
// You will add to this interface the implementations as you complete them.
type Feed interface {
	CheckOrder([]int64) (string, bool) // Do not remove this function!!! Needed by testing!
	Add(body string, timestamp int64)
	Remove(timestamp int64) bool
	Contains(timestamp int64) bool
	String() string
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
