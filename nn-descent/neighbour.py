# Neighbour objects represent the neighbours of a point.
from vertecie import Vertecie


class Neighbour:
    def __init__(self, vert: Vertecie, distance: float):
        self.vert = vert
        self.distance = distance
        self.flag = True

    def __hash__(self):
        return hash(self.vert.id)

    def __eq__(self, other):
        return self.vert.id == other.vert.id

    def __lt__(self, other):
        return self.distance < other.distance
