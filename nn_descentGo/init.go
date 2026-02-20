package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

func initGraph(filepath string, N, D, K int) {
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
		N:                  N,
		K:                  K,
		Dim:                D,
		Data:               make([]float64, N*D),
		NeighborsID:        make([]int, N*K),
		ReverseNeighborsID: make([][]int, N),
		Flags:              make([]bool, N*K),
		Distances:          make([]float64, N*K),
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

	for I := 0; I < N; I++ {
		IdList := getKRandomNumbers(N, K, I)
		for J := 0; J < K; J++ {
			Graph.NeighborsID[I*K+J] = IdList[J]
			Graph.Distances[I*K+J] = euclideanDistance(Graph.Data[I*D:(I+1)*D], Graph.Data[IdList[J]*D:(IdList[J]+1)*D])
			Graph.ReverseNeighborsID[IdList[J]] = append(Graph.ReverseNeighborsID[IdList[J]], I)
		}
	}

}
