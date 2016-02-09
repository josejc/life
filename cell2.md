# Cells

The main function control the data of the world's state

                            // World divide in M zones and concurrently calculate the next state of this zone
##data
world map[point]int         // points of the cells are alife
visit map[point]int         // point for calculate the next state (cells life and her neighbours)
c[M] channels of point      // a channel for every go routine for send the point to calculate next state
sol[M] channels of map      // the solution of the zone (map[point]int)

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

PROBLEMS

    -All the goroutines read the structure data world
    -With this solution if all the cells are in one zone don't improve the time of execution, the best improve is 
    when the cells are uniform distribution in the world
    -How match the borderlines in the different zones of the world

```
