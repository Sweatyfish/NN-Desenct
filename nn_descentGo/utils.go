package main

import (
	"fmt"
	"math/rand"

	mapset "github.com/deckarep/golang-set/v2"
)

func getNandDFromFilename(filename string) (int32, int32) {
	var N, Dimensions int32
	_, err := fmt.Sscanf(filename, "data-N_%Dimensions-D_%Dimensions.csv", &N, &Dimensions)
	if err != nil {
		panic(err)
	}
	return N, Dimensions
}

func getVertex(V int32) []float32 {
	return graph.Data[V*Dimensions : (V+1)*Dimensions]
}

func CosineDistance(Vertex1, Vertex2 []float32) float32 {
	var dot float32
	for i := 0; i < len(Vertex1); i++ {
		dot += Vertex1[i] * Vertex2[i]
	}
	return 1 - dot
}

func CosineDistanceBatchN(NeighbourList []int32) []float32 {

	n := len(NeighbourList)
	out := make([]float32, n*(n-1)/2)
	idx := 0
	for i := 0; i < n; i++ {
		offseta := Dimensions * NeighbourList[i]
		a := graph.Data[offseta : offseta+Dimensions]
		for j := i + 1; j < n; j++ {
			offsetb := Dimensions * NeighbourList[j]
			b := graph.Data[offsetb : offsetb+Dimensions]
			var sum float32
			for l := int32(0); l < Dimensions; l += 8 {
				sum +=
					a[l+0]*b[l+0] +
						a[l+1]*b[l+1] +
						a[l+2]*b[l+2] +
						a[l+3]*b[l+3] +
						a[l+4]*b[l+4] +
						a[l+5]*b[l+5] +
						a[l+6]*b[l+6] +
						a[l+7]*b[l+7]
			}
			out[idx] = 1 - sum
			idx++
		}
	}
	return out
}

func CosineDistanceBatchNM(NewNeighbourlist []int32, OldNeighbour []int32) []float32 {
	Dimensions := Dimensions
	n := len(NewNeighbourlist)
	m := len(OldNeighbour)
	out := make([]float32, n*m)
	idx := 0
	for i := 0; i < n; i++ {
		offseta := Dimensions * NewNeighbourlist[i]
		a := graph.Data[offseta : offseta+Dimensions]
		for j := 0; j < m; j++ {
			offsetb := Dimensions * OldNeighbour[j]
			b := graph.Data[offsetb : offsetb+Dimensions]
			var sum float32
			for l := int32(0); l < Dimensions; l += 8 {
				sum +=
					a[l+0]*b[l+0] +
						a[l+1]*b[l+1] +
						a[l+2]*b[l+2] +
						a[l+3]*b[l+3] +
						a[l+4]*b[l+4] +
						a[l+5]*b[l+5] +
						a[l+6]*b[l+6] +
						a[l+7]*b[l+7]
			}
			out[idx] = 1 - sum
			idx++
		}
	}
	return out
}

func contains(slice []int32, num int32) bool {
	for _, v := range slice {
		if v == num {
			return true
		}
	}
	return false
}

func getKRandomNumbers(N, K, Alpha int32) []int32 {
	randomNumbers := make([]int32, int32(0), K)
	for int32(len(randomNumbers)) < K {
		num := int32(rand.Intn(int(N)))
		if num != Alpha && !contains(randomNumbers, num) {
			randomNumbers = append(randomNumbers, num)
		}
	}
	return randomNumbers
}

func getNeighbor(V int32) []NeighborTuple {
	graph.Locks[V].Lock()
	list := graph.NeighborsID[V*graph.K : (V+1)*graph.K]
	graph.Locks[V].Unlock()
	return list
}

func getReverseNeighbor(V int32) []int32 {
	return *graph.ReverseNeighbors[V].Load()
}

func sampleKRandomNeighbors(Set mapset.Set[int32], rho float32) mapset.Set[int32] {
	SampledSet := mapset.NewSet[int32]()
	for neighbor := range Set.Iter() {
		if rand.Float32() < rho {
			SampledSet.Add(neighbor)
		}
	}
	return SampledSet
}

func getWorstNeighborInfo(vertex int32) neighborInfo {
	var worstN neighborInfo
	for i := vertex * graph.K; i < (vertex+1)*graph.K; i++ {
		if graph.Distances[i] > worstN.distance {
			worstN.distance = graph.Distances[i]
			worstN.id = neighborID(graph.NeighborsID[i])
			worstN.index = i
		}
	}
	return worstN
}

func getWorstNeighborInfoBatch(VertexList []int32) []neighborInfo {
	worstNList := make([]neighborInfo, len(VertexList))
	for i, v := range VertexList {
		graph.Locks[v].Lock()
		start := v * graph.K
		end := (v + 1) * graph.K
		var worst neighborInfo
		for j := start; j < end; j++ {
			if graph.Distances[j] > worst.distance {
				worst.distance = graph.Distances[j]
				worst.id = neighborID(graph.NeighborsID[j])
				worst.index = j
			}
		}
		graph.Locks[v].Unlock()
		worstNList[i] = worst
	}
	return worstNList
}

