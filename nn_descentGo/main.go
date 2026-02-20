package main

/*Set Filname to the corresponding data set you want to load*/
var filename string = "data-N:1000-D:10.csv"
var filepath string = "Data/" + filename

var K = 10

type Graph struct {
	N, K, Dim          int
	Data               []float64 /*Prolly needs changing*/
	NeighborsID        []int
	ReverseNeighborsID [][]int
	Flags              []bool
	Distances          []float64
}

func main() {
	N, D := getNandDFromFilename(filename)
	initGraph(filepath, N, D, K)
}
