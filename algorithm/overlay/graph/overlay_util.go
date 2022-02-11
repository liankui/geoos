package graph

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/coordtransform"
	"github.com/spatial-go/geoos/space"
	"math"
)

const (
	SAFE_ENV_BUFFER_FACTOR = 0.1
	SAFE_ENV_GRID_FACTOR   = 3
)

// Utility methods for overlay processing.
type OverlayUtil struct {
}

// Computes a clipping envelope for overlay input geometries. The clipping envelope
// encloses all geometry line segments which might participate in the overlay,
// with a buffer to account for numerical precision (in particular, rounding due to
// a precision model. The clipping envelope is used in both the RingClipper and in the LineLimiter.
// Some overlay operations (i.e. and OverlayNG#SYMDIFFERENCE cannot use clipping as an
// optimization, since the result envelope is the full extent of the two input geometries.
// In this case the returned envelope is null to indicate this.
// Params:
//		opCode – the overlay op code
//		inputGeom – the input geometries
//		pm – the precision model being used
// Returns:
//		an envelope for clipping and line limiting, or null if no clipping is performed
func (o *OverlayUtil) clippingEnvelope(opCode int, inputGeom *InputGeometry, pm *PrecisionModel) *envelope.Envelope {
	resultEnv := o.resultEnvelope(opCode, inputGeom, pm)
	if resultEnv == nil {
		return nil
	}
	clipEnv := RobustClipEnvelopeComputer.getEnvelope(
		inputGeom.getGeometry(0),
		inputGeom.getGeometry(1),
		resultEnv)

	safeEnv := o.safeEnv(clipEnv, pm)
	return safeEnv
}

// Computes an envelope which covers the extent of the result of a given overlay
// operation for given inputs. The operations which have a result envelope smaller
// than the extent of the inputs are:
//		OverlayNG.INTERSECTION: result envelope is the intersection of the input envelopes
//		OverlayNG.DIFERENCE: result envelope is the envelope of the A input geometry
// Otherwise, null is returned to indicate full extent.
// Params:
//		opCode –
//		inputGeom –
//		pm –
//Returns:
//		the result envelope, or null if the full extent
func (o *OverlayUtil) resultEnvelope(opCode int, inputGeom *InputGeometry, pm *PrecisionModel) *envelope.Envelope {
	overlapEnv := new(envelope.Envelope)
	switch opCode {
	case INTERSECTION:
		// use safe envelopes for intersection to ensure they contain rounded coordinates
		envA := o.safeEnv(inputGeom.getEnvelope(0), pm)
		envB := o.safeEnv(inputGeom.getEnvelope(1), pm)
		overlapEnv = envA.Intersection(envB)
	case DIFFERENCE:
		overlapEnv = o.safeEnv(inputGeom.getEnvelope(0), pm)
	}
	// return null for UNION and SYMDIFFERENCE to indicate no clipping
	return overlapEnv
}

// Determines a safe geometry envelope for clipping, taking into account the precision model being used.
// Params:
//		env – a geometry envelope
//		pm – the precision model
// Returns:
//		a safe envelope to use for clipping
func (o *OverlayUtil) safeEnv(env *envelope.Envelope, pm *PrecisionModel) *envelope.Envelope {
	envExpandDist := o.safeExpandDistance(env, pm)
	safeEnv := env.Copy()
	safeEnv.ExpandBy(envExpandDist)
	return safeEnv
}

// safeExpandDistance...
func (o *OverlayUtil) safeExpandDistance(env *envelope.Envelope, pm *PrecisionModel) float64 {
	var envExpandDist float64
	if o.isFloating(pm) {
		// if PM is FLOAT then there is no scale factor, so add 10%
		minSize := math.Min(env.Height(), env.Width())
		// heuristic to ensure zero-width envelopes don't cause total clipping
		if minSize <= 0.0 {
			minSize = math.Max(env.Height(), env.Width())
		}
		envExpandDist = SAFE_ENV_BUFFER_FACTOR * minSize
	} else {
		// if PM is fixed, add a small multiple of the grid size
		gridSize := 1.0 / pm.Scale
		envExpandDist = SAFE_ENV_GRID_FACTOR * gridSize
	}
	return envExpandDist
}

// isFloating A null-handling wrapper for PrecisionModel.isFloating()
func (o *OverlayUtil) isFloating(pm *PrecisionModel) bool {
	if pm == nil {
		return true
	}
	return pm.isFloating()
}

// toLines...
func (o *OverlayUtil) toLines(graph *OverlayGraph, isOutputEdges bool) space.Geometry {
	lines := make([]matrix.LineMatrix, 0)
	for _, edge := range graph.edges {
		includeEdge := isOutputEdges || edge.isInResultArea
		if !includeEdge {
			continue
		}
		pts := edge.getCoordinatesOriented()

		tmp := make([][]float64, 0)
		for _, pt := range pts {
			tmp = append(tmp, pt)
		}
		var line matrix.LineMatrix = tmp

		//line.setUserData(labelForResult(edge) )
		lines = append(lines, line)
	}

	// todo 暂时使用geoos的构建方式
	var trans coordtransform.Transformer
	lineMatrices := trans.TransformMultiLineString(lines)

	var linesTmp space.MultiLineString
	for _, lineMatrix := range lineMatrices {
		linesTmp = append(linesTmp, space.LineString(lineMatrix))
	}
	return linesTmp
}
