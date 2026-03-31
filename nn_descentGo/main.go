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

// max n = 3001496
var k = int32(15)
var n = int32(800000)
var delta = 0.001
var numThreads = 12
var rho float32 = 0.5
var benchmarking = false
var benchmarkingReal = true
var timeMeasure = false
var checkMemory = false

// write pca dimensions
var PCAcase = 3

// 0 = 384 dimension
// 1 = 320 dimension
// 2 = 160 dimension
// 3 = 144 dimension
// 4 = 128 dimension
// 5 = 120 dimension
// 6 = 80  dimension
var Dimensions int32

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
	graph                   Graph
	newNeighborsFoundAtomic int64
)

var counter int64

type neighborInfo struct {
	id       int32
	distance float32
	index    int32
}

func NNDecent(c chan int32) {
	oldNeighbors := mapset.NewSet[int32]()
	newNeighbors := mapset.NewSet[int32]()
	//oldPrime := mapset.NewSet[int32]()
	newPrime := mapset.NewSet[int32]()

	for true {
		V := <-c
		addToNewNeighbours := int32(0)
		oldNeighbors.Clear()
		newNeighbors.Clear()
		//oldPrime.Clear()
		newPrime.Clear()

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
		/*for neighbor := range oldNeighbors.Iter() {
			for _, reverseNeighbor := range getReverseNeighbor(neighbor) {
				oldPrime.Add(reverseNeighbor)
			}
		}*/

		newNeighbors = newNeighbors.Union(sampleKRandomNeighbors(newPrime, rho))
		//oldNeighbors = oldNeighbors.Union(sampleKRandomNeighbors(oldPrime, rho))
		//oldPrime.Clear()
		newPrime.Clear()
		newNeighboursList := newNeighbors.ToSlice()
		oldNeighboursList := oldNeighbors.ToSlice()
		oldNeighbors.Clear()
		newNeighbors.Clear()

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

		atomic.AddInt64(&newNeighborsFoundAtomic, int64(addToNewNeighbours))

		atomic.AddInt64(&counter, 1)
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

	graph = initGraph(n, k)
	atomic.StoreInt64(&newNeighborsFoundAtomic, -1)
	c := make(chan int32)
	for i := 0; i < numThreads; i++ {
		go NNDecent(c)
	}

	iterations := 0
	totalTimeStart := time.Now()
	for atomic.LoadInt64(&newNeighborsFoundAtomic) > int64(delta*float64(k)*float64(n)) || atomic.LoadInt64(&newNeighborsFoundAtomic) == -1 {
		fmt.Println("Iteration number:", iterations)

		atomic.StoreInt64(&newNeighborsFoundAtomic, 0)

		start := time.Now()
		for i := int32(0); i < n; i++ {
			c <- i
		}

		end := time.Now()
		if timeMeasure {
			fmt.Println("Time taken to process all vertices in this iteration:", end.Sub(start))
		}

		for atomic.LoadInt64(&counter) != int64(n) {
			time.Sleep(50 * time.Millisecond)
		}
		atomic.StoreInt64(&counter, 0)

		if timeMeasure {
			fmt.Println("Time taken to update reverse neighbors:", end.Sub(start))
		}

		iterations++
		fmt.Println("New neighbors found in this iteration:", atomic.LoadInt64(&newNeighborsFoundAtomic))
	}
	fmt.Println("n: ", n, " k: ", k, " Dimensions: ", Dimensions, " threads: ", numThreads, " rho: ", rho, " delta: ", delta)

	start := time.Now()
	fmt.Println(time.Since(totalTimeStart))

	if benchmarking {
		fmt.Println("Calculating accuracy...")

	}

	var groundTruth [][]int
	if benchmarkingReal {
		groundTruth = loadGroundTruth("../../data/groundtruth.i32", n, k)
		accuracy := benchmarkNew(graph, groundTruth)
		fmt.Println("Calculated Accuracy is:", accuracy, "%")
		fmt.Println("Accuracy considering n:", accuracy*(float32(3001496)/float32(n)), "%")
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
