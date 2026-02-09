# NN-Descent

This project implements a variation of the **NN-Descent** algorithm for approximate nearest neighbor search.

Currently, we use:

- `matplotlib.pyplot` to visualize the graph.
- `make_blobs` from `sklearn` to generate synthetic data points.
- `Bencmark` Custom package the measures accuracy based on brute-force
- `scipy.spatial.distance` to measure "distance" between two points

Still missing from paper
`Local Join`
`Incremental Search`
`Sampling`
`Iteration optimization` with heap or something simular

## Version 1

In **version 1**, we implemented a very simple and poorly optimized version of the algorithm.  
It runs for a fixed number of iterations without any advanced optimization or convergence criteria.

You can find performance benchmarks and comparisons with other versions in  
[`benchmarks.md`](./benchmarks.md).
