# Each point in space is represented as a Vertecie object.
# Store coordinates as a NumPy array for faster numeric operations.
import numpy as np
from typing import Iterable


class Vertecie:
    def __init__(self, coordinates: Iterable[float], id):
        self.coordinates = np.asarray(coordinates)
        self.id = id

    def __hash__(self):
        return hash(self.id)

    def __eq__(self, other):
        return isinstance(other, Vertecie) and self.id == other.id