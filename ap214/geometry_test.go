package ap214_test

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/lestrrat-3d/step"
	"github.com/lestrrat-3d/step/ap214"
	"github.com/stretchr/testify/require"
)

func TestCurvedGeometryRecords(t *testing.T) {
	t.Parallel()
	f := ap214.NewFile(
		step.Header{
			Description:   []string{"curved geometry"},
			Name:          "curved.step",
			Timestamp:     time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC),
			Authors:       []string{"Example"},
			Organizations: []string{"Example"},
		},
		ap214.CartesianPoint(5, "origin", [3]step.Real{0, 0, 0}),
		ap214.Direction(6, "z", [3]step.Real{0, 0, 1}),
		ap214.Direction(7, "x", [3]step.Real{1, 0, 0}),
		ap214.Axis2Placement3D(8, "frame", 5, 6, 7),
		ap214.ToroidalSurface(1, "torus", 8, 4, 1),
		ap214.SphericalSurface(2, "sphere", 8, 2.5),
		ap214.Ellipse(3, "ellipse", 8, step.Real(2*math.Sqrt2), 2),
	)
	encoded, err := f.Marshal()
	require.NoError(t, err)
	for _, line := range []string{
		"#1=TOROIDAL_SURFACE('torus',#8,4.,1.);",
		"#2=SPHERICAL_SURFACE('sphere',#8,2.5);",
		"#3=ELLIPSE('ellipse',#8,2.8284271247461903,2.);",
	} {
		require.Truef(t, bytes.Contains(encoded, []byte(line+"\n")), "missing record %s", line)
	}
}
