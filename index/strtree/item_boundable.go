package strtree

import (
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
)

// IBoundable
type IBoundable interface {
	getBounds() []*envelope.Envelope
	getItem() interface{}
}

// ItemBoundable leaf node bounds.
// Boundable wrapper for a non-Boundable spatial object. Used internally by AbstractSTRtree.
type ItemBoundable struct {
	Bounds []*envelope.Envelope `json:"bounds"`
	Item   interface{}          `json:"item"`
}

func (i *ItemBoundable) getBounds() []*envelope.Envelope {
	if i != nil {
		return i.Bounds
	}
	return nil
}

func (i *ItemBoundable) getItem() interface{} {
	if i != nil {
		return i.Item
	}
	return nil
}
