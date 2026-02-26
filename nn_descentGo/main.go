package main

import (
	"sync"
	"sync/atomic"
	"time"
)

/*Set Filname to the corresponding data set you want to load*/
var filename string = "data-N_1000-D_10.csv"
var filepath string = "Data/" + filename

/* amount of neighbors to be considered for each point, can be changed to any number you want*/
var K = 10
var Delta = 0.001
var numThreads = 4

/* amount of verticies each lock resides over */

type NeighborTuple struct {
	Isnew bool
	Id    int
}

type Graph struct {
	N, K, Dim              int
	Data                   []float64 /*Prolly needs changing*/
	NeighborsID            []NeighborTuple
	ReverseNeighbors       []atomic.Pointer[[]NeighborTuple]
	Distances              []float64
	FreezeNeighbors        []NeighborTuple
	FreezeReverseNeighbors []atomic.Pointer[[]NeighborTuple]
}

var (
	graph             Graph
	NewneigboursLock  sync.Mutex
	Newneigboursfound int
	Counterlock       sync.Mutex
	Counter           int
)

func NNDecent(c chan int) {
	/* Iterate over all verticies in the graph, and for each vertex, iterate over its neighbors and reverse neighbors to find potential new neighbors. If a new neighbor is found, add it to the graph and mark it as new. */

}

func main() {
	N, D := getNandDFromFilename(filename)
	graph = initGraph(filepath, N, D, K)

	c := make(chan int)

	for i := 0; i < numThreads; i++ {

		go NNDecent(c)

	}

	for float64(Newneigboursfound) > Delta*float64(K)*float64(N) {
		NewneigboursLock.Lock()
		Newneigboursfound = 0
		NewneigboursLock.Unlock()
		for i := 0; i < N; i++ {
			c <- i

		}
		for Counter != N {
			time.Sleep(50 * time.Millisecond)
		}
	}

}
