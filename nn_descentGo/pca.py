from sklearn.decomposition import PCA
import numpy as np

X = np.load("../../data/train.npy")

# normalize first

X = X / np.linalg.norm(X, axis=1, keepdims=True)
# Write amount of dimensions wanted
D = 128

pca = PCA(n_components=D)
X_reduced = pca.fit_transform(X)

# normalize again
X_reduced = X_reduced / np.linalg.norm(X_reduced, axis=1, keepdims=True)
# change the number after 'reduced'
np.save(f"../../data/reduced_{D}.npy", X_reduced)

