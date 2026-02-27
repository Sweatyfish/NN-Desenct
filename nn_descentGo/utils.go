package main

import (
	"fmt"
	"math"
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
func getvertex(V int, graph Graph) []float64 {
	return graph.Data[V*graph.Dim : (V+1)*graph.Dim]
}

/* Function to calculate the Euclidean distance between two vectors*/
func euclideanDistance(vec1, vec2 []float64) float64 {

	var total float64 = 0
	for i := range vec1 {
		diff := vec1[i] - vec2[i]
		total += diff * diff
	}
	return math.Sqrt(total)
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
func getneighbour(V int, graph Graph) []NeighborTuple {
	list := graph.NeighborsID[V*graph.K : (V+1)*graph.K]
	for i := V * graph.K; i < (V+1)*graph.K; i++ {
		graph.NeighborsID[i].Isnew = false
	}
	return list
}

func getreverseneighbour(V int, graph Graph) []NeighborTuple {
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

func tryInsert(Vertex1, Vertex2 int, distance float64) {
	/*Compare the distance with the distances of the current neighbors and update the neighbor list if necessary.*/

}
