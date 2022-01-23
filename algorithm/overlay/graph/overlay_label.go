package graph

// A structure recording the topological situation for an edge in a topology graph
// used during overlay processing. A label contains the topological Locations for
// one or two input geometries to an overlay operation. An input geometry may be
// either a Line or an Area. The label locations for each input geometry are populated
// with the s for the edge Positions when they are created or once they are computed by
// topological evaluation. A label also records the (effective) dimension of each input
// geometry. For area edges the role (shell or hole) of the originating ring is recorded,
// to allow determination of edge handling in collapse cases.
// In an OverlayGraph a single label is shared between the two oppositely-oriented
type OverlayLabel struct {

}
