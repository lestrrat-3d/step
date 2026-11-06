package ap214

import "github.com/lestrrat-3d/step"

// ColourRGB defines a colour by red, green, and blue channels in [0, 1].
func ColourRGB(id uint64, name string, r, g, b step.Real) step.Entity {
	return step.Entity{ID: id, Name: "COLOUR_RGB", Parameters: []step.Value{step.String(name), r, g, b}}
}

// FillAreaStyleColour fills an area with a colour.
func FillAreaStyleColour(id uint64, name string, colour step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "FILL_AREA_STYLE_COLOUR", Parameters: []step.Value{step.String(name), colour}}
}

// FillAreaStyle groups fill styles such as FillAreaStyleColour.
func FillAreaStyle(id uint64, name string, styles ...step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "FILL_AREA_STYLE", Parameters: []step.Value{step.String(name), refs(styles)}}
}

// SurfaceStyleFillArea applies a fill area style to a surface.
func SurfaceStyleFillArea(id uint64, fillArea step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "SURFACE_STYLE_FILL_AREA", Parameters: []step.Value{fillArea}}
}

// SurfaceSideStyle groups surface style elements such as SurfaceStyleFillArea.
func SurfaceSideStyle(id uint64, name string, styles ...step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "SURFACE_SIDE_STYLE", Parameters: []step.Value{step.String(name), refs(styles)}}
}

// SurfaceStyleUsage applies a side style to the given side of a surface. The
// side is an enumeration literal without dots: "POSITIVE", "NEGATIVE", or
// "BOTH".
func SurfaceStyleUsage(id uint64, side string, style step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "SURFACE_STYLE_USAGE", Parameters: []step.Value{step.Enumeration(side), style}}
}

// PresentationStyleAssignment collects styles such as SurfaceStyleUsage.
func PresentationStyleAssignment(id uint64, styles ...step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "PRESENTATION_STYLE_ASSIGNMENT", Parameters: []step.Value{refs(styles)}}
}

// StyledItem assigns presentation styles to a geometric item such as a solid.
func StyledItem(id uint64, name string, styles []step.Reference, item step.Reference) step.Entity {
	return step.Entity{ID: id, Name: "STYLED_ITEM", Parameters: []step.Value{step.String(name), refs(styles), item}}
}

// MechanicalDesignGeometricPresentationRepresentation collects styled items
// in a context.
func MechanicalDesignGeometricPresentationRepresentation(
	id uint64, name string, context step.Reference, items ...step.Reference,
) step.Entity {
	return step.Entity{ID: id, Name: "MECHANICAL_DESIGN_GEOMETRIC_PRESENTATION_REPRESENTATION", Parameters: []step.Value{
		step.String(name), refs(items), context,
	}}
}
