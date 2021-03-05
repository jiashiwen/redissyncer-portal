package commons

import (
	"fmt"
	"testing"
)

func TestSortMapByValue(t *testing.T) {
	testmap := make(map[string]int64)

	testmap["a"] = 23
	testmap["b"] = 44
	testmap["c"] = 5
	testmap["d"] = 21
	testmap["e"] = 4
	testmap["f"] = 56

	parelist := SortMapByValue(testmap, false)
	parelist2 := SortMapByValue(testmap, true)

	fmt.Println(parelist)
	fmt.Println(parelist2)

}
