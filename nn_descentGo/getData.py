import h5py
import numpy as np

# Load the HDF5 file
f = h5py.File("../../data/benchmark-dev-gooaq.h5", "r")

# Inspect the keys
print("Keys in HDF5 file:", list(f.keys()))

# Get the 'train' dataset
train_data = f["train"][:]

# Print its shape
print("Train dataset shape:", train_data.shape)

# Example: print the first 5 rows
print("First 5 entries:\n", train_data[:5])

# Example: print the first row entirely
print("First row (all 384 dimensions):\n", train_data[0])

# Example: print a specific element, e.g., row 100, dimension 50
print("Row 100, dimension 50:", train_data[100, 50])

# Example: print last row
print("Last row:", train_data[-1])
# np.save("train.npy", train_data)