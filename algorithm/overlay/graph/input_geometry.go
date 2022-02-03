package graph

import (
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/space"
)

// Manages the input geometries for an overlay operation.
// The second geometry is allowed to be null, to support
// for instance precision reduction.
type InputGeometry struct {
	geom [2]space.Geometry

}

// getGeometry...
func (i *InputGeometry) getGeometry(geomIndex int) space.Geometry {
	return i.geom[geomIndex]
}

// getEnvelope...
func (i *InputGeometry) getEnvelope(geomIndex int) *envelope.Envelope {
	return i.geom[geomIndex].GetEnvelopeInternal()
}
