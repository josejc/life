# Cells

The main function control the data of the world's state

##data

world map[point]int         // points of the cells are alife
visit map[point]int       // points for check the neighbours (the life cells and her neighbours)
c[M*N] channels of point      // a channel for every go routine N for send the point to calculate next state
life[M*N] channels of map       // return the map of points with the cell life in the next state


##pseudocode
```
main 
    initializations
    create M*N go routines (nextstate)
    execute M*N // waiting for points for calculate nextstate
    for t=0; t<T_SIMUL; t++
        for p range world
            visit=p and neighbours
        endfor
        i=0
        // every go routine is waiting for a point to calculate the next state
        for p range visit
            p -> c[i]
            i++
        endfor
        // the next world is the union of all cells alife
        for j=0 to i
            solve += <- life[j]
        endfor
        world = solve
    endfor
endmain

go routine nextstate
    // nextstate - i
    p <- c[i]
    life = neifhbours(p)
    if life
        life[i] <- p
    endif
    // if not life need a point X,X for don't block the main
endnextstate
```

PROBLEMS

    -All the goroutines read the structure data world
    -If wait for a point life and the point is dead need a point X,X for this case
    -With this structure... improve the deterministic solution? I think NO,
        the go routines working concurrently but main is send all the points and recollect
        the solution in deterministic :p


```
