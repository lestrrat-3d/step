package ap214

import "github.com/lestrrat-3d/step"

// VertexPoint associates a topological vertex with a geometric point.
func VertexPoint(id uint64, name string, point step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "VERTEX_POINT", Parameters: []step.Value{step.String(name), point}}
}

// EdgeCurve associates an edge's endpoints with its geometric curve.
func EdgeCurve(id uint64, name string, start, end, curve step.Reference, sameSense bool) step.Entity {
	return step.Entity{ID: id, Name: "EDGE_CURVE", Parameters: []step.Value{
		step.String(name), start, end, curve, logical(sameSense),
	}}
}

// OrientedEdge uses derived endpoint attributes (*) and a referenced edge.
func OrientedEdge(id uint64, name string, edge step.Reference, orientation bool) step.Entity {
	return step.Entity{ID: id, Name: "ORIENTED_EDGE", Parameters: []step.Value{
		step.String(name), step.Omitted{}, step.Omitted{}, edge, logical(orientation),
	}}
}

// EdgeLoop is an ordered circuit of oriented edges.
func EdgeLoop(id uint64, name string, edges ...step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "EDGE_LOOP", Parameters: []step.Value{
		step.String(name), refs(edges),
	}}
}

// FaceOuterBound marks a face's outer loop.
func FaceOuterBound(id uint64, name string, loop step.Reference, orientation bool) step.Entity {
	return step.Entity{ID: id, Name: "FACE_OUTER_BOUND", Parameters: []step.Value{
		step.String(name), loop, logical(orientation),
	}}
}

// FaceBound marks a face's inner or other loop.
func FaceBound(id uint64, name string, loop step.Reference, orientation bool) step.Entity {
	return step.Entity{ID: id, Name: "FACE_BOUND", Parameters: []step.Value{
		step.String(name), loop, logical(orientation),
	}}
}

// AdvancedFace bounds an analytic surface by face bounds.
func AdvancedFace(id uint64, name string, surface step.Reference, sameSense bool, bounds ...step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "ADVANCED_FACE", Parameters: []step.Value{
		step.String(name), refs(bounds), surface, logical(sameSense),
	}}
}

// ClosedShell collects the faces of a closed shell.
func ClosedShell(id uint64, name string, faces ...step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "CLOSED_SHELL", Parameters: []step.Value{
		step.String(name), refs(faces),
	}}
}

// ManifoldSolidBrep names the outer shell of a solid.
func ManifoldSolidBrep(id uint64, name string, outer step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "MANIFOLD_SOLID_BREP", Parameters: []step.Value{
		step.String(name), outer,
	}}
}
