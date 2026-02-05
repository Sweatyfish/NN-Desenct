from datasets import load_dataset
import h5py
import numpy as np

# Open the .h5 file
f = h5py.File("benchmark-dev-gooaq.h5", "r")

# Inspect keys
print(list(f.keys()))

# train_embeddings = f["train"][:]
# print(train_embeddings)


# ds = load_dataset("sentence-transformers/gooaq")
# # print(ds)
# print(ds["train"][:5])

# Get the training embeddings
train = f["train"][:]  # shape (N, 384) or (384, N)

# Let's check the shape
print(train.shape)

# If shape is (384, N), transpose so each row is a vector
if train.shape[0] == 384:
    train = train.T

# Take the first vector
first_vector = train[0]  # shape (384,)
print("First 384-dim vector:")
print(first_vector)