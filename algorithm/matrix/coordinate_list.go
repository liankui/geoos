package matrix

import (
	"reflect"
)

type CoordinateList []Matrix

// RemoveRepeatedPoints If the coordinate array argument has repeated points,
// constructs a new array containing no repeated points. Otherwise, returns the argument.
// Params:
//		coord – an array of coordinates
// Returns:
//		the array with repeated coordinates removed
func (c CoordinateList) RemoveRepeatedPoints(coord []Matrix) []Matrix {
	if !c.hasRepeatedPoints(coord) {
		return coord
	}
	coordList := c.CoordinateList(coord, false)
	return coordList
}

// hasRepeatedPoints Tests whether []Matrix returns true
// for any two consecutive Coordinates in the given array.
func (c CoordinateList) hasRepeatedPoints(coord []Matrix) bool {
	for i := 1; i < len(coord); i++ {
		if reflect.DeepEqual(coord[i-1], coord[i]) {
			return true
		}
	}
	return false
}

// Constructs a new list from an array of Coordinates, allowing caller to specify
// if repeated points are to be removed.
// Params:
//		coord – the array of coordinates to load into the list
//		allowRepeated – if false, repeated points are removed
func (c CoordinateList) CoordinateList(coord []Matrix, allowRepeated bool) CoordinateList {
	return c.AddWithDirection(coord, allowRepeated, true)
}

// AddWithDirection Adds an array of coordinates to the list.
// Params:
//		coord – The coordinates
//		allowRepeated – if set to false, repeated coordinates are collapsed
//		direction – if false, the array is added in reverse order
func (c CoordinateList) AddWithDirection(coord []Matrix, allowRepeated, direction bool) CoordinateList {
	if direction {
		for i := 0; i < len(coord); i++ {
			c = c.AddToEndList(coord[i], allowRepeated)
		}
	} else {
		for i := len(coord) - 1; i >= 0; i-- {
			c = c.AddToEndList(coord[i], allowRepeated)
		}
	}
	return c
}

// AddToEndList Adds a coordinate to the end of the list.
// Params:
//		coord – The coordinates
//		allowRepeated – if set to false, repeated coordinates are collapsed
func (c CoordinateList) AddToEndList(coord Matrix, allowRepeated bool) CoordinateList {
	if !allowRepeated {
		if len(c) > 0 {
			last := c[len(c)-1]
			if last.Equals(coord) {
				return c
			}
		}
	}
	return append(c, coord)
}

// CloseRing Ensure this coordList is a ring, by adding the start point if necessary
func (c CoordinateList) CloseRing() {
	if len(c) > 0 {
		duplicate := c[0]
		c.AddToEndList(duplicate, false)
	}
}

// ToPts...
func (c CoordinateList) ToPts() []Matrix {
	return c
}

// ToPts...
func (c CoordinateList) ToLineString() LineMatrix {
	tmp := make([][]float64, 0)
	for _, matrix := range c {
		tmp = append(tmp, matrix)
	}
	return tmp
}

// ToCoordinateArray Creates an array containing the coordinates in this list, oriented
// in the given direction (forward or reverse).
// Params:
//		isForward – true if the direction is forward, false for reverse
// Returns:
//		an oriented array of coordinates
func (c CoordinateList) ToCoordinateArray(isForward bool) CoordinateList {
	if isForward {
		return c[0].Bound()
	}
	// construct reversed array
	pts := make([]Matrix, len(c))
	for i := 0; i < len(c); i++ {
		pts[i] = pts[len(c)-i-1]	// todo 是否数组越界
	}
	return pts
}
