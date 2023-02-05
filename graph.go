package toposort

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var (
	// ErrCircular is raised when a cyclic relationship has been found.
	ErrCircular = errors.New("cyclic")
	// ErrMultipleRoots is raised when a graph contains multiple root nodes.
	ErrMultipleRoots = errors.New("multiple roots")
	// ErrInvalidName is raised when a name format couldn't be validated.
	ErrInvalidName = errors.New("invalid name")
)

type Vertex struct {
	afters []string
	id     string
}

// tsort sorts the given graph topologically.
func tsort(g map[string]*Vertex) (sorted []string, recursive map[string]bool, recursion []string) {
	sorted = []string{}
	visited := make(map[string]bool)
	recursive = make(map[string]bool) // keys caught in a recursive chain
	recursion = []string{}            // recursion paths for printing out in the error messages

	var visit func(id string, ancestors []string)

	visit = func(id string, ancestors []string) {
		vertex := g[id]
		if _, ok := visited[id]; ok {
			return
		}
		ancestors = append(ancestors, id)
		visited[id] = true
		for _, afterID := range vertex.afters {
			if sliceContainsString(ancestors, afterID) {
				recursive[id] = true
				for _, id := range ancestors {
					recursive[id] = true
				}
				recursion = append(recursion, append([]string{id}, ancestors...)...)
			} else {
				visit(afterID, ancestors[:])
			}
		}
		sorted = append([]string{id}, sorted...)
	}

	for k := range g {
		visit(k, []string{})
	}

	return
}

type Graph struct {
	data      map[string]*Vertex // graph itself
	ids       map[string]string  // original IDs
	sorted    []string           // toposorted keys
	recursive map[string]bool    // recursive keys
	recursion []string           // recursion paths
}

func NewGraph(data map[string]string) (*Graph, error) {
	relations, ids, err := buildRelations(data)
	if err != nil {
		return nil, err
	}

	vertices := make(map[string]*Vertex)

	for c, p := range relations {
		if _, ok := vertices[c]; !ok {
			vertices[c] = &Vertex{id: c}
		}
		if _, ok := vertices[p]; !ok {
			vertices[p] = &Vertex{id: p}
		}
		vertices[p].afters = append(vertices[p].afters, c)
	}

	g := new(Graph)
	g.sorted, g.recursive, g.recursion = tsort(vertices)
	g.ids = ids
	g.data = vertices

	if err := validateGraph(g); err != nil {
		return nil, err
	}

	return g, nil
}

// SortedIDs returns the sorted IDs in the original format.
func (g *Graph) SortedIDs() []string {
	ret := []string{}
	for _, k := range g.sorted {
		ret = append(ret, g.ids[k])
	}
	return ret
}

// validateGraph checks a graph for recursive paths and multiple root nodes.
func validateGraph(g *Graph) (err MultiError) {
	var visit func(id string)

	length := 0

	visit = func(id string) {
		v := g.data[id]
		for _, afterID := range v.afters {
			length += 1
			visit(afterID)
		}
	}

	roots := []string{}
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

	recursions, start, ko := [][]string{}, 0, ""
	for i, k := range g.recursion {
		if k == ko {
			recursions = append(recursions, g.recursion[start:i+1])
			start = i + 1
			ko = ""
		} else if ko == "" {
			ko = k
		}
	}

	// add all cyclic dependency errors to the multierror instance
	for _, xs := range recursions {
		err = append(err, fmt.Errorf("%w: %s", ErrCircular, strings.Join(xs, " -> ")))
	}

	// add multiple roots error after that if found any
	if len(roots) > 1 {
		names := []string{}
		for _, k := range roots {
			names = append(names, g.ids[k])
		}
		err = append(err, fmt.Errorf("%w: %s", ErrMultipleRoots, strings.Join(names, ", ")))
	}

	return
}

// buildRelations validates IDs in the initially provided relations map.
// It returns a new relations map with all of the IDs lowercased, and
// a 2nd index with the original IDs.
func buildRelations(data map[string]string) (map[string]string, map[string]string, error) {
	relations := make(map[string]string, len(data))
	ids := make(map[string]string)

	var err MultiError

	for k, v := range data {
		sk := strings.ToLower(k)
		sv := strings.ToLower(v)

		relations[sk] = sv

		if _, ok := ids[sk]; !ok {
			if !isAlpha(k) || len(k) < 2 { // names can only contain letters and the length must be < 2
				err = append(err, fmt.Errorf("%w: \"%s\"", ErrInvalidName, k))
			}
			ids[sk] = k
		}
		if _, ok := ids[sv]; !ok {
			if !isAlpha(v) || len(v) < 2 {
				err = append(err, fmt.Errorf("%w: \"%s\"", ErrInvalidName, v))
			}
			ids[sv] = v
		}
	}

	if err != nil {
		return nil, nil, err
	}

	return relations, ids, nil
}

func isAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func sliceContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
