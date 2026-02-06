import matplotlib.pyplot as plt
from sklearn.datasets import make_blobs
from typing import Dict, List
import random
import scipy.spatial.distance as sci


k = 10
n = 300
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
# NNG is a dictionary that maps each Vertecie to a list of its Neighbour objects
NNG: Dict[Vertecie, List[Neighbour]] = {}

#Returns a set of Neighbours neighbours
def getNeighboursNeighbour(point: Vertecie):
    listOfNeighbours = NNG[point]

    setOfNN = set()
    for neighbour in listOfNeighbours:
        for nNeigbhour in NNG[neighbour.vert]:
            if nNeigbhour.vert.id != point.id:
                setOfNN.add(nNeigbhour)
    return (setOfNN)
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
    for vert in NNG:
        setOfNN = getNeighboursNeighbour(vert)
        
        vert.candidates = sorted(NNG[vert],
                                key=lambda nb: nb.distance[0] if hasattr(nb.distance, '__iter__') else nb.distance, reverse=True)
        # For every neighbours neighbours
        for nn in setOfNN:
            # For every candidate
            for i in range (len(vert.candidates)):
                nn.distance = getDistance(vert.coordinates, nn.vert.coordinates)
                if nn.distance < vert.candidates[i].distance and nn.vert.id not in [c.vert.id for c in vert.candidates]: 
                    # print(nn.vert.id, nn.distance)
                    # print(vert.candidates[i].vert.id, vert.candidates[i].distance)
                    # print("---")
                    vert.candidates[i] = nn
                    vert.candidates.sort(key=lambda nb: nb.distance[0] if hasattr(nb.distance, '__iter__') else nb.distance, reverse=True)
                    break
    # for vert in NNG:
        NNG[vert] = vert.candidates
    return NNG
# Main function that generates the NNG, iterates over it for a number of iterations, and then draws the final NNG.
def main():
    NNG, X, y = getNNG()
    for i in range (10):
        print(i)
        NNG = iterate(NNG)

    draw(NNG, X, y)
    

main()