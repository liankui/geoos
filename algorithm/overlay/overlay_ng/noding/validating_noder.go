package noding

import "fmt"

// A wrapper for Noders which validates the output arrangement is correctly noded.
// An arrangement of line segments is fully noded if there is no line segment
// which has another segment intersecting its interior. If the noding is not correct,
// a org.locationtech.jts.geom.TopologyException is thrown with details of the first
// invalid location found. See Also: FastNodingValidator
type ValidatingNoder struct {
	noder   Noder
	nodedSS []SegmentString
}

// Creates a noding validator wrapping the given Noder
// Params:
//		noder – the Noder to validate
func NewValidatingNoder(noder Noder) *ValidatingNoder {
	return &ValidatingNoder{
		noder: noder,
	}
}

// ComputeNodes Checks whether the output of the wrapped noder is fully noded.
// Throws an exception if it is not.
func (v *ValidatingNoder) ComputeNodes(segStrings interface{}) {
	fmt.Println("====computeNodes9")
	v.noder.ComputeNodes(segStrings)
	fmt.Println("====computeNodes9-2")
	nodedSS := v.noder.GetNodedSubstrings()
	fmt.Println("====computeNodes9-3,nodedSS=", nodedSS)
	//v.validate() todo
}

// validate...
func (v *ValidatingNoder) validate() {

}

// GetNodedSubstrings...
func (v *ValidatingNoder) GetNodedSubstrings() interface{} {
	fmt.Println("GetNodedSubstrings 9")
	return v.nodedSS
}
