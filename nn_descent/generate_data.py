import matplotlib.pyplot as plt
from sklearn.datasets import make_blobs
from typing import Dict, List
import random
import scipy.spatial.distance as sci
from Benchmark import evaluate_accuracy
import heapq
import time

bencmark_Result = True
drawFlag = True
dimensions = 2
k = 5
n = 1000

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
        self.flag = True

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
    X, y = make_blobs(n_samples=n, centers=3, n_features=dimensions,
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

def getReverseNNG(NNG):
    R = {v: [] for v in NNG}

    for u in NNG:
        for neigh in NNG[u]:
            v = neigh.vert
            R[v].append(Neighbour(u,0.0))
    return R

def try_insert(heap, vert, dist):
    ids = {n.vert.id for n in heap}
    if vert.id in ids:
        return 0

    if len(heap) < k:
        heap.append(Neighbour(vert, dist))
        return 1

    worst = max(heap, key=lambda x: x.distance)
    if dist < worst.distance:
        heap.remove(worst)
        heap.append(Neighbour(vert, dist))
        return 1

    return 0
#Takes dictionary of Vertecie to list of Neighbour and a Vertecie and returns the list of Neighbour for that Vertecie
def getNeighbours(vert, NNG):
    return NNG[vert]

def iterate(NNG):
    R = getReverseNNG(NNG)
    counter = 0
    for vert in NNG:
        Neighbour = getNeighbours(vert, NNG)
        old_neighbors = []
        new_neighbors = []
        

        '''Mark sampled items in B[v] as false;'''
            
        for neigh in Neighbour:
            """old[v] ←− all items in B[v] with a false flag"""
            if not neigh.flag:
                old_neighbors.append(neigh)
            else:
                """new[v] ←− ρK items in B[v] with a true flag"""
                new_neighbors.append(neigh)
                neigh.flag = False
                
        
        Neighbour_rev = getNeighbours(vert, R)
        old_rev = []
        new_rev = []
        
        '''old′ ← Reverse(old), new′ ← Reverse(new)'''
        for neigh in Neighbour_rev:
            if not neigh.flag:
                old_rev.append(neigh)
            else:
                new_rev.append(neigh)
                neigh.flag = False

        ''' old[v] ←− old[v] ∪ Sample(old′ [v], ρK)
            new[v] ←− new[v] ∪ Sample(new′ [v], ρK)'''
        
        old_neighbors = set(old_neighbors)
        old_neighbors.update(old_rev)
        new_neighbors = set(new_neighbors)
        new_neighbors.update(new_rev)
        
        """ c ←− c + UpdateNN(B[u1], hu2, l, true)
            c ←− c + UpdateNN(B[u2], hu1, l, true)"""
        for i in new_neighbors:
            for j in new_neighbors:
                if i.vert.id < j.vert.id:
                    dist = getDistance(i.vert.coordinates, j.vert.coordinates)
                    counter += try_insert(NNG[i.vert], j.vert, dist)
                    counter += try_insert(NNG[j.vert], i.vert, dist)
            for j in old_neighbors:
                    if i.vert.id < j.vert.id:
                        dist = getDistance(i.vert.coordinates, j.vert.coordinates)
                        counter += try_insert(NNG[i.vert], j.vert, dist)
                        counter += try_insert(NNG[j.vert], i.vert, dist)

        # for i in range(len(general)):
        #     for j in range(i+1, len(general)):
        #         p = general[i].vert
        #         q = general[j].vert
        #         if (general[i].flag or general[j].flag):
        #             dist = getDistance(p.coordinates, q.coordinates)
        #             counter += try_insert(NNG[p], q, dist)
        #             counter += try_insert(NNG[q], p, dist)
        #     general[i].flag=False
    return NNG, counter


# Draws the NNG by plotting lines between points and their neighbours, and coloring the points according to their cluster labels
def draw(NNG : Dict[Vertecie, List[Neighbour]], X, y):
    # Plot all lines from points to neighbors
    drawLinesValue = 100/n
    if drawLinesValue > 1:
        drawLinesValue = 1
    drawCirclesValue = 10000/n
    for vert in NNG:
        for neighbor in NNG[vert]:
            nx, ny = neighbor.vert.coordinates
            #Draws lines
            plt.plot([vert.coordinates[0], nx], [vert.coordinates[1], ny], 'r-', alpha=drawLinesValue)
    
    #Draws points
    plt.scatter(X[:, 0], X[:, 1], c=y, cmap='viridis', s=drawCirclesValue)

    plt.xlabel('Feature 1')
    plt.ylabel('Feature 2')
    plt.title('Generated Clusters')
    plt.colorbar(label='Cluster')
    plt.show()
# Iterates over the NNG and updates the candidates for each point based on the distances to its neighbours and their neighbours.    
# #plt.savefig("output.png")

# Main function that generates the NNG, iterates over it for a number of iterations, and then draws the final NNG.
def main():
    NNG, X, y = getNNG()

    newNeighborsFound = 0
    iterationcounter = 0
    start = time.time()
    while True:
        NNG, newNeighborsFound = iterate(NNG)

        if newNeighborsFound == 0 or newNeighborsFound < Delta * n * k or iterationcounter >= Iterationcelling:
            break
        print("Iteration number: ",iterationcounter + 1)
        print("new neighbors found: ", newNeighborsFound)
        iterationcounter += 1
    end = time.time()

    print("Time from start of first iteraion to end of last:")
    print(end - start)

    if bencmark_Result:
        accuracy = evaluate_accuracy(NNG, k)
        print(f"Accuracy: {accuracy:.4f}")
    if drawFlag:
        draw(NNG, X, y)
    

main()