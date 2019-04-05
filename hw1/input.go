//
// This program is an example of retrieving input from the user
//
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {

	//A Scanner value is needed to read input from the console
	scanner := bufio.NewScanner(os.Stdin)

	//Prompt to the user to enter in
	fmt.Println("What is your name?")

	//Scan() reads a line and removes the newline character from the end.
	//Scan() returns "true" if there is a line and "false" when there's no
	// more input. The result text can be retrieved by calling the
	//Text() method.
	if scanner.Scan() {
		name := scanner.Text()
		fmt.Println("Your name is ", name)
		fmt.Println("What is your age?")
		if scanner.Scan() {
			//Converts Atoi a string to int
			age, _ := strconv.Atoi(scanner.Text())
			fmt.Println("Your Age is", age)
		} else {
			fmt.Println("Nothing was scanned")
		}

	} else {
		fmt.Println("Nothing was scanned")
	}
}
