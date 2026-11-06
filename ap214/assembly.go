package ap214

import "github.com/lestrrat-3d/step"

// ShapeRepresentation associates geometric items with a context. An
// assembly's shape representation lists the placement axes of its occurrences.
func ShapeRepresentation(id uint64, name string, context step.Reference, items ...step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "SHAPE_REPRESENTATION", Parameters: []step.Value{
		step.String(name), refs(items), context,
	}}
}

// NextAssemblyUsageOccurrence records one use of the related product
// definition inside the relating one. Its reference designator is null.
func NextAssemblyUsageOccurrence(
	id uint64, identifier, name, description string, relating, related step.Reference,
) step.Entity {
	return step.Entity{ID: id, Name: "NEXT_ASSEMBLY_USAGE_OCCURRENCE", Parameters: []step.Value{
		step.String(identifier), step.String(name), step.String(description), relating, related, step.Null{},
	}}
}

// ItemDefinedTransformation maps the axis placement item1 onto item2. The
// first item lies in the child representation and the second in the parent.
func ItemDefinedTransformation(id uint64, name, description string, item1, item2 step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "ITEM_DEFINED_TRANSFORMATION", Parameters: []step.Value{
		step.String(name), step.String(description), item1, item2,
	}}
}

// ShapeRepresentationRelationshipWithTransformation relates rep1 (the child)
// to rep2 (the parent) through a transformation such as
// ItemDefinedTransformation. It is the complex instance of
// REPRESENTATION_RELATIONSHIP, REPRESENTATION_RELATIONSHIP_WITH_TRANSFORMATION,
// and SHAPE_REPRESENTATION_RELATIONSHIP, with components in alphabetical order.
func ShapeRepresentationRelationshipWithTransformation(
	id uint64, name, description string, rep1, rep2, transformation step.Reference,
) step.Entity {
	return step.Entity{ID: id, Components: []step.Record{
		{Name: "REPRESENTATION_RELATIONSHIP", Parameters: []step.Value{
			step.String(name), step.String(description), rep1, rep2,
		}},
		{Name: "REPRESENTATION_RELATIONSHIP_WITH_TRANSFORMATION", Parameters: []step.Value{transformation}},
		{Name: "SHAPE_REPRESENTATION_RELATIONSHIP"},
	}}
}

// ContextDependentShapeRepresentation attaches a representation relationship
// to the product definition shape of an assembly usage occurrence.
func ContextDependentShapeRepresentation(id uint64, relationship, definitionShape step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "CONTEXT_DEPENDENT_SHAPE_REPRESENTATION", Parameters: []step.Value{
		relationship, definitionShape,
	}}
}
