# step

`github.com/lestrrat-3d/step` defines and writes ISO 10303-21 clear-text STEP
files. The root package handles Part 21 headers, entity records, values, and
serialization. The `ap214` subpackage supplies constructors for a focused set
of `AUTOMOTIVE_DESIGN` records.

The packages do not parse files, validate a complete AP214 model, or convert a
CAD body into a STEP solid. See the executable examples in `examples/` and the
scope in `docs/step-design.md`.
