import h5py
import numpy as np

f = h5py.File("../../data/benchmark-dev-gooaq.h5")
np.save("train.npy", f["train"][:])