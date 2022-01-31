package graph

const (
	SYM_UNKNOWN  = '#'
	SYM_BOUNDARY = 'B'
	SYM_COLLAPSE = 'C'
	SYM_LINE     = 'L'
	/**
	 * The dimension of an input geometry which is not known
	 */
	DIM_UNKNOWN = -1
	/**
	 * The dimension of an edge which is not part of a specified input geometry.
	 */
	DIM_NOT_PART = DIM_UNKNOWN

	/**
	 * The dimension of an edge which is a line.
	 */
	DIM_LINE = 1

	/**
	 * The dimension for an edge which is part of an input Area geometry boundary.
	 */
	DIM_BOUNDARY = 2

	/**
	 * The dimension for an edge which is a collapsed part of an input Area geometry boundary.
	 * A collapsed edge represents two or more line segments which have the same endpoints.
	 * They usually are caused by edges in valid polygonal geometries
	 * having their endpoints become identical due to precision reduction.
	 */
	DIM_COLLAPSE = 3

	/**
	 * Indicates that the location is currently unknown
	 */
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
