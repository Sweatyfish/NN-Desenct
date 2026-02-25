package main

import (
	"sync"
	"sync/atomic"
)

/*Set Filname to the corresponding data set you want to load*/
var filename string = "data-N_1000-D_10.csv"
var filepath string = "Data/" + filename

/* amount of neighbors to be considered for each point, can be changed to any number you want*/
var K = 10

/* amount of verticies each lock resides over */
type revNeighborTuple struct {
	Id  int
	New bool
}

type Graph struct {
	N, K, Dim        int
	Data             []float64 /*Prolly needs changing*/
	NeighborsID      []int
	ReverseNeighbors []atomic.Pointer[[]revNeighborTuple]
	Flags            []bool
	Distances        []float64
}

var (
	graph       Graph
	counterLock sync.Mutex
	counter     int
)

func main() {
	N, D := getNandDFromFilename(filename)
	graph = initGraph(filepath, N, D, K)

}
