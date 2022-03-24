package noding

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"math"
)

type Type string

const (
	/**
	 * Fixed Precision indicates that coordinates have a fixed number of decimal places.
	 * The number of decimal places is determined by the log10 of the Scale factor.
	 */
	FIXED Type = "FIXED"
	/**
	 * Floating precision corresponds to the standard Java
	 * double-precision floating-point representation, which is
	 * based on the IEEE-754 standard
	 */
	FLOATING Type = "FLOATING"
	/**
	 * Floating single precision corresponds to the standard Java
	 * single-precision floating-point representation, which is
	 * based on the IEEE-754 standard
	 */
	FLOATING_SINGLE Type = "FLOATING SINGLE"
	/**
	 *  The maximum precise value representable in a double. Since IEE754
	 *  double-precision numbers allow 53 bits of mantissa, the value is equal to
	 *  2^53 - 1.  This provides <i>almost</i> 16 decimal digits of precision.
	 */
	maximumPreciseValue = 9007199254740992.0
)

/*
Specifies the precision model of the Coordinates in a Geometry. In other words,
specifies the grid of allowable points for a Geometry. A precision model may be
floating (FLOATING or FLOATING_SINGLE), in which case normal floating-point value semantics apply.
For a FIXED precision model the makePrecise(Coordinate) method allows rounding
a coordinate to a "precise" value; that is, one whose precision is known exactly.
Coordinates are assumed to be precise in geometries. That is, the coordinates are assumed
to be rounded to the precision model given for the geometry. All internal operations assume
that coordinates are rounded to the precision model. Constructive methods (such as boolean
operations) always round computed coordinates to the appropriate precision model.
Three types of precision model are supported:
	FLOATING - represents full double precision floating point. This is the default precision model used in JTS
	FLOATING_SINGLE - represents single precision floating point.
	FIXED - represents a model with a fixed number of decimal places. A Fixed Precision Model is specified by a Scale factor. The Scale factor specifies the size of the grid which numbers are rounded to. Input coordinates are mapped to fixed coordinates according to the following equations:
		jtsPt.x = round( (inputPt.x * Scale ) / Scale
		jtsPt.y = round( (inputPt.y * Scale ) / Scale
For example, to specify 3 decimal places of precision, use a Scale factor of 1000.
To specify -3 decimal places of precision (i.e. rounding to the nearest 1000), use a Scale factor of 0.001.
It is also supported to specify a precise grid size by providing it as a negative Scale factor.
This allows setting a precise grid size rather than using a fractional Scale, which provides
more accurate and robust rounding. For example, to specify rounding to the nearest 1000 use
a Scale factor of -1000.
Coordinates are represented internally as Java double-precision values. Java uses the IEEE-394
floating point standard, which provides 53 bits of precision. (Thus the maximum precisely
representable integer is 9,007,199,254,740,992 - or almost 16 decimal digits of precision).
*/
type PrecisionModel struct {
	Name string
	// The type of PrecisionModel this represents.
	ModelType Type
	// The Scale factor which determines the number of decimal places in fixed precision.
	Scale float64
	// If non-zero, the precise grid size specified. In this case, the Scale is also valid and
	// is computed from the grid size. If zero, the Scale is used to compute the grid size where needed.
	GridSize float64
}

// NewDefaultPrecisionModel...
func NewDefaultPrecisionModel() *PrecisionModel {
	return &PrecisionModel{
		ModelType: FLOATING,
	}
}

func NewPrecisionModelByScale(scale float64) *PrecisionModel {
	p := PrecisionModel{
		ModelType: FIXED,
	}
	p.setScale(scale)
	return &p
}

// IsFloating Tests whether the precision model supports floating point
func (p *PrecisionModel) IsFloating() bool {
	if p.ModelType == FLOATING || p.ModelType == FLOATING_SINGLE {
		return true
	}
	return false
}

// setScale Sets the multiplying factor used to obtain a precise coordinate.
// This method is private because PrecisionModel is an immutable (value) type.
func (p *PrecisionModel) setScale(scale float64) {
	if scale < 0 {
		p.GridSize = math.Abs(scale)
		p.Scale = 1.0 / p.GridSize
	} else {
		p.Scale = math.Abs(scale)
		p.GridSize = 0.0
	}
}

// makePrecise Rounds a numeric value to the PrecisionModel grid. Asymmetric Arithmetic Rounding
// is used, to provide uniform rounding behaviour no matter where the number is on the number line.
func (p *PrecisionModel) makePrecise(val float64) float64 {
	if p.ModelType == FLOATING_SINGLE {
		return val
	}
	if p.ModelType == FIXED {
		if p.GridSize > 0 {
			return math.Round(val/p.GridSize) * p.GridSize
		} else {
			return math.Round(val*p.Scale) / p.Scale
		}
	}
	// modelType == FLOATING - no rounding necessary
	return val
}

// MakePrecise Rounds a Coordinate to the PrecisionModel grid.
func (p *PrecisionModel) MakePrecise(coord matrix.Matrix) {
	// optimization for full precision
	if p.ModelType == FLOATING {
		return
	}
	coord[0] = p.makePrecise(coord[0])
	coord[1] = p.makePrecise(coord[1])
}
