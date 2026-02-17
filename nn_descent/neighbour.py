# Neighbour objects are used to represent the neighbours of a point, and they store the neighbour's Vertecie object and the distance to that neighbour
from typing import Dict, List
from vertecie import Vertecie


class Neighbour: 
    def __init__(self, vert : Vertecie, distance : float):
        self.vert = vert
        self.distance = distance
        self.flag = True

    def __hash__(self):
        return hash(self.vert.id)
    def __eq__(self, other):
        return self.vert.id == other.vert.id
    def __lt__(self, other):
        return self.distance < other.distance
