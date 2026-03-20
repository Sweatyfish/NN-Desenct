package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func benchmark(graph Graph) float32 {
	totalCorrect := int64(0)
	totalPossible := graph.N * graph.K
	var wg sync.WaitGroup
	for i := 0; i < graph.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			foundNN := graph.NeighborsID[i*graph.K : (i+1)*graph.K]
			trueNN := getTrueNN(i)
			correct := 0
			for _, neighbor := range foundNN {
				for _, trueNeighbor := range trueNN {
					if neighborID(neighbor) == trueNeighbor {
						correct++
						break
					}
				}
			}
			atomic.AddInt64(&totalCorrect, int64(correct))
		}(i)
		if i%1000 == 0 {
			fmt.Println(i, "/", graph.N)
		}
	}
	wg.Wait()
	return float32(totalCorrect) / float32(totalPossible) * 100
}

func getTrueNN(vertex int) []int {
	truenn := make([]int, 0)
	distanceList := make([]float32, 0)
	for i := 0; i < graph.N; i++ {
		if i != vertex {
			calculatedDistance := CosineDistance(
				graph.Data[vertex*graph.Dim:(vertex+1)*graph.Dim],
				graph.Data[i*graph.Dim:(i+1)*graph.Dim],
			)
			if len(truenn) <= graph.K {
				truenn = append(truenn, i)
				distanceList = append(distanceList, calculatedDistance)
				continue
			}
			max := findmax(distanceList)
			if calculatedDistance < max[0] {
				distanceList[int(max[1])] = calculatedDistance
				truenn[int(max[1])] = i
			}
		}
	}
	return truenn
}

func findmax(slice []float32) []float32 {
	max := []float32{0, 0}
	for i, v := range slice {
		if v > max[0] {
			max[1] = float32(i)
			max[0] = v
		}
	}
	return max
}
