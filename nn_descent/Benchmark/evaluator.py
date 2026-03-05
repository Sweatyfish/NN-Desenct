import numpy as np

# Returns the distance between two coordinates using NumPy (vectorized)
def getDistance(point1, point2):
    return float(np.linalg.norm(np.asarray(point1) - np.asarray(point2)))

# Evaluates the accuracy of the found NNG based on brute-force the nearest neighbors and returns the accuracy
# It takes the "current" NNG datatype and the value of k as input if NNG changes, structure the function breaks
def evaluate_accuracy(NNG, k, original_data):
    print("Evaluating accuracy...")
    vertices = list(NNG.keys())
    correct = 0
    total = 0
    length = len(vertices) / 100
    counter = 0

    # For each vertex, compute distances to all points vectorized and pick top-k
    for vertex, neighbors in NNG.items():
        counter += 1
        if (counter % 100 == 0):
            print(counter / 100, "/", length)

        # Vectorized distances from this vertex to all data points
        v = np.asarray(original_data[vertex.id])
        all_dists = np.linalg.norm(np.asarray(original_data) - v, axis=1)
        # exclude self
        all_dists[vertex.id] = np.inf

        # get indices of k smallest distances
        if k < len(all_dists):
            idx = np.argpartition(all_dists, k)[:k]
            true_nn_ids = set(int(i) for i in idx)
        else:
            true_nn_ids = set(i for i in range(len(all_dists)) if i != vertex.id)

        predicted_nn_ids = {nb.vert.id for nb in sorted(neighbors, key=lambda nb: nb.distance)[:k]}

        correct += len(predicted_nn_ids & true_nn_ids)
        total += k

    return correct / total if total > 0 else 0.0