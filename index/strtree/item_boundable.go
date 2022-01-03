package strtree

import (
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
)

// IBoundable
type IBoundable interface {
	getBounds() []*envelope.Envelope
	getItem() interface{}
}

// Boundable internal node bounds
type Boundable struct {
	bounds []*envelope.Envelope
}

func (b *Boundable) getBounds() []*envelope.Envelope {
	return b.bounds
}

func (b *Boundable) getItem() interface{} {
	return nil
}

// ItemBoundable leaf node bounds.
// Boundable wrapper for a non-Boundable spatial object. Used internally by AbstractSTRtree.
type ItemBoundable struct {
	bounds []*envelope.Envelope
	item   interface{}	// 标识叶子节点，可以为任意值
}

func (i *ItemBoundable) getBounds() []*envelope.Envelope {
	return i.bounds
}

func (i *ItemBoundable) getItem() interface{} {
	return i.item
}