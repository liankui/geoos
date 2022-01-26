package matrix

// A lightweight class used to store coordinates on the 2-dimensional Cartesian plane.
// It is distinct from Point, which is a subclass of Geometry. Unlike objects of type
// Point (which contain additional information such as an envelope, a precision model,
// and spatial reference system information), a Coordinate only contains ordinate values
// and accessor methods.
// Coordinates are two-dimensional points, with an additional Z-ordinate. If an
// Z-ordinate value is not specified or not defined, constructed coordinates have a
// Z-ordinate of NaN (which is also the value of NULL_ORDINATE). The standard comparison
// functions ignore the Z-ordinate. Apart from the basic accessor functions, JTS supports
// only specific operations involving the Z-ordinate.
// Implementations may optionally support Z-ordinate and M-measure values as appropriate
// for a CoordinateSequence. Use of getZ() and getM() accessors, or getOrdinate(int) are recommended.
type Coordinate Matrix

// equals2D Returns whether the planar projections of the two Coordinates are equal.
// Params:
//		other – a Coordinate with which to do the 2D comparison.
// Returns:
//		true if the x- and y-coordinates are equal; the z-coordinates do not have to be equal.
func (c Coordinate) equals2D(other Coordinate) bool {
	if len(other) < 2 {
		return false
	}
	if c[0] != other[0] {
		return false
	}
	if c[1] != other[1] {
		return false
	}
	return true
}

