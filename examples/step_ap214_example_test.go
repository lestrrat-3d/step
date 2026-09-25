package examples_test

import (
	"fmt"
	"time"

	"github.com/lestrrat-3d/step"
	"github.com/lestrrat-3d/step/ap214"
)

func Example_step_ap214() {
	header := step.Header{
		Description:   []string{"AP214 geometry"},
		Name:          "geometry.step",
		Timestamp:     time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC),
		Authors:       []string{"Example"},
		Organizations: []string{"Example"},
	}
	f := ap214.NewFile(header,
		ap214.CartesianPoint(1, "origin", [3]step.Real{0, 0, 0}),
		ap214.Direction(2, "z", [3]step.Real{0, 0, 1}),
		ap214.Direction(3, "x", [3]step.Real{1, 0, 0}),
		ap214.Axis2Placement3D(4, "frame", 1, 2, 3),
		ap214.Plane(5, "plane", 4),
	)
	data, err := f.Marshal()
	if err != nil {
		fmt.Printf("failed to write AP214: %s\n", err)
		return
	}
	fmt.Print(string(data))
	// Output:
	// ISO-10303-21;
	// HEADER;
	// FILE_DESCRIPTION(('AP214 geometry'),'2;1');
	// FILE_NAME('geometry.step','2026-09-25T00:00:00Z',('Example'),('Example'),'','','');
	// FILE_SCHEMA(('AUTOMOTIVE_DESIGN'));
	// ENDSEC;
	// DATA;
	// #1=CARTESIAN_POINT('origin',(0.,0.,0.));
	// #2=DIRECTION('z',(0.,0.,1.));
	// #3=DIRECTION('x',(1.,0.,0.));
	// #4=AXIS2_PLACEMENT_3D('frame',#1,#2,#3);
	// #5=PLANE('plane',#4);
	// ENDSEC;
	// END-ISO-10303-21;
}
