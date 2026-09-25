package step_test

import (
	"bytes"
	"errors"
	"io"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/lestrrat-3d/step"
	"github.com/stretchr/testify/require"
)

func sampleFile() step.File {
	return step.File{
		Header: step.Header{
			Description:         []string{"sample"},
			Name:                "sample.step",
			Timestamp:           time.Date(2026, 9, 25, 12, 30, 0, 0, time.UTC),
			Authors:             []string{"A. Author"},
			Organizations:       []string{"Example"},
			PreprocessorVersion: "step package",
			OriginatingSystem:   "test",
			Schemas:             []string{"EXAMPLE_SCHEMA"},
		},
		Entities: []step.Entity{
			{ID: 2, Name: "POINT", Parameters: []step.Value{step.String("origin"), step.List{step.Real(0), step.Real(1e20), step.Real(-0.5)}}},
			{ID: 1, Name: "PAIR", Parameters: []step.Value{step.Reference(2), step.Typed{Type: "LENGTH", Value: step.Real(2.5)}, step.Enumeration("T"), step.Binary("0AF"), step.Null{}, step.Omitted{}, step.Integer(-7)}},
		},
	}
}

func TestFileWrite(t *testing.T) {
	t.Parallel()
	f := sampleFile()
	want := "ISO-10303-21;\n" +
		"HEADER;\n" +
		"FILE_DESCRIPTION(('sample'),'2;1');\n" +
		"FILE_NAME('sample.step','2026-09-25T12:30:00Z',('A. Author'),('Example'),'step package','test','');\n" +
		"FILE_SCHEMA(('EXAMPLE_SCHEMA'));\n" +
		"ENDSEC;\n" +
		"DATA;\n" +
		"#2=POINT('origin',(0.,1.E+20,-0.5));\n" +
		"#1=PAIR(#2,LENGTH(2.5),.T.,\"0AF\",$,*,-7);\n" +
		"ENDSEC;\n" +
		"END-ISO-10303-21;\n"

	var b bytes.Buffer
	require.NoError(t, f.Write(&b))
	require.Equal(t, want, b.String())
	encoded, err := f.Marshal()
	require.NoError(t, err)
	require.Equal(t, want, string(encoded))
	var again bytes.Buffer
	require.NoError(t, f.Write(&again))
	require.Equal(t, b.String(), again.String())
}

func TestFileStringEscaping(t *testing.T) {
	t.Parallel()
	f := sampleFile()
	f.Entities[0].Parameters[0] = step.String("O'Brien\\日本\n😀")
	encoded, err := f.Marshal()
	require.NoError(t, err)
	require.Contains(t, string(encoded), "'O''Brien\\\\\\X2\\65E5\\X0\\\\X2\\672C\\X0\\\\X\\0A\\X4\\0001F600\\X0\\'")
	require.NotContains(t, string(encoded), "日本")
}

func TestFileRejectsInvalidModelBeforeWriting(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		change func(*step.File)
	}{
		{"no schema", func(f *step.File) { f.Header.Schemas = nil }},
		{"bad schema", func(f *step.File) { f.Header.Schemas = []string{"schema);DATA;"} }},
		{"bad schema OID", func(f *step.File) { f.Header.Schemas = []string{"AUTOMOTIVE_DESIGN { 1 0 10303 );"} }},
		{"multiple schemas", func(f *step.File) { f.Header.Schemas = []string{"EXAMPLE_SCHEMA", "OTHER_SCHEMA"} }},
		{"zero timestamp", func(f *step.File) { f.Header.Timestamp = time.Time{} }},
		{"invalid timestamp year", func(f *step.File) { f.Header.Timestamp = time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC) }},
		{"long header", func(f *step.File) { f.Header.Name = strings.Repeat("a", 257) }},
		{"zero ID", func(f *step.File) { f.Entities[0].ID = 0 }},
		{"duplicate ID", func(f *step.File) { f.Entities[1].ID = 2 }},
		{"bad entity name", func(f *step.File) { f.Entities[0].Name = "POINT);DATA;" }},
		{"mixed complex record", func(f *step.File) { f.Entities[0].Components = []step.Record{{Name: "POINT"}} }},
		{"bad complex component", func(f *step.File) {
			f.Entities[0] = step.Entity{ID: 2, Components: []step.Record{{Name: "BAD);"}}}
		}},
		{"dangling reference", func(f *step.File) { f.Entities[1].Parameters[0] = step.Reference(99) }},
		{"nil value", func(f *step.File) { f.Entities[1].Parameters[0] = nil }},
		{"nonfinite real", func(f *step.File) { f.Entities[0].Parameters[1] = step.Real(math.Inf(1)) }},
		{"bad enumeration", func(f *step.File) { f.Entities[1].Parameters[2] = step.Enumeration("T);#3=BAD(") }},
		{"enumeration with underscore", func(f *step.File) { f.Entities[1].Parameters[2] = step.Enumeration("NOT_SET") }},
		{"bad binary", func(f *step.File) { f.Entities[1].Parameters[3] = step.Binary("4FF") }},
		{"bad UTF-8", func(f *step.File) { f.Entities[0].Parameters[0] = step.String(string([]byte{0xff})) }},
		{"long string", func(f *step.File) { f.Entities[0].Parameters[0] = step.String(strings.Repeat("x", 32768)) }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := sampleFile()
			tc.change(&f)
			var b bytes.Buffer
			err := f.Write(&b)
			require.Error(t, err)
			require.Empty(t, b.String())
		})
	}
}

func TestFileComplexEntityAndSchemaOID(t *testing.T) {
	t.Parallel()
	f := sampleFile()
	f.Header.Schemas = []string{"AUTOMOTIVE_DESIGN { 1 0 10303 214 3 1 1 }"}
	f.Entities[0] = step.Entity{ID: 2, Components: []step.Record{
		{Name: "NAMED_UNIT", Parameters: []step.Value{step.Omitted{}}},
		{Name: "SI_UNIT", Parameters: []step.Value{step.Enumeration("MILLI"), step.Enumeration("METRE")}},
	}}
	encoded, err := f.Marshal()
	require.NoError(t, err)
	require.Contains(t, string(encoded), "FILE_SCHEMA(('AUTOMOTIVE_DESIGN { 1 0 10303 214 3 1 1 }'));\n")
	require.Contains(t, string(encoded), "#2=(NAMED_UNIT(*)SI_UNIT(.MILLI.,.METRE.));\n")
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestFileWriteError(t *testing.T) {
	t.Parallel()
	err := sampleFile().Write(failingWriter{})
	require.ErrorIs(t, err, io.ErrClosedPipe)
	require.True(t, errors.Is(err, io.ErrClosedPipe))
}
