package graph

import "log"

// Performs merging on the noded edges of the input geometries. Merging takes place
// on edges which are coincident (i.e. have the same coordinate list, modulo direction).
// The following situations can occur:
//		Coincident edges from different input geometries have their labels combined
//		Coincident edges from the same area geometry indicate a topology collapse.
//			In this case the topology locations are "summed" to provide a final assignment of side location
//		Coincident edges from the same linear geometry can simply be merged using the same ON location
// The merging attempts to preserve the direction of linear edges if possible (which
// is the case if there is no other coincident edge, or if all coincident edges have
// the same direction). This ensures that the overlay output line direction will be as
// consistent as possible with input lines.
// The merger also preserves the order of the edges in the input. This means that for
// polygon-line overlay the result lines will be in the same order as in the input
// (possibly with multiple result lines for a single input line).
type EdgeMerger struct {}

func (e *EdgeMerger) merge(edges []*Edge) []*Edge {
	// use a list to collect the final edges, to preserve order
	mergedEdges := make([]*Edge, 0)
	edgeMap := make(map[*EdgeKey]*Edge)

	for _, edge := range edges {
		edgeKey := new(EdgeKey)
		edgeKey.create(edge)
		baseEdge, ok := edgeMap[edgeKey]
		if !ok {
			// this is the first (and maybe only) edge for this line
			edgeMap[edgeKey] = edge
			mergedEdges = append(mergedEdges, edge)
		} else {
			// found an existing edge
			if len(baseEdge.pts) != len(edge.pts) {
				log.Printf("Merge of edges of different sizes - probable noding error.\n")
			}
			baseEdge.merge(edge) // 边缘的合并算法
		}
	}
	return mergedEdges
}
