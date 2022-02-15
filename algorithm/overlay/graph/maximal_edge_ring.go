package graph

import "log"

const (
	STATE_FIND_INCOMING = 1
	STATE_LINK_OUTGOING = 2
)

type MaximalEdgeRing struct {
	startEdge *OverlayEdge
}

// NewMaximalEdgeRing...
func NewMaximalEdgeRing(e *OverlayEdge) *MaximalEdgeRing {
	maximalEdgeRing := new(MaximalEdgeRing)
	maximalEdgeRing.attachEdges(e)
	return maximalEdgeRing
}

// attachEdges...
func (m *MaximalEdgeRing) attachEdges(startEdge *OverlayEdge) {
	edge := startEdge
	for edge != startEdge {
		if edge == nil {
			log.Printf("Ring edge is null\n")
			return
		}
		if edge.maxEdgeRing == m {
			log.Printf("Ring edge visited twice at %v %v\n", edge.origin, edge.origin)
			return
		}
		if edge.nextResultMaxEdge == nil {
			log.Printf("Ring edge missing at (%v)\n", edge.DirectionPt())
			return
		}
		edge.maxEdgeRing = m
		edge = edge.nextResultMaxEdge
	}
}

// linkResultAreaMaxRingAtNode Traverses the star of edges originating at a node and links
// consecutive result edges together into maximal edge rings. To link two edges the resultNextMax
// pointer for an incoming result edge is set to the next outgoing result edge.
// Edges are linked when:
//		they belong to an area (i.e. they have sides)
//		they are marked as being in the result
// Edges are linked in CCW order (which is the order they are linked in the underlying graph).
// This means that rings have their face on the Right (in other words, the topological location of
// the face is given by the RHS label of the DirectedEdge). This produces rings with CW orientation.
// PRECONDITIONS:
//		- This edge is in the result
//		- This edge is not yet linked
//		- The edge and its sym are NOT both marked as being in the result
func (m *MaximalEdgeRing) linkResultAreaMaxRingAtNode(nodeEdge *OverlayEdge) {
	if !nodeEdge.isInResultArea {
		log.Printf("Attempt to link non-result edge")
		return
	}

	/**
	 * Since the node edge is an out-edge,
	 * make it the last edge to be linked
	 * by starting at the next edge.
	 * The node edge cannot be an in-edge as well,
	 * but the next one may be the first in-edge.
	 */
	endOut := nodeEdge.oNextOE()
	currOut := endOut
	state := STATE_FIND_INCOMING
	currResultIn := new(OverlayEdge)
	for currOut != endOut {
		/**
		 * If an edge is linked this node has already been processed
		 * so can skip further processing
		 */
		if currResultIn != nil && currResultIn.nextResultMaxEdge != nil {
			return
		}

		switch state {
		case STATE_FIND_INCOMING:
			currIn := currOut.symOE()
			if !currIn.isInResultArea {
				break
			}
			currResultIn = currIn
			state = STATE_LINK_OUTGOING
		case STATE_LINK_OUTGOING:
			if !currOut.isInResultArea {
				break
			}
			// link the in edge to the out edge
			if currResultIn != nil {
				currResultIn.nextResultMaxEdge = currOut
			}
			state = STATE_FIND_INCOMING
		}
		currOut = currOut.oNextOE()
	}
	if state == STATE_LINK_OUTGOING {
		log.Printf("no outgoing edge found, coordinate=%v\n", nodeEdge.origin)
	}
}

// buildMinimalRings...
func (m *MaximalEdgeRing) buildMinimalRings() []*OverlayEdgeRing {
	m.linkMinimalRings()
	minEdgeRings := make([]*OverlayEdgeRing, 0)
	e := m.startEdge
	for e != m.startEdge {
		if e.edgeRing == nil {
			minEr := NewOverlayEdgeRing(e)
			minEdgeRings = append(minEdgeRings, minEr)
		}
		e = e.nextResultMaxEdge
	}
	return minEdgeRings
}

// linkMinimalRings...
func (m *MaximalEdgeRing) linkMinimalRings() {
	e := m.startEdge
	for e != m.startEdge {
		m.linkMinRingEdgesAtNode(e, m)
		e = e.nextResultMaxEdge
	}
}

// linkMinRingEdgesAtNode Links the edges of a MaximalEdgeRing around this node
// into minimal edge rings (OverlayEdgeRings). Minimal ring edges are linked in the
// opposite orientation (CW) to the maximal ring. This changes self-touching rings
// into a two or more separate rings, as per the OGC SFS polygon topology semantics.
// This relinking must be done to each max ring separately, rather than all the node
// result edges, since there may be more than one max ring incident at the node.
func (m *MaximalEdgeRing) linkMinRingEdgesAtNode(nodeEdge *OverlayEdge, maxRing *MaximalEdgeRing) {
	// The node edge is an out-edge, so it is the first edge linked with the next CCW in-edge
	endOut := nodeEdge
	currMaxRingOut := endOut
	currOut := endOut.oNextOE()

	for currOut != endOut {
		if m.isAlreadyLinked(currOut.symOE(), maxRing) {
			return
		}
		if currMaxRingOut == nil {
			currMaxRingOut = m.selectMaxOutEdge(currOut, maxRing)
		} else {
			currMaxRingOut = m.linkMaxInEdge(currOut, currMaxRingOut, maxRing)
		}
		currOut = currOut.oNextOE()
	}
	if currMaxRingOut != nil {
		log.Printf("Unmatched edge found during min-ring linking %v\n", nodeEdge.origin)
	}
}

// isAlreadyLinked Tests if an edge of the maximal edge ring is already linked into
// a minimal OverlayEdgeRing. If so, this node has already been processed earlier
// in the maximal edgering linking scan.
// Params:
//		edge – an edge of a maximal edgering
//		maxRing – the maximal edgering
// Returns:
//		true if the edge has already been linked into a minimal edgering.
func (m *MaximalEdgeRing) isAlreadyLinked(edge *OverlayEdge, maxRing *MaximalEdgeRing) bool {
	return edge.maxEdgeRing == maxRing && edge.nextResultEdge != nil
}

// selectMaxOutEdge...
func (m *MaximalEdgeRing) selectMaxOutEdge(currOut *OverlayEdge, maxEdgeRing *MaximalEdgeRing) *OverlayEdge {
	if currOut.maxEdgeRing == maxEdgeRing {
		return currOut
	}
	return nil
}

// linkMaxInEdge...
func (m *MaximalEdgeRing) linkMaxInEdge(currOut, currMaxRingOut *OverlayEdge, maxEdgeRing *MaximalEdgeRing) *OverlayEdge {
	currIn := currOut.symOE()
	if currIn.maxEdgeRing != maxEdgeRing {
		return currMaxRingOut
	}
	currIn.nextResultEdge = currMaxRingOut
	return nil
}
