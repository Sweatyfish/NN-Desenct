package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
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

	Graph := Graph{
		N:                N,
		K:                K,
		Dim:              D,
		Data:             make([]float64, N*D),
		NeighborsID:      make([]int, N*K),
		ReverseNeighbors: make([]atomic.Pointer[[]revNeighborTuple], N),
		Flags:            make([]bool, N*K),
		Distances:        make([]float64, N*K),
	}
	/* Insert Vector data into graph*/
	for i, row := range records {
		for j, value := range row {
			Graph.Data[i*D+j], err = strconv.ParseFloat(value, 64)
			if err != nil {
				panic(err)
			}
		}
	}

	/* Set all flags to true*/
	for i, _ := range Graph.Flags {
		Graph.Flags[i] = true
	}
	/* Initialize neighbors, distances and reverse neighbors*/
	for I := 0; I < N; I++ {
		IdList := getKRandomNumbers(N, K, I)
		for J := 0; J < K; J++ {
			Graph.NeighborsID[I*K+J] = IdList[J]
			Graph.Distances[I*K+J] = euclideanDistance(Graph.Data[I*D:(I+1)*D], Graph.Data[IdList[J]*D:(IdList[J]+1)*D])
			if Graph.ReverseNeighbors[IdList[J]].Load() == nil {
				Graph.ReverseNeighbors[IdList[J]].Store(&[]revNeighborTuple{{Id: I, New: true}})
			} else {
				revPointer := Graph.ReverseNeighbors[IdList[J]].Load()
				*revPointer = append(*revPointer, revNeighborTuple{Id: I, New: true})
				Graph.ReverseNeighbors[IdList[J]].Store(revPointer)
			}
		}
	}
	fmt.Println("Graph Initialized")

	return Graph
}
