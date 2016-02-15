# Cells

The main function control the data of the world's state

##Data

X is the numbre of goroutines for concurrently chech if the points are alife in the next state

    * world map[point]int       ;points of the cells are alife
    * visit map[point]int       ;point for calculate the next state (cells life and her neighbours)
    * c[X] channels of point  ;a channel for every go routine for send the point to calculate next state
    * sol[X] channels of map  ;the point life in the next state calculates for the gorotuine 'x' (map[point]int)

##Pseudocode
```
main 
    initializations
    create M go routines (nextstate)
    execute M // waiting for points for calculate nextstate
    for t=0; t<T_SIMUL; t++
        for p range world
            visit=p and neighbours
        endfor
        // every go routine is waiting for a point to calculate the next state
        i=0
        for p range visit
            p -> c[i]
            i = (i+1)%c
        endfor
        // send special point for finish and return the solve
        for i=0 to X
            specialp -> c[i]
            // the next world is the union of all cells alife in the different zones
            solve += <- sol[i]
        endfor
        world = solve
    endfor
endmain

go routine nextstate
    // nextstate - x
    loop
        solzone = empty
        for p <- c[x]; p != specialp; p<-c[x]
            life = neifhbours(p)        // Only read in data structure world, and change in time
            if life
                solzone[p] = life
            endif
            // if not life the point is not added
        endfor
        sol[x] <- solzone
    endloop
endnextstate
```

##Problems

    * All the goroutines read the structure data world
        No problem, because only read, the main function control with channels when goroutines running
        Nothing to do
    * How match the borderlines in the different zones of the world
        No problem, the main function divide the world in zones but the goroutines read all the world, no borderlines
    * Be careful, to selection X
        if X too large, goroutines do nothing
        if X is small, doesn't profit of concurrency

