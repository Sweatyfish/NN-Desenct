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
	// 0  = 16  dimension
// 1  = 32  dimension
// 2  = 64  dimension
// 3  = 96  dimension
// 4  = 128 dimension
// 5  = 160 dimension
// 6  = 192 dimension
// 7  = 224 dimension
// 8  = 256 dimension
// 9  = 288 dimension
// 10 = 320 dimension
// 11 = 352 dimension
// 12 = 384 dimension
	var (
		data []float32
		err  error
	)

	switch PCAcase {
	case 0:
		Dimensions = int32(16)
		data, err = loadNpyFirstN("../../../../mnt/large_storage/2026_nndescent/reduced_16.npy", N, Dimensions)
	case 1:
		Dimensions = int32(32)
		data, err = loadNpyFirstN("../../../../mnt/large_storage/2026_nndescent/reduced_32.npy", N, Dimensions)
	case 2:
		Dimensions = int32(64)
		data, err = loadNpyFirstN("../../../../mnt/large_storage/2026_nndescent/reduced_64.npy", N, Dimensions)
	case 3:
		Dimensions = int32(96)
		data, err = loadNpyFirstN("../../../../mnt/large_storage/2026_nndescent/reduced_96.npy", N, Dimensions)
	case 4:
		Dimensions = int32(128)
		data, err = loadNpyFirstN("../../../../mnt/large_storage/2026_nndescent/reduced_128.npy", N, Dimensions)
	case 5:
		Dimensions = int32(160)
		data, err = loadNpyFirstN("../../../../mnt/large_storage/2026_nndescent/reduced_160.npy", N, Dimensions)
	case 6:
		Dimensions = int32(192)
		data, err = loadNpyFirstN("../../../../mnt/large_storage/2026_nndescent/reduced_192.npy", N, Dimensions)
	case 7:
		Dimensions = int32(224)
		data, err = loadNpyFirstN("../../../../mnt/large_storage/2026_nndescent/reduced_224.npy", N, Dimensions)
	case 8:
		Dimensions = int32(256)
		data, err = loadNpyFirstN("../../../../mnt/large_storage/2026_nndescent/reduced_256.npy", N, Dimensions)
	case 9:
		Dimensions = int32(288)
		data, err = loadNpyFirstN("../../../../mnt/large_storage/2026_nndescent/reduced_288.npy", N, Dimensions)
	case 10:
		Dimensions = int32(320)
		data, err = loadNpyFirstN("../../../../mnt/large_storage/2026_nndescent/reduced_320.npy", N, Dimensions)
	case 11:
		Dimensions = int32(352)
		data, err = loadNpyFirstN("../../../../mnt/large_storage/2026_nndescent/reduced_352.npy", N, Dimensions)
	case 12:
		Dimensions = int32(384)
		data, err = loadNpyFirstN("../../../../mnt/large_storage/2026_nndescent/train.npy", N, Dimensions)
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
