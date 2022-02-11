package graph

import "github.com/spatial-go/geoos/space"

// Extracts Point resultants from an overlay graph created by an Intersection operation
// between non-Point inputs. Points may be created during intersection if lines or areas
// touch one another at single points. Intersection is the only overlay operation which
// can result in Points from non-Point inputs.
// Overlay operations where one or more inputs are Points are handled via a different code path.
type IntersectionPointBuilder struct {
	graph *OverlayGraph
	points []space.Point

}

// NewIntersectionPointBuilder...
func NewIntersectionPointBuilder() *IntersectionPointBuilder {
	return &IntersectionPointBuilder{

	}
}