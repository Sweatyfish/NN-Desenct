package main

import (
	"fmt"
	"math/rand"

	mapset "github.com/deckarep/golang-set/v2"
)

func getNandDFromFilename(filename string) (int, int) {
	var N, D int
	_, err := fmt.Sscanf(filename, "data-N_%d-D_%d.csv", &N, &D)
	if err != nil {
		panic(err)
	}
	return N, D
}
func getVertex(V int) []float32 {
	return graph.Data[V*graph.Dim : (V+1)*graph.Dim]
}

/* Function to calculate the Euclidean distance between two vectors*/
func euclideanDistance(vec1, vec2 []float32) float32 {

	var total float32 = 0
	for i := range vec1 {
		diff := vec1[i] - vec2[i]
		total += diff * diff
	}
	return (float32((float32(total))))
}

func CosineDistance(Vertex1, Vertex2 []float32) float32 {
	var dot float32

	for i := 0; i < graph.Dim; i++ {
		dot += Vertex1[i] * Vertex2[i]
	}

	return 1 - dot
}

func CosineDistanceBatch(VertexID int, NeighbourList []int) []float32 {
	dim := graph.Dim
	out := make([]float32, len(NeighbourList))
	a := graph.Data[VertexID*dim : (VertexID+1)*dim]
	for j := range NeighbourList {
		offset := dim * NeighbourList[j]
		b := graph.Data[offset : offset+dim]
		var sum float32

		for i := 0; i < dim; i += 8 {
			sum +=
				a[i+0]*b[i+0] +
					a[i+1]*b[i+1] +
					a[i+2]*b[i+2] +
					a[i+3]*b[i+3] +
					a[i+4]*b[i+4] +
					a[i+5]*b[i+5] +
					a[i+6]*b[i+6] +
					a[i+7]*b[i+7]
		}

		out[j] = 1 - sum
	}
	return out
}

/* Helper function to check if a slice contains a specific number*/
func contains(slice []int, num int) bool {
	for _, v := range slice {
		if v == num {
			return true
		}
	}
	return false
}

/* Helper function to get K unique random numbers from 0 to N-1, excluding Alpha*/
func getKRandomNumbers(N, K, Alpha int) []int {
	randomNumbers := make([]int, 0, K)

	for len(randomNumbers) < K {
		num := rand.Intn(N)
		if num != Alpha && !contains(randomNumbers, num) {
			randomNumbers = append(randomNumbers, num)
		}
	}
	return randomNumbers
}
func getNeighbor(V int) []NeighborTuple {
	graph.Locks[V].Lock()
	list := graph.NeighborsID[V*graph.K : (V+1)*graph.K]
	graph.Locks[V].Unlock()
	return list

}

func getReverseNeighbor(V int) []int {
	return *graph.ReverseNeighbors[V].Load()
}
func sampleKRandomNeighbors(Set mapset.Set[int], rho float32) mapset.Set[int] {
	SampledSet := mapset.NewSet[int]()
	for neighbor := range Set.Iter() {
		if rand.Float32() < rho {
			SampledSet.Add(neighbor)
		}
	}
	return SampledSet
}

func tryInsert(Vertex1, Vertex2 int, distance float32) int {
	if Vertex1 == Vertex2 {
		return 0
	}
	skipvertex1 := false
	skipvertex2 := false

	//Check to find the current neighbor with the longest distance for both vertices
	//Type of this variable is [neighborID, Distance, Placement in neighbor list]
	LongestNeighborVertex1 := make([]float32, 3)
	LongestNeighborVertex2 := make([]float32, 3)
	inserted := 0
	for i := Vertex1 * graph.K; i < (Vertex1+1)*graph.K; i++ {
		if graph.NeighborsID[i].Id == Vertex2 {
			skipvertex1 = true
			break
		}
		if LongestNeighborVertex1[1] == 0.0 || graph.Distances[i] > LongestNeighborVertex1[1] {
			LongestNeighborVertex1[0] = float32(graph.NeighborsID[i].Id)
			LongestNeighborVertex1[1] = graph.Distances[i]
			LongestNeighborVertex1[2] = float32(i)
		}
	}
	for i := Vertex2 * graph.K; i < (Vertex2+1)*graph.K; i++ {
		if graph.NeighborsID[i].Id == Vertex1 {
			skipvertex2 = true
			break
		}
		if LongestNeighborVertex2[1] == 0.0 || graph.Distances[i] > LongestNeighborVertex2[1] {
			LongestNeighborVertex2[0] = float32(graph.NeighborsID[i].Id)
			LongestNeighborVertex2[1] = graph.Distances[i]
			LongestNeighborVertex2[2] = float32(i)
		}
	}

	//If the new distance is smaller than the longest distance, we can insert the new neighbor
	if !skipvertex1 && distance < LongestNeighborVertex1[1] {
		//Update the reverse neighbors of the removed neighbor and the new neighbor for vertex 1
		removeReverseNeighbor(Vertex1, int(LongestNeighborVertex1[0]))
		InsertNewReverseNeighbor(Vertex1, Vertex2)
		//We need to lock the vertex before modifying its neighbors
		graph.Locks[Vertex1].Lock()
		//Insert the new neighbor
		graph.NeighborsID[int(LongestNeighborVertex1[2])] = NeighborTuple{isNew: true, Id: Vertex2}
		//Insert the new distance
		graph.Distances[int(LongestNeighborVertex1[2])] = distance
		//Unlock the vertex after modification
		graph.Locks[Vertex1].Unlock()
		inserted++

	}
	//Same Process for the second vertex
	if !skipvertex2 && distance < LongestNeighborVertex2[1] {
		removeReverseNeighbor(Vertex2, int(LongestNeighborVertex2[0]))
		InsertNewReverseNeighbor(Vertex2, Vertex1)
		graph.Locks[Vertex2].Lock()
		graph.NeighborsID[int(LongestNeighborVertex2[2])] = NeighborTuple{isNew: true, Id: Vertex1}
		graph.Distances[int(LongestNeighborVertex2[2])] = distance
		graph.Locks[Vertex2].Unlock()
		inserted++

	}
	return inserted
}

// Removes vertex1 from the reverse neighbor list of vertex2
func removeReverseNeighbor(Vertex1, Vertex2 int) {
	//We need to lock the vertex before modifying its reverse neighbors
	graph.Locks[Vertex2].Lock()
	revPointer := graph.FreezeReverseNeighbors[Vertex2].Load()
	for i, neighbor := range *revPointer {
		if neighbor == Vertex1 {
			//Remove the neighbor by swapping it with the last element and truncating the slice
			(*revPointer)[i] = (*revPointer)[len(*revPointer)-1]
			*revPointer = (*revPointer)[:len(*revPointer)-1]

			break
		}
	}
	graph.FreezeReverseNeighbors[Vertex2].Store(revPointer)
	graph.Locks[Vertex2].Unlock()
}

// Inserts vertex1 into the reverse neighbor list of vertex2
func InsertNewReverseNeighbor(Vertex1, Vertex2 int) {
	//We need to lock the vertex before modifying its reverse neighbors
	graph.Locks[Vertex2].Lock()
	//Insert the new neighbor by appending it to the slice
	revPointer := graph.FreezeReverseNeighbors[Vertex2].Load()
	*revPointer = append(*revPointer, Vertex1)
	graph.FreezeReverseNeighbors[Vertex2].Store(revPointer)
	//Unlock the vertex after modification
	graph.Locks[Vertex2].Unlock()
}
