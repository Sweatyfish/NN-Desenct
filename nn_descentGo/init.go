package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

func loadNpyFirstN(filename string, n int32, Dimensions int32) ([]float32, error) {
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
	data := make([]float32, n*Dimensions)
	err = binary.Read(f, binary.LittleEndian, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func initGraph(N, K int32) Graph {

	fmt.Println("Initializing Graph...")
	var (
		data []float32
		err  error
	)
	switch PCAcase {
	case 0:
		Dimensions = int32(384)
		data, err = loadNpyFirstN("../../data/train.npy", N, Dimensions)

	case 1:
		Dimensions = int32(320)
		data, err = loadNpyFirstN("../../data/reduced_320.npy", N, Dimensions)

	case 2:
		Dimensions = int32(160)
		data, err = loadNpyFirstN("../../data/reduced_160.npy", N, Dimensions)

	}

	if err != nil {
		panic(err)
	}

	graph := Graph{
		N:                N,
		K:                K,
		Dim:              Dimensions,
		Data:             data,
		NeighborsID:      make([]NeighborTuple, N*K),
		ReverseNeighbors: make([]atomic.Pointer[[]int32], N),
		Distances:        make([]float32, N*K),
		Locks:            make([]sync.Mutex, N),
	}

	for i := int32(0); i < N; i++ {
		empty := make([]int32, 0)
		graph.ReverseNeighbors[i].Store(&empty)
	}

	for I := int32(0); I < N; I++ {
		IdList := getKRandomNumbers(N, K, I)
		for J := int32(0); J < K; J++ {
			if IdList[J] == I {
				continue
			}
			graph.NeighborsID[I*K+J] = makeNeighbor(IdList[J], true)
			graph.Distances[I*K+J] = CosineDistance(
				graph.Data[I*Dimensions:(I+1)*Dimensions],
				graph.Data[IdList[J]*Dimensions:(IdList[J]+1)*Dimensions],
			)
			revPointer := graph.ReverseNeighbors[IdList[J]].Load()
			*revPointer = append(*revPointer, I)
			graph.ReverseNeighbors[IdList[J]].Store(revPointer)
		}
	}

	fmt.Println("Graph Initialized")
	return graph
}
