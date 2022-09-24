package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// put log file location here
const LogFile = "orders.log"

var (
	CountMap     = make(map[string]int)
	EaterFoodMap = make(map[string][]string)
	readFile     *os.File
)

// Will check if we have any duplicate eater_id and food_id, input from the log file
func FoundDuplicate(fs *bufio.Scanner) {
	defer closeFile()
	for fs.Scan() {
		//eater,food => eater
		eater := strings.Split(fs.Text(), ",")[0]
		food := strings.Split(fs.Text(), ",")[1]
		// check if we already have the eater id ...
		if _, ok := EaterFoodMap[eater]; ok {
			// check if already had the food previously ...
			for _, f := range EaterFoodMap[eater] {
				if food == f {
					fmt.Printf("Error: Same eater same food again!!!(%v -> %v)\n", eater, food)
					return
				}
			}
		}
		// Record the unique combination to EaterFoodMap
		EaterFoodMap[eater] = append(EaterFoodMap[eater], food)
	}
}

// Will track how many times a food_id is ordered ,input from the log file
func CountFoodItem(wg *sync.WaitGroup, fs *bufio.Scanner) {
	defer closeFile()
	for fs.Scan() {
		//eater,food
		food := strings.Split(fs.Text(), ",")[1]
		CountMap[food] = CountMap[food] + 1
	}
	wg.Done()

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func closeFile() {
	readFile.Close()
}

// get fs for the log file
func getFileScanner(filePath string) *bufio.Scanner {
	readFile, err := os.Open(filePath)
	checkError(err)
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	return fileScanner
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	FoundDuplicate(getFileScanner(LogFile))
	go CountFoodItem(wg, getFileScanner(LogFile))
	wg.Wait()
	fmt.Println(CountMap)
}
