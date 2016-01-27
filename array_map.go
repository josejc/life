package main

import "fmt"

func main() {
	//m := make([]map[string]int, 2) // Slice (len=2) of map
	var m [2]map[string]int // Array [2] of map

	//n := make(map[string]int) // Equal to create map in the next line
	n := map[string]int{} // Map is empty set

	// the code section ;)
	n["a"] = 1
	n["b"] = 2
	m[0] = n
	n = make(map[string]int)
	n["A"] = 100
	m[1] = n
	fmt.Println(m)
}
