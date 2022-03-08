package overlay_ng

import (
	"github.com/spatial-go/geoos/space"
	"log"
)

type PolygonBuilder struct {
	shellList          []*OverlayEdgeRing
	freeHoleList       []*OverlayEdgeRing
	isEnforcePolygonal bool // default=true
}

// NewPolygonBuilder...
func NewPolygonBuilder(resultAreaEdges []*OverlayEdge) *PolygonBuilder {
	polygonBuilder := PolygonBuilder{
		isEnforcePolygonal: true,
	}
	polygonBuilder.buildRings(resultAreaEdges)
	return &polygonBuilder
}

// getPolygons...
func (p *PolygonBuilder) getPolygons() []space.Polygon {
	return p.computePolygons(p.shellList)
}

// computePolygons...
func (p *PolygonBuilder) computePolygons(shellList []*OverlayEdgeRing) []space.Polygon {
	resultPolyList := make([]space.Polygon, 0)
	// add Polygons for all shells
	for _, er := range shellList {
		poly := er.toPolygon()
		resultPolyList = append(resultPolyList, poly)
	}
	return resultPolyList
}

// buildRings...
func (p *PolygonBuilder) buildRings(resultAreaEdges []*OverlayEdge) {
	p.linkResultAreaEdgesMax(resultAreaEdges)
	maxRings := p.buildMaximalRings(resultAreaEdges)
	p.buildMinimalRings(maxRings)
	p.placeFreeHoles(p.shellList, p.freeHoleList)
}

// linkResultAreaEdgesMax...
func (p *PolygonBuilder) linkResultAreaEdgesMax(resultEdges []*OverlayEdge) {
	var maximalEdgeRing MaximalEdgeRing
	for _, edge := range resultEdges {
		maximalEdgeRing.linkResultAreaMaxRingAtNode(edge)
	}
}

// buildMaximalRings...
func (p *PolygonBuilder) buildMaximalRings(edges []*OverlayEdge) []*MaximalEdgeRing {
	edgeRings := make([]*MaximalEdgeRing, 0)
	for _, e := range edges {
		if e.isInResultArea && e.label.isBoundaryEither() {
			// if this edge has not yet been processed
			if e.maxEdgeRing == nil {
				er := NewMaximalEdgeRing(e)
				edgeRings = append(edgeRings, er)
			}
		}
	}
	return edgeRings
}

// buildMinimalRings...
func (p *PolygonBuilder) buildMinimalRings(maxRings []*MaximalEdgeRing) {
	for _, erMax := range maxRings {
		minRings := erMax.buildMinimalRings()
		p.assignShellsAndHoles(minRings)
	}
}

// placeFreeHoles Place holes have not yet been assigned to a shell. These "free" holes
// should all be properly contained in their parent shells, so it is safe to use the
// findEdgeRingContaining method. (This is the case because any holes which are NOT
// properly contained (i.e. are connected to their parent shell) would have formed part of a
// MaximalEdgeRing and been handled in a previous step).
func (p *PolygonBuilder) placeFreeHoles(shellList []*OverlayEdgeRing, freeHoleList []*OverlayEdgeRing) {
	for _, hole := range freeHoleList {
		// only place this hole if it doesn't yet have a shell
		if hole.shell == nil {
			shell := hole.findEdgeRingContaining(shellList)
			// only when building a polygon-valid result
			if p.isEnforcePolygonal && shell == nil {
				log.Printf("unable to assign free hole to a shell %v\n", hole.ringPts[0])
			}
			hole.setShell(shell)
		}
	}
}

// assignShellsAndHoles...
func (p *PolygonBuilder) assignShellsAndHoles(minRings []*OverlayEdgeRing) {
	/**
	 * Two situations may occur:
	 * - the rings are a shell and some holes
	 * - rings are a set of holes
	 * This code identifies the situation
	 * and places the rings appropriately
	 */
	shell := p.findSingleShell(minRings)
	if shell != nil {
		p.assignHoles(shell, minRings)
		p.shellList = append(p.shellList, shell)
	} else {
		// all rings are holes; their shell will be found later
		p.freeHoleList = append(p.freeHoleList, minRings...)
	}
}

// findSingleShell Finds the single shell, if any, out of a list of minimal rings derived
// from a maximal ring. The other possibility is that they are a set of (connected) holes,
// in which case no shell will be found.
// Returns:
//		the shell ring, if there is one or null, if all rings are holes
func (p *PolygonBuilder) findSingleShell(edgeRings []*OverlayEdgeRing) *OverlayEdgeRing {
	shellCount := 0
	shell := new(OverlayEdgeRing)
	for _, er := range edgeRings {
		if !er.isHole {
			shell = er
			shellCount++
		}
	}
	if shellCount > 1 {
		log.Printf("found two shells in EdgeRing list\n")
	}
	return shell
}

// assignHoles For the set of minimal rings comprising a maximal ring, assigns the holes
// to the shell known to contain them. Assigning the holes directly to the shell serves two purposes:
//		it is faster than using a point-in-polygon check later on.
//		it ensures correctness, since if the PIP test was used the point chosen might lie on the shell,
//			which might return an incorrect result from the PIP test
func (p *PolygonBuilder) assignHoles(shell *OverlayEdgeRing , edgeRings []*OverlayEdgeRing) {
	for _, er := range edgeRings {
		if er.isHole {
			er.shell = shell
		}
	}
}
