package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const usage = "Usage: generate <num_of_ops> <add_percent> <remove_percent> <contains_percent> <string_percent>" +
	"\t <num_of_ops> = the number of operations you want to generate\n" +
	"\t <add_percent> = the percentage of add operations you want to generate. Integer [0,100]\n" +
	"\t <remove_percent> = the percentage of remove operations you want to generate. Integer [0,100]\n" +
	"\t <contains_percent> = the percentage of contains operations you want to generate. Integer [0,100]\n" +
	"\t <string_percent> = the percentage of string operations you want to generate. Integer [0,100]\n" +
	"These percentages neeed to add up to 100.\n" +
	"Sample Runs:\n" +
	"\t./generate 100 50 50 0 0 > 100.txt -- Generates 100 operations with 50% being add and 50% being remove and saves it to a file\n" +
	"\t./generate 10000 25 25 25 25 > 10000.txt -- Generates 10000 operations with an equal distribution of operations and saves it to a file\n"

func genString(numOfAdds, reqIDCount int) ([]string, int) {

	var tasks []string

	for i := 0; i < numOfAdds; i++ {
		tasks = append(tasks, fmt.Sprintf("{{%v},{%v}}", "STRING", reqIDCount))
		reqIDCount++
	}
	return tasks, reqIDCount
}
func genRemove(numOfOps int, timeStamps []int64, reqIDCount int) ([]string, int) {

	var tasks []string
	var tIndex int
	var rightNow int64
	for i := 0; i < numOfOps; i++ {
		if tIndex < len(timeStamps) && i%2 == 0 {
			rightNow = timeStamps[tIndex]
			tIndex++
		} else {
			rightNow = time.Now().Unix() + int64(i)
		}
		tasks = append(tasks, fmt.Sprintf("{{%v},{%v},{%v}}", "REMOVE", reqIDCount, rightNow))
		reqIDCount++
	}
	return tasks, reqIDCount
}
func genContains(numOfOps int, timeStamps []int64, reqIDCount int) ([]string, int) {

	var tasks []string
	var tIndex int
	var rightNow int64
	for i := 0; i < numOfOps; i++ {
		if tIndex < len(timeStamps) && i%2 == 0 {
			rightNow = timeStamps[tIndex]
			tIndex++
		} else {
			rightNow = time.Now().Unix() + int64(i)
		}
		tasks = append(tasks, fmt.Sprintf("{{%v},{%v},{%v}}", "CONTAINS", reqIDCount, rightNow))
		reqIDCount++
	}
	return tasks, reqIDCount
}
func genAdd(numOfOps, reqIDCount int) ([]string, []int64, int) {

	var tasks []string
	var timestamps []int64
	for i := 0; i < numOfOps; i++ {
		rightNow := time.Now().Unix() + int64(i)
		time.Sleep(time.Millisecond)
		timestamps = append(timestamps, rightNow)
		tasks = append(tasks, fmt.Sprintf("{{%v},{%v},{%v},{%v}}", "ADD", reqIDCount, i, rightNow))
		reqIDCount++
	}
	return tasks, timestamps, reqIDCount
}
func main() {

	if len(os.Args) != 6 {
		fmt.Println(usage)
	} else {
		numOfOps, _ := strconv.Atoi(os.Args[1])
		aPer, _ := strconv.Atoi(os.Args[2])
		rPer, _ := strconv.Atoi(os.Args[3])
		cPer, _ := strconv.Atoi(os.Args[4])
		sPer, _ := strconv.Atoi(os.Args[5])
		if aPer+rPer+cPer+sPer != 100 {
			fmt.Println(usage)
			return
		}
		calculateAmount := func(per int) int {
			return int(float32(numOfOps) * (float32(per) / 100))
		}
		remaining := numOfOps - calculateAmount(aPer) - calculateAmount(rPer) - calculateAmount(sPer) - calculateAmount(cPer)
		tasks, timestamps, reqCount := genAdd(calculateAmount(aPer)+remaining, 0)
		cTasks, reqCount2 := genContains(calculateAmount(cPer), timestamps, reqCount)
		tasks = append(tasks, cTasks...)
		rTasks, reqCount3 := genRemove(calculateAmount(rPer), timestamps, reqCount2)
		tasks = append(tasks, rTasks...)
		sTasks, _ := genString(calculateAmount(sPer), reqCount3)
		tasks = append(tasks, sTasks...)
		for _, str := range tasks {
			fmt.Println(str)
		}
		fmt.Println("{DONE}")
	}

}
