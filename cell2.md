# Cells

The main function control the data of the world's state

##data

The world is divided in (x-subdivision horizontal, y-subdivision vertical)
    * world map[point]int       ;points of the cells are alife
    * visit map[point]int       ;point for calculate the next state (cells life and her neighbours)
    * c[x*y] channels of point  ;a channel for every go routine for send the point to calculate next state
    * sol[x*y] channels of map  ;the solution of the zone (map[point]int)

##pseudocode
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
        for p range visit
            zone = zoneof(p)
            p -> c[zone]
        endfor
        // send special point for finish and return the solve
        for i to M
            specialp -> c[i]
            // the next world is the union of all cells alife in the different zones
            solve += <- life[j]
        endfor
        world = solve
    endfor
endmain

go routine nextstate
    // nextstate - zone
    loop
        for p <- c[zone]; p != specialp; p<-c[zone]
            life = neifhbours(p)        // Only read in data structure world, and change in time
            if life
                solzone[p] = life
            endif
            // if not life the point is not added
        endfor
        sol[M] <- solzone
    endloop
endnextstate

//The world divides x in horizontal and y in vertial, know the zone and the channel of goroutine for control this zone
function zoneworld(point) zone
    d=M div x   // M the horizontal size of the world, x the subdivision
    hzone=-1
    for i=0; i=x-1; i++
        if p.x pertain [i*d,((i+1)*d)-1]
            hzone=i
            exit for
        endif
    endfor
    if hzone=-1
        hzone=x
    endif
    d=N div y   // N the vertical size of the world, y the subdivisions
    vzone=-1
    for i=0; i=y-1; i++
        if p.y pertain [i*d,((i+1)*d)-1]
            vzone=i
            exit for
        endif
    endfor
    if vzone=-1
        vzone=y
    endif
    // 2d -> 1d
    return (vzone*x+hzone)
```

##Problems

    * All the goroutines read the structure data world
        No problem, because only read, the main function control with channels when goroutines running
    * With this solution if all the cells are in one zone don't improve the time of execution, the best improve is 
        when the cells are uniform distribution in the world
        Nothing to do
    * How match the borderlines in the different zones of the world
        No problem, the main function divide the world in zones but the goroutines read all the world, no borderlines

