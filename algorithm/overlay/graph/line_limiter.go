package graph

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
)

// Limits the segments in a list of segments to those which intersect an envelope.
// This creates zero or more sections of the input segment sequences, containing
// only line segments which intersect the limit envelope. Segments are not clipped,
// since that can move line segments enough to alter topology, and it happens in
// the overlay in any case. This can substantially reduce the number of vertices
// which need to be processed during overlay.
// This optimization is only applicable to Line geometries, since it does not maintain
// the closed topology of rings. Polygonal geometries are optimized using the RingClipper.
type LineLimiter struct {
	limitEnv    *envelope.Envelope
	ptList      []matrix.Matrix
	lastOutside matrix.Matrix
}
