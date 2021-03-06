// Package MAIN for proves the functions to include in a next package
// for implementation Cellular Automatas or Conway's Life
// This version is sequential without the characteristics of concurrency
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Constants in CAPITALS
const (
	M       = 100 // Rows
	N       = 100 // Columns
	H       = 3   // History
	T_SIMUL = 100 // Max. time of simulation (end if static condition or socillation)
	X       = 20  // Number of goroutines to create, and wait for calculate the next state of points
)

type Point struct {
	x, y int
}

type World struct {
	Matrix [H]map[Point]int
	static bool // World is static?
	T      int  // Actual time
	X      int  // Number of Goroutines
}

var w World
var punts [X]chan<- Point
var sols [X]<-chan map[Point]int

// printw print world -> array of [][]cells in time t
func printw() {
	m := map[Point]int{} // Map is empty set
	m = w.Matrix[w.T]
	fmt.Println("---")
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			if m[Point{i, j}] == 0 {
				fmt.Print(".")
			} else {
				fmt.Print("*")
			}
		}
		fmt.Println()
	}
}

// randomw generate a random initial state
func randomw() {
	var p Point
	m := map[Point]int{}
	// A 25% of cells are alive
	for i := 0; i < (M * N / 4); i++ {
		p.x = rand.Intn(M)
		p.y = rand.Intn(N)
		if m[p] == 0 {
			m[p] = 1
		}
	}
	w.Matrix[0] = m
}

// initw read a file .LIF and configure the world with it
func initw(f *os.File) bool {
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
	var p Point
	m := map[Point]int{}
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
			p.x = x
			for cells := int(r.Size()); cells > 0; cells-- {
				p.y = y
				switch b {
				case '.':
					{
						//m[p] = 0
					}
				case '*':
					{
						m[p] = 1
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
	w.Matrix[0] = m
	return true
	// NOTE: ignoring potential errors from input.Err()
}

// oscilt2 compare Actual (t) = Past (t - 2) for know if the system is oscillator
func oscilt2(t int) bool {
	oscil := true
	if t < 2 {
		return false
	}
	pt := (t - 2) % H // Past time
	m_at := w.Matrix[w.T]
	m_pt := w.Matrix[pt]
	if !reflect.DeepEqual(m_at, m_pt) {
		oscil = false
	}
	return oscil
}

//
func nextConcurrently() (chan<- Point, <-chan map[Point]int) {
	c_punts := make(chan Point)       // Can only read from
	c_sol := make(chan map[Point]int) // Can only write to
	go func() {                       // We launch the goroutine from inside the function.
		var p Point
		p_end := Point{M, N}
		for {
			for p = <-c_punts; p != p_end; p = <-c_punts {
				// wait to change the actual time ;)
			}
			p_sol := map[Point]int{} // Map points alife
			m_at := w.Matrix[w.T]
			for p = <-c_punts; p != p_end; p = <-c_punts {
				nxt := neighbours(m_at, p.x, p.y)
				if nxt == 1 {
					p_sol[p] = 1
				}
				if w.static && (nxt != m_at[p]) {
					w.static = false
				}
			}
			c_sol <- p_sol
		}
	}()
	return c_punts, c_sol
}

// nextw compute for all the world the next state of the cells
func nextw() {
	var paux Point

	w.static = true
	at := w.T          // Actual time
	nt := (at + 1) % H // Next time
	m_at := w.Matrix[at]
	m_nt := map[Point]int{}
	m_s := map[Point]int{} // Partial solution offer by goroutine
	m_v := map[Point]int{} // Map points for calculate the next state
	for p, v := range m_at {
		if v == 1 { // Only add the points for visit the cells alife and her neighbours
			top, bottom, left, right := pneighbours(p.x, p.y)
			paux = Point{top, left}
			m_v[paux] = 1
			paux = Point{top, p.y}
			m_v[paux] = 1
			paux = Point{top, right}
			m_v[paux] = 1
			paux = Point{p.x, left}
			m_v[paux] = 1
			paux = p
			m_v[paux] = 1
			paux = Point{p.x, right}
			m_v[paux] = 1
			paux = Point{bottom, left}
			m_v[paux] = 1
			paux = Point{bottom, p.y}
			m_v[paux] = 1
			paux = Point{bottom, right}
			m_v[paux] = 1
		}
	}
	i := 0
	count := 0
	for p, _ := range m_v {
		if count < w.X {
			punts[i] <- Point{M, N}
		}
		punts[i] <- p
		i = (i + 1) % w.X
		count++
	}
	for i = count; i < w.X; i++ {
		punts[i] <- Point{M, N}
	}
	for i = 0; i < w.X; i++ {
		punts[i] <- Point{M, N}
		m_s = <-sols[i]
		for k, v := range m_s {
			m_nt[k] = v
		}
	}
	w.Matrix[nt] = m_nt
}

// pneighbours return the positions of the neighbours of (i,j)
func pneighbours(i int, j int) (int, int, int, int) {
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
	return top, bottom, left, right
}

// neighbours calculate the next state of the cells
func neighbours(m map[Point]int, i int, j int) int {
	var nb int // number of neifhbours life

	top, bottom, left, right := pneighbours(i, j)
	nb = m[Point{top, left}] + m[Point{top, j}] + m[Point{top, right}]
	nb += m[Point{i, left}] + m[Point{i, right}]
	nb += m[Point{bottom, left}] + m[Point{bottom, j}] + m[Point{bottom, right}]
	if (nb == 2) && (m[Point{i, j}] == 1) {
		return 1
	}
	if nb == 3 {
		return 1
	}
	return 0
}

// main function for run and test the implementation of the functions
func main() {
	start := time.Now()
	run := true
	w.X = X
	randomPtr := flag.Bool("random", false, "Initialize the world with a random state")
	filePtr := flag.String("file", "name.lif", "File name .lif")
	nGRPtr := flag.Int("x", X, "Number of goroutines for calculate next world")
	flag.Parse()
	switch {
	case *randomPtr:
		{
			randomw()
		}
	case *filePtr != "name.lif":
		{
			f, err := os.Open(*filePtr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				run = false
			} else {
				run = initw(f)
				f.Close()
			}
		}
	default:
		{
			fmt.Println("Use $life -h for help")
			run = false
		}
	}
	if *nGRPtr > X {
		fmt.Fprintf(os.Stderr, "ERROR: The higher number of GoRoutines is %v\n", X)
		run = false
	} else if *nGRPtr < 1 {
		fmt.Fprintf(os.Stderr, "Error: The minimal number of GoRoutines is 1\n")
		run = false
	}
	if run {
		w.X = *nGRPtr
		fmt.Println("---X", w.X)
		for i := 0; i < w.X; i++ {
			punts[i], sols[i] = nextConcurrently()
		}
		for t := 0; t < T_SIMUL; t++ {
			fmt.Println("Time:", t)
			w.T = t % H
			printw()
			// Check the world before calculate the next
			if oscilt2(t) {
				fmt.Println("End simulation, the system is oscillator with period=2")
				t = T_SIMUL
			}
			nextw()
			if w.static {
				fmt.Println("End simulation, the system is static.")
				t = T_SIMUL
			}
		}
	}
	elapsed := time.Since(start)
	log.Printf("Total time: %s", elapsed)
}
