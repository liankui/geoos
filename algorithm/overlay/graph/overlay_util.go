package graph

import (
	"github.com/spatial-go/geoos/algorithm/overlay/graph/noding"
	"log"
	"math"
	"reflect"

	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/space"
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
func (o *OverlayUtil) clippingEnvelope(opCode int, inputGeom *InputGeometry, pm *noding.PrecisionModel) *envelope.Envelope {
	resultEnv := o.resultEnvelope(opCode, inputGeom, pm)
	if resultEnv == nil {
		return nil
	}

	// todo 暂时未到这里
	//clipEnv := RobustClipEnvelopeComputer.getEnvelope(
	//	inputGeom.getGeometry(0),
	//	inputGeom.getGeometry(1),
	//	resultEnv)

	//safeEnv := o.safeEnv(clipEnv, pm)
	//return safeEnv

	return nil
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
func (o *OverlayUtil) resultEnvelope(opCode int, inputGeom *InputGeometry, pm *noding.PrecisionModel) *envelope.Envelope {
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
func (o *OverlayUtil) safeEnv(env *envelope.Envelope, pm *noding.PrecisionModel) *envelope.Envelope {
	envExpandDist := o.safeExpandDistance(env, pm)
	safeEnv := env.Copy()
	safeEnv.ExpandBy(envExpandDist)
	return safeEnv
}

// safeExpandDistance...
func (o *OverlayUtil) safeExpandDistance(env *envelope.Envelope, pm *noding.PrecisionModel) float64 {
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
func (o *OverlayUtil) isFloating(pm *noding.PrecisionModel) bool {
	if pm == nil {
		return true
	}
	return pm.IsFloating()
}

// createEmptyResult Creates an empty result geometry of the appropriate dimension,
// based on the given overlay operation and the dimensions of the inputs. The created
// geometry is an atomic geometry, not a collection (unless the dimension is -1,
// in which case a GEOMETRYCOLLECTION EMPTY is created.)
// Params:
//		dim – the dimension of the empty geometry to create
//		geomFact – the geometry factory being used for the operation
// Returns:
//		an empty atomic geometry of the appropriate dimension
func (o *OverlayUtil) createEmptyResult(dim int) space.Geometry {
	var result space.Geometry
	switch dim {
	case 0:
		result = space.Point{}
	case 1:
		result = space.LineString{}
	case 2:
		result = space.Polygon{}
	case -1:
		result = space.Collection{}
	default:
		log.Printf("Unable to determine overlay result geometry dimension\n")
	}
	return result
}

// resultDimension Computes the dimension of the result of applying the given operation
// to inputs with the given dimensions. This assumes that complete collapse does not occur.
// The result dimension is computed according to the following rules:
// 		OverlayNG.INTERSECTION - result has the dimension of the lowest input dimension
//		OverlayNG.UNION - result has the dimension of the highest input dimension
//		OverlayNG.DIFFERENCE - result has the dimension of the left-hand input
//		OverlayNG.SYMDIFFERENCE - result has the dimension of the highest input dimension
//			(since the Symmetric Difference is the Union of the Differences).
// Params:
//		opCode – the overlay operation
//		dim0 – dimension of the LH input
//		dim1 – dimension of the RH input
// Returns:
//		the dimension of the result
func (o *OverlayUtil) resultDimension(opCode, dim0, dim1 int) int {
	resultDimension := -1
	switch opCode {
	case INTERSECTION:
		resultDimension = min(dim0, dim1)
	case UNION:
		resultDimension = max(dim0, dim1)
	case DIFFERENCE:
		resultDimension = dim0
	case SYMDIFFERENCE:
		/**
		 * This result is chosen because
		 * <pre>
		 * SymDiff = Union( Diff(A, B), Diff(B, A) )
		 * </pre>
		 * and Union has the dimension of the highest-dimension argument.
		 */
		resultDimension = max(dim0, dim1)
	}
	return resultDimension
}

// createResultGeometry Creates an overlay result geometry for homogeneous or mixed components.
// Params:
//		resultPolyList – the list of result polygons (may be empty or null)
//		resultLineList – the list of result lines (may be empty or null)
//		resultPointList – the list of result points (may be empty or null)
// Returns:
//		a geometry structured according to the overlay result semantics
func (o *OverlayUtil) createResultGeometry(resultPolyList []space.Polygon, resultLineList []space.LineString,
	resultPointList []space.Point) space.Geometry {

	geomList := make([]space.Geometry, 0)
	// element geometries of the result are always in the order A,L,P
	if resultPolyList != nil {
		for _, polygon := range resultPolyList {
			geomList = append(geomList, polygon)
		}
	}
	if resultLineList != nil {
		for _, lineString := range resultLineList {
			geomList = append(geomList, lineString)
		}
	}
	if resultPointList != nil {
		for _, point := range resultPointList {
			geomList = append(geomList, point)
		}
	}

	return o.buildGeometry(geomList)
}

// toLines...
func (o *OverlayUtil) toLines(graph *OverlayGraph, isOutputEdges bool) space.Geometry {
	lines := make([]space.LineString, 0)
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
		var line space.LineString = tmp

		// todo line.setUserData(labelForResult(edge) )
		lines = append(lines, line)
	}

	tmpLines := make([]space.Geometry, 0)
	for _, line := range lines {
		tmpLines = append(tmpLines, line)
	}
	return o.buildGeometry(tmpLines)
}

// buildGeometry Build an appropriate Geometry, MultiGeometry, or GeometryCollection to
// contain the Geometrys in it. For example:
//		If geomList contains a single Polygon, the Polygon is returned.
//		If geomList contains several Polygons, a MultiPolygon is returned.
//		If geomList contains some Polygons and some LineStrings, a GeometryCollection is returned.
//		If geomList is empty, an empty GeometryCollection is returned
// Note that this method does not "flatten" Geometries in the input, and hence if any
// MultiGeometries are contained in the input a GeometryCollection containing them will be returned.
// Params:
//		geomList – the Geometrys to combine
// Returns:
//		a Geometry of the "smallest", "most type-specific" class that can contain the elements of geomList .
func (o *OverlayUtil) buildGeometry(geomList []space.Geometry) space.Geometry {
	// Determine some facts about the geometries in the list
	var geomClass interface{}
	isHeterogeneous := false
	hasGeometryCollection := false
	for _, geom := range geomList {
		partClass := reflect.TypeOf(geom)
		if geomClass == nil {
			geomClass = partClass
		}
		if partClass != geomClass {
			isHeterogeneous = true
		}
		switch geom.(type) {
		case space.Collection:
			hasGeometryCollection = true
		}
	}

	// Now construct an appropriate geometry to return
	// for the empty geometry, return an empty GeometryCollection
	if geomClass == nil {
		return space.Collection{}
	}
	if isHeterogeneous || hasGeometryCollection {
		return space.Collection(geomList)
	}

	// at this point we know the collection is hetereogenous.
	// Determine the type of the result from the first Geometry in the list
	// this should always return a geometry, since otherwise an empty collection would have already been returned
	geom0 := geomList[0]
	isCollection := len(geomList) > 1
	if isCollection {
		switch geom0.(type) {
		case space.Polygon:
			tmp := make([]space.Polygon, 0)
			for _, geom := range geomList {
				tmp = append(tmp, geom.(space.Polygon))
			}
			return space.MultiPolygon(tmp)
		case space.LineString:
			tmp := make([]space.LineString, 0)
			for _, geom := range geomList {
				tmp = append(tmp, geom.(space.LineString))
			}
			return space.MultiLineString(tmp)
		case space.Point:
			tmp := make([]space.Point, 0)
			for _, geom := range geomList {
				tmp = append(tmp, geom.(space.Point))
			}
			return space.MultiPoint(tmp)
		}
	}
	return geom0
}
