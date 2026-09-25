package examples_test

import (
	"fmt"
	"time"

	"github.com/lestrrat-3d/step"
)

func Example_step_file() {
	f := step.File{
		Header: step.Header{
			Description:   []string{"one point"},
			Name:          "point.step",
			Timestamp:     time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC),
			Authors:       []string{"Example"},
			Organizations: []string{"Example"},
			Schemas:       []string{"EXAMPLE_SCHEMA"},
		},
		Entities: []step.Entity{
			{ID: 1, Name: "CARTESIAN_POINT", Parameters: []step.Value{
				step.String("origin"), step.List{step.Real(0), step.Real(0), step.Real(0)},
			}},
		},
	}
	data, err := f.Marshal()
	if err != nil {
		fmt.Printf("failed to write STEP: %s\n", err)
		return
	}
	fmt.Print(string(data))
	// Output:
	// ISO-10303-21;
	// HEADER;
	// FILE_DESCRIPTION(('one point'),'2;1');
	// FILE_NAME('point.step','2026-09-25T00:00:00Z',('Example'),('Example'),'','','');
	// FILE_SCHEMA(('EXAMPLE_SCHEMA'));
	// ENDSEC;
	// DATA;
	// #1=CARTESIAN_POINT('origin',(0.,0.,0.));
	// ENDSEC;
	// END-ISO-10303-21;
}
