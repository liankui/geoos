package strtree

import (
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
)

type STRtree struct {
}

func centreX(e envelope.Envelope) float64 {
	return avg(e.MinX, e.MaxX)
}

func centreY(e envelope.Envelope) float64 {
	return avg(e.MinY, e.MaxY)
}

func avg(a, b float64) float64 {
	return (a + b) / 2.0
}

func intersects(a,b *envelope.Envelope) bool {
	return a.IsIntersects(b)
}