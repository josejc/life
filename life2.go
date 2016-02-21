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
	M       = 99  // Rows
	N       = 99  // Columns
	H       = 3   // History
	T_SIMUL = 100 // Max. time of simulation (end if static condition or socillation)
	zx      = 3   // Horizontal subdivisions in the world
	zy      = 3   // Vertical subdivisions in the world
)

type Point struct {
	x, y int
}

type World struct {
	Matrix [H]map[Point]int
}

// printw print world -> array of [][]cells in time t
func printw(w World, t int) {
	m := map[Point]int{} // Map is empty set
	m = w.Matrix[t]
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
func randomw(w *World) {
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
func initw(f *os.File, w *World) bool {
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
func oscilt2(w *World, t int) bool {
	oscil := true
	if t < 2 {
		return false
	}
	at := t % H       // Actual time
	pt := (t - 2) % H // Past time
	m_at := w.Matrix[at]
	m_pt := w.Matrix[pt]
	if !reflect.DeepEqual(m_at, m_pt) {
		oscil = false
	}
	return oscil
}

// nextw compute for all the world the next state of the cells
func nextw(w *World, t int) bool {
	var paux Point

	static := true
	at := t % H       // Actual time
	nt := (t + 1) % H // Next time
	m_at := w.Matrix[at]
	m_nt := map[Point]int{}
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
	for p, _ := range m_v {
		// Now calculate the zone in the world were is p
		hzone := 0
		vzone := 0
		d := M / zx
	hloop:
		for i := 0; i < zx; i++ {
			// if p is inside [id, (i+1)d -1]
			if (p.x >= (i * d)) && (p.x < (i+1)*d) {
				hzone = i
				break hloop
			}
		}
		d = N / zy
	vloop:
		for i := 0; i < zy; i++ {
			// if p is inside [id, (i+1)d -1]
			if (p.y >= (i * d)) && (p.y < (i+1)*d) {
				vzone = i
				break vloop
			}
		}
		zone := vzone*zx + hzone
		fmt.Println("Point", p, "Zone h:", hzone, "Zone v:", vzone, "Zone 2d:", zone)
		//TODO: Goroutines and send the point to the gorotuine calculate the next state in this zone
		static = nextstate(static, p, m_at, m_nt)
	}
	w.Matrix[nt] = m_nt
	return static
}

// nextstate calculate the next state and return if the systems is static
// The maps are reference types ;) don´t necessary a pointer ;)
func nextstate(static bool, p Point, m_at map[Point]int, m_nt map[Point]int) bool {
	nxt := neighbours(m_at, p.x, p.y)
	if nxt == 1 {
		m_nt[p] = 1
	}
	if static && (nxt != m_at[p]) {
		static = false
	}
	return static
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
	var world World

	start := time.Now()
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
	elapsed := time.Since(start)
	log.Printf("Total time: %s", elapsed)
}
