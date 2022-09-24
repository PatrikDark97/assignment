package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

type Sort struct {
	Key   string
	Value int
}
type SortList []Sort

func (p SortList) Len() int           { return len(p) }
func (p SortList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p SortList) Less(i, j int) bool { return p[i].Value > p[j].Value }

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
	fmt.Println("Top Food items: ")
	//sort the CountMap and display the top 3 for <3 show all
	CountSlice := make(SortList, len(CountMap))
	i := 0
	for k, v := range CountMap {
		CountSlice[i] = Sort{k, v}
		i++
	}
	sort.Sort(CountSlice)

	if len(CountSlice) > 0 && len(CountSlice) < 3 {
		for i, v := range CountSlice {
			fmt.Printf("%v %v\n", i+1, v.Key)
		}
	} else {
		for i, v := range CountSlice {
			if i > 2 {
				break
			}
			fmt.Printf("%v %v\n", i+1, v.Key)
		}
	}
}
