

n = 3 * 1000000
k = 15
Dataset = n*384*(8*8)
Neigbours_reverse = 2*(n*k*(4*8))
VertexId = n*(4*8)
Distance = n*k*(8*8)
Freeze = n*k*k*(4*8)
Boolflag = n*k*k
Pointers = n*(8*8)
print("Dataset: Gb", Dataset/8000000000)
print("Neigbour_reverse: Gb", Neigbours_reverse/8000000000)
print("VertexId: Gb", VertexId/8000000000)
print("Distance: Gb", Distance/8000000000)
print("Freeze: Gb", Freeze/8000000000)
print("Boolflag: Gb", Boolflag/8000000000)
print("Pointers: Gb", Pointers/8000000000)
print("Total Memory Usage: Gb", (Dataset + Neigbours_reverse + VertexId + Distance + Freeze + Boolflag + Pointers)/8000000000)
