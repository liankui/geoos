package overlay_ng

import "github.com/spatial-go/geoos/space"

// Internal class which encapsulates the runtime switch to use OverlayNG,
// and some additional extensions for optimization and GeometryCollection handling.
// This class allows the Geometry overlay methods to be switched between the
// original algorithm and the modern OverlayNG codebase via a system property jts.overlay.
type GeometryOverlay struct {}

// overlay...
func (g GeometryOverlay) overlay(a, b space.Geometry, opCode int) space.Geometry {
	var overlayNGRobust OverlayNGRobust
	return overlayNGRobust.overlay(a, b, opCode)
}

// union...
func (g GeometryOverlay) union(a, b space.Geometry) space.Geometry {
	// handle empty geometry cases
	if a.IsEmpty() || b.IsEmpty() {
		if a.IsEmpty() && b.IsEmpty() {
			var op OverlayOp
			emptyResult, _ := op.CreateEmptyResult(UNION, a, b)
			return emptyResult
		}
		if a.IsEmpty() {
			return b
		}
		if b.IsEmpty() {
			return a
		}
	}

	return g.overlay(a, b, UNION)
}