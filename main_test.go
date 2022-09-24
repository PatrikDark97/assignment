package main

import (
	"reflect"
	"sync"
	"testing"
)

func TestCountFoodItem(t *testing.T) {
	WantCountMap := make(map[string]int)
	WantCountMap["Food_1"] = 6
	WantCountMap["Food_10"] = 1
	WantCountMap["Food_1000"] = 4
	WantCountMap["Food_2"] = 3

	wg := new(sync.WaitGroup)
	wg.Add(1)
	CountFoodItem(wg, getFileScanner("ordersTest.log"))
	if !reflect.DeepEqual(WantCountMap, CountMap) {
		t.Errorf("Not equal")
	}
}

// this one is self documenting test
func ExampleFoundDuplicate() {
	FoundDuplicate(getFileScanner("ordersTest.log"))
	//Output:
	//Error: Same eater same food again!!!(Cus_101 -> Food_1)
}
