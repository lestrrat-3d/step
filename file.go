package step

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Header holds the three mandatory Part 21 header records. The caller must
// supply a timestamp; no field is silently filled from the current clock.
type Header struct {
	Description         []string
	Name                string
	Timestamp           time.Time
	Authors             []string
	Organizations       []string
	PreprocessorVersion string
	OriginatingSystem   string
	Authorization       string
	Schemas             []string
}

// Record is one named component of a complex entity instance.
type Record struct {
	Name       string
	Parameters []Value
}

// Entity is an instance in the file's unnamed DATA section. A simple instance
// uses Name and Parameters. A complex instance uses Components instead. ID must
// be positive and unique; forward references are allowed.
type Entity struct {
	ID         uint64
	Name       string
	Parameters []Value
	Components []Record
}

// File is a single-data-section ISO 10303-21 exchange file. Entity order is
// retained. The caller must not mutate it while Write or Marshal is running.
type File struct {
	Header   Header
	Entities []Entity
}

// Write validates and writes the complete file. Validation errors leave w
// untouched. I/O errors may leave a partial file and are returned as errors.
func (f File) Write(w io.Writer) error {
	if w == nil {
		return fmt.Errorf("step: nil writer")
	}
	if err := f.validate(); err != nil {
		return fmt.Errorf("step: %w", err)
	}
	bw := bufio.NewWriter(w)
	e := emitter{w: bw}
	e.line("ISO-10303-21;\nHEADER;\n")
	e.line("FILE_DESCRIPTION(")
	e.stringList(f.Header.Description)
	e.line(",'2;1');\nFILE_NAME(")
	e.string(f.Header.Name)
	e.line(",")
	e.string(f.Header.Timestamp.UTC().Format(time.RFC3339Nano))
	e.line(",")
	e.stringList(f.Header.Authors)
	e.line(",")
	e.stringList(f.Header.Organizations)
	for _, s := range []string{f.Header.PreprocessorVersion, f.Header.OriginatingSystem, f.Header.Authorization} {
		e.line(",")
		e.string(s)
	}
	e.line(");\nFILE_SCHEMA(")
	e.stringList(f.Header.Schemas)
	e.line(");\nENDSEC;\nDATA;\n")
	for _, entity := range f.Entities {
		e.line("#")
		e.line(strconv.FormatUint(entity.ID, 10))
		e.line("=")
		if len(entity.Components) == 0 {
			e.record(Record{Name: entity.Name, Parameters: entity.Parameters})
		} else {
			e.line("(")
			for _, component := range entity.Components {
				e.record(component)
			}
			e.line(")")
		}
		e.line(";\n")
	}
	e.line("ENDSEC;\nEND-ISO-10303-21;\n")
	if e.err != nil {
		return fmt.Errorf("step: write failed: %w", e.err)
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("step: write failed: %w", err)
	}
	return nil
}

