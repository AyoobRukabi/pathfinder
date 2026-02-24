## Step 1. Extracting the Route (Path Reconstruction)
To get the actual route (e.g., Waterloo -> Victoria -> Euston) instead of just the distance,
we must track where train came from during the A* or Dijkstra search.

**Best Practice: Parent Array (`cameFrom`).**

1. Create a slice `parent := make([]int, numStations)`.
2. Whenever we update a shorter distance to neighbor `v` from current station `u`, we record it: `parent[v] = u`.
3. When we reach the target, we backtrack from the target to the start using this array to build the path.

## Step 2. Finding Multiple Routes for Trains

Since we cannot occupy the same station or track at the same time.
We need Disjoint Paths (paths that share no vertices or edges other than the start and end).

**Best Practice: Suurballe’s Algorithm (or Bhandari's Algorithm).**

[Suurballe’s Algorithm](https://en.wikipedia.org/wiki/Suurballe%27s_algorithm)

How it works:
1. It finds the shortest path, reverses the direction of its edges (making their weights negative), and runs the search again.
2. If the second search uses a reversed edge, it "cancels out" that segment,
magically untangling the routes to give you the absolute optimal pair of disjoint paths.

## Step 3. Evaluating Efficiency (Distributing `x` Trains)

Once we have `k` paths with lengths (measured in turns/edges) of `L_1, L_2, ..., L_k`,
we need to distribute `X` trains among them to minimize the total turns.

You do not simply divide the trains equally. You must balance the "load".
A train should only be sent down a longer path if the shorter path is too congested.

The formula for the number of turns `T` it takes to send `X` trains down a single path of length `L` is:
```text
T = L + X - 1
```

To optimize across multiple paths, we iterate through your paths (sorted shortest to longest).

We add paths to our active pool as long as using the next path actually reduces the overall number of turns.