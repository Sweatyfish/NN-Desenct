package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

func loadNpyFirstN(filename string, n int, D int) ([]float32, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the 128-byte header (typical .npy v1.0-1.1)
	header := make([]byte, 128)
	_, err = f.Read(header)
	if err != nil {
		return nil, err
	}

	// Allocate only n*D
	data := make([]float32, n*D)

	// Read first n rows directly
	err = binary.Read(f, binary.LittleEndian, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func initGraph(N, D, K int) Graph {
	fmt.Println("Initializing Graph...")

	data, err := loadNpyFirstN("../../data/train.npy", N, D)
	if err != nil {
		panic(err)
	}

	graph := Graph{
		N:                      N,
		K:                      K,
		Dim:                    D,
		Data:                   make([]float32, N*D),
		NeighborsID:            make([]NeighborTuple, N*K),
		ReverseNeighbors:       make([]atomic.Pointer[[]NeighborTuple], N),
		Distances:              make([]float32, N*K),
		Locks:                  make([]sync.Mutex, N),
		FreezeReverseNeighbors: make([]atomic.Pointer[[]NeighborTuple], N),
	}
	// Insert vector data into graph
	copy(graph.Data, data)

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
			graph.NeighborsID[I*K+J] = NeighborTuple{isNew: true, Id: IdList[J]}
			graph.Distances[I*K+J] = euclideanDistance(
				graph.Data[I*D:(I+1)*D],
				graph.Data[IdList[J]*D:(IdList[J]+1)*D],
			)

			// Add reverse neighbor
			revPointer := graph.ReverseNeighbors[IdList[J]].Load()
			*revPointer = append(*revPointer, NeighborTuple{isNew: true, Id: I})
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
