package ap214

import "github.com/lestrrat-3d/step"

// ApplicationContext names the application for product definitions.
func ApplicationContext(id uint64, application string) step.Entity {
	return step.Entity{ID: id, Name: "APPLICATION_CONTEXT", Parameters: []step.Value{step.String(application)}}
}

// ApplicationProtocolDefinition identifies the AP214 protocol and its year.
func ApplicationProtocolDefinition(id uint64, status string, year int64, application step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "APPLICATION_PROTOCOL_DEFINITION", Parameters: []step.Value{
		step.String(status), step.String("automotive_design"), step.Integer(year), application,
	}}
}

// ProductContext defines a discipline within an application context.
func ProductContext(id uint64, name string, application step.Reference, discipline string) step.Entity {
	return step.Entity{ID: id, Name: "PRODUCT_CONTEXT", Parameters: []step.Value{
		step.String(name), application, step.String(discipline),
	}}
}

// Product defines a product and its applicable contexts.
func Product(id uint64, identifier, name, description string, contexts ...step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "PRODUCT", Parameters: []step.Value{
		step.String(identifier), step.String(name), step.String(description), refs(contexts),
	}}
}

// ProductRelatedProductCategory groups products under a category such as part.
func ProductRelatedProductCategory(id uint64, name string, products ...step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "PRODUCT_RELATED_PRODUCT_CATEGORY", Parameters: []step.Value{
		step.String(name), step.Null{}, refs(products),
	}}
}

// ProductDefinitionFormation identifies a version of a product.
func ProductDefinitionFormation(id uint64, identifier, description string, product step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "PRODUCT_DEFINITION_FORMATION", Parameters: []step.Value{
		step.String(identifier), step.String(description), product,
	}}
}

// ProductDefinitionContext names a product life-cycle context.
func ProductDefinitionContext(id uint64, name string, application step.Reference, lifeCycle string) step.Entity {
	return step.Entity{ID: id, Name: "PRODUCT_DEFINITION_CONTEXT", Parameters: []step.Value{
		step.String(name), application, step.String(lifeCycle),
	}}
}

// ProductDefinition joins a product formation to a definition context.
func ProductDefinition(id uint64, identifier, description string, formation, context step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "PRODUCT_DEFINITION", Parameters: []step.Value{
		step.String(identifier), step.String(description), formation, context,
	}}
}

// ProductDefinitionShape describes the shape of a product definition.
func ProductDefinitionShape(id uint64, name, description string, definition step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "PRODUCT_DEFINITION_SHAPE", Parameters: []step.Value{
		step.String(name), step.String(description), definition,
	}}
}
