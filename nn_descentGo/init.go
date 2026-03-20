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
	header := make([]byte, 128)
	_, err = f.Read(header)
	if err != nil {
		return nil, err
	}
	data := make([]float32, n*D)
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
		Data:                   data,
		NeighborsID:            make([]NeighborTuple, N*K),
		ReverseNeighbors:       make([]atomic.Pointer[[]int], N),
		Distances:              make([]float32, N*K),
		Locks:                  make([]sync.Mutex, N),
		FreezeReverseNeighbors: make([]atomic.Pointer[[]int], N),
	}

	for i := 0; i < N; i++ {
		empty := make([]int, 0)
		graph.ReverseNeighbors[i].Store(&empty)
		emptyFreeze := make([]int, 0)
		graph.FreezeReverseNeighbors[i].Store(&emptyFreeze)
	}

	for I := 0; I < N; I++ {
		IdList := getKRandomNumbers(N, K, I)
		for J := 0; J < K; J++ {
			if IdList[J] == I {
				continue
			}
			graph.NeighborsID[I*K+J] = makeNeighbor(IdList[J], true)
			graph.Distances[I*K+J] = CosineDistance(
				graph.Data[I*D:(I+1)*D],
				graph.Data[IdList[J]*D:(IdList[J]+1)*D],
			)
			revPointer := graph.ReverseNeighbors[IdList[J]].Load()
			*revPointer = append(*revPointer, I)
			graph.ReverseNeighbors[IdList[J]].Store(revPointer)
		}
	}

	for i := 0; i < N; i++ {
		ptr := graph.ReverseNeighbors[i].Load()
		copySlice := make([]int, len(*ptr))
		copy(copySlice, *ptr)
		graph.FreezeReverseNeighbors[i].Store(&copySlice)
	}

	fmt.Println("Graph Initialized")
	return graph
}
