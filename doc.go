// Package step defines and writes single-section ISO 10303-21 STEP files.
//
// It models the Part 21 exchange syntax, not an application protocol or CAD
// geometry. Callers supply the schema and entities. The package has no
// dependencies outside the Go standard library and does not import a CAD
// engine.
package step
