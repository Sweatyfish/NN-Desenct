import os
import h5py
import numpy as np

# Load the HDF5 file
h5_path = "../../../../mnt/large_storage/2026_nndescent/benchmark-dev-gooaq.h5"
f = h5py.File(h5_path, "r")


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

out_path = os.path.join(os.path.dirname(h5_path), "train.npy")
np.save(out_path, train_data)
print("Saved to:", out_path)
