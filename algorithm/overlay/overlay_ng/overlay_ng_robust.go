package overlay_ng

import (
	"github.com/spatial-go/geoos/algorithm/overlay/overlay_ng/noding"
	"github.com/spatial-go/geoos/space"
)

type OverlayNGRobust struct {
	
}

// overlay Overlay two geometries, using heuristics to ensure computation completes
// correctly. In practice the heuristics are observed to be fully correct.
func (o *OverlayNGRobust) overlay(geom0, geom1 space.Geometry, opCode int) space.Geometry {
	/**
	 * First try overlay with a FLOAT noder, which is fast and causes least
	 * change to geometry coordinates
	 * By default the noder is validated, which is required in order
	 * to detect certain invalid noding situations which otherwise
	 * cause incorrect overlay output.
	 */
	//var ov OverlayNG
	//overlay := ov.overlay(geom0, geom1, opCode)
	//if overlay != nil {
	//	return overlay
	//}

	/**
	 * On failure retry using snapping noding with a "safe" tolerance.
	 * if this throws an exception just let it go,
	 * since it is something that is not a TopologyException
	 */
	// todo overlaySnapTries

	// On failure retry using snap-rounding with a heuristic scale factor (grid size).
	// todo overlaySR
	result := o.overlaySR(geom0, geom1, opCode)
	if result != nil {
		return result
	}

	return nil
}

func (o *OverlayNGRobust) overlaySR(geom0, geom1 space.Geometry, opCode int) space.Geometry {
	var p PrecisionUtil
	safeScale := p.safeScaleForTwoGeom(geom0, geom1)
	pmSafe := noding.NewPrecisionModelByScale(safeScale)
	var ov OverlayNG
	return ov.overlaySR(geom0, geom1, opCode, pmSafe)
}