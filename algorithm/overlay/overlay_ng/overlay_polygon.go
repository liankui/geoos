package overlay_ng

import (
	"github.com/spatial-go/geoos/algorithm"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/algorithm/overlay"
	"github.com/spatial-go/geoos/algorithm/relate"
)

// PolygonOverlay  Computes the overlay of two geometries,either or both of which may be nil.
type PolygonOverlay struct {
	*overlay.PointOverlay
	subjectPlane, clippingPlane *overlay.Plane
}

// isPolyInHole
func (p *PolygonOverlay) isPolyInHole(polyMatrix matrix.PolygonMatrix) (inHole bool) {
	subjectPoly, ok := p.Subject.(matrix.PolygonMatrix)
	if !ok {
		return
	}
	for i := range subjectPoly {
		if i == 0 {
			continue
		}
		subjectMatrix := matrix.PolygonMatrix{subjectPoly[i]}
		inter := envelope.Bound(polyMatrix.Bound()).IsIntersects(envelope.Bound(subjectMatrix.Bound()))
		im := relate.IM(subjectMatrix, polyMatrix, inter)
		if mark := im.IsContains(); mark {
			inHole = true
			break
		}
	}
	if inHole {
		return
	}
	clippingPoly, ok := p.Clipping.(matrix.PolygonMatrix)
	if !ok {
		return
	}
	for i := range clippingPoly {
		if i == 0 {
			continue
		}
		clippingMatrix := matrix.PolygonMatrix{clippingPoly[i]}
		inter := envelope.Bound(polyMatrix.Bound()).IsIntersects(envelope.Bound(clippingMatrix.Bound()))
		im := relate.IM(clippingMatrix, polyMatrix, inter)
		if mark := im.IsContains(); mark {
			inHole = true
			break
		}
	}
	return
}

// Union  Computes the Union of two geometries,either or both of which may be nil.
func (p *PolygonOverlay) Union() (matrix.Steric, error) {
	//if res, ok := p.unionCheck(); !ok {
	//	return res, nil
	//}
	if ps, ok := p.Subject.(matrix.PolygonMatrix); ok {
		if pc, ok := p.Clipping.(matrix.PolygonMatrix); ok {
			inter := envelope.Bound(ps.Bound()).IsIntersects(envelope.Bound(pc.Bound()))
			im := relate.IM(ps, pc, inter)
			if mark := im.IsCovers(); mark {
				return ps, nil
			}
			if mark := im.IsCoveredBy(); mark {
				return pc, nil
			}
			if mark := !im.IsIntersects(); mark {
				return matrix.Collection{p.Subject.(matrix.PolygonMatrix), p.Clipping.(matrix.PolygonMatrix)}, nil
			}
			if mark, _ := im.Matches("FF**0****"); mark {
				return matrix.Collection{p.Subject.(matrix.PolygonMatrix), p.Clipping.(matrix.PolygonMatrix)}, nil
			}
			//
			//cpo := &ComputeMergeOverlay{p}
			//
			//cpo.prepare()
			//_, exitingPoints := cpo.Weiler()
			//
			//result := ToPolygonMatrix(cpo.ComputePolygon(exitingPoints, cpo))
			var result matrix.Steric
			return result, nil
		}
	}
	return nil, algorithm.ErrNotMatchType
}

// Intersection  Computes the Intersection of two geometries,either or both of which may be nil.
func (p *PolygonOverlay) Intersection() (matrix.Steric, error) {
	op := &OverlayOp{}

	// special case: if one input is empty ==> empty
	if p.Subject.IsEmpty() || p.Clipping.IsEmpty() {
		//return op.CreateEmptyResult(INTERSECTION, p.Subject, p.Clipping)
	}

	//switch p.Subject.(type) {
	//case matrix.Collection:
	//
	//}

	return op.Overlay(p.Subject, p.Clipping, INTERSECTION)
}

// Difference returns a geometry that represents that part of geometry A that does not intersect with geometry B.
// One can think of this as GeometryA - Intersection(A,B).
// If A is completely contained in B then an empty geometry collection is returned.
func (p *PolygonOverlay) Difference() (matrix.Steric, error) {
	//if res, ok := p.differenceCheck(); !ok {
	//	return res, nil
	//}
	if poly, ok := p.Subject.(matrix.PolygonMatrix); ok {
		if c, ok := p.Clipping.(matrix.PolygonMatrix); ok {

			inter := envelope.Bound(poly.Bound()).IsIntersects(envelope.Bound(c.Bound()))
			im := relate.IM(poly, c, inter)
			if mark := im.IsCoveredBy(); mark {
				return matrix.PolygonMatrix{}, nil

			}
			if mark, _ := im.Matches("212FF1FF2"); mark {
				poly = append(poly, c...)
				return poly, nil
			}

			//cpo := &ComputeMainOverlay{p}
			//
			//cpo.prepare()
			//_, exitingPoints := cpo.Weiler()
			//result := ToPolygonMatrix(cpo.ComputePolygon(exitingPoints, cpo))
			var result matrix.PolygonMatrix
			return result, nil
		}
	}
	return nil, algorithm.ErrNotMatchType
}

// DifferenceReverse returns a geometry that represents reverse that part of geometry A that does not intersect with geometry B .
// One can think of this as GeometryB - Intersection(A,B).
// If B is completely contained in A then an empty geometry collection is returned.
func (p *PolygonOverlay) DifferenceReverse() (matrix.Steric, error) {
	newPoly := &PolygonOverlay{PointOverlay: &overlay.PointOverlay{Subject: p.Clipping, Clipping: p.Subject}}
	return newPoly.Difference()
}

// SymDifference returns a geometry that represents the portions of A and B that do not intersect.
// It is called a symmetric difference because SymDifference(A,B) = SymDifference(B,A).
// One can think of this as Union(geomA,geomB) - Intersection(A,B).
func (p *PolygonOverlay) SymDifference() (matrix.Steric, error) {

	result := matrix.Collection{}
	if res, err := p.Difference(); err == nil && !res.IsEmpty() {
		result = append(result, res)
	}
	if res, err := p.DifferenceReverse(); err == nil && !res.IsEmpty() {
		result = append(result, res)
	}
	return result, nil
}



