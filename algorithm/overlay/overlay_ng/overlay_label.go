package overlay_ng

import "github.com/spatial-go/geoos/algorithm/calc"

const (
	SYM_UNKNOWN  = '#'
	SYM_BOUNDARY = 'B'
	SYM_COLLAPSE = 'C'
	SYM_LINE     = 'L'
	// The dimension of an input geometry which is not known
	DIM_UNKNOWN = -1
	// The dimension of an edge which is not part of a specified input geometry.
	DIM_NOT_PART = DIM_UNKNOWN
	// The dimension of an edge which is a line.
	DIM_LINE = 1
	// The dimension for an edge which is part of an input Area geometry boundary.
	DIM_BOUNDARY = 2
	/**
	 * The dimension for an edge which is a collapsed part of an input Area geometry boundary.
	 * A collapsed edge represents two or more line segments which have the same endpoints.
	 * They usually are caused by edges in valid polygonal geometries
	 * having their endpoints become identical due to precision reduction.
	 */
	DIM_COLLAPSE = 3
	// Indicates that the location is currently unknown
	LOC_UNKNOWN = -1
)

// A structure recording the topological situation for an edge in a topology graph
// used during overlay processing. A label contains the topological Locations for
// one or two input geometries to an overlay operation. An input geometry may be
// either a Line or an Area. The label locations for each input geometry are populated
// with the s for the edge Positions when they are created or once they are computed by
// topological evaluation. A label also records the (effective) dimension of each input
// geometry. For area edges the role (shell or hole) of the originating ring is recorded,
// to allow determination of edge handling in collapse cases.
// In an OverlayGraph a single label is shared between the two oppositely-oriented
type OverlayLabel struct {
	aDim      int // DIM_NOT_PART
	aIsHole   bool
	aLocLeft  int // LOC_UNKNOWN
	aLocRight int // LOC_UNKNOWN
	aLocLine  int // LOC_UNKNOWN

	bDim      int // DIM_NOT_PART
	bIsHole   bool
	bLocLeft  int // LOC_UNKNOWN
	bLocRight int // LOC_UNKNOWN
	bLocLine  int // LOC_UNKNOWN
}

// Initializes the label for an edge which is the collapse of part of the boundary
// of an Area input geometry. The location of the collapsed edge relative to the
// parent area geometry is initially unknown. It must be determined from the topology
// of the overlay graph.
// Params:
// 		index – the index of the parent input geometry
//		isHole – whether the dominant edge role is a hole or a shell
func (o *OverlayLabel) initCollapse(index int, isHole bool) {
	if index == 0 {
		o.aDim = DIM_COLLAPSE
		o.aIsHole = isHole
	} else {
		o.bDim = DIM_COLLAPSE
		o.bIsHole = isHole
	}
}

// Initializes the label for an input geometry which is an Area boundary.
// Params:
//		index – the input index of the parent geometry
//		locLeft – the location of the left side of the edge
//		locRight – the location of the right side of the edge
//		isHole – whether the edge role is a hole or a shell
func (o *OverlayLabel) initBoundary(index, locLeft, locRight int, isHole bool) {
	if index == 0 {
		o.aDim = DIM_BOUNDARY
		o.aIsHole = isHole
		o.aLocLeft = locLeft
		o.aLocRight = locRight
		o.aLocLine = calc.ImInterior
	} else {
		o.bDim = DIM_BOUNDARY
		o.bIsHole = isHole
		o.bLocLeft = locLeft
		o.bLocRight = locRight
		o.bLocLine = calc.ImInterior
	}
}

// initLine Initializes the label for an input geometry which is a Line.
// Params:
//		index – the index of the parent input geometry
func (o *OverlayLabel) initLine(index int) {
	if index == 0 {
		o.aDim = DIM_LINE
		o.aLocLine = LOC_UNKNOWN
	} else {
		o.bDim = DIM_LINE
		o.bLocLine = LOC_UNKNOWN
	}
}

// initNotPart Initializes the label for an edge which is not part of an input geometry.
// Params:
//		index – the index of the input geometry
func (o *OverlayLabel) initNotPart(index int) {
	if index == 0 {
		o.aDim = DIM_NOT_PART
	} else {
		o.bDim = DIM_NOT_PART
	}
}

// isBoundary Tests if a label is for an edge which is in the boundary of a source geometry.
// Collapses are not reported as being in the boundary.
// Params:
//		index – the index of the input geometry
// Returns:
//		true if the label is a boundary for the source
func (o *OverlayLabel) isBoundary(index int) bool {
	if index == 0 {
		return o.aDim == DIM_BOUNDARY
	}
	return o.bDim == DIM_BOUNDARY
}

// isBoundaryEither Tests if a label is for an edge which is in the boundary of either source geometry.
func (o *OverlayLabel) isBoundaryEither() bool {
	return o.aDim == DIM_BOUNDARY || o.bDim == DIM_BOUNDARY
}

