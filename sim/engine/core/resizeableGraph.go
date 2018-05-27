package core

import (
	"fmt"
)

// Defines a graph that supports adding / removing nodes and edges.
// Does not support edges pointing to the same node, or more than int32.max node additions / deletions
type ResizeableGraph struct {
	// If there is a node [v], maps to the other vertex [w], with cost [c]
	edges        map[int]map[int]int64
	newNodeIndex int
	vertexCount  int
}

func NewResizeableGraph() *ResizeableGraph {
	return &ResizeableGraph{
		edges:        make(map[int]map[int]int64, 0),
		newNodeIndex: 0,
		vertexCount:  0}
}

// These two methods implement graph operations as per https://github.com/yourbasic/graph/blob/master/mutable.go

// Returns the number of vertices in the graph.
func (g *ResizeableGraph) Order() int {
	return g.vertexCount
}

// Calls the do function for each neighbor w of v,
// with c equal to the cost of the edge from v to w.
// If do returns true, Visit returns immediately,
// skipping any remaining neighbors, and returns true.
//
// The iteration order is not specified and is not guaranteed
// to be the same every time.
// It is safe to delete, but not to add, edges adjacent to v
// during a call to this method.
func (g *ResizeableGraph) Visit(v int, do func(w int, c int64) bool) bool {
	if _, ok := g.edges[v]; !ok {
		panic(fmt.Sprintf("Attempted to visit node %v, which doesn't exist (newNodeIndex: %v. node count: %v)!", v, g.newNodeIndex, g.vertexCount))
	}

	for w, c := range g.edges[v] {
		// We don't delete edges when we delete nodes, so verify they still exist
		if _, ok := g.edges[w]; ok {
			if do(w, c) {
				return true
			}
		} else {
			// Delete the edge while we're traversing to ensure cleanup
			delete(g.edges[v], w)
		}
	}

	return false
}

// Returns the number of outward directed edges from v.
//func (g *ResizeableGraph) Degree(v int) int {
//	return -1 // TODO: Make this work with node deletion
// return len(g.edges[v])
//}

// Tells if there is an edge from v to w.
func (g *ResizeableGraph) Edge(v, w int) bool {
	ok := false
	if _, ok = g.edges[v]; ok {
		if _, ok = g.edges[w]; ok {
			_, ok = g.edges[v][w]
		}
	}
	return ok

}

// Returns the cost of an edge from v to w, or -1 if no such edge exists.
func (g *ResizeableGraph) Cost(v, w int) int64 {
	if g.Edge(v, w) {
		return g.edges[v][w]
	}

	return -1
}

// Adds a new node with no edges, returning the node id.
func (g *ResizeableGraph) AddNode() int {
	g.edges[g.newNodeIndex] = make(map[int]int64)
	g.newNodeIndex++
	g.vertexCount++

	return g.newNodeIndex - 1
}

// Deletes a node, doing nothing if the node was already missing
func (g *ResizeableGraph) DeleteNode(v int) {
	if _, ok := g.edges[v]; ok {
		delete(g.edges, v)
		g.vertexCount--
	}
}

// Inserts or updates a directed edge from v to w with cost c.
func (g *ResizeableGraph) AddOrUpdateEdgeCost(v, w int, c int64) {
	if v == w {
		panic(fmt.Sprintf("We do not support links to the same node. Node: %v (newNodeIndex: %v. node count: %v)!", w, g.newNodeIndex, g.vertexCount))
	}

	if _, ok := g.edges[v]; ok {
		if _, ok := g.edges[w]; ok {
			g.edges[v][w] = c
		} else {
			panic(fmt.Sprintf("Vertex '%v' does not exist (newNodeIndex: %v. node count: %v)!", w, g.newNodeIndex, g.vertexCount))
		}
	} else {
		panic(fmt.Sprintf("Vertex '%v' does not exist (newNodeIndex: %v. node count: %v)!", v, g.newNodeIndex, g.vertexCount))
	}
}

// Inserts or updates undirected edges from v to w with cost c.
func (g *ResizeableGraph) AddOrUpdateEdgeCostBoth(v, w int, c int64) {
	g.AddOrUpdateEdgeCost(v, w, c)
	g.AddOrUpdateEdgeCost(w, v, c)
}

// Delete removes an edge from v to w.

func (g *ResizeableGraph) Delete(v, w int) {
	if _, ok := g.edges[v]; ok {
		delete(g.edges[v], w)
	} else {
		panic(fmt.Sprintf("Vertex '%v' does not exist (newNodeIndex: %v. node count: %v)!", v, g.newNodeIndex, g.vertexCount))
	}
}

// DeleteBoth removes all edges between v and w.

func (g *ResizeableGraph) DeleteBoth(v, w int) {
	if v == w {
		panic(fmt.Sprintf("We do not support deleting links to the same node. Node: %v (newNodeIndex: %v. node count: %v)!", w, g.newNodeIndex, g.vertexCount))
	}

	g.Delete(v, w)
	g.Delete(w, v)
}
