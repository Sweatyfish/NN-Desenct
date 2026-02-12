import matplotlib.pyplot as plt
from sklearn.datasets import make_blobs
from typing import Dict, List
import random
import scipy.spatial.distance as sci
from Benchmark import evaluate_accuracy
import heapq

bencmark_Result = False
k = 4
n = 25000
#Much higher delta than the original paper, Might be dataset size migth be that we don't have reverse and or local join
#Paper delta is 0.001
Delta = 0.02
#The maximum number of iterations allowed
Iterationcelling = 50

def _heapreplace_max(heap, item):
    """Maxheap version of a heappop followed by a heappush."""
    returnitem = heap[0]    # raises appropriate IndexError if heap is empty
    heap[0] = item
    heapq._siftup_max(heap, 0)
    return returnitem

# Each point in space is represented as a Vertecie object, which has a list of its neighbours
# (as Neighbour objects) and a list of candidates (also as Neighbour objects)
class Vertecie:
    def __init__(self, coordinates: List[float], id):
        self.coordinates = coordinates
        self.id = id
        self.candidates = []
# Neighbour objects are used to represent the neighbours of a point, and they store the neighbour's Vertecie object and the distance to that neighbour
class Neighbour: 
    def __init__(self, vert : Vertecie, distance : float):
        self.vert = vert
        self.distance = distance

    def __hash__(self):
        return hash(self.vert.id)
    def __eq__(self, other):
        return self.vert.id == other.vert.id
    def __lt__(self, other):
        return self.distance < other.distance
# NNG is a dictionary that maps each Vertecie to a list of its Neighbour objects
NNG: Dict[Vertecie, List[Neighbour]] = {}

#Returns a set of Neighbours neighbours
def getNeighboursNeighbour(point: Vertecie):
    direct_neighbors = NNG[point]
    direct_ids = {n.vert.id for n in direct_neighbors}

    setOfNN = set()
    for neighbour in direct_neighbors:
        for nn in NNG[neighbour.vert]:
            if nn.vert.id != point.id and nn.vert.id not in direct_ids:
                setOfNN.add(nn)

    return setOfNN

# Gets k random points from the dataset, excluding the point itself, and returns their indices
def getKRandomPoints(k, point_index, data_list):
    """Get n random neighbor indices excluding the point itself"""
    other_indices = [i for i in range(len(data_list)) if i != point_index]
    return random.sample(other_indices, k)
# Gets the distance between two points using scipy's pdist function, which computes the pairwise distance 
# between two points in a N dimension space
def getDistance (point1, point2):
    Dist = sci.pdist([point1, point2], 'euclidean')
    # print (Dist)
    return Dist

# Generates the NNG by creating a dataset using make_blobs, initializing the Vertecie objects and their neighbours, and storing them in the NNG dictionary
def getNNG():
    X, y = make_blobs(n_samples=n, centers=3, n_features=2,
                    random_state=0)
    vert_list = []

    # Runs over all points
    for i in range(len(X)):
        # Adds points to vert_list
        temp = X[i].tolist()
        vert_list.append(Vertecie(temp, i))

    # Runs over all points in vert_list and assigns them neighbours
    for vert in vert_list:
        #Gets random neigbours
        random_neighbors = getKRandomPoints(k, vert.id, vert_list)
        # print(random_neighbors)
        neighbor_list = []

        #Initializes neigbours (with hardcoded distance)
        for index in random_neighbors:
            # print("looking at points", vert.coordinates, vert_list[index].coordinates)
            neighbor_list.append(Neighbour(vert_list[index], getDistance(vert.coordinates, vert_list[index].coordinates)))
        NNG[vert] = neighbor_list

    return NNG, X, y
# Draws the NNG by plotting lines between points and their neighbours, and coloring the points according to their cluster labels
def draw(NNG : Dict[Vertecie, List[Neighbour]], X, y):
    # Plot all lines from points to neighbors
    for vert in NNG:
        for neighbor in NNG[vert]:
            nx, ny = neighbor.vert.coordinates
            #Draws lines
            plt.plot([vert.coordinates[0], nx], [vert.coordinates[1], ny], 'r-', alpha=0.1)
    
    #Draws points
    plt.scatter(X[:, 0], X[:, 1], c=y, cmap='viridis', s=10)

    plt.xlabel('Feature 1')
    plt.ylabel('Feature 2')
    plt.title('Generated Clusters')
    plt.colorbar(label='Cluster')
    plt.show()
# Iterates over the NNG and updates the candidates for each point based on the distances to its neighbours and their neighbours.    #plt.savefig("output.png")


def iterate(NNG):
    counter = 0
    for vert in NNG:
        setOfNN = getNeighboursNeighbour(vert)
        
        # Initialize candidates from current neighbors
        vert.candidates = list(NNG[vert])
        # Heapify in-place (returns None, modifies vert.candidates)
        heapq._heapify_max(vert.candidates)
        
        # Save the original neighbors to get a count of newly added neighbors
        original_ids = {c.vert.id for c in vert.candidates}
        candidate_ids = set(original_ids)
        
        # For every neighbours neighbours
        for nn in setOfNN:
            # Skip if already a candidate
            if nn.vert.id in candidate_ids:
                continue    

            dist = getDistance(vert.coordinates, nn.vert.coordinates)
            new_nn = Neighbour(nn.vert, dist)

            if dist < vert.candidates[0].distance:
                worst = vert.candidates[0]
                _heapreplace_max(vert.candidates, new_nn)
                candidate_ids.discard(worst.vert.id)
                candidate_ids.add(new_nn.vert.id)


        # Here newly added neighbors are counted by taking the difference between the original neighbors and the new candidates
        new_ids = {c.vert.id for c in vert.candidates}
        new_neighbors = len(new_ids - original_ids)
        counter += new_neighbors
    # Unfrezze graph and set new neighbors
    for vert in NNG:
        NNG[vert] = vert.candidates
    return NNG, counter
# Main function that generates the NNG, iterates over it for a number of iterations, and then draws the final NNG.
def main():
    NNG, X, y = getNNG()

    newNeighborsFound = 0
    iterationcounter = 0
    while True:
        NNG, newNeighborsFound = iterate(NNG)

        if newNeighborsFound == 0 or newNeighborsFound < Delta * n * k or iterationcounter >= Iterationcelling:
            break
        print("Iteration number: ",iterationcounter)
        print("new neighbors found: ", newNeighborsFound)
        iterationcounter += 1

    if bencmark_Result:
        accuracy = evaluate_accuracy(NNG, k)
        print(f"Accuracy: {accuracy:.4f}")
    # draw(NNG, X, y)
    

main()