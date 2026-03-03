package main

// this function is used to benchmark the accuracy based on the graph it recieves from main
//this could also easily be multi-threaded but it would need structural changes in main "Rasmus" I might add this if necessary
func benchmark(graph Graph) float64 {
	totalCorrect := 0
	totalPossible := graph.N * graph.K

	for i := 0; i < graph.N; i++ {
		foundNN := getneighbour(i)
		trueNN := getTrueNN(graph, i)

		// convert trueNN to set for fast lookup
		trueSet := make(map[int]bool)
		for _, v := range trueNN {
			trueSet[v] = true
		}

		for _, v := range foundNN {
			if trueSet[v.Id] {
				totalCorrect++
			}
		}
	}

	accuracy := float64(totalCorrect) / float64(totalPossible) * 100
	return accuracy
}

func getTrueNN(graph Graph, i int) []int {
	trueNN := make([]int, 0)
	distanceOfNN := make([]float64, 0)
	for j := 0; j < graph.N; j++ {
		if i != j {
			distance := euclideanDistance(getvertex(i), getvertex(j))
			LongestDistArr := getLongestDistance(distanceOfNN)
			if len(trueNN) < graph.K {
				trueNN = append(trueNN, j)
				distanceOfNN = append(distanceOfNN, distance)
			} else if distance < LongestDistArr[1] {
				trueNN[int(LongestDistArr[0])] = j
				distanceOfNN[int(LongestDistArr[0])] = distance
			}
		}
	}
	return trueNN
}

//For further optimization this could be a heap operation or dynamically maintained once found
//This returns a slice [placement,distance] correlating to the current neighbor which is the furthest away
func getLongestDistance(DistSlice []float64) []float64 {
	longestdist := make([]float64, 2)
	for i := 0; i < len(DistSlice); i++ {
		if DistSlice[i] > longestdist[1] {
			longestdist[1] = DistSlice[i]
			longestdist[0] = float64(i)

		}
	}
	return longestdist
}
