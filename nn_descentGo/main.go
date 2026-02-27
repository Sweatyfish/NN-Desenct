package main

import (
	"sync"
	"sync/atomic"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

/*Set Filname to the corresponding data set you want to load*/
var filename string = "data-N_1000-D_10.csv"
var filepath string = "Data/" + filename

/* amount of neighbors to be considered for each point, can be changed to any number you want*/
var K = 10
var Delta = 0.001
var numThreads = 4
var rho float32 = 0.5

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
	FreezeReverseNeighbors []atomic.Pointer[[]NeighborTuple]
	Locks                  []sync.Mutex
}

var (
	graph              Graph
	NewneighboursLock  sync.Mutex
	Newneighboursfound int
	Counterlock        sync.Mutex
	Counter            int
)

func NNDecent(c chan int) {

	oldneighbours := mapset.NewSet[int]()
	newneighbours := mapset.NewSet[int]()
	oldprime := mapset.NewSet[int]()
	newprime := mapset.NewSet[int]()
	for true {
		V := <-c
		oldneighbours.Clear()
		newneighbours.Clear()
		oldprime.Clear()
		newprime.Clear()

		/*We go through all the neighbours and check whether they are new or old neighbours */
		for _, neighbor := range getneighbour(V, graph) {
			if neighbor.Isnew {
				newneighbours.Add(neighbor.Id)
			} else {
				oldneighbours.Add(neighbor.Id)
			}
		}
		/*Iterate over the reverse neighbors of V and add them to the corresponding sets based on whether they are new or old neighbors.*/
		/*We take reverseneighbours and add them to two seperate lists*/
		for Vertex := range newneighbours.Union(oldneighbours).Iter() {
			for _, neighbor := range getreverseneighbour(Vertex, graph) {
				if neighbor.Isnew {
					newprime.Add(neighbor.Id)
				} else {
					oldprime.Add(neighbor.Id)
				}

			}
			newneighbours = newneighbours.Union(sampleKRandomNeighbors(newprime, rho))
			oldneighbours = oldneighbours.Union(sampleKRandomNeighbors(oldprime, rho))
			newneighbourslist := newneighbours.ToSlice()
			oldneighbourslist := oldneighbours.ToSlice()

			for i := 0; i < len(newneighbourslist); i++ {
				for j := i + 1; j < len(newneighbourslist); j++ {

					distance = euclideanDistance(getvertex(newneighbourslist[i], graph), getvertex(newneighbourslist[j], graph))
					tryInsert()
				}
			}

		}
	}
}

func main() {
	N, D := getNandDFromFilename(filename)
	graph = initGraph(filepath, N, D, K)

	c := make(chan int)

	for i := 0; i < numThreads; i++ {

		go NNDecent(c)

	}

	for float64(Newneigboursfound) > Delta*float64(K)*float64(N) {
		NewneighboursLock.Lock()
		Newneigboursfound = 0
		NewneighboursLock.Unlock()
		for i := 0; i < N; i++ {
			c <- i

		}
		for Counter != N {
			time.Sleep(50 * time.Millisecond)
		}
	}

}
