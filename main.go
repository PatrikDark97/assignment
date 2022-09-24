package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const LogFile = "orders.log"

var CountMap = make(map[string]int)

func CountFoodItem(fs *bufio.Scanner) {
	for fs.Scan() {
		food := strings.Split(fs.Text(), ",")[1]
		CountMap[food] = CountMap[food] + 1
	}
}

func main() {
	readFile, err := os.Open(LogFile)
	defer readFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	CountFoodItem(fileScanner)

	fmt.Println(CountMap)
}
