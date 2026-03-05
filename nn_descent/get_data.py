from datasets import load_dataset
import h5py
import numpy as np

# Open the .h5 file
# f = h5py.File("../../data/benchmark-dev-gooaq.h5", "r")

# # Get the training embeddings
# train = f["train"][:]  # shape (N, 384) or (384, N)

# # Let's check the shape
# print(train.shape)
# print(train.shape[0])

# for entry in range (5):
#     print ("Entry:",entry + 1)
#     print(train[entry])

def getData(n):
    f = h5py.File("../../data/benchmark-dev-gooaq.h5", "r")
    data = f["train"][:n]
    return data

