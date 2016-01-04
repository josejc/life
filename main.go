package main

import "fmt"

// Constants in CAPITALS
const M = 5  // Rows
const N = 10 // Columns

func printw(w [M][N][2]bool, t int) {
	for i, m := range w {
		for j := 0; j < len(m); j++ {
			fmt.Print("[", i, ",", j, ":", m[j][t], "]")
		}
		fmt.Println()
	}
}

func main() {
	var world [M][N][2]bool

	fmt.Println("World 0")
	printw(world, 0)
	fmt.Println("World 1")
	printw(world, 1)
}
