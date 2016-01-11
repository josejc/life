package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Constants in CAPITALS
const (
	M       = 101 // Rows
	N       = 101 // Columns
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

func initw(f *os.File, w *[M][N][2]int) bool {
	var r *strings.Reader
	var b byte
	var x, y, oldy int

	input := bufio.NewScanner(f)
	input.Scan()
	if input.Text() != "#Life 1.05" {
		fmt.Fprintf(os.Stderr, "ERROR: The file for initialization the world is not a valid .LIF format\n")
		return false
	}
header:
	// Read header of .LIF
	for input.Scan() {
		r = strings.NewReader(input.Text())
		b, _ = r.ReadByte()
		if b != '#' {
			fmt.Println(input.Text())
		} else {
			b, _ = r.ReadByte()
			switch b {
			case 'D':
				{
					fmt.Println("Description")
				}
			case 'N':
				{
					fmt.Println("Rules Conway R 23/3")
				}
			case 'R':
				{
					fmt.Fprintf(os.Stderr, "ERROR: 'R' option not implemented\n")
					return false
				}
			case 'P':
				{
					s := strings.Split(input.Text(), " ")
					x, _ = strconv.Atoi(s[1])
					y, _ = strconv.Atoi(s[2])
					x += (M / 2)
					y += (N / 2)
					oldy = y
					break header // Exit loop, now only blocks of position and cells
				}
			default:
				{
					fmt.Fprintf(os.Stderr, "ERROR: Option in header not implemented\n")
					return false
				}
			}
		}
	}
	// Read patterns and positions
	for input.Scan() {
		r = strings.NewReader(input.Text())
		b, _ = r.ReadByte()
		if b == '#' {
			b, _ = r.ReadByte()
			if b == 'P' {
				s := strings.Split(input.Text(), " ")
				x, _ = strconv.Atoi(s[1])
				y, _ = strconv.Atoi(s[2])
				x += (M / 2)
				y += (N / 2)
				oldy = y
			} else {
				fmt.Fprintf(os.Stderr, "ERROR: Expected Position or blocks not config parameters\n")
				return false
			}
		} else {
			for cells := int(r.Size()); cells > 0; cells-- {
				if b == '.' || b == '*' {
					if b == '.' {
						w[x][y][0] = 0
					} else {
						w[x][y][0] = 1
					}
				} else {
					fmt.Fprintf(os.Stderr, "ERROR: Character not valid, only '.' or '*'\n")
					return false
				}
				b, _ = r.ReadByte()
				y++
			}
		}
		x++
		y = oldy
	}
	return true
	// NOTE: ignoring potential errors from input.Err()
}

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

	run := true
	files := os.Args[1:]
	if len(files) == 0 {
		// World initialization
		figure(0, 6, 2, &world)
		//figure(1, 4, 4, &world)
		figure(2, 0, 0, &world)
	} else {
		// Only open the first argument, one file
		arg := files[0]
		f, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			run = false
		} else {
			run = initw(f, &world)
			f.Close()
		}
	}
	if run {
		for t := 0; t < T_SIMUL; t++ {
			fmt.Println("Time:", t)
			printw(world, t%2)
			if nextw(&world, t) {
				fmt.Println("End simulation, the system is static.")
				t = T_SIMUL
			}
		}
	}
}
