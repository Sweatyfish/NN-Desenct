package main

// this function is used to benchmark the accuracy based on the graph it recieves from main
func benchmark(graph Graph) {
	for i := 0; i < graph.N; i++ {
		foundNN := getneighbour(i, graph)
		trueNN := getTrueNN(graph, i)

	}

}
func getTrueNN(graph Graph, i int) []int {
	trueNN := make([]int, graph.K)
	distanceOfNN := make([]float64, graph.K)
	for j := 0; j < graph.N; j++ {
		if i != j {
			distance := euclideanDistance(getvertex(i, graph), getvertex(j, graph))
			if len(trueNN) < graph.K {
				trueNN = append(trueNN, j)
				distanceOfNN = append(distanceOfNN, distance)
			} else if distance < distanceOfNN[graph.K-1] {
				trueNN[graph.K-1] = j
				distanceOfNN[graph.K-1] = distance
			}

		}
	}
}
