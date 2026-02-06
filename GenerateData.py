import matplotlib.pyplot as plt
from sklearn.datasets import make_blobs
from typing import Dict, List
import random

k = 2

class Vertecie:
    def __init__(self, coordinates: List[float]):
        self.coordinates = coordinates

class Neighbour: 
    def __init__(self, vert : Vertecie, distance : float):
        self.vert = vert
        self.distance = distance

NNG: Dict[Vertecie, List[Neighbour]] = {}



def getNRandomPoints(n, point_index, data_list):
    """Get n random neighbor indices excluding the point itself"""
    other_indices = [i for i in range(len(data_list)) if i != point_index]
    return random.sample(other_indices, n)

def getDistance (point1, point2):
    return 0


def getNNG():
    X, y = make_blobs(n_samples=30, centers=3, n_features=2,
                    random_state=0)

    for i in range(len(X)):
        # Gets point and converts to Vertecie
        point_coords = X[i].tolist()
        point = Vertecie(point_coords)

        #Gets random neigbours
        random_neighbors = getNRandomPoints(k, i, X)
        neighbor_list = []

        #Initializes neigbours (with hardcoded distance)
        for idx in random_neighbors:
            neighbor_coords = X[idx].tolist()
            neighbor_vert = Vertecie(neighbor_coords)
            neighbor_list.append(Neighbour(neighbor_vert, 5.0))  # hardcoded distance of 5.0
        NNG[point] = neighbor_list

    return NNG, X, y

def draw(NNG : Dict[Vertecie, List[Neighbour]], X, y):
    # Plot all lines from points to neighbors
    for vert in NNG:
        for neighbor in NNG[vert]:
            nx, ny = neighbor.vert.coordinates
            #Draws lines
            plt.plot([vert.coordinates[0], nx], [vert.coordinates[1], ny], 'r-', alpha=0.03)
    
    #Draws points
    plt.scatter(X[:, 0], X[:, 1], c=y, cmap='viridis', s=10)
    plt.xlabel('Feature 1')
    plt.ylabel('Feature 2')
    plt.title('Generated Clusters')
    plt.colorbar(label='Cluster')
    plt.show()

def main():
    NNG, X, y = getNNG()
    draw(NNG, X, y)

main()