import random
from typing import Dict, List
from Benchmark import evaluate_accuracy
import time
from get_data import getData
from neighbour import Neighbour
from vertecie import Vertecie
from sklearn.decomposition import PCA
import numpy as np

bencmark_Result = True
usePCA = False
PCAdimensions = 32
k = 10
n = 4000
#sample rate
rho = 0.5

# NNG is a dictionary that maps each Vertecie to a list of its Neighbour objects
NNG: Dict[Vertecie, List[Neighbour]] = {}

#Much higher delta than the original paper, Might be dataset size migth be that we don't have reverse and or local join
#Paper delta is 0.001
Delta = 0.001
#The maximum number of iterations allowed
Iterationceiling = 50

# Gets k random points from the dataset, excluding the point itself, and returns their indices
def getKRandomPoints(k, point_index, data_list):
    """Get n random neighbor indices excluding the point itself"""
    # Use numpy choice for faster sampling on large arrays
    if k <= 0:
        return []
    n_items = len(data_list)
    all_idx = np.arange(n_items)
    other_idx = np.delete(all_idx, point_index)
    # np.random.choice returns ndarray; convert to Python list of ints
    return list(np.random.choice(other_idx, size=k, replace=False))

# Gets the distance between two points using scipy's pdist function, which computes the pairwise distance 
# between two points in a N dimension space
def getDistance(point1, point2):
    """Return Euclidean distance as a float (faster than pdist here)."""
    a = np.asarray(point1)
    b = np.asarray(point2)
    return float(np.linalg.norm(a - b))

# Generates the NNG by creating a dataset using make_blobs, initializing the Vertecie objects and their neighbours, and storing them in the NNG dictionary
def getNNG():
    data = getData(n)
    pca = PCA(n_components=PCAdimensions) 

    # Fit and transform your data
    data_reduced = pca.fit_transform(data)

    # You can also check how much variance is explained
    print(f"Explained variance ratio: {pca.explained_variance_ratio_.sum():.4f}")

    if usePCA:
        data = data_reduced
    vert_list = []

    # Runs over all points
    for i in range(len(data)):
        # Adds points to vert_list
        vert_list.append(Vertecie(data[i], i))

    # Runs over all points in vert_list and assigns them neighbours
    for vert in vert_list:
        #Gets random neigbours
        random_neighbors = getKRandomPoints(k, vert.id, vert_list)
        neighbor_list = []

        #Initializes neigbours (with hardcoded distance)
        for index in random_neighbors:
            neighbor_list.append(Neighbour(vert_list[index], getDistance(vert.coordinates, vert_list[index].coordinates)))
        NNG[vert] = neighbor_list

    return NNG

def getReverseNNG(NNG):
    R = {v: [] for v in NNG}

    for u in NNG:
        for neigh in NNG[u]:
            v = neigh.vert
            R[v].append(Neighbour(u,0.0))
    return R

def try_insert(u1_vert, u2_vert, dist):
    u1_heap = NNG[u1_vert]
    ids = {n.vert.id for n in u1_heap}
    if u2_vert.id in ids:
        return 0
    
    u3_worst = max(u1_heap, key=lambda x: x.distance)
    if dist < u3_worst.distance:
        
        u1_heap.remove(u3_worst)
        u1_heap.append(Neighbour(u2_vert, dist))
        
        RNNG[u2_vert].append(Neighbour(u1_vert,0.0))
        RNNG[u3_worst.vert].remove(next((n for n in RNNG[u3_worst.vert] if n.vert.id == u1_vert.id), None))
        return 1

    return 0
#Takes dictionary of Vertecie to list of Neighbour and a Vertecie and returns the list of Neighbour for that Vertecie
def getNeighbours(vert, NNG):
    return NNG[vert]

def sample_ref(neigh_list):
    sample_size = int(rho * len(neigh_list))
    if sample_size == 0:
        return []
    # Use numpy to pick sample indices for speed, then return the objects
    indices = np.arange(len(neigh_list))
    chosen = np.random.choice(indices, size=sample_size, replace=False)
    return_list = [neigh_list[i] for i in chosen]
    for neigh in return_list:
        neigh.flag = False
    return return_list
    
def iterate(NNG, RNNG):
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
                
        
        '''old′ ← Reverse(old), new′ ← Reverse(new)'''
        # Collect reverse neighbours by looking up RNNG for each vertex in old/new
        oldPrime = []
        newPrime = []
        try:
            for neigh in old_neighbors:
                oldPrime.extend(RNNG.get(neigh.vert, []))
            for neigh in new_neighbors:
                newPrime.extend(RNNG.get(neigh.vert, []))
        except NameError:
            # If RNNG is not in scope, fall back to empty reverse lists
            oldPrime = []
            newPrime = []

        ''' old[v] ←− old[v] ∪ Sample(old′ [v], ρK)
            new[v] ←− new[v] ∪ Sample(new′ [v], ρK)'''
        
        old_neighbors = set(old_neighbors)
        old_neighbors.update(sample_ref(oldPrime))
        new_neighbors = set(new_neighbors)
        new_neighbors.update(sample_ref(newPrime))
        
        """ c ←− c + UpdateNN(B[u1], hu2, l, true)
            c ←− c + UpdateNN(B[u2], hu1, l, true)"""
        for i in new_neighbors:
            for j in new_neighbors:
                if i.vert.id < j.vert.id:
                    dist = getDistance(i.vert.coordinates, j.vert.coordinates)
                    counter += try_insert(i.vert, j.vert, dist)
                    counter += try_insert(j.vert, i.vert, dist)
            for j in old_neighbors:
                    if i.vert.id < j.vert.id:
                        dist = getDistance(i.vert.coordinates, j.vert.coordinates)
                        counter += try_insert(i.vert, j.vert, dist)
                        counter += try_insert(j.vert, i.vert, dist)
    return NNG, counter


# Main function that generates the NNG, iterates over it for a number of iterations, and then draws the final NNG.
def main():
    original_data = getData(n)
    
    NNG = getNNG()
    global RNNG 
    RNNG = getReverseNNG(NNG)
    newNeighborsFound = 0
    iterationcounter = 0
    start = time.time()
    while True:
        NNG, newNeighborsFound = iterate(NNG,RNNG)

        if newNeighborsFound == 0 or newNeighborsFound < Delta * n * k or iterationcounter >= Iterationceiling:
            break
        print("Iteration number: ",iterationcounter + 1)
        print("new neighbors found: ", newNeighborsFound)
        iterationcounter += 1
    end = time.time()

    print("Time from start of first iteraion to end of last:")
    print(end - start)

    if bencmark_Result:
        accuracy = evaluate_accuracy(NNG, k, original_data)
        print("n:", n, "k:", k, "Dimensions:",original_data.shape[1], "Rho:", rho, "Delta:", Delta, "Iterations:", iterationcounter, "PCA:", usePCA)
        print(f"Accuracy: {accuracy:.4f}")
        print("time", end - start)

main()