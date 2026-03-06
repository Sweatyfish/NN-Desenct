package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

/* amount of neighbors to be considered for each point, can be changed to any number you want*/
var k = 15
var n = 10000
var delta = 0.001
var numThreads = 8
var rho float32 = 0.5
var benchmarking = false
var timeMeasure = false

/* amount of verticies each lock resides over */

type NeighborTuple struct {
	isNew bool
	Id    int
}

type Graph struct {
	N, K, Dim              int
	Data                   []float32 /*Prolly needs changing*/
	NeighborsID            []NeighborTuple
	ReverseNeighbors       []atomic.Pointer[[]int]
	Distances              []float32
	FreezeReverseNeighbors []atomic.Pointer[[]int]
	Locks                  []sync.Mutex
}

var (
	graph Graph
	//This int counts the amount of new neighbors found in each iteration, it is used to check the stopping condition
	//Remeber to use the associated lock when modifying this variable
	newNeighborsLock  sync.Mutex
	newNeighborsFound int
	//This is used to count the amount of vertices that have been processed in each iteration, it is used to make sure all vertices are processed before starting a new iteration
	//Remeber to use the associated lock when modifying this variable
	counterLock sync.Mutex
	counter     int
)

func NNDecent(c chan int) {
	//Instantiate the sets of old and new neighbors
	oldNeighbors := mapset.NewSet[int]()
	newNeighbors := mapset.NewSet[int]()
	oldPrime := mapset.NewSet[int]()
	newPrime := mapset.NewSet[int]()

	for true {
		// Wait for a id to be sent on the channel and then process it
		V := <-c
		//Reset all the sets for a new iteration
		oldNeighbors.Clear()
		newNeighbors.Clear()
		oldPrime.Clear()
		newPrime.Clear()
		//This will add up to be the amount of new neighbors this iteration on this vertex found
		addToNewNeighbours := 0

		/*We go through all the neighbours and check whether they are new or old neighbours */
		for _, neighbor := range getNeighbor(V) {
			if neighbor.Id != V {
				if neighbor.isNew {
					newNeighbors.Add(neighbor.Id)
				} else {
					oldNeighbors.Add(neighbor.Id)
				}
			}
		}

		/*Iterate over the reverse neighbors of V and add them to the corresponding sets based on whether they are new or old neighbors.*/
		/*We take reverseneighbours and add them to two seperate lists*/
		for neighbor := range newNeighbors.Iter() {
			for _, reverseNeighbor := range getReverseNeighbor(neighbor) {
				newPrime.Add(reverseNeighbor)
			}
		}
		for neighbor := range oldNeighbors.Iter() {
			for _, reverseNeighbor := range getReverseNeighbor(neighbor) {
				oldPrime.Add(reverseNeighbor)
			}
		}

		//Union operations on the 4 sets with sampling, also making them into a slice
		newNeighbors = newNeighbors.Union(sampleKRandomNeighbors(newPrime, rho))
		oldNeighbors = oldNeighbors.Union(sampleKRandomNeighbors(oldPrime, rho))
		newNeighboursList := newNeighbors.ToSlice()
		oldNeighboursList := oldNeighbors.ToSlice()

		//Set all current neighbors and reverse neighbors to old neighbors for the next iteration needs to be done before we begin changing neighbors
		neighbors := getNeighbor(V)
		for i := range neighbors {
			neighbors[i].isNew = false
		}

		//The same for the reverse neighbors
		//This (!IMPORTANT!) This is SUBOPTIMAL, but currently i dont now of a better way to do this, since no thread has the full picture of which are new or old reverse neighbors
		//The main loop checking neighbors neighbors against each other
		DistancesNew := CosineDistanceBatch(newNeighboursList)
		idx := 0
		for i := 0; i < len(newNeighboursList); i++ {

			for j := i + 1; j < len(newNeighboursList); j++ {
				addToNewNeighbours += tryInsert(newNeighboursList[i], newNeighboursList[j], DistancesNew[idx])
				idx++

			}

			/*
				for j := i + 1; j < len(newNeighboursList); j++ {

					distance := CosineDistance(getVertex(newNeighboursList[i]), getVertex(newNeighboursList[j]))
					//println("Distance between", newneighbourslist[i], "and", newneighbourslist[j], "is", distance)
					addToNewNeighbours += tryInsert(newNeighboursList[i], newNeighboursList[j], distance)
				}
			*/
			/*
				if len(oldNeighboursList) < 1 {
					continue

				}



				for d := 0; d < len(DistancesOld); d++ {
					println("kom nu")

					addToNewNeighbours += tryInsert(newNeighboursList[i], oldNeighboursList[d], DistancesOld[d])
				}
			*/

			for j := 0; j < len(oldNeighboursList); j++ {
				distance := CosineDistance(getVertex(newNeighboursList[i]), getVertex(oldNeighboursList[j]))
				addToNewNeighbours += tryInsert(newNeighboursList[i], oldNeighboursList[j], distance)
			}

		}

		//Adding the amount of new neighbors found to the total amount of new neighbors found in this iteration
		newNeighborsLock.Lock()
		newNeighborsFound += addToNewNeighbours
		newNeighborsLock.Unlock()

		//Adding a finished vertex to the counter lock
		counterLock.Lock()
		counter++
		counterLock.Unlock()
	}
}

