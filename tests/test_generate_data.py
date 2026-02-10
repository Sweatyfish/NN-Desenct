import pytest
from nn_descent.generate_data import (
    getKRandomPoints,
    getDistance,
    getNNG,
    getNeighboursNeighbour,
    iterate,
    Vertecie,
    Neighbour,
    NNG,
)

# ------------------------
# getKRandomPoints
# ------------------------
# Hard-coded NNG for testing
def create_test_NNG():
    # Create 4 vertices with 2D coordinates
    v0 = Vertecie([0, 0], 0)
    v1 = Vertecie([0, 1], 1)
    v2 = Vertecie([1, 0], 2)
    v3 = Vertecie([1, 1], 3)

    # Manually define neighbors with distances
    NNG_test = {
        v0: [Neighbour(v1, getDistance(v0.coordinates, v1.coordinates)),
             Neighbour(v2, getDistance(v0.coordinates, v2.coordinates))],
        v1: [Neighbour(v0, getDistance(v1.coordinates, v0.coordinates)),
             Neighbour(v3, getDistance(v1.coordinates, v3.coordinates))],
        v2: [Neighbour(v0, getDistance(v2.coordinates, v0.coordinates)),
             Neighbour(v3, getDistance(v2.coordinates, v3.coordinates))],
        v3: [Neighbour(v1, getDistance(v3.coordinates, v1.coordinates)),
             Neighbour(v2, getDistance(v3.coordinates, v2.coordinates))]
    }

    return NNG_test, [v0, v1, v2, v3]


# Hard-coded 12-vertex NNG for testing
def create_test_NNG_large():
    # Create 12 vertices in a 3x4 grid
    v0 = Vertecie([0,0], 0)
    v1 = Vertecie([1,0], 1)
    v2 = Vertecie([2,0], 2)
    v3 = Vertecie([3,0], 3)
    v4 = Vertecie([0,1], 4)
    v5 = Vertecie([1,1], 5)
    v6 = Vertecie([2,1], 6)
    v7 = Vertecie([3,1], 7)
    v8 = Vertecie([0,2], 8)
    v9 = Vertecie([1,2], 9)
    v10 = Vertecie([2,2], 10)
    v11 = Vertecie([3,2], 11)

    vertices = [v0,v1,v2,v3,v4,v5,v6,v7,v8,v9,v10,v11]

    # Manually define neighbors (next 2 vertices, wrapping around)
    NNG_test = {
        v0: [Neighbour(v1,getDistance(v0.coordinates,v1.coordinates)),
             Neighbour(v2,getDistance(v0.coordinates,v2.coordinates))],

        v1: [Neighbour(v2,getDistance(v1.coordinates,v2.coordinates)),
             Neighbour(v3,getDistance(v1.coordinates,v3.coordinates))],

        v2: [Neighbour(v3,getDistance(v2.coordinates,v3.coordinates)),
             Neighbour(v4,getDistance(v2.coordinates,v4.coordinates))],

        v3: [Neighbour(v4,getDistance(v3.coordinates,v4.coordinates)),
             Neighbour(v5,getDistance(v3.coordinates,v5.coordinates))],

        v4: [Neighbour(v5,getDistance(v4.coordinates,v5.coordinates)),
             Neighbour(v6,getDistance(v4.coordinates,v6.coordinates))],

        v5: [Neighbour(v6,getDistance(v5.coordinates,v6.coordinates)),
             Neighbour(v7,getDistance(v5.coordinates,v7.coordinates))],

        v6: [Neighbour(v7,getDistance(v6.coordinates,v7.coordinates)),
             Neighbour(v8,getDistance(v6.coordinates,v8.coordinates))],

        v7: [Neighbour(v8,getDistance(v7.coordinates,v8.coordinates)),
             Neighbour(v9,getDistance(v7.coordinates,v9.coordinates))],

        v8: [Neighbour(v9,getDistance(v8.coordinates,v9.coordinates)),
             Neighbour(v10,getDistance(v8.coordinates,v10.coordinates))],

        v9: [Neighbour(v10,getDistance(v9.coordinates,v10.coordinates)),
             Neighbour(v11,getDistance(v9.coordinates,v11.coordinates))],

        v10: [Neighbour(v11,getDistance(v10.coordinates,v11.coordinates)),
              Neighbour(v0,getDistance(v10.coordinates,v0.coordinates))],

        v11: [Neighbour(v0,getDistance(v11.coordinates,v0.coordinates)),
              Neighbour(v1,getDistance(v11.coordinates,v1.coordinates))]
    }

    return NNG_test, vertices




def test_getKRandomPoints_length():
    data = list(range(10))
    result = getKRandomPoints(3, 0, data)
    assert len(result) == 3

def test_getKRandomPoints_excludes_self():
    data = list(range(10))
    result = getKRandomPoints(5, 4, data)
    assert 4 not in result


# ------------------------
# getDistance
# ------------------------

def test_getDistance_zero():
    d = getDistance([0, 0], [0, 0])
    assert d[0] == 0

def test_getDistance_known():
    d = getDistance([0, 0], [3, 4])
    assert pytest.approx(d[0]) == 5.0


# ------------------------
# getNNG
# ------------------------

def test_getNNG_structure():
    nng, X, y = getNNG()

    # number of vertices
    assert len(nng) > 0

    # each vertex has k neighbors
    for vert, neighbors in nng.items():
        assert len(neighbors) == 5   # k = 5
        for nb in neighbors:
            assert isinstance(nb, Neighbour)
            assert nb.vert != vert


# ------------------------
# getNeighboursNeighbour
# ------------------------

def test_neighbours_of_neighbours_not_self():
    nng, _, _ = getNNG()
    vert = list(nng.keys())[0]
    nn = getNeighboursNeighbour(vert)

    for n in nn:
        assert n.vert.id != vert.id


# ------------------------
# iterate
# ------------------------

def test_iterate_returns_same_size_graph():
    nng, _, _ = getNNG()
    new_nng, counter = iterate(nng)

    assert len(new_nng) == len(nng)
    assert isinstance(counter, int)

def test_neighbours_of_neighbours_hardcoded():
    # Use the hardcoded NNG
    nng, vertices = create_test_NNG()
    v0 = vertices[0]

    # Patch the global NNG used by getNeighboursNeighbour
    from nn_descent import generate_data
    generate_data.NNG = nng

    nn = getNeighboursNeighbour(v0)

    # v0's neighbors: v1, v2
    # neighbors of neighbors: v0's neighbors' neighbors = [v0,v3,v0,v3] -> excluding v0 -> {v3}
    nn_ids = {n.vert.id for n in nn}
    assert nn_ids == {3}, f"Expected neighbors-of-neighbors {3}, got {nn_ids}"

def test_neighbours_of_neighbours_large():
    nng, vertices = create_test_NNG_large()
    v0 = vertices[0]

    # Patch global NNG
    from nn_descent import generate_data
    generate_data.NNG = nng

    nn = getNeighboursNeighbour(v0)

    # v0 neighbors: vertices[1] and vertices[2]
    # neighbors-of-neighbors = union of neighbors of 1 and 2, excluding v0
    expected_ids = {3, 4}  # vertex 1 neighbors: 2,3 ; vertex 2 neighbors: 3,4 ; exclude 0 -> {3,4}
    nn_ids = {n.vert.id for n in nn}
    assert nn_ids == expected_ids, f"Expected {expected_ids}, got {nn_ids}"
