// Package operation define valid func for geometries.
package graph

import (
	"github.com/spatial-go/geoos/algorithm"
	"github.com/spatial-go/geoos/algorithm/matrix"
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

/*
CreateEmptyResult Creates an empty result geometry of the appropriate dimension,
based on the given overlay operation and the dimensions of the inputs.
The created geometry is always an atomic geometry, not a collection.
*/
func (o *OverlayOp) CreateEmptyResult(overlayOpCode int, a, b matrix.Steric) (matrix.Steric, error) {
	resultDim := o.resultDimension(overlayOpCode, a, b)
	result, _ := createEmpty(resultDim)
	return result, nil
}

func (o *OverlayOp) resultDimension(opCode int, g0, g1 matrix.Steric) int {
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

func createEmpty(dimension int) (matrix.Steric, error) {
	switch dimension {
	case -1:
		//return createGeometryCollection()
	case 0:
		//return createPoint()
	case 1:
		//return createLineString()
	case 2:
		//return createPolygon()
	default:
		return nil, algorithm.ErrInvalidDimension(dimension)
	}
	return nil, nil
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func (o *OverlayOp) Overlay(a, b matrix.Steric, opCode int) (matrix.Steric, error) {


	return nil, nil
}
