# Cells

Next make a second implementation using concurrency for calculate the next state every cell independent and synchronization in the main function.

Ideas for a implementation in a package 'cell' the functions of cell and run concurrently

##Data of Cell
* 8 channels for communication with neighbours (send/receive state)
* position (x,y) for a representation text/graphic
* state (life/dead)
* 1 channel for synchronization (receive ticks) 

##Neighbours:
* Top (left, middle, right)
* Middle (left, right)
* Bottom (left, middle, right)

##pseudocode
```
main 
    initializations
    for t=0; t<T_SIMUL; t++
        send tick (sync channel shared?)
        send print (order channel? shared?)
        test static world?
        test oscillator with period=2?
    endfor
endmain

cell_loop
    wait tick
    for neighbours
        // Deadlocks?
        state_n += <-channel neighbour
        state_cell -> channel neighbour
    endfor
    if state_n = 3
        state_cell = 1
    endif
    if state_n = 2 and...
        state_cell = 1
    endif
    state_cell = 0
    wait order? the same channel for ticks and orders?
endcell_loop
```

Now I need know how channels work ;)
