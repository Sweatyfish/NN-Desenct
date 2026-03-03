package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

func initGraph(filepath string, N, D, K int) Graph {
	fmt.Println("Initializing Graph...")
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	graph := Graph{
		N:                      N,
		K:                      K,
		Dim:                    D,
		Data:                   make([]float64, N*D),
		NeighborsID:            make([]NeighborTuple, N*K),
		ReverseNeighbors:       make([]atomic.Pointer[[]NeighborTuple], N),
		Distances:              make([]float64, N*K),
		Locks:                  make([]sync.Mutex, N),
		FreezeReverseNeighbors: make([]atomic.Pointer[[]NeighborTuple], N),
	}

	// Insert vector data into graph
	for i, row := range records {
		for j, value := range row {
			graph.Data[i*D+j], err = strconv.ParseFloat(value, 64)
			if err != nil {
				panic(err)
			}
		}
	}

	// Initialize all atomic pointers to empty slices
	for i := 0; i < N; i++ {
		empty := make([]NeighborTuple, 0)
		graph.ReverseNeighbors[i].Store(&empty)

		emptyFreeze := make([]NeighborTuple, 0)
		graph.FreezeReverseNeighbors[i].Store(&emptyFreeze)
	}

	// Initialize neighbors and distances
	for I := 0; I < N; I++ {
		IdList := getKRandomNumbers(N, K, I)
		for J := 0; J < K; J++ {
			if IdList[J] == I {
				continue
			}

			// Set neighbor
			graph.NeighborsID[I*K+J] = NeighborTuple{Isnew: true, Id: IdList[J]}
			graph.Distances[I*K+J] = euclideanDistance(
				graph.Data[I*D:(I+1)*D],
				graph.Data[IdList[J]*D:(IdList[J]+1)*D],
			)

			// Add reverse neighbor
			revPointer := graph.ReverseNeighbors[IdList[J]].Load()
			*revPointer = append(*revPointer, NeighborTuple{Isnew: true, Id: I})
			graph.ReverseNeighbors[IdList[J]].Store(revPointer)
		}
	}

	// Deep copy ReverseNeighbors into FreezeReverseNeighbors
	for i := 0; i < N; i++ {
		ptr := graph.ReverseNeighbors[i].Load()
		copySlice := make([]NeighborTuple, len(*ptr))
		copy(copySlice, *ptr)
		graph.FreezeReverseNeighbors[i].Store(&copySlice)
	}

	fmt.Println("Graph Initialized")
	return graph
}
