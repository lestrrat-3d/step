package ap214

import "github.com/lestrrat-3d/step"

// Millimetre defines the SI length unit used by the context's coordinates.
func Millimetre(id uint64) step.Entity {
	return step.Entity{ID: id, Components: []step.Record{
		{Name: "LENGTH_UNIT"},
		{Name: "NAMED_UNIT", Parameters: []step.Value{step.Omitted{}}},
		{Name: "SI_UNIT", Parameters: []step.Value{step.Enumeration("MILLI"), step.Enumeration("METRE")}},
	}}
}

// Radian defines the SI plane-angle unit.
func Radian(id uint64) step.Entity {
	return step.Entity{ID: id, Components: []step.Record{
		{Name: "NAMED_UNIT", Parameters: []step.Value{step.Omitted{}}},
		{Name: "PLANE_ANGLE_UNIT"},
		{Name: "SI_UNIT", Parameters: []step.Value{step.Null{}, step.Enumeration("RADIAN")}},
	}}
}

// Steradian defines the SI solid-angle unit.
func Steradian(id uint64) step.Entity {
	return step.Entity{ID: id, Components: []step.Record{
		{Name: "NAMED_UNIT", Parameters: []step.Value{step.Omitted{}}},
		{Name: "SI_UNIT", Parameters: []step.Value{step.Null{}, step.Enumeration("STERADIAN")}},
		{Name: "SOLID_ANGLE_UNIT"},
	}}
}

// GeometricRepresentationContext defines a 3D context with assigned units.
// The unit references normally include one length, one plane-angle, and one
// solid-angle unit, such as Millimetre, Radian, and Steradian.
func GeometricRepresentationContext(id uint64, identifier, contextType string, units ...step.Reference) step.Entity {
	return step.Entity{ID: id, Components: []step.Record{
		{Name: "GEOMETRIC_REPRESENTATION_CONTEXT", Parameters: []step.Value{step.Integer(3)}},
		{Name: "GLOBAL_UNIT_ASSIGNED_CONTEXT", Parameters: []step.Value{refs(units)}},
		{Name: "REPRESENTATION_CONTEXT", Parameters: []step.Value{step.String(identifier), step.String(contextType)}},
	}}
}

// AdvancedBrepShapeRepresentation associates B-rep solids with a context.
func AdvancedBrepShapeRepresentation(id uint64, name string, context step.Reference, items ...step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "ADVANCED_BREP_SHAPE_REPRESENTATION", Parameters: []step.Value{
		step.String(name), refs(items), context,
	}}
}

// ShapeDefinitionRepresentation attaches a representation to a product shape.
func ShapeDefinitionRepresentation(id uint64, definition, representation step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "SHAPE_DEFINITION_REPRESENTATION", Parameters: []step.Value{
		definition, representation,
	}}
}
