package strtree

// Boundable wrapper for a non-Boundable spatial object. Used internally by AbstractSTRtree.
type ItemBoundable struct {
	Bounds interface{}
	Item   interface{}
}
