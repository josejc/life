package main

import "fmt"

// Constants in CAPITALS
const (
	M       = 10 // Rows
	N       = 10 // Columns
	T_SIMUL = 5
)

func printw(w [M][N][2]int, t int) {
	for _, m := range w {
		for j := 0; j < len(m); j++ {
			//fmt.Print("[", i, ",", j, ":", m[j][t], "]")
			fmt.Print(m[j][t])
		}
		fmt.Println()
	}
}

// Initialization figures in the world
// x,y is the corner on top left of the figure
// TODO: Test the dimension of the figure in the world and the border cases
// beware only put 'true' in the cells and not test the other
func figure(figure int, x int, y int, w *[M][N][2]int) {
	switch figure {
	case 0: // Block
		w[x][y][0] = 1
		w[x+1][y][0] = 1
		w[x][y+1][0] = 1
		w[x+1][y+1][0] = 1
	case 1: // Blinker
		w[x][y][0] = 1
		w[x][y+1][0] = 1
		w[x][y+2][0] = 1
	case 2: // Slider
		w[x][y+1][0] = 1
		w[x+1][y+2][0] = 1
		w[x+2][y][0] = 1
		w[x+2][y+1][0] = 1
		w[x+2][y+2][0] = 1
	}
}

// TODO: Read the world from file

func main() {
	var world [M][N][2]int

	figure(0, 0, 0, &world)
	figure(1, 4, 4, &world)
	for t := 0; t < T_SIMUL; t++ {
		fmt.Println("Time:", t)
		printw(world, t%2)
	}
}
