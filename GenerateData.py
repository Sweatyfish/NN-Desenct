import matplotlib.pyplot as plt

from sklearn.datasets import make_blobs
X, y = make_blobs(n_samples=30, centers=3, n_features=2,
                  random_state=0)
print(X.shape)
print(X)

plt.scatter(X[:, 0], X[:, 1], c=y, cmap='viridis', s=100)
plt.xlabel('Feature 1')
plt.ylabel('Feature 2')
plt.title('Generated Clusters')
plt.colorbar(label='Cluster')
plt.show()