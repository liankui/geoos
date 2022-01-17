package graph

import "github.com/spatial-go/geoos/space"

type GeometryOverlay struct {
	
}

func (g GeometryOverlay) overlay(a, b space.Geometry, opCode int) space.Geometry {
	var overlayNGRobust OverlayNGRobust
	return overlayNGRobust.overlay(a, b, opCode)
}

func (g GeometryOverlay) union(a, b space.Geometry) space.Geometry {
	// todo handle empty geometry cases


	return g.overlay(a, b, UNION)
}