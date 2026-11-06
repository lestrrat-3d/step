package ap214_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/lestrrat-3d/step"
	"github.com/lestrrat-3d/step/ap214"
	"github.com/stretchr/testify/require"
)

func testHeader() step.Header {
	return step.Header{
		Description:   []string{"assembly"},
		Name:          "assembly.step",
		Timestamp:     time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC),
		Authors:       []string{"Example"},
		Organizations: []string{"Example"},
	}
}

func requireLines(t *testing.T, f step.File, lines []string) {
	t.Helper()
	encoded, err := f.Marshal()
	require.NoError(t, err)
	for _, line := range lines {
		require.Truef(t, bytes.Contains(encoded, []byte(line+"\n")), "missing record %s", line)
	}
}

func TestAssemblyRecords(t *testing.T) {
	t.Parallel()
	f := ap214.NewFile(testHeader(),
		ap214.Millimetre(1),
		ap214.Radian(2),
		ap214.Steradian(3),
		ap214.GeometricRepresentationContext(4, "3D", "model", 1, 2, 3),
		ap214.CartesianPoint(5, "", [3]step.Real{0, 0, 0}),
		ap214.Direction(6, "", [3]step.Real{0, 0, 1}),
		ap214.Direction(7, "", [3]step.Real{1, 0, 0}),
		ap214.Axis2Placement3D(8, "", 5, 6, 7),
		ap214.Axis2Placement3D(9, "", 5, 6, 7),
		ap214.ApplicationContext(10, "automotive design"),
		ap214.ProductDefinitionContext(11, "part definition", 10, "design"),
		ap214.ProductDefinition(12, "design", "", 100, 11),
		ap214.ProductDefinition(13, "design", "", 100, 11),
		ap214.AdvancedBrepShapeRepresentation(14, "child", 4, 8),
		ap214.ShapeRepresentation(15, "assembly", 4, 8, 9),
		ap214.NextAssemblyUsageOccurrence(16, "NAUO1", "left", "", 12, 13),
		ap214.ProductDefinitionShape(17, "placement", "", 16),
		ap214.ItemDefinedTransformation(18, "", "", 8, 9),
		ap214.ShapeRepresentationRelationshipWithTransformation(19, "", "", 14, 15, 18),
		ap214.ContextDependentShapeRepresentation(20, 19, 17),
		ap214.ProductDefinitionFormation(100, "1", "", 101),
		ap214.Product(101, "p", "p", ""),
	)
	requireLines(t, f, []string{
		"#15=SHAPE_REPRESENTATION('assembly',(#8,#9),#4);",
		"#16=NEXT_ASSEMBLY_USAGE_OCCURRENCE('NAUO1','left','',#12,#13,$);",
		"#18=ITEM_DEFINED_TRANSFORMATION('','',#8,#9);",
		"#19=(REPRESENTATION_RELATIONSHIP('','',#14,#15)" +
			"REPRESENTATION_RELATIONSHIP_WITH_TRANSFORMATION(#18)SHAPE_REPRESENTATION_RELATIONSHIP());",
		"#20=CONTEXT_DEPENDENT_SHAPE_REPRESENTATION(#19,#17);",
	})
	require.NoError(t, f.Write(&bytes.Buffer{}))
}

func TestAssemblyReferenceIntegrity(t *testing.T) {
	t.Parallel()
	f := ap214.NewFile(testHeader(),
		ap214.ItemDefinedTransformation(1, "", "", 8, 9),
		ap214.ShapeRepresentationRelationshipWithTransformation(2, "", "", 3, 4, 1),
	)
	var out bytes.Buffer
	require.Error(t, f.Write(&out))
	require.Zero(t, out.Len())
}

func TestPresentationRecords(t *testing.T) {
	t.Parallel()
	f := ap214.NewFile(testHeader(),
		ap214.ColourRGB(1, "", 1, 0.5, 0),
		ap214.FillAreaStyleColour(2, "", 1),
		ap214.FillAreaStyle(3, "", 2),
		ap214.SurfaceStyleFillArea(4, 3),
		ap214.SurfaceSideStyle(5, "", 4),
		ap214.SurfaceStyleUsage(6, "BOTH", 5),
		ap214.PresentationStyleAssignment(7, 6),
		ap214.CartesianPoint(8, "", [3]step.Real{0, 0, 0}),
		ap214.StyledItem(9, "color", []step.Reference{7}, 8),
		ap214.Millimetre(10),
		ap214.GeometricRepresentationContext(11, "3D", "model", 10),
		ap214.MechanicalDesignGeometricPresentationRepresentation(12, "", 11, 9),
		ap214.ClosedShell(13, "shell"),
		ap214.ClosedShell(14, "void"),
		ap214.BrepWithVoids(15, "solid", 13, 14),
		ap214.OrientedClosedShell(16, "", 14, false),
	)
	requireLines(t, f, []string{
		"#1=COLOUR_RGB('',1.,0.5,0.);",
		"#2=FILL_AREA_STYLE_COLOUR('',#1);",
		"#3=FILL_AREA_STYLE('',(#2));",
		"#4=SURFACE_STYLE_FILL_AREA(#3);",
		"#5=SURFACE_SIDE_STYLE('',(#4));",
		"#6=SURFACE_STYLE_USAGE(.BOTH.,#5);",
		"#7=PRESENTATION_STYLE_ASSIGNMENT((#6));",
		"#9=STYLED_ITEM('color',(#7),#8);",
		"#12=MECHANICAL_DESIGN_GEOMETRIC_PRESENTATION_REPRESENTATION('',(#9),#11);",
		"#15=BREP_WITH_VOIDS('solid',#13,(#14));",
		"#16=ORIENTED_CLOSED_SHELL('',*,#14,.F.);",
	})
	require.NoError(t, f.Write(&bytes.Buffer{}))
}
