package main

import "fmt"

// Constants in CAPITALS
const (
	M       = 10 // Rows
	N       = 10 // Columns
	T_SIMUL = 100
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

// Compute for all the world the next state of the cells
func nextw(w *[M][N][2]int, t int) bool {

	static := true
	at := t % 2       // Actual time
	nt := (t + 1) % 2 // Next time
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			w[i][j][nt] = neighbours(w, i, j, at)
			if static && (w[i][j][nt] != w[i][j][at]) {
				static = false
			}
		}
	}
	return static
}

// Compute the next state interacts with is neighbours
func neighbours(w *[M][N][2]int, i int, j int, t int) int {
	var nb int // number of neifhbours life

	top := i - 1
	bottom := i + 1
	left := j - 1
	right := j + 1
	if top == -1 {
		top = M - 1
	}
	if bottom == M {
		bottom = 0
	}
	if left == -1 {
		left = N - 1
	}
	if right == N {
		right = 0
	}
	nb = w[top][left][t] + w[top][j][t] + w[top][right][t]
	nb += w[i][left][t] + w[i][right][t]
	nb += w[bottom][left][t] + w[bottom][j][t] + w[bottom][right][t]
	//fmt.Println(top, bottom, left, right)
	//fmt.Println(i, j, t, "nb:", nb, "w:", w[i][j][t])
	if (nb == 2) && (w[i][j][t] == 1) {
		return 1
	}
	if nb == 3 {
		return 1
	}
	return 0
}

// TODO: End simulation oscillators period=2?
func main() {
	var world [M][N][2]int

	figure(0, 6, 2, &world)
	//figure(1, 4, 4, &world)
	figure(2, 0, 0, &world)
	for t := 0; t < T_SIMUL; t++ {
		fmt.Println("Time:", t)
		printw(world, t%2)
		if nextw(&world, t) {
			fmt.Println("End simulation, the system is static.")
			t = T_SIMUL
		}
	}
}
