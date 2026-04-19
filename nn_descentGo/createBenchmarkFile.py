import h5py
import numpy as np
import os

h5_path = "../../../../mnt/large_storage/2026_nndescent/allknn-benchmark-dev-gooaq.h5"

f = h5py.File(h5_path, "r")
knns = f['knns'][:]

# Remove self-index (first column per vector)
groundtruth_knns = knns[:, 1:16]
groundtruth_knns -= 1  # convert 1-based → 0-based

# Save next to the .h5 file
out_path = os.path.join(os.path.dirname(h5_path), "groundtruth.i32")

groundtruth_knns.astype(np.int32).tofile(out_path)

print("Saved ground truth for Go benchmark:", groundtruth_knns.shape)
print("Saved to:", out_path)
