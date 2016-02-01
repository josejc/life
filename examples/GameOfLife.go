package main

//Imports
import (
	"rand"
	"time"
	"http"
	"fmt"
	"strconv"
)

//Struct template
type Iteration struct {
	board [][]bool
	row, col int
	time int
}

//This handles blank page responses
func handle(writer http.ResponseWriter, request *http.Request) {
	run(writer,0,0,0)
}

//This is what sends the Post data to the primarry "run" function
func handleTest(writer http.ResponseWriter, request *http.Request) {
	BoardSize, _ := strconv.Atoi(request.FormValue("BoardSize"))
	NumberThreads, _ := strconv.Atoi(request.FormValue("NumberThreads"))
	NumberIterations, _ := strconv.Atoi(request.FormValue("NumberIterations"))
	run(writer, BoardSize, NumberThreads, NumberIterations)
	http.Redirect(writer, request, "/", http.StatusFound)
}

//This is the Iteration constructor
func (iteration *Iteration) Init(size int) *Iteration {
	iteration.row = size
	iteration.col = size
	iteration.board = make([][]bool, size)
	for rowIndex := 0; rowIndex < size; rowIndex++ {
		iteration.board[rowIndex] = make([]bool, size)
		for colIndex := 0; colIndex < size; colIndex++ {
			if rand.Float32() < 0.3 {
				iteration.board[rowIndex][colIndex] = true
			}
		}
	}
	iteration.time = 0
	return iteration
}

//This is used by Print() to make thigns look preatty
func displayState(state bool) (str string) {
	if state {
		str = "[*]"
	} else {
		str = "[ ]"
	}
	return str
}

//This displays the board 
func (iteration *Iteration) Print() {
	for rowIndex := 0; rowIndex < iteration.row; rowIndex++ {
		for colIndex := 0; colIndex < iteration.col; colIndex++ {
			print(displayState(iteration.board[rowIndex][colIndex]))
		}
		println()
	}
	println()
}

//This is used to count the state of the selected cell's neighbors
func (iteration *Iteration) count(rowIndex, colIndex int) (state int) {
	if rowIndex < 0 {
		rowIndex += iteration.row
	}
	if iteration.row <= rowIndex {
		rowIndex -= iteration.row
	}
	if colIndex < 0 {
		colIndex += iteration.col
	}
	if iteration.col <= colIndex {
		colIndex -= iteration.col
	}
	if iteration.board[rowIndex][colIndex] {
		state = 1
	} else {
		state = 0
	}
	return state
}

//This determines what the cell will be next iteration
func (iteration *Iteration) calculateState(rowIndex, colIndex int) (state bool) {
	count := iteration.count(rowIndex - 1, colIndex - 1) +
		 iteration.count(rowIndex - 1, colIndex    ) +
		 iteration.count(rowIndex - 1, colIndex + 1) +
		 iteration.count(rowIndex    , colIndex - 1) +
		 iteration.count(rowIndex    , colIndex    ) +
		 iteration.count(rowIndex    , colIndex + 1) +
		 iteration.count(rowIndex + 1, colIndex - 1) +
		 iteration.count(rowIndex + 1, colIndex    ) +
		 iteration.count(rowIndex + 1, colIndex + 1)
	switch count {
		case 3:
			state = true
		case 4:
			state = iteration.board[rowIndex][colIndex]
		default:
			state = false
	}
	return state
}

//This is where each worker thread does is business
func workOnPortion(iteration *Iteration, channel chan<- *Iteration, next *Iteration,threadId int, NumberThreads int, BoardSize int) {
	workOffset := (BoardSize / NumberThreads)
	workStart := (workOffset)*(threadId - 1)
	workEnd := (workOffset)*(threadId)

	for rowIndex := workStart; rowIndex < workEnd; rowIndex++ {
		for colIndex := workStart; colIndex < workEnd; colIndex++ {
			next.board[rowIndex][colIndex] = iteration.calculateState(rowIndex, colIndex)
		} 
	}
	
	next.time = iteration.time +1
	channel <- next
	iteration = next
}

//This is the master thread that is in charge of managing each worker thread
func generate_gen(iteration *Iteration, channel chan<- *Iteration, quit chan bool, NumberIterations int, NumberThreads int, BoardSize int) {
	count:= 0
	for count < NumberIterations {
		next := new(Iteration).Init(iteration.row)
		for threadId:= 1; threadId < NumberThreads+1; threadId++ {
			go workOnPortion(iteration, channel, next, threadId, NumberThreads, BoardSize)
		}
		next.Print()
		count = count + 1
	}
	//Exit since were done
	quit <- true	
}

//Required by the web frontend
func init() {
	http.HandleFunc("/", handle)
	http.HandleFunc("/test", handleTest)
}

//This wrtied the header information for all http rqeuests
func writeHeader(writer http.ResponseWriter, BoardSize int, NumberThreads int, NumberIterations int) {
	fmt.Fprint(writer, "<html><body>")
	fmt.Fprint(writer, "Starting Game of Life with [")
	fmt.Fprint(writer, NumberIterations)
	fmt.Fprint(writer, "] Iterations with a Board Size of [")
	fmt.Fprint(writer, BoardSize)
	fmt.Fprint(writer, "] And [")
	fmt.Fprint(writer, NumberThreads)
	fmt.Fprint(writer, "] Threads! \n<br/>")
}

//This builds the form that is used to collect data
func writeForm(writer http.ResponseWriter) {
	fmt.Fprint(writer, "<form action=\"/test\" method=\"/post\">")
	fmt.Fprint(writer, "Board Size: <input type=\"text\" name=\"BoardSize\" /><br/>")
	fmt.Fprint(writer, "Number Of Threads: <input type=\"text\" name=\"NumberThreads\" /><br/>")
	fmt.Fprint(writer, "Number Of Iterations: <input type=\"text\" name=\"NumberIterations\" /><br/>")
	fmt.Fprint(writer, "<input type=\"submit\" value=\"Test!\"></form>")
}

//This is the "main" function that is called each time data is submitted
func run(writer http.ResponseWriter, BoardSize int, NumberThreads int, NumberIterations int) {
	writeForm(writer)
	writeHeader(writer, BoardSize, NumberThreads, NumberIterations)
	startTime := time.Nanoseconds()
	game := new(Iteration).Init(BoardSize)
	channel := make(chan *Iteration)
	quit := make(chan bool)
	go generate_gen(game, channel, quit, NumberIterations, NumberThreads, BoardSize)
	endTime := time.Nanoseconds() - startTime
	fmt.Fprint(writer, "<br/>Total Time:")
	fmt.Fprint(writer, endTime)
	fmt.Fprint(writer, "ns\n</body></html>")
	<-quit
}