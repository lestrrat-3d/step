package ap214_test

import (
	"strings"
	"testing"
	"time"

	"github.com/lestrrat-3d/step"
	"github.com/lestrrat-3d/step/ap214"
	"github.com/stretchr/testify/require"
)

func TestAP214Records(t *testing.T) {
	t.Parallel()
	entities := []step.Entity{
		ap214.Millimetre(1),
		ap214.Radian(2),
		ap214.Steradian(3),
		ap214.GeometricRepresentationContext(4, "3D", "model", 1, 2, 3),
		ap214.CartesianPoint(5, "origin", [3]step.Real{0, 0, 0}),
		ap214.Direction(6, "z", [3]step.Real{0, 0, 1}),
		ap214.Direction(7, "x", [3]step.Real{1, 0, 0}),
		ap214.Axis2Placement3D(8, "frame", 5, 6, 7),
		ap214.Plane(9, "plane", 8),
		ap214.CylindricalSurface(10, "cylinder", 8, 2),
		ap214.Circle(11, "circle", 8, 2),
		ap214.Vector(12, "vector", 6, 4),
		ap214.Line(13, "line", 5, 12),
		ap214.VertexPoint(14, "vertex", 5),
		ap214.EdgeCurve(15, "edge", 14, 14, 13, true),
		ap214.OrientedEdge(16, "forward", 15, false),
		ap214.EdgeLoop(17, "loop", 16),
		ap214.FaceOuterBound(18, "outer", 17, true),
		ap214.FaceBound(19, "inner", 17, false),
		ap214.AdvancedFace(20, "face", 9, true, 18, 19),
		ap214.ClosedShell(21, "shell", 20),
		ap214.ManifoldSolidBrep(22, "brep", 21),
		ap214.AdvancedBrepShapeRepresentation(23, "shape", 4, 22),
		ap214.ApplicationContext(24, "automotive design"),
		ap214.ApplicationProtocolDefinition(25, "international standard", 2003, 24),
		ap214.ProductContext(26, "mechanical parts", 24, "mechanical"),
		ap214.Product(27, "part-1", "Part", "example", 26),
		ap214.ProductRelatedProductCategory(28, "part", 27),
		ap214.ProductDefinitionFormation(29, "1", "", 27),
		ap214.ProductDefinitionContext(30, "part definition", 24, "design"),
		ap214.ProductDefinition(31, "design", "", 29, 30),
		ap214.ProductDefinitionShape(32, "shape", "", 31),
		ap214.ShapeDefinitionRepresentation(33, 32, 23),
	}
	header := step.Header{
		Description:   []string{"AP214 records"},
		Name:          "records.step",
		Timestamp:     time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC),
		Authors:       []string{"Example"},
		Organizations: []string{"Example"},
		Schemas:       []string{"WRONG_SCHEMA"},
	}
	f := ap214.NewFile(header, entities...)
	require.Equal(t, []string{ap214.Schema}, f.Header.Schemas)
	require.Equal(t, []string{"WRONG_SCHEMA"}, header.Schemas)
	encoded, err := f.Marshal()
	require.NoError(t, err)
	out := string(encoded)
	require.Contains(t, out, "FILE_SCHEMA(('AUTOMOTIVE_DESIGN'));\n")
	for _, line := range []string{
		"#1=(LENGTH_UNIT()NAMED_UNIT(*)SI_UNIT(.MILLI.,.METRE.));",
		"#2=(NAMED_UNIT(*)PLANE_ANGLE_UNIT()SI_UNIT($,.RADIAN.));",
		"#3=(NAMED_UNIT(*)SI_UNIT($,.STERADIAN.)SOLID_ANGLE_UNIT());",
		"#4=(GEOMETRIC_REPRESENTATION_CONTEXT(3)GLOBAL_UNIT_ASSIGNED_CONTEXT((#1,#2,#3))REPRESENTATION_CONTEXT('3D','model'));",
		"#5=CARTESIAN_POINT('origin',(0.,0.,0.));",
		"#6=DIRECTION('z',(0.,0.,1.));",
		"#8=AXIS2_PLACEMENT_3D('frame',#5,#6,#7);",
		"#9=PLANE('plane',#8);",
		"#10=CYLINDRICAL_SURFACE('cylinder',#8,2.);",
		"#11=CIRCLE('circle',#8,2.);",
		"#12=VECTOR('vector',#6,4.);",
		"#13=LINE('line',#5,#12);",
		"#14=VERTEX_POINT('vertex',#5);",
		"#15=EDGE_CURVE('edge',#14,#14,#13,.T.);",
		"#16=ORIENTED_EDGE('forward',*,*,#15,.F.);",
		"#17=EDGE_LOOP('loop',(#16));",
		"#18=FACE_OUTER_BOUND('outer',#17,.T.);",
		"#19=FACE_BOUND('inner',#17,.F.);",
		"#20=ADVANCED_FACE('face',(#18,#19),#9,.T.);",
		"#21=CLOSED_SHELL('shell',(#20));",
		"#22=MANIFOLD_SOLID_BREP('brep',#21);",
		"#23=ADVANCED_BREP_SHAPE_REPRESENTATION('shape',(#22),#4);",
		"#24=APPLICATION_CONTEXT('automotive design');",
		"#25=APPLICATION_PROTOCOL_DEFINITION('international standard','automotive_design',2003,#24);",
		"#26=PRODUCT_CONTEXT('mechanical parts',#24,'mechanical');",
		"#27=PRODUCT('part-1','Part','example',(#26));",
		"#28=PRODUCT_RELATED_PRODUCT_CATEGORY('part',$,(#27));",
		"#29=PRODUCT_DEFINITION_FORMATION('1','',#27);",
		"#30=PRODUCT_DEFINITION_CONTEXT('part definition',#24,'design');",
		"#31=PRODUCT_DEFINITION('design','',#29,#30);",
		"#32=PRODUCT_DEFINITION_SHAPE('shape','',#31);",
		"#33=SHAPE_DEFINITION_REPRESENTATION(#32,#23);",
	} {
		require.Truef(t, strings.Contains(out, line+"\n"), "missing record %s", line)
	}
}
