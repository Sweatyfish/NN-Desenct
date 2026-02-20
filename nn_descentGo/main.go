package main

import (
	"sync"
)

/*Set Filname to the corresponding data set you want to load*/
var filename string = "data-N:1000-D:10.csv"
var filepath string = "Data/" + filename

/* amount of neighbors to be considered for each point, can be changed to any number you want*/
var K = 10

/* amount of verticies each lock resides over */
var lockSize = 128

type Graph struct {
	N, K, Dim          int
	Data               []float64 /*Prolly needs changing*/
	NeighborsID        []int
	ReverseNeighborsID [][]int
	Flags              []bool
	Distances          []float64
	Locks              []sync.Mutex
}

var (
	Graph_old Graph
	Graph_new Graph
)

func main() {
	N, D := getNandDFromFilename(filename)
	Graph_old = initGraph(filepath, N, D, K, lockSize)
	Graph_new = Graph_old
}
