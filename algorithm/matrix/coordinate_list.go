package matrix

import "reflect"

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
	for i := 1; i < len(coord); i++ { // todo 顺序问题
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
	tmp := c
	if direction {
		for i := 0; i < len(coord); i++ {
			tmp = tmp.AddToEndList(coord[i], allowRepeated)
		}
	} else {
		for i := len(coord) - 1; i >= 0; i-- {
			tmp = tmp.AddToEndList(coord[i], allowRepeated)
		}
	}
	return tmp
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
	tmp := c
	tmp = append(tmp, coord)
	return tmp
}
