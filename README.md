# NN-Descent

This project implements a variation of the **NN-Descent** algorithm for approximate k nearest neighbor search.
This repo contains both a working python implementation and a golang implementation, all aimed at solving the Sisap 2025 Task 2

This is part of an academical product currently being develop in colaberation with students and supervisers at the ITU-University of Copenhagen

go tool pprof cpu.prof
(pprof) top
(pprof) list NNDecent


var checkMemory = true
if checkMemory {
		f, err := os.Create("mem.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		runtime.GC() // important: get up-to-date heap
		pprof.WriteHeapProfile(f)
	}

go tool pprof mem.prof
top

Currently, we use:

- `matplotlib.pyplot` to visualize the graph.
- `make_blobs` from `sklearn` to generate synthetic data points.
- `Bencmark` Custom package the measures accuracy based on brute-force
- `scipy.spatial.distance` to measure "distance" between two points
- `numpy` for benchmark and general python improvements

You need to run:

- `pip install pytest pytest-cov`
- `pip install scikit-learn`
- `pip install matplotlib`
- `pip install scipy`
- `pip install numpy`

