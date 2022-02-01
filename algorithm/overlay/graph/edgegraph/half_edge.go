package edgegraph

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"reflect"
)

type IHalfEdge interface {
	ToString() string
	DirectionPt() matrix.Matrix
}

// Represents a directed component of an edge in an EdgeGraph. HalfEdges link
// vertices whose locations are defined by Coordinates. HalfEdges start at an
// origin vertex, and terminate at a destination vertex. HalfEdges always occur
// in symmetric pairs, with the sym() method giving access to the oppositely-oriented
// component. HalfEdges and the methods on them form an edge algebra, which can
// be used to traverse and query the topology of the graph formed by the edges.
// To support graphs where the edges are sequences of coordinates each edge may also
// have a direction point supplied. This is used to determine the ordering of the edges
// around the origin. HalfEdges with the same origin are ordered so that the ring of
// edges formed by them is oriented CCW.
// By design HalfEdges carry minimal information about the actual usage of the graph
// they represent. They can be subclassed to carry more information if required.
// HalfEdges form a complete and consistent data structure by themselves, but an
// EdgeGraph is useful to allow retrieving edges by vertex and edge location,
// as well as ensuring edges are created and linked appropriately.
type HalfEdge struct {
	IHalfEdge
	orig      matrix.Matrix
	sym, next *HalfEdge
}

// Link Links this edge with its sym (opposite) edge.
// This must be done for each pair of edges created.
// Params:
//		sym – the sym edge to link.
func (h *HalfEdge) Link(sym IHalfEdge) {
	h.setSym(sym.(*HalfEdge))
	sym.(*HalfEdge).setSym(h)
	// set next ptrs for a single segment
	h.setNext(sym.(*HalfEdge))
	sym.(*HalfEdge).setNext(h)
}

// Sym Gets the symmetric pair edge of this edge.
// Returns:
//		the symmetric pair edge
func (h *HalfEdge) Sym(e *HalfEdge) *HalfEdge {
	return h.sym
}

// setSym Sets the symmetric (opposite) edge to this edge.
// Params:
//		e – the sym edge to set
func (h *HalfEdge) setSym(e *HalfEdge) {
	h.sym = e
}

// setNext the next edge CCW around the destination vertex of this edge.
// Params:
//		e – the next edge
func (h *HalfEdge) setNext(e *HalfEdge) {
	h.next = e
}

// oNext Gets the next edge CCW around the origin of this edge, with the
// same origin. If the origin vertex has degree 1 then this is the edge itself.
// e.oNext() is equal to e.sym().next()
// Returns:
//		the next edge around the origin
func (h *HalfEdge) oNext() *HalfEdge {
	return h.sym.next
}

// insert inserts an edge into the ring of edges around the origin vertex
// of this edge, ensuring that the edges remain ordered CCW. The inserted
// edge must have the same origin as this edge.
// Params:
//		eAdd – the edge to insert
func (h *HalfEdge) Insert(eAdd *HalfEdge) {
	// If this is only edge at origin, insert it after this
	if h.oNext() == h {
		// set linkage so ring is correct
		h.insertAfter(eAdd)
		return
	}
	// Scan edges until insertion point is found
	ePrev := h.insertionEdge(eAdd)
	ePrev.insertAfter(eAdd)
}

// insertionEdge Finds the insertion edge for a edge being added to this origin,
// ensuring that the star of edges around the origin remains fully CCW.
// Params:
//		eAdd – the edge being added
// Returns:
//		the edge to insert after
func (h *HalfEdge) insertionEdge(eAdd *HalfEdge) *HalfEdge {
	ePrev := h
	for ePrev != h {
		eNext := ePrev.oNext()
		/**
		 * Case 1: General case,
		 * with eNext higher than ePrev.
		 *
		 * Insert edge here if it lies between ePrev and eNext.
		 */
		if eNext.compareTo(ePrev) > 0 &&
			eAdd.compareTo(ePrev) >= 0 &&
			eAdd.compareTo(eNext) <= 0 {
			return ePrev
		}
		/**
		 * Case 2: Origin-crossing case,
		 * indicated by eNext <= ePrev.
		 *
		 * Insert edge here if it lies
		 * in the gap between ePrev and eNext across the origin.
		 */
		if eNext.compareTo(ePrev) <= 0 &&
			(eAdd.compareTo(eNext) <= 0 || eAdd.compareTo(ePrev) >= 0) {
			return ePrev
		}
		ePrev = eNext
	}
	return nil
}

// insertAfter Insert an edge with the same origin after this one.
// Assumes that the inserted edge is in the correct position around the ring.
// Params:
//		e – the edge to insert (with same origin)
func (h *HalfEdge) insertAfter(e *HalfEdge) {
	if !reflect.DeepEqual(h.orig, e.orig) {
		return
	}
	save := h.oNext()
	h.sym.setNext(e)
	e.sym.setNext(save)
}

// compareTo Compares edges which originate at the same vertex based on the
// angle they make at their origin vertex with the positive X-axis. This allows
// sorting edges around their origin vertex in CCW order.
func (h *HalfEdge) compareTo(obj interface{}) int {
	e := obj.(*HalfEdge)
	return h.compareAngularDirection(e)
}

// compareAngularDirection Implements the total order relation:
// The angle of edge a is greater than the angle of edge b, where the angle of
// an edge is the angle made by the first segment of the edge with the positive x-axis.
// When applied to a list of edges originating at the same point, this produces a CCW
// ordering of the edges around the point.
// Using the obvious algorithm of computing the angle is not robust, since the angle
// calculation is susceptible to roundoff error. A robust algorithm is:
// First, compare the quadrants the edge vectors lie in. If the quadrants are different,
// it is trivial to determine which edge has a greater angle.
// if the vectors lie in the same quadrant, the Orientation.index(Coordinate, Coordinate,
// Coordinate) function can be used to determine the relative orientation of the vectors.
func (h *HalfEdge) compareAngularDirection(e *HalfEdge) int {
	// todo
	return 0
}

// DirectionPt...
func (h *HalfEdge) DirectionPt() matrix.Matrix {
	return h.sym.orig
}

// ToString Provides a string representation of a HalfEdge.
// Returns:
//		a string representation
func (h *HalfEdge) ToString() string {
	return fmt.Sprintf("HE(%v %v, %v %v)",
		h.orig[0], h.orig[0], h.sym.orig[0], h.sym.orig[0])
}
