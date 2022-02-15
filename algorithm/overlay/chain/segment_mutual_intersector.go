package chain

import "github.com/spatial-go/geoos/algorithm/matrix"

// SegmentMutualIntersector Intersects two sets of SegmentString.
type SegmentMutualIntersector struct {
	SegmentMutual matrix.LineMatrix
}

// Process the given collection of SegmentStrings and the set of indexed segments.
func (s *SegmentMutualIntersector) Process(segStrings matrix.LineMatrix, segInt Intersector) {
	monoChains, testChains := []*MonotoneChain{}, []*MonotoneChain{}

	monoChains = s.AddToMonoChains(s.SegmentMutual, "subject", monoChains)

	testChains = s.AddToMonoChains(segStrings, "test", testChains)

	s.IntersectChains(monoChains, testChains, segInt)
}

// AddToMonoChains ...
func (s *SegmentMutualIntersector) AddToMonoChains(segMatrix matrix.LineMatrix, context interface{}, monoChains []*MonotoneChain) []*MonotoneChain {
	segChains := ChainsContext(segMatrix, context)

	monoChains = append(monoChains, segChains...)
	return monoChains
}

// IntersectChains ...
func (s *SegmentMutualIntersector) IntersectChains(monoChains []*MonotoneChain, testChains []*MonotoneChain, segInt Intersector) {
	overlapAction := &SegmentOverlapAction{MonotoneChainOverlapAction: &MonotoneChainOverlapAction{}, Si: segInt}

	for _, queryChain := range monoChains {
		for _, testChain := range testChains {
			queryChain.ComputeOverlaps(testChain, overlapAction)
			if segInt.IsDone() {
				return
			}
		}
	}
}

// SegmentOverlapAction implement OverlapAction.
type SegmentOverlapAction struct {
	*MonotoneChainOverlapAction
	Si Intersector
}

// NewSegmentOverlapAction...
func NewSegmentOverlapAction(si Intersector) *SegmentOverlapAction {
	return &SegmentOverlapAction{
		Si: si,
	}
}

// Overlap This function can be overridden if the original chains are needed.
func (s *SegmentOverlapAction) Overlap(mc1 *MonotoneChain, start1 int, mc2 *MonotoneChain, start2 int) {
	ss1 := mc1.Edge
	ss2 := mc2.Edge
	s.Si.ProcessIntersections(ss1, start1, ss2, start2)
}
