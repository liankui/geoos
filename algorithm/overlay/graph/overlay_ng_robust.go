package graph

import "github.com/spatial-go/geoos/space"

type OverlayNGRobust struct {
	
}

func (o OverlayNGRobust) overlay(geom0, geom1 space.Geometry, opCode int) space.Geometry {
	/**
	 * First try overlay with a FLOAT noder, which is fast and causes least
	 * change to geometry coordinates
	 * By default the noder is validated, which is required in order
	 * to detect certain invalid noding situations which otherwise
	 * cause incorrect overlay output.
	 */
	var ov OverlayNG
	return ov.overlay(geom0, geom1, opCode)
}