// isBoundarySingleton Tests whether a label is for an edge which is a boundary of one
// geometry and not part of the other.
// Returns:
//		true if the edge is a boundary singleton
func (o *OverlayLabel) isBoundarySingleton() bool {
	if o.aDim == DIM_BOUNDARY && o.bDim == DIM_NOT_PART {
		return true
	}
	if o.bDim == DIM_BOUNDARY && o.aDim == DIM_NOT_PART {
		return true
	}
	return false
}

// isLine Tests whether at least one of the sources is a Line.
// Returns:
//		true if at least one source is a line
func (o *OverlayLabel) isLine() bool {
	return o.aDim == DIM_LINE || o.bDim == DIM_LINE
}

// isLineByIndex Tests whether a source is a Line.
// Params:
//		index – the index of the input geometry
// Returns:
//		true if the input is a Line
func (o *OverlayLabel) isLineByIndex(index int) bool {
	if index == 0 {
		return o.aDim == DIM_LINE
	}
	return o.bDim == DIM_LINE
}

// isLineInArea Tests if a line edge is inside a source geometry (i.e. it has location Location.INTERIOR).
// Params:
//		index – the index of the input geometry
// Returns:
//		true if the line is inside the source geometry
func (o *OverlayLabel) isLineInArea(index int) bool {
	if index == 0 {
		return o.aLocLine == calc.ImInterior
	}
	return o.bLocLine == calc.ImInterior
}

// isBoundaryBoth Tests if a label is for an edge which is in the boundary of both source geometries.
// Returns:
//		true if the label is a boundary for both sources
func (o *OverlayLabel) isBoundaryBoth() bool {
	return o.aDim == DIM_BOUNDARY && o.bDim == DIM_BOUNDARY
}

// isBoundaryCollapse Tests if the label is a collapsed edge of one area and is a
// (non-collapsed) boundary edge of the other area.
// Returns:
//		true if the label is for a collapse coincident with a boundary
func (o *OverlayLabel) isBoundaryCollapse() bool {
	if o.isLine() {
		return false
	}
	return !o.isBoundaryBoth()
}

// isBoundaryTouch Tests if a label is for an edge where two area touch along their boundary.
// Returns:
//		true if the edge is a boundary touch
func (o *OverlayLabel) isBoundaryTouch() bool {
	return o.isBoundaryBoth() &&
		o.getLocation(0, calc.SideRight, true) != o.getLocation(1, calc.SideRight, true)
}

// isCollapse Tests if an edge is a Collapse for a source geometry.
// Params:
//		index – the index of the input geometry
// Returns:
//		true if the label indicates the edge is a collapse for the source
func (o *OverlayLabel) isCollapse(index int) bool {
	dimension := 0
	if index == 0 {
		dimension = o.aDim
	}
	dimension = o.bDim
	return dimension == DIM_COLLAPSE
}

// isInteriorCollapse Tests if a label is a Collapse has location Location.INTERIOR,
// to at least one source geometry.
// Returns:
//		true if the label is an Interior Collapse to a source geometry
func (o *OverlayLabel) isInteriorCollapse() bool {
	if o.aDim == DIM_COLLAPSE && o.aLocLine == calc.ImInterior {
		return true
	}
	if o.bDim == DIM_COLLAPSE && o.bLocLine == calc.ImInterior {
		return true
	}
	return false
}

// isCollapseAndNotPartInterior Tests if a label is a Collapse and NotPart with location
// Location.INTERIOR for the other geometry.
// Returns:
//		true if the label is a Collapse and a NotPart with Location Interior
func (o *OverlayLabel) isCollapseAndNotPartInterior() bool {
	if o.aDim == DIM_COLLAPSE && o.bDim == DIM_NOT_PART && o.bLocLine == calc.ImInterior {
		return true
	}
	if o.bDim == DIM_COLLAPSE && o.aDim == DIM_NOT_PART && o.aLocLine == calc.ImInterior {
		return true
	}
	return false
}

// getLineLocation Gets the line location for a source geometry.
// Params:
//		index – the index of the input geometry
// Returns:
//		the line location for the source
func (o *OverlayLabel) getLineLocation(index int) int {
	if index == 0 {
		return o.aLocLine
	}
	return o.bLocLine
}

// getLocation Gets the location for a Position of an edge of a source for an edge with given orientation.
// Params:
//		index – the index of the source geometry
//		position – the position to get the location for
//		isForward – true if the orientation of the containing edge is forward
// Returns:
//		the location of the oriented position in the source
func (o *OverlayLabel) getLocation(index, position int, isForward bool) int {
	if index == 0 {
		switch position {
		case calc.SideLeft:
			if isForward {
				return o.aLocLeft
			}
			return o.aLocRight
		case calc.SideRight:
			if isForward {
				return o.aLocRight
			}
			return o.aLocLeft
		case calc.SideOn:
			return o.aLocLine
		}
	}
	// index == 1
	switch position {
	case calc.SideLeft:
		if isForward {
			return o.bLocLeft
		}
		return o.bLocRight
	case calc.SideRight:
		if isForward {
			return o.bLocRight
		}
		return o.bLocLeft
	case calc.SideOn:
		return o.bLocLine
	}
	return LOC_UNKNOWN
}
