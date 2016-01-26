# Game of Life

2nd project, implement 'Game of Life' in Go

Try step by step, everyday a little...

* main.go:          Run ;)

The idea is do a first implementation with array of cells and main function for compute the next state of the cells and loop n cycles for evualuate the simulation.

Next make a second implementation using concurrency for calculate the next state every cell independent and synchronization in the main function.

World of cells:

Is array of M,N: M,N>=3 M=rows,N=cols

Beware of border cases:

* M=0 the first row, neighboring top cells will be M-1
* M=M-1 the last row, neighboring below cells will be 0
* N=0 the first column, neighboring left cells will be N-1
* N=N-1 the last column, neighboring right cells will be 0

The corner cells:

* (0,0), (0,N-1), (M-1,0), (M-1,N-1) the neighboring corner are the other corner cells

With this structure we have a spherical world ;) in 2d rectangular array of cells

For evaluation in next time(tick) the world of cells have to [2], for iterate between and grow the simulation.

--
The constant H is for history of save the different steps in the simulation (now use H=3)

When [0]=[1] the system is static and stop the simulation
When [0]=[2] the system is oscyllatos with t=2 and stop the simulation

--

Different programs for the same:

* life0: Implementation in Go with in mind C programming
* life1: Implementation with data structures of Go (Slices of cells and matrix a map of cells)
* life2: Implmentation a concurrent version with Go routines and channels
* life3: Implementation with graphics (first study another programs/games for what graphic library use)

--

## Graphic library? SDL, GL or another basic for create animation

First, install the libraries ;) in the SDL case

$sudo dnf install SDL2-devel SDL2 SDL2_mixer SDL2_ttf SDL2_image

Second, install the package "https://github.com/veandco/go-sdl2"

Now I need learn to use :p

More information:
 [Wikipedia](https://en.wikipedia.org/wiki/Conway's_Game_of_Life "https://en.wikipedia.org/wiki/Conway's_Game_of_Life")

Bonda...

For run: "$go run *.go" //With predefined world in the code

"$go build main.go"
"$./main file.LIF"  // Read .LIF file with the pattern defined in file
