// Package MAIN for proves the functions to include in a next package
// for implementation Cellular Automatas or Conway's Life
// This version is sequential without the characteristics of concurrency
package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// Constants in CAPITALS
const (
	M       = 101 // Rows
	N       = 101 // Columns
	H       = 3   // History
	T_SIMUL = 100
)

type Point struct {
	x, y int
}

type World struct {
	Cells  [][H]Point
	Matrix [H]map[Point]int
}

// printw print world -> array of [][]cells in time t
func printw(w World, t int) {
	fmt.Println(w.Cells)
}

// randomw generate a random initial state
func randomw(w *World) {
	var p Point
	// A 25% of cells are alive
	for i := 0; i < (M * N / 4); i++ {
		p.x = rand.Intn(M)
		p.y = rand.Intn(N)
		w.Matrix[p][0] = 1
	}
}

/*---
// initw read a file .LIF and configure the world with it
func initw(f *os.File, w *[M][N][H]int) bool {
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
				switch b {
				case '.':
					{
						w[x][y][0] = 0
					}
				case '*':
					{
						w[x][y][0] = 1
					}
				default:
					{
						fmt.Fprintf(os.Stderr, "ERROR: Character not valid, only '.' or '*'\n")
						return false
					}
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

// oscilt2 compare Actual (t) = Past (t - 2) for know if the system is oscillator
func oscilt2(w *[M][N][H]int, t int) bool {
	oscil := true
	if t < 2 {
		return false
	}
	at := t % H       // Actual time
	pt := (t - 2) % H // Past time
loop:
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			if w[i][j][at] != w[i][j][pt] {
				oscil = false
				break loop
			}
		}
	}
	return oscil
}

// nextw compute for all the world the next state of the cells
func nextw(w *[M][N][H]int, t int) bool {

	static := true
	at := t % H       // Actual time
	nt := (t + 1) % H // Next time
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

// neighbours calculate the next state of the cells
func neighbours(w *[M][N][H]int, i int, j int, t int) int {
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
---*/

// main function for run and test the implementation of the functions
func main() {
	var world World

	/*
		run := true
		randomPtr := flag.Bool("random", false, "Initialize the world with a random state")
		filePtr := flag.String("file", "name.lif", "File name .lif")
		flag.Parse()
		switch {
		case *randomPtr:
			{
				randomw(&world)
			}
		case *filePtr != "name.lif":
			{
				f, err := os.Open(*filePtr)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
					run = false
				} else {
					run = initw(f, &world)
					f.Close()
				}
			}
		default:
			{
				fmt.Println("Use $life -h for help")
				run = false
			}
		}
		if run {
			for t := 0; t < T_SIMUL; t++ {
				fmt.Println("Time:", t)
				printw(world, t%H)
				// Check the world before calculate the next
				if oscilt2(&world, t) {
					fmt.Println("End simulation, the system is oscillator with period=2")
					t = T_SIMUL
				}
				if nextw(&world, t) {
					fmt.Println("End simulation, the system is static.")
					t = T_SIMUL
				}
			}
		}
	*/

	randomw(&world)
	printw(world, 0)

}
