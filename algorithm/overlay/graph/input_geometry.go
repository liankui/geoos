package graph

import (
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/space"
)

// Manages the input geometries for an overlay operation.
// The second geometry is allowed to be null, to support
// for instance precision reduction.
type InputGeometry struct {
	geom        [2]space.Geometry
	isCollapsed [2]bool
}

// getGeometry...
func (i *InputGeometry) getGeometry(geomIndex int) space.Geometry {
	return i.geom[geomIndex]
}

// getEnvelope...
func (i *InputGeometry) getEnvelope(geomIndex int) *envelope.Envelope {
	return i.geom[geomIndex].GetEnvelopeInternal()
}

// setCollapsed...
func (i *InputGeometry) setCollapsed(geomIndex int, isGeomCollapsed bool) {
	i.isCollapsed[geomIndex] = isGeomCollapsed
}

// getDimension...
func (i *InputGeometry) getDimension(index int) int {
	if i.geom[index] == nil {
		return -1
	}
	return i.geom[index].Dimensions()
}

// getAreaIndex...
func (i *InputGeometry) getAreaIndex() int {
	if i.getDimension(0) == 2 {
		return 0
	}
	if i.getDimension(1) == 2 {
		return 1
	}
	return -1
}
