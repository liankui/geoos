package overlay_ng

import (
	"fmt"
	"github.com/spatial-go/geoos/encoding/wkt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnion(t *testing.T) {
	a, _ := wkt.UnmarshalString("POLYGON ((0 0, 10 0, 10 10, 0 10, 0 0))")
	b, _ := wkt.UnmarshalString("POLYGON ((5 5, 15 5, 15 15, 5 15, 5 5))")
	fmt.Printf("a=%+v\n", a)
	fmt.Printf("b=%+v\n", b)

	var ol GeometryOverlay
	result := ol.union(a, b)
	fmt.Printf("[result]=%+v\n", result)
	want, _ := wkt.UnmarshalString("POLYGON ((10 0, 0 0, 0 10, 5 10, 5 15, 15 15, 15 5, 10 5, 10 0))")
	assert.Equal(t, result, want)
}