func insert(v1Id, v2Id int32, nInfo neighborInfo, distance float32) (int32, neighborInfo) {
	if v1Id == v2Id {
		return 0, nInfo
	}
	var secondWorst neighborInfo
	for i := v1Id * graph.K; i < (v1Id+1)*graph.K; i++ {
		if v2Id == neighborID(graph.NeighborsID[i]) {
			return 0, nInfo
		}
		if secondWorst.distance < graph.Distances[i] && nInfo.id != neighborID(graph.NeighborsID[i]) {
			secondWorst.distance = graph.Distances[i]
			secondWorst.id = neighborID(graph.NeighborsID[i])
			secondWorst.index = i
		}
	}
	if secondWorst.distance < distance {
		secondWorst.id = v2Id
		secondWorst.distance = distance
		secondWorst.index = nInfo.index
	}
	graph.Locks[v1Id].Lock()
	graph.NeighborsID[nInfo.index] = makeNeighbor(v2Id, true)
	graph.Distances[nInfo.index] = distance
	graph.Locks[v1Id].Unlock()
	removeReverseNeighbor(v1Id, nInfo.id)
	InsertNewReverseNeighbor(v1Id, v2Id)
	return 1, secondWorst
}

func insertNoreturn(v1Id, v2Id int32, distance float32) int32 {
	if v1Id == v2Id {
		return 0
	}
	var Worst neighborInfo
	for i := v1Id * graph.K; i < (v1Id+1)*graph.K; i++ {
		if v2Id == neighborID(graph.NeighborsID[i]) {
			return 0
		}
		if Worst.distance < graph.Distances[i] {
			Worst.distance = graph.Distances[i]
			Worst.id = neighborID(graph.NeighborsID[i])
			Worst.index = i
		}
	}
	if Worst.distance > distance {
		graph.Locks[v1Id].Lock()
		graph.NeighborsID[Worst.index] = makeNeighbor(v2Id, true)
		graph.Distances[Worst.index] = distance
		graph.Locks[v1Id].Unlock()
		removeReverseNeighbor(v1Id, Worst.id)
		InsertNewReverseNeighbor(v1Id, v2Id)
		return 1
	}
	return 0
}

func tryInsert(Vert1 int32, Vertex2 int32, distance float32) int32 {
	Vertex1 := Vert1
	if Vertex1 == Vertex2 {
		return int32(0)
	}
	skipvertex1 := false
	skipvertex2 := false
	LongestNeighborVertex1 := make([]float32, 3)
	LongestNeighborVertex2 := make([]float32, 3)
	inserted := int32(0)
	for i := Vertex1 * graph.K; i < (Vertex1+1)*graph.K; i++ {
		if neighborID(graph.NeighborsID[i]) == Vertex2 {
			skipvertex1 = true
			break
		}
		if LongestNeighborVertex1[1] == 0.0 || graph.Distances[i] > LongestNeighborVertex1[1] {
			LongestNeighborVertex1[0] = float32(neighborID(graph.NeighborsID[i]))
			LongestNeighborVertex1[1] = graph.Distances[i]
			LongestNeighborVertex1[2] = float32(i)
		}
	}
	for i := Vertex2 * graph.K; i < (Vertex2+1)*graph.K; i++ {
		if neighborID(graph.NeighborsID[i]) == Vertex1 {
			skipvertex2 = true
			break
		}
		if LongestNeighborVertex2[1] == 0.0 || graph.Distances[i] > LongestNeighborVertex2[1] {
			LongestNeighborVertex2[0] = float32(neighborID(graph.NeighborsID[i]))
			LongestNeighborVertex2[1] = graph.Distances[i]
			LongestNeighborVertex2[2] = float32(i)
		}
	}
	if !skipvertex1 && distance < LongestNeighborVertex1[1] {
		removeReverseNeighbor(Vertex1, int32(LongestNeighborVertex1[0]))
		InsertNewReverseNeighbor(Vertex1, Vertex2)
		graph.Locks[Vertex1].Lock()
		graph.NeighborsID[int32(LongestNeighborVertex1[2])] = makeNeighbor(Vertex2, true)
		graph.Distances[int32(LongestNeighborVertex1[2])] = distance
		graph.Locks[Vertex1].Unlock()
		inserted++
	}
	if !skipvertex2 && distance < LongestNeighborVertex2[1] {
		removeReverseNeighbor(Vertex2, int32(LongestNeighborVertex2[0]))
		InsertNewReverseNeighbor(Vertex2, Vertex1)
		graph.Locks[Vertex2].Lock()
		graph.NeighborsID[int32(LongestNeighborVertex2[2])] = makeNeighbor(Vertex1, true)
		graph.Distances[int32(LongestNeighborVertex2[2])] = distance
		graph.Locks[Vertex2].Unlock()
		inserted++
	}
	return inserted
}

func removeReverseNeighbor(Vertex1, Vertex2 int32) {
	graph.Locks[Vertex2].Lock()
	revPointer := graph.ReverseNeighbors[Vertex2].Load()
	for i, neighbor := range *revPointer {
		if neighbor == Vertex1 {
			(*revPointer)[i] = (*revPointer)[len(*revPointer)-1]
			*revPointer = (*revPointer)[:len(*revPointer)-1]
			break
		}
	}
	graph.ReverseNeighbors[Vertex2].Store(revPointer)
	graph.Locks[Vertex2].Unlock()
}

func InsertNewReverseNeighbor(Vertex1, Vertex2 int32) {
	graph.Locks[Vertex2].Lock()
	revPointer := graph.ReverseNeighbors[Vertex2].Load()
	*revPointer = append(*revPointer, Vertex1)
	graph.ReverseNeighbors[Vertex2].Store(revPointer)
	graph.Locks[Vertex2].Unlock()
}
