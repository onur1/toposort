package toposort

import (
	"errors"
	"fmt"
)

var (
	// ErrCircular is raised when a cyclic relationship has been found.
	ErrCircular = errors.New("cyclic")
	// ErrMultipleRoots is raised when a graph contains multiple root nodes.
	ErrMultipleRoots = errors.New("multiple roots")
)

type Vertex[K comparable] struct {
	afters []K
	id     K
}

// tsort sorts the given graph topologically.
func tsort[K comparable](g map[K]*Vertex[K]) (sorted []K, recursive map[K]bool, recursion []K) {
	sorted = []K{}
	visited := make(map[K]bool)
	recursive = make(map[K]bool) // keys caught in a recursive chain
	recursion = []K{}            // recursion paths for printing out in the error messages

	var visit func(id K, ancestors []K)

	visit = func(id K, ancestors []K) {
		vertex := g[id]
		if _, ok := visited[id]; ok {
			return
		}
		ancestors = append(ancestors, id)
		visited[id] = true
		for _, afterID := range vertex.afters {
			if sliceContains(ancestors, afterID) {
				recursive[id] = true
				for _, id := range ancestors {
					recursive[id] = true
				}
				recursion = append(recursion, append([]K{id}, ancestors...)...)
			} else {
				visit(afterID, ancestors[:])
			}
		}
		sorted = append([]K{id}, sorted...)
	}

	for k := range g {
		visit(k, []K{})
	}

	return
}

type Graph[K comparable] struct {
	data      map[K]*Vertex[K] // graph itself
	sorted    []K              // toposorted keys
	recursive map[K]bool       // recursive keys
	recursion []K              // recursion paths
}

func Sort[K comparable](relations map[K]K) ([]K, error) {
	vertices := make(map[K]*Vertex[K])

	for c, p := range relations {
		if _, ok := vertices[c]; !ok {
			vertices[c] = &Vertex[K]{id: c}
		}
		if _, ok := vertices[p]; !ok {
			vertices[p] = &Vertex[K]{id: p}
		}
		vertices[p].afters = append(vertices[p].afters, c)
	}

	g := new(Graph[K])
	g.sorted, g.recursive, g.recursion = tsort(vertices)
	g.data = vertices

	if err := validateGraph(g); err != nil {
		return nil, err
	}

	return g.sorted, nil
}

// validateGraph checks a graph for recursive paths and multiple root nodes.
func validateGraph[K comparable](g *Graph[K]) (err MultiError) {
	var visit func(id K)

	length := 0

	visit = func(id K) {
		v := g.data[id]
		for _, afterID := range v.afters {
			length += 1
			visit(afterID)
		}
	}

	roots := []K{}
	var o int
	for _, id := range g.sorted {
		o = length
		length = 0
		if !g.recursive[id] { // avoid stack overflow
			visit(id)
			if length > o {
				// if the length of the dependencies is increased,
				// that means we are traversing a new tree.
				roots = append(roots, id)
			}
		}
	}

	null := *(new(K))

	recursions, start, ko := [][]K{}, 0, null
	for i, k := range g.recursion {
		if k == ko {
			recursions = append(recursions, g.recursion[start:i+1])
			start = i + 1
			ko = null
		} else if ko == null {
			ko = k
		}
	}

	// add all cyclic dependency errors to the multierror instance
	for _, xs := range recursions {
		// TODO
		// err = append(err, fmt.Errorf("%w: %s", ErrCircular, strings.Join(xs, " -> ")))
		err = append(err, fmt.Errorf("%w: %v", ErrCircular, xs))
	}

	// add multiple roots error after that if found any
	if len(roots) > 1 {
		// TODO
		// err = append(err, fmt.Errorf("%w: %s", ErrMultipleRoots, strings.Join(names, ", ")))
		err = append(err, fmt.Errorf("%w: %v", ErrMultipleRoots, roots))
	}

	return
}

func sliceContains[K comparable](s []K, e K) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
