package ap214

import "github.com/lestrrat-3d/step"

// CartesianPoint describes a point in the representation context's length
// unit. AP214 allows other coordinate dimensions; this constructor is 3D.
func CartesianPoint(id uint64, name string, xyz [3]step.Real) step.Entity {
	return step.Entity{ID: id, Name: "CARTESIAN_POINT", Parameters: []step.Value{
		step.String(name), step.List{xyz[0], xyz[1], xyz[2]},
	}}
}

// Direction describes direction ratios. The ratios must not all be zero.
func Direction(id uint64, name string, ratios [3]step.Real) step.Entity {
	return step.Entity{ID: id, Name: "DIRECTION", Parameters: []step.Value{
		step.String(name), step.List{ratios[0], ratios[1], ratios[2]},
	}}
}

// Axis2Placement3D places an origin and two directions. AP214 permits absent
// directions, but this constructor uses explicit axis and reference direction.
func Axis2Placement3D(id uint64, name string, location, axis, refDirection step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "AXIS2_PLACEMENT_3D", Parameters: []step.Value{
		step.String(name), location, axis, refDirection,
	}}
}

// Plane defines an infinite plane from a 3D placement.
func Plane(id uint64, name string, placement step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "PLANE", Parameters: []step.Value{step.String(name), placement}}
}

// CylindricalSurface defines a cylinder; radius uses the context's length unit.
func CylindricalSurface(id uint64, name string, placement step.Reference, radius step.Real) step.Entity {
	return step.Entity{ID: id, Name: "CYLINDRICAL_SURFACE", Parameters: []step.Value{
		step.String(name), placement, radius,
	}}
}

// Circle defines a circle; radius uses the context's length unit.
func Circle(id uint64, name string, placement step.Reference, radius step.Real) step.Entity {
	return step.Entity{ID: id, Name: "CIRCLE", Parameters: []step.Value{
		step.String(name), placement, radius,
	}}
}

// Vector combines a direction with a magnitude in the context's length unit.
func Vector(id uint64, name string, orientation step.Reference, magnitude step.Real) step.Entity {
	return step.Entity{ID: id, Name: "VECTOR", Parameters: []step.Value{
		step.String(name), orientation, magnitude,
	}}
}

// Line defines a line from a point and a vector.
func Line(id uint64, name string, point, direction step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "LINE", Parameters: []step.Value{
		step.String(name), point, direction,
	}}
}
