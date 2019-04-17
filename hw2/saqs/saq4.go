package main

import (
	"fmt"
	"sync"
)

// printIDAndShared function prints the gorountine's id, increments
// a shared variables and prints it.
func printIDAndShared(goID int,
	sharedVarMutex *sync.Mutex,
	wg *sync.WaitGroup,
	sharedVar *int) {

	// FILO deferment will execute this last
	defer wg.Done()

	sharedVarMutex.Lock()
	defer sharedVarMutex.Unlock()

	if *sharedVar == goID {
		*sharedVar++
		fmt.Printf("Gorountine id = %v, sharedVar = %v\n", goID, *sharedVar)
	} else {
		wg.Add(1)
		go printIDAndShared(goID, sharedVarMutex, wg, sharedVar)
	}
}

func main() {

	//We will spawn 10 threads (i.e., 10 goroutines) in this program
	numOfThreads := 10

	//This variable is a shared variable between the threads
	//Each thread will increment this resource in a synchornized fashion.
	//The last thread incrementing it being the main thread.
	var sharedVar int

	//A mutex that is needed to allow sequential access to the shared variable
	//Look at the Lecture 2 slides that explain the need for the mutex
	// Note that when loc
	var mux sync.Mutex

	/** From the Go Documentation:

	A WaitGroup waits for a collection of goroutines to finish.
	The main goroutine calls Add to set the number of goroutines to wait for.
	Then each of the goroutines runs and calls Done when finished.
	At the same time, Wait can be used to block until all goroutines
	have finished.

	https://golang.org/pkg/sync/#WaitGroup

	**/
	var wg sync.WaitGroup

	// Spawn a goroutine (aka thread). It's sometimes useful
	// to give a thread a unique identifier so you can use this
	// value id to determine how to decompose work to a thread.
	for goID := 0; goID < numOfThreads; goID++ {

		// Increment the WaitGroup counter.
		wg.Add(1)

		/****** IMPORTANT ******
		Why am I passing pointers (*sync.Mutex, *sync.WaitGroup, etc.) into this function?
		Go by default passes arguments by value (i.e. it's going to be a copy) of
		the value. Thus, if I didn't pass a pointer to these values then
		each thread will have their own duplicate copy of them, which defeats
		the point of having a .

		Be mindful of this on in your homework and future assignments!
		***/
		go printIDAndShared(goID, &mux, &wg, &sharedVar)
	}
	//The main goroutine needs to wait for all the goroutines to end
	// before exiting the program.
	wg.Wait()

	//Increment the shared variable one last time
	sharedVar++

	//Print the final variable value and quit
	fmt.Printf("Main Goroutine, sharedVar = %v\nDone\n", sharedVar)
}