// Marshal returns the complete Part 21 file as bytes.
func (f File) Marshal() ([]byte, error) {
	var b bytes.Buffer
	if err := f.Write(&b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (f File) validate() error {
	h := f.Header
	if len(h.Description) == 0 || len(h.Authors) == 0 || len(h.Organizations) == 0 {
		return fmt.Errorf("description, authors, and organizations each require at least one entry")
	}
	if len(h.Schemas) != 1 {
		return fmt.Errorf("the unnamed data section requires exactly one schema")
	}
	if h.Name == "" {
		return fmt.Errorf("file name is empty")
	}
	if h.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is zero")
	}
	if year := h.Timestamp.UTC().Year(); year < 1 || year > 9999 {
		return fmt.Errorf("timestamp year is outside 0001–9999")
	}
	for _, field := range []struct {
		name  string
		value string
	}{
		{"name", h.Name},
		{"timestamp", h.Timestamp.UTC().Format(time.RFC3339Nano)},
		{"preprocessor version", h.PreprocessorVersion},
		{"originating system", h.OriginatingSystem},
		{"authorization", h.Authorization},
	} {
		if err := validateHeaderString(field.value, 256); err != nil {
			return fmt.Errorf("%s: %w", field.name, err)
		}
	}
	for _, group := range []struct {
		name   string
		values []string
	}{
		{"description", h.Description},
		{"author", h.Authors},
		{"organization", h.Organizations},
	} {
		for i, value := range group.values {
			if err := validateHeaderString(value, 256); err != nil {
				return fmt.Errorf("%s %d: %w", group.name, i, err)
			}
		}
	}
	if !validSchemaIdentifier(h.Schemas[0]) {
		return fmt.Errorf("invalid schema identifier %q", h.Schemas[0])
	}
	if err := validateHeaderString(h.Schemas[0], 1024); err != nil {
		return fmt.Errorf("schema: %w", err)
	}
	ids := make(map[uint64]struct{}, len(f.Entities))
	for i, entity := range f.Entities {
		if entity.ID == 0 {
			return fmt.Errorf("entity %d: ID must be positive", i)
		}
		if _, exists := ids[entity.ID]; exists {
			return fmt.Errorf("entity %d: duplicate ID #%d", i, entity.ID)
		}
		ids[entity.ID] = struct{}{}
		if len(entity.Components) == 0 {
			if !validKeyword(entity.Name, true) {
				return fmt.Errorf("entity #%d: invalid name %q", entity.ID, entity.Name)
			}
			continue
		}
		if entity.Name != "" || len(entity.Parameters) != 0 {
			return fmt.Errorf("entity #%d: complex instance cannot also have a simple record", entity.ID)
		}
		for j, component := range entity.Components {
			if !validKeyword(component.Name, true) {
				return fmt.Errorf("entity #%d component %d: invalid name %q", entity.ID, j, component.Name)
			}
		}
	}
	for _, entity := range f.Entities {
		records := entity.Components
		if len(records) == 0 {
			records = []Record{{Name: entity.Name, Parameters: entity.Parameters}}
		}
		for _, record := range records {
			for i, value := range record.Parameters {
				if err := validateValue(value, ids, 0); err != nil {
					return fmt.Errorf("entity #%d %s parameter %d: %w", entity.ID, record.Name, i, err)
				}
			}
		}
	}
	return nil
}

func validateHeaderString(s string, maxRunes int) error {
	if utf8.RuneCountInString(s) > maxRunes {
		return fmt.Errorf("more than %d characters", maxRunes)
	}
	_, err := quoteString(s)
	return err
}

func validSchemaIdentifier(s string) bool {
	name, suffix, hasSuffix := strings.Cut(s, " ")
	if !validKeyword(name, false) {
		return false
	}
	if !hasSuffix {
		return true
	}
	suffix = strings.TrimSpace(suffix)
	if len(suffix) < 3 || suffix[0] != '{' || suffix[len(suffix)-1] != '}' {
		return false
	}
	arcs := strings.Fields(suffix[1 : len(suffix)-1])
	if len(arcs) < 2 {
		return false
	}
	for _, arc := range arcs {
		for i := range len(arc) {
			if arc[i] < '0' || arc[i] > '9' {
				return false
			}
		}
	}
	return true
}

type emitter struct {
	w   io.Writer
	err error
}

func (e *emitter) line(s string) {
	if e.err != nil {
		return
	}
	_, e.err = io.WriteString(e.w, s)
}

func (e *emitter) string(s string) {
	quoted, _ := quoteString(s) // File.validate checked every string.
	e.line(quoted)
}

func (e *emitter) stringList(values []string) {
	e.line("(")
	for i, value := range values {
		if i > 0 {
			e.line(",")
		}
		e.string(value)
	}
	e.line(")")
}

func (e *emitter) record(record Record) {
	e.line(record.Name)
	e.line("(")
	for i, value := range record.Parameters {
		if i > 0 {
			e.line(",")
		}
		e.value(value)
	}
	e.line(")")
}

func (e *emitter) value(v Value) {
	switch x := v.(type) {
	case String:
		e.string(string(x))
	case Integer:
		e.line(strconv.FormatInt(int64(x), 10))
	case Real:
		e.line(realToken(x))
	case Reference:
		e.line("#")
		e.line(strconv.FormatUint(uint64(x), 10))
	case Enumeration:
		e.line(".")
		e.line(string(x))
		e.line(".")
	case Binary:
		e.line("\"")
		e.line(string(x))
		e.line("\"")
	case List:
		e.line("(")
		for i, value := range x {
			if i > 0 {
				e.line(",")
			}
			e.value(value)
		}
		e.line(")")
	case Typed:
		e.line(x.Type)
		e.line("(")
		e.value(x.Value)
		e.line(")")
	case Null:
		e.line("$")
	case Omitted:
		e.line("*")
	}
}
