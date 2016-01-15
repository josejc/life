# Examples of code

## Multiple goroutines listening on one chanell

For sincronization and read the ticks sending the main function

A couple of rules of thumb that should make things feel much more straightforward.

* **prefer using formal arguments for the channels** you pass to go-routines instead of accessing channels in global scope. You can get more compiler checking this way, and better modularity too.
* **avoid both reading and writing on the same channel in a particular go-routine** (including the 'main' one). Otherwise, deadlock is a much greater risk.

 Here's an alternative version of your program, applying these two guidelines. This case demonstrates many writers & one reader on a channel:

```
c := make(chan string)

for i := 1; i <= 5; i++ {
   go func(i int, co chan<- string) {
      for j := 1; j <= 5; j++ {
         co <- fmt.Sprintf("hi from %d.%d", i, j)
      }
   }(i, c)
}
for i := 1; i <= 25; i++ {
   fmt.Println(<-c)
}
```

It creates the five go-routines writing to a single channel, each one writing five times. The main go-routine reads all twenty five messages - you may notice that the order they appear in is often not sequential (i.e. the concurrency is evident).

This example demonstrates a feature of Go channels: it is possible to have multiple writers sharing one channel; Go will interleave the messages automatically.

The same applies for one writer and multiple readers on one channel, as seen in the second example here:

```
c := make(chan int)
var w sync.WaitGroup
w.Add(5)

for i := 1; i <= 5; i++ {
   go func(i int, ci <-chan int) {
      j := 1
      for v := range ci {
         time.Sleep(time.Millisecond)
         fmt.Printf("%d.%d got %d\n", i, j, v)
         j += 1
      }
      w.Done()
   }(i, c)
}

for i := 1; i <= 25; i++ {
   c <- i
}
close(c)
w.Wait()
```

This second example includes a wait imposed on the main goroutine, which would otherwise exit promptly and cause the other five goroutines to be terminated early (thanks to olov for this correction).

In both examples, no buffering was needed. It is generally a good principle to view buffering as a performance enhancer only. If your program does not deadlock without buffers, it won't deadlock with buffers either (but the converse is not always true). So, as **another rule of thumb, start without buffering then add it later as needed**.

 [Stackoverflow](http://stackoverflow.com/questions/15715605/multiple-goroutines-listening-on-one-channel "http://stackoverflow.com/questions/15715605/multiple-goroutines-listening-on-one-channel")

Bonda...

