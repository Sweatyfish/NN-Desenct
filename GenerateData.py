import matplotlib.pyplot as plt
from sklearn.datasets import make_blobs
from typing import Dict, List
import random

k = 5

Neighbours: Dict[int, List[tuple]] = {}


def getNRandomPoints(n, point_index, data_list):
    """Get n random neighbor indices excluding the point itself"""
    other_indices = [i for i in range(len(data_list)) if i != point_index]
    return random.sample(other_indices, n)


X, y = make_blobs(n_samples=30, centers=3, n_features=2,
                  random_state=0)
# print(X.shape)
# print(X)

for i in range(len(X)):
    point = X[i]
    random_neighbors = getNRandomPoints(k, i, X)
    Neighbours[i] = [X[idx] for idx in random_neighbors]
    print(f"Point {i}: {point}")
print(Neighbours)

for i in range(len(X)):
    point = X[i]
    neighbors = Neighbours[i]
    for neighbor in neighbors:
        plt.plot([point[0], neighbor[0]], [point[1], neighbor[1]], 'r-', alpha=0.5)

plt.scatter(X[:, 0], X[:, 1], c=y, cmap='viridis', s=100)
plt.xlabel('Feature 1')
plt.ylabel('Feature 2')
plt.title('Generated Clusters')
plt.colorbar(label='Cluster')
plt.show()