package overlay_ng

import (
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/space"
	"math"
)

type PrecisionUtilType int

const MAX_ROBUST_DP_DIGITS PrecisionUtilType = 14

// Functions for computing precision model scale factors that ensure robust geometry operations.
// In particular, these can be used to automatically determine appropriate scale factors for
// operations using limited-precision noding (such as OverlayNG).
type PrecisionUtil struct {
}

// Computes a safe scale factor for a numeric value. A safe scale factor ensures that rounded number has no more than MAX_ROBUST_DP_DIGITS digits of precision.
func (p *PrecisionUtil) safeScale(value float64) float64 {
	return p.precisionScale(value, int(MAX_ROBUST_DP_DIGITS))
}

// safeScale Computes a safe scale factor for two geometries. A safe scale factor ensures
// that the rounded ordinates have no more than MAX_ROBUST_DP_DIGITS digits of precision.
func (p *PrecisionUtil) safeScaleForTwoGeom(a, b space.Geometry) float64 {
	maxBnd := p.maxBoundMagnitude(a.ComputeEnvelopeInternal())
	if b != nil {
		maxBndB := p.maxBoundMagnitude(b.ComputeEnvelopeInternal())
		maxBnd = math.Max(maxBnd, maxBndB)
	}
	scale := p.safeScale(maxBnd)
	return scale
}

// maxBoundMagnitude Determines the maximum magnitude (absolute value) of the bounds of an
// of an envelope. This is equal to the largest ordinate value which must be accommodated by a scale factor.
func (p *PrecisionUtil) maxBoundMagnitude(env *envelope.Envelope) float64 {
	return maxInFour(
		math.Abs(env.MaxX),
		math.Abs(env.MaxY),
		math.Abs(env.MinX),
		math.Abs(env.MinY),
	)
}

// maxInFour...
func maxInFour(v1, v2, v3, v4 float64) float64 {
	max := v1
	if v2 > max {
		max = v2
	}
	if v3 > max {
		max = v3
	}
	if v4 > max {
		max = v4
	}
	return max
}

// precisionScale Computes the scale factor which will produce a given number of
// digits of precision (significant digits) when used to round the given number.
func (p *PrecisionUtil) precisionScale(value float64, precisionDigits int) float64 {
	magnitude := (int)(math.Log(value)/math.Log(10) + 1.0)
	precDigits := precisionDigits - magnitude

	scaleFactor := math.Pow(10.0, float64(precDigits))
	return scaleFactor
}
