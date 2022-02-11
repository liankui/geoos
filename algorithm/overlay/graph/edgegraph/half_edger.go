package edgegraph

type HalfEdger interface {
	//ToString() string
	//DirectionPt() matrix.Matrix

	Link(sym *HalfEdge)
	ONext() *HalfEdge
	Sym() *HalfEdge
	Insert(eAdd *HalfEdge)
}

// HalfEdgerLink...
func HalfEdgerLink(h, sym HalfEdger) {
	h.Link(sym.(*HalfEdge))
}

// HalfEdgerONext...
func HalfEdgerONext(h HalfEdger) HalfEdger {
	return h.ONext()
}

// HalfEdgerSym...
func HalfEdgerSym(h HalfEdger) HalfEdger {
	return h.ONext()
}

// HalfEdgerInsert...
func HalfEdgerInsert(h, eAdd HalfEdger) {
	h.Insert(eAdd.(*HalfEdge))
}
