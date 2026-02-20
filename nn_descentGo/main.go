package main

type Vertices struct {
	coordinate []float64 /*64 bit numbers for now (Prolly needs changing)*/
	Id         int
	neighbors  []int
}
