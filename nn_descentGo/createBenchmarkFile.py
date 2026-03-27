import h5py
import numpy as np

f = h5py.File("../../data/allknn-benchmark-dev-gooaq.h5", "r")
knns = f['knns'][:]

# Remove self-index (first column per vector)
groundtruth_knns = knns[:, 1:16]  # columns 1 to 15
groundtruth_knns -= 1  # convert 1-based → 0-based

# Save to binary for Go
groundtruth_knns.astype(np.int32).tofile("groundtruth.i32")
print("Saved ground truth for Go benchmark:", groundtruth_knns.shape)