func main() {
	// Start CPU profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer func() {
		pprof.StopCPUProfile()
		f.Close()
	}()
	graph = initGraph(n, 384, k)
	//Instastiate to -1 for entering the first loop
	newNeighborsFound = -1
	c := make(chan int)
	//Start the NNDescent algorithm with the specified amount of threads
	for i := 0; i < numThreads; i++ {
		go NNDecent(c)
	}
	iterations := 0
	//This master threads will send vertex ids to the worker threads and check the stopping condition after each iteration, it will also update the reverse neighbors with the freeze reverse neighbors after each iteration
	//Updating the reverse neighbors could be multithreaded as well, but would require deeper changes to the code, so I decided to keep it single threaded for now
	//The way we currently keep track of which neighbors to switch should also be a heap currently is not optimal
	totalTimeStart := time.Now()
	for float64(newNeighborsFound) > delta*float64(k)*float64(n) || newNeighborsFound == -1 {
		fmt.Println("Iteration number:", iterations)

		newNeighborsLock.Lock()
		newNeighborsFound = 0
		newNeighborsLock.Unlock()

		start := time.Now()
		for i := 0; i < n; i++ {
			c <- i
		}

		//Wait until all vertices have been processed before procedding
		for counter != n {
			time.Sleep(50 * time.Millisecond)
		}

		end := time.Now()
		if timeMeasure {
			fmt.Println("Time taken to process all vertices in this iteration:", end.Sub(start))
		}

		//Reset the counter for the next iteration
		counterLock.Lock()
		counter = 0
		counterLock.Unlock()
		//We load all data from the freeze reverse neighbors into the reverse neighbor
		start = time.Now()
		for i := 0; i < n; i++ {
			//We need to make a copy of the slice because otherwise we would be modifying the same slice for all vertices that share the same reverse neighbor
			ListToCopy := *graph.FreezeReverseNeighbors[i].Load()
			newSlice := make([]int, len(ListToCopy))
			copy(newSlice, ListToCopy)
			graph.ReverseNeighbors[i].Store(&newSlice)
		}

		end = time.Now()
		if timeMeasure {
			fmt.Println("Time taken to update reverse neighbors:", end.Sub(start))
		}

		iterations++
		fmt.Println("New neighbors found in this iteration:", newNeighborsFound)
	}

	start := time.Now()
	fmt.Println(time.Since(totalTimeStart))

	if benchmarking {
		fmt.Println("Calculating accuracy...")
		accuracy := benchmark(graph)
		fmt.Println("Calculated Accuracy is:", accuracy, "%")
	}

	end := time.Now()
	if timeMeasure {
		fmt.Println("Time taken to benchmark:", end.Sub(start))
	}
}
