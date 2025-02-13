package wkt

import (
	"testing"

	"github.com/spatial-go/geoos/space"
)

func TestMarshalString(t *testing.T) {
	type args struct {
		geom space.Geometry
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "marshal string", args: args{space.LineString{{50, 100}, {50, 200}}},
			want: "LINESTRING(50 100,50 200)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MarshalString(tt.args.geom); got != tt.want {
				t.Errorf("MarshalString() = %v, want %v", got, tt.want)
			}
		})
	}
}
