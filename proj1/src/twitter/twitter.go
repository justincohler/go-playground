package main

import (
	"bufio"
	"feed"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	cADD      = "ADD"
	cREMOVE   = "REMOVE"
	cCONTAINS = "CONTAINS"
	cSTRING   = "STRING"
)

// Command is a generic Interface for Feed Types.
type Command interface {
}

// Add puts the given body onto the feed.
type Add struct {
	reqID     int64
	body      string
	timestamp int64
}

// Remove removes a given reqID from the feed.
type Remove struct {
	reqID     int64
	timestamp int64
}

// Contains returns the existence of a given reqID on the feed.
type Contains struct {
	reqID     int64
	timestamp int64
}

// String prints the details of a given reqID on the feed.
type String struct {
	reqID int64
}

func parseLine(line string) (string, int64, string, int64) {
	line = strings.Replace(line, "{", "", -1)
	line = strings.Replace(line, "}", "", -1)
	args := strings.Split(line, ",")
	commandName := args[0]
	reqID, _ := strconv.ParseInt(args[1], 10, 64)

	var timestamp int64
	var body string
	if len(args) == 3 {
		timestamp, _ = strconv.ParseInt(args[2], 10, 64)
	} else if len(args) == 4 {
		body = args[2]
		timestamp, _ = strconv.ParseInt(args[4], 10, 64)
	}

	return commandName, reqID, body, timestamp
}

func executeLine(feed *feed.Feed, line string) {
	commandName, reqID, body, timestamp := parseLine(line)

	var res bool
	var status string

	switch commandName {
	case cADD:
		(*feed).Lock()
		defer feed.Unlock()
		feed.Add(body, timestamp)
		fmt.Println("{{", reqID, "}, {SUCCESS}}")
	case cREMOVE:
		feed.Lock()
		defer feed.Unlock()

		if feed.Remove(timestamp) {
			status = "SUCCESS"
		} else {
			status = "FAILED"
		}
		fmt.Println("{{", reqID, "}, {", status, "}}")

	case cCONTAINS:
		feed.RLock()
		defer feed.RUnlock()

		if feed.Contains(timestamp) {
			status = "YES"
		} else {
			status = "NO"
		}
		fmt.Println("{{", reqID, "}, {", status, "}}")

	case cSTRING:
		feed.RLock()
		defer feed.RUnlock()
		fmt.Println("{{", reqID, "}, {", feed.String(), "}}")
	}
}

func main() {
	var feed feed.Feed

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		go executeLine(feed, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
