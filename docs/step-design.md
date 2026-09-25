# STEP file package

`github.com/lestrrat-3d/step` defines and writes the common, single-data-section
ISO 10303-21 clear-text exchange file. The root package imports only the Go
standard library and does not depend on a CAD engine.

The package models the mandatory `FILE_DESCRIPTION`, `FILE_NAME`, and
`FILE_SCHEMA` header records and simple or complex entity instances in one unnamed
`DATA` section. Values cover strings, integers, reals, enumerations, binary literals,
references, aggregates, typed parameters, omitted parameters, and nulls. It writes
implementation level `2;1`, using ASCII string control directives for non-ASCII
characters. The unnamed section requires exactly one schema, with an optional
numeric object identifier. The caller supplies the timestamp; writing never
reads the clock.

`File.Write` validates the entire model before touching the writer. It rejects
invalid tokens, non-finite reals, duplicate entity IDs, and references to IDs
absent from the file. It checks Part 21 syntax and the header's basic EXPRESS
shape, but does not claim that entity names, attribute types, or relationships
conform to a particular application protocol schema. `step` does not parse
files, model geometry, or export CAD bodies.

`ap214` owns the `AUTOMOTIVE_DESIGN` schema identity and pure constructors
for a focused set of AP214 records: product linkage, 3D analytic geometry,
advanced B-rep topology, and the geometric representation context. It depends
only on `step`, and `step` does not import it. Its constructors fix record names,
field order, and Part 21 encodings; they do not certify EXPRESS WHERE rules,
topological validity, or full AP214 schema conformance. A caller supplies the
records and their references, and uses `step.File.Write` to serialize them.
