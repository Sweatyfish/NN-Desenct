package main

import (
	"fmt"
	"math"
	"math/rand"
)

func getNandDFromFilename(filename string) (int, int) {
	var N, D int
	_, err := fmt.Sscanf(filename, "data-N:%d-D:%d.csv", &N, &D)
	if err != nil {
		panic(err)
	}
	return N, D
}

func euclideanDistance(vec1, vec2 []float64) float64 {

	var total float64 = 0
	for i := range vec1 {
		diff := vec1[i] - vec2[i]
		total += diff * diff
	}
	return math.Sqrt(total)
}

func contains(slice []int, num int) bool {
	for _, v := range slice {
		if v == num {
			return true
		}
	}
	return false
}

func getKRandomNumbers(N, K, Alpha int) []int {
	randomNumbers := make([]int, K)
	for len(randomNumbers) < K {
		num := rand.Intn(N)
		if !contains(randomNumbers, num) && num != Alpha {
			randomNumbers = append(randomNumbers, num)
		}
	}
	return randomNumbers
}
