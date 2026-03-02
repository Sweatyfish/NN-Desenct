package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

/*Set Filname to the corresponding data set you want to load*/
var filename string = "data-N_5000-D_80.csv"
var filepath string = "Data/" + filename

/* amount of neighbors to be considered for each point, can be changed to any number you want*/
var K = 10
var Delta = 0.001
var numThreads = 4
var rho float32 = 0.5
var benchmarking = true

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
	graph Graph
	//This int counts the amount of new neighbors found in each iteration, it is used to check the stopping condition
	//Remeber to use the associated lock when modifying this variable
	NewneighboursLock  sync.Mutex
	Newneighboursfound int
	//This is used to count the amount of vertices that have been processed in each iteration, it is used to make sure all vertices are processed before starting a new iteration
	//Remeber to use the associated lock when modifying this variable
	Counterlock sync.Mutex
	Counter     int
)

func NNDecent(c chan int) {
	//Instantiate the sets of old and new neighbors
	oldneighbours := mapset.NewSet[int]()
	newneighbours := mapset.NewSet[int]()
	oldprime := mapset.NewSet[int]()
	newprime := mapset.NewSet[int]()
	for true {
		// Wait for a id to be sent on the channel and then process it
		V := <-c
		//Reset all the sets for a new iteration
		oldneighbours.Clear()
		newneighbours.Clear()
		oldprime.Clear()
		newprime.Clear()
		//This will add up to be the amount of new neighbors this iteration on this vertex found
		addToNewNeighbors := 0

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
			//Union operations on the 4 sets with sampling, also making them into a slice
			newneighbours = newneighbours.Union(sampleKRandomNeighbors(newprime, rho))
			oldneighbours = oldneighbours.Union(sampleKRandomNeighbors(oldprime, rho))
			newneighbourslist := newneighbours.ToSlice()
			oldneighbourslist := oldneighbours.ToSlice()

			//Set all current neighbors and reverse neighbors to old neighbors for the next iteration needs to be done before we begin changing neighbors
			for _, neighbor := range getneighbour(Vertex, graph) {
				neighbor.Isnew = false
			}

			//The same for the reverse neighbors
			//This (!IMPORTANT!) This is SUBOPTIMAL, but currently i dont now of a better way to do this, since no thread has the full picture of which are new or old reverse neighbors
			graph.Locks[Vertex].Lock()

			slicePtr := graph.FreezeReverseNeighbors[Vertex].Load()
			if slicePtr != nil {
				slice := *slicePtr
				for i := range slice {
					slice[i].Isnew = false
				}
			}

			graph.Locks[Vertex].Unlock()
			//The main loop checking neighbors neighbors against each other
			for i := 0; i < len(newneighbourslist); i++ {
				for j := i + 1; j < len(newneighbourslist); j++ {

					distance := euclideanDistance(getvertex(newneighbourslist[i], graph), getvertex(newneighbourslist[j], graph))
					addToNewNeighbors += tryInsert(graph, newneighbourslist[i], newneighbourslist[j], distance)
				}
				for j := 0; j < len(oldneighbourslist); j++ {
					distance := euclideanDistance(getvertex(newneighbourslist[i], graph), getvertex(oldneighbourslist[j], graph))
					addToNewNeighbors += tryInsert(graph, newneighbourslist[i], oldneighbourslist[j], distance)
				}
			}
			//Adding the amount of new neighbors found to the total amount of new neighbors found in this iteration
			NewneighboursLock.Lock()
			Newneighboursfound += addToNewNeighbors
			NewneighboursLock.Unlock()

		}
		//Adding a finished vertex to the counter lock
		Counterlock.Lock()
		Counter++
		Counterlock.Unlock()

	}
}

func main() {
	N, D := getNandDFromFilename(filename)
	graph = initGraph(filepath, N, D, K)
	//Instastiate to -1 for entering the first loop
	Newneighboursfound = -1
	c := make(chan int)
	//Start the NNDescent algorithm with the specified amount of threads
	for i := 0; i < numThreads; i++ {

		go NNDecent(c)

	}
	Iterations := 0
	//This master threads will send vertex ids to the worker threads and check the stopping condition after each iteration, it will also update the reverse neighbors with the freeze reverse neighbors after each iteration
	//Updating the reverse neighbors could be multithreaded as well, but would require deeper changes to the code, so I decided to keep it single threaded for now
	//The way we currently keep track of which neighbors to switch should also be a heap currently is not optimal
	for float64(Newneighboursfound) > Delta*float64(K)*float64(N) || Newneighboursfound == -1 {
		fmt.Println("Iteration number:", Iterations)
		NewneighboursLock.Lock()
		Newneighboursfound = 0
		NewneighboursLock.Unlock()
		for i := 0; i < N; i++ {
			c <- i

		}
		//Wait until all vertices have been processed before procedding
		for Counter != N {
			time.Sleep(50 * time.Millisecond)
		}
		//Reset the counter for the next iteration
		Counterlock.Lock()
		Counter = 0
		Counterlock.Unlock()
		//We load all data from the freeze reverse neighbors into the reverse neighbor
		for i := 0; i < N; i++ {
			//We need to make a copy of the slice because otherwise we would be modifying the same slice for all vertices that share the same reverse neighbor
			ListToCopy := *graph.FreezeReverseNeighbors[i].Load()
			newSlice := make([]NeighborTuple, len(ListToCopy))
			copy(newSlice, ListToCopy)
			graph.ReverseNeighbors[i].Store(&newSlice)
		}

		Iterations++
		fmt.Println("New neighbors found in this iteration:", Newneighboursfound)
	}

	if benchmarking {
		fmt.Println("Calculating accuracy...")
		accuracy := benchmark(graph)
		fmt.Println("Calculated Accuracy is:", accuracy)
	}

}
