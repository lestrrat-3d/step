package ap214

import "github.com/lestrrat-3d/step"

// Schema is the AP214 EXPRESS schema name in FILE_SCHEMA.
const Schema = "AUTOMOTIVE_DESIGN"

// NewFile creates a Part 21 file using AP214's schema. The caller supplies
// header metadata and entity IDs; Header.Schemas is set to Schema. The entity
// slice is copied, so appending to the original slice cannot change the file.
func NewFile(header step.Header, entities ...step.Entity) step.File {
	header.Schemas = []string{Schema}
	return step.File{Header: header, Entities: append([]step.Entity(nil), entities...)}
}

func refs(ids []step.Reference) step.List {
	values := make(step.List, len(ids))
	for i, id := range ids {
		values[i] = id
	}
	return values
}

func logical(value bool) step.Enumeration {
	if value {
		return step.Enumeration("T")
	}
	return step.Enumeration("F")
}
