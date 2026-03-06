package main

// this function is used to benchmark the accuracy based on the graph it recieves from main
// this could also easily be multi-threaded but it would need structural changes in main "Rasmus" I might add this if necessary
func benchmark(graph Graph) float32 {
	totalCorrect := 0
	totalPossible := graph.N * graph.K

	for i := 0; i < graph.N; i++ {
		foundNN := graph.NeighborsID[i*graph.K : (i+1)*graph.K]
		trueNN := getTrueNN(i)

		for _, neighbor := range foundNN {
			for _, trueNeighbor := range trueNN {
				if neighbor.Id == trueNeighbor {
					totalCorrect++
					break
				}
			}
		}
	}

	accuracy := float32(totalCorrect) / float32(totalPossible) * 100
	return accuracy
}

func getTrueNN(vertex int) []int {
	truenn := make([]int, 0)
	distanceList := make([]float32, 0)
	for i := 0; i < graph.N; i++ {
		if i != vertex {
			calculatedDistance := euclideanDistance(graph.Data[vertex*graph.Dim:(vertex+1)*graph.Dim], graph.Data[i*graph.Dim:(i+1)*graph.Dim])
			if len(truenn) <= graph.K {
				truenn = append(truenn, i)
				distanceList = append(distanceList, calculatedDistance)
				continue
			}
			max := findmax(distanceList)
			if calculatedDistance < max[0] {
				// replace the max distance with the new one
				distanceList[int(max[1])] = calculatedDistance
				truenn[int(max[1])] = i
			}
		}
	}
	return truenn
}

// returns in format [distance,Index]
func findmax(slice []float32) []float32 {
	max := []float32{0, 0} // [placement, distance]
	for i, v := range slice {
		if v > max[0] {
			max[1] = float32(i)
			max[0] = v
		}
	}
	return max
}
