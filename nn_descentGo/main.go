package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

var k = int32(15)
var n = int32(3000000)
var delta = 0.001
var numThreads = 8
var rho float32 = 0.5
var benchmarking = false
var timeMeasure = false
var checkMemory = false

// NeighborTuple is packed into a single int32 to save memory.
// Negative value = new neighbor, positive = old neighbor.
// ID is stored as abs(value). 0 is reserved, so IDs are 1-indexed internally.
type NeighborTuple = int32

func makeNeighbor(id int32, isNew bool) NeighborTuple {
	if isNew {
		return -int32(id + 1)
	}
	return int32(id + 1)
}

func neighborIsNew(t NeighborTuple) bool {
	return t < 0
}

func neighborID(t NeighborTuple) int32 {
	if t < 0 {
		return int32(-t) - 1
	}
	return int32(t) - 1
}

func setOld(t *NeighborTuple) {
	if *t < 0 {
		*t = -*t
	}
}

type Graph struct {
	N, K, Dim        int32
	Data             []float32
	NeighborsID      []NeighborTuple
	ReverseNeighbors []atomic.Pointer[[]int32]
	Distances        []float32
	Locks            []sync.Mutex
}

var (
	graph             Graph
	newNeighborsLock  sync.Mutex
	newNeighborsFound int32
	counterLock       sync.Mutex
	counter           int32
)

type neighborInfo struct {
	id       int32
	distance float32
	index    int32
}

func NNDecent(c chan int32) {
	oldNeighbors := mapset.NewSet[int32]()
	newNeighbors := mapset.NewSet[int32]()
	oldPrime := mapset.NewSet[int32]()
	newPrime := mapset.NewSet[int32]()

	for true {
		V := <-c
		oldNeighbors.Clear()
		newNeighbors.Clear()
		oldPrime.Clear()
		newPrime.Clear()
		addToNewNeighbours := int32(0)

		for _, neighbor := range getNeighbor(V) {
			id := neighborID(neighbor)
			if id != V {
				if neighborIsNew(neighbor) {
					newNeighbors.Add(id)
				} else {
					oldNeighbors.Add(id)
				}
			}
		}

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

		newNeighbors = newNeighbors.Union(sampleKRandomNeighbors(newPrime, rho))
		oldNeighbors = oldNeighbors.Union(sampleKRandomNeighbors(oldPrime, rho))
		newNeighboursList := newNeighbors.ToSlice()
		oldNeighboursList := oldNeighbors.ToSlice()

		neighbors := getNeighbor(V)
		for i := range neighbors {
			setOld(&neighbors[i])
		}

		NxNMatrix := CosineDistanceBatchN(newNeighboursList)
		NxOMatrix := CosineDistanceBatchNM(newNeighboursList, oldNeighboursList)

		idx1 := int32(0)
		idx2 := int32(0)
		added := int32(0)

		var worstPrimaryList []neighborInfo
		newNlen := int32(len(newNeighboursList))
		oldNlen := int32(len(oldNeighboursList))
		worstPrimaryList = getWorstNeighborInfoBatch(newNeighboursList)
		for i := int32(0); i < newNlen; i++ {
			for j := i + int32(1); j < newNlen; j++ {
				dist := NxNMatrix[idx1]
				if worstPrimaryList[i].distance > dist {
					added, worstPrimaryList[i] = insert(newNeighboursList[i], newNeighboursList[j], worstPrimaryList[i], dist)
					addToNewNeighbours += added
				}
				if worstPrimaryList[j].distance > dist {
					added, worstPrimaryList[j] = insert(newNeighboursList[j], newNeighboursList[i], worstPrimaryList[j], dist)
					addToNewNeighbours += added
				}
				idx1++
			}
			for j := int32(0); j < oldNlen; j++ {
				dist := NxOMatrix[idx2]
				idx2++
				if worstPrimaryList[i].distance > dist {
					added, worstPrimaryList[i] = insert(newNeighboursList[i], oldNeighboursList[j], worstPrimaryList[i], dist)
					addToNewNeighbours += added
				}
			}
		}

		newNeighborsLock.Lock()
		newNeighborsFound += addToNewNeighbours
		newNeighborsLock.Unlock()

		counterLock.Lock()
		counter++
		counterLock.Unlock()
	}
}

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer func() {
		pprof.StopCPUProfile()
		f.Close()
	}()

	graph = initGraph(n, int32(384), k)
	newNeighborsFound = -1
	c := make(chan int32)
	for i := 0; i < numThreads; i++ {
		go NNDecent(c)
	}

	iterations := 0
	totalTimeStart := time.Now()
	for float64(newNeighborsFound) > delta*float64(k)*float64(n) || newNeighborsFound == -1 {
		fmt.Println("Iteration number:", iterations)

		newNeighborsLock.Lock()
		newNeighborsFound = 0
		newNeighborsLock.Unlock()

		start := time.Now()
		for i := int32(0); i < n; i++ {
			c <- i
		}

		for counter != n {
			time.Sleep(50 * time.Millisecond)
		}

		end := time.Now()
		if timeMeasure {
			fmt.Println("Time taken to process all vertices in this iteration:", end.Sub(start))
		}

		counterLock.Lock()
		counter = 0
		counterLock.Unlock()

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
	if checkMemory {
		f, err := os.Create("mem.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		runtime.GC() // important: get up-to-date heap
		pprof.WriteHeapProfile(f)
	}
}
