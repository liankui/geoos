// Package operation define valid func for geometries.
package overlay_ng

import (
	"github.com/spatial-go/geoos/algorithm"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/space"
)

const (
	INTERSECTION  = 1
	UNION         = 2
	DIFFERENCE    = 3
	SYMDIFFERENCE = 4
)

// OverlayOp computes the geometric overlay of two Geometries.
type OverlayOp struct {
	//matrix.Steric
}

// CreateEmptyResult Creates an empty result geometry of the appropriate dimension,
// based on the given overlay operation and the dimensions of the inputs. The created
// geometry is always an atomic geometry, not a collection.
// The empty result is constructed using the following rules:
//		INTERSECTION - result has the dimension of the lowest input dimension
//		UNION - result has the dimension of the highest input dimension
//		DIFFERENCE - result has the dimension of the left-hand input
//		SYMDIFFERENCE - result has the dimension of the highest input dimension (since
//			the symmetric Difference is the union of the differences).
// Params:
//		overlayOpCode – the code for the overlay operation being performed
//		a – an input geometry
//		b – an input geometry
// Returns:
//		an empty atomic geometry of the appropriate dimension
func (o *OverlayOp) CreateEmptyResult(overlayOpCode int, a, b space.Geometry) (space.Geometry, error) {
	resultDim := o.resultDimension(overlayOpCode, a, b)
	// Handles resultSDim = -1, although should not happen
	result, _ := o.createEmpty(resultDim)
	return result, nil
}

// resultDimension...
func (o *OverlayOp) resultDimension(opCode int, g0, g1 space.Geometry) int {
	dim0 := g0.Dimensions()
	dim1 := g1.Dimensions()

	resultDimension := -1
	switch opCode {
	case INTERSECTION:
		resultDimension = min(dim0, dim1)
	case UNION:
		resultDimension = max(dim0, dim1)
	case DIFFERENCE:
		resultDimension = dim0
	case SYMDIFFERENCE:
		// SymDiff = Union(Diff(A, B), Diff(B, A)
		resultDimension = max(dim0, dim1)
	}
	return resultDimension
}

// createEmpty Creates an empty atomic geometry of the given dimension. If passed
// a dimension of -1 will create an empty GeometryCollection.
// Params:
//		dimension – the required dimension (-1, 0, 1 or 2)
// Returns:
//		an empty atomic geometry of given dimension
func (o *OverlayOp) createEmpty(dimension int) (space.Geometry, error) {
	switch dimension {
	case -1:
		return space.Collection{}, nil
	case 0:
		return space.Point{}, nil
	case 1:
		return space.LineString{}, nil
	case 2:
		return space.Polygon{}, nil
	default:
		return nil, algorithm.ErrInvalidDimension(dimension)
	}
}

// max return int
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

// min return int
func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// Overlay...
func (o *OverlayOp) Overlay(a, b matrix.Steric, opCode int) (matrix.Steric, error) {


	return nil, nil
}
