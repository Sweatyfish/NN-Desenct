# Each point in space is represented as a Vertecie object, which has a list of its neighbours
# (as Neighbour objects) and a list of candidates (also as Neighbour objects)
from typing import List


class Vertecie:
    def __init__(self, coordinates: List[float], id):
        self.coordinates = coordinates
        self.id = id
    
    def __hash__(self):
        return hash(self.id)

    def __eq__(self, other):
        return isinstance(other, Vertecie) and self.id == other.id