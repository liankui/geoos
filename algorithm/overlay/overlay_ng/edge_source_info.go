package overlay_ng

import "github.com/spatial-go/geoos/algorithm/calc"

// Records topological information about an edge representing a piece of
// linework (lineString or polygon ring) from a single source geometry.
// This information is carried through the noding process (which may result
// in many noded edges sharing the same information object).
// It is then used to populate the topology info fields in Edges (possibly via merging).
// That information is used to construct the topology graph OverlayLabels.
type EdgeSourceInfo struct {
	index      int
	dim        int // default=-999
	isHole     bool
	depthDelta int
}

func NewEdgeSourceInfo(index, depthDelta int, isHole bool) *EdgeSourceInfo {
	return &EdgeSourceInfo{
		index:      index,
		dim:        calc.ImA,
		isHole:     isHole,
		depthDelta: depthDelta,
	}
}
