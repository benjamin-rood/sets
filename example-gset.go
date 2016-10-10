package main

import (
	"fmt"
	"github.com/benjamin-rood/sets/gset"
)

func main() {
	set := gset.New()
	integer := 34
	str := "apple"
	type point struct {
		x, y int
	}
	p := point{3, -1}
	fmt.Println(set.Insert(integer)) // true
	fmt.Println(set.Insert(34))      // false
	fmt.Println(set.Insert("apple")) // true
	fmt.Println(set.Cardinality())
	fmt.Println(set.Insert(str)) // false
	fmt.Println(set.Cardinality())
	fmt.Println(set.Insert(p)) // true
	fmt.Println(set.Insert(p)) // false
	p = point{0, 10}           //
	fmt.Println(set.Insert(p)) // true
	fmt.Println(set.Cardinality())
	fmt.Println(gset.In(p, set)) // true
}
