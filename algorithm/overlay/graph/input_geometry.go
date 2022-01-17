package graph

import "github.com/spatial-go/geoos/space"

// Manages the input geometries for an overlay operation.
// The second geometry is allowed to be null, to support
// for instance precision reduction.
type InputGeometry struct {
	geom [2]space.Geometry

}
