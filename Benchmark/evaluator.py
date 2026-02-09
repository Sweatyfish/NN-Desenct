import scipy.spatial.distance as sci
import heapq

# Returns the distance between two coordinates using scipy's pdist function
def getDistance (point1, point2):
    Dist = sci.pdist([point1, point2], 'euclidean')
    # print (Dist)
    return Dist

# Evaluates the accuracy of the found NNG based on brute-force the nearest neighbors and returns the accuracy
# It takes the "current" NNG datatype and the value of k as input if NNG changes, structure the function breaks
def evaluate_accuracy(NNG, k):
    print("Evaluating accuracy...")
    vertices = list(NNG.keys())
    correct = 0
    total = 0


    #This head bruteforces the nearest neighbors for each vertex and compares them to the predicted neighbors in the NNG
    for vertex, neighbors in NNG.items():
        heap = []

        for other in vertices:
            if other.id == vertex.id:
                continue
            dist = getDistance(vertex.coordinates, other.coordinates)

            if len(heap) < k:
                heapq.heappush(heap, (-dist, other.id))
            else:

                if dist < -heap[0][0]:
                    heapq.heappushpop(heap, (-dist, other.id))

        true_nn_ids = {nid for _, nid in heap}

        predicted_nn_ids = {nb.vert.id for nb in sorted(neighbors, key=lambda nb: nb.distance)[:k]}

        #Take the intersection of the predicted and true nearest neighbors and count how many are correct
        correct += len(predicted_nn_ids & true_nn_ids)
        #The total number of neighbors evaluated is the number of vertices times k
        total += k 

    return correct / total if total > 0 else 0.0