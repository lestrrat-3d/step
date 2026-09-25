package step

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Value is a Part 21 entity parameter. Its concrete forms are defined here.
type Value interface{ stepValue() }

// String is a Part 21 string value.
type String string

// Integer is a Part 21 integer value.
type Integer int64

// Real is a Part 21 real value.
type Real float64

// Reference names another entity in the same file.
type Reference uint64

// Enumeration is an unqualified enumeration item, for example T or F.
type Enumeration string

// Binary is the contents of a Part 21 binary literal, without quotes. Its
// first digit (0–3) states the number of unused bits; the rest are hex digits.
type Binary string

// List is a possibly empty, possibly nested Part 21 aggregate.
type List []Value

// Typed is a parameter whose EXPRESS type is stated explicitly.
type Typed struct {
	Type  string
	Value Value
}

// Null represents an absent optional value ($).
type Null struct{}

// Omitted represents a value derived by the schema (*).
type Omitted struct{}

func (String) stepValue()      {}
func (Integer) stepValue()     {}
func (Real) stepValue()        {}
func (Reference) stepValue()   {}
func (Enumeration) stepValue() {}
func (Binary) stepValue()      {}
func (List) stepValue()        {}
func (Typed) stepValue()       {}
func (Null) stepValue()        {}
func (Omitted) stepValue()     {}

const maxValueDepth = 256

func validateValue(v Value, ids map[uint64]struct{}, depth int) error {
	if depth > maxValueDepth {
		return fmt.Errorf("aggregate nesting exceeds %d levels", maxValueDepth)
	}
	switch x := v.(type) {
	case String:
		_, err := quoteString(string(x))
		return err
	case Integer:
		return nil
	case Real:
		if math.IsNaN(float64(x)) || math.IsInf(float64(x), 0) {
			return fmt.Errorf("real must be finite")
		}
		return nil
	case Reference:
		if _, ok := ids[uint64(x)]; !ok {
			return fmt.Errorf("reference #%d has no entity", x)
		}
		return nil
	case Enumeration:
		if !validEnumeration(string(x)) {
			return fmt.Errorf("invalid enumeration %q", x)
		}
		return nil
	case Binary:
		if !validBinary(string(x)) {
			return fmt.Errorf("invalid binary literal %q", x)
		}
		return nil
	case List:
		for i, item := range x {
			if err := validateValue(item, ids, depth+1); err != nil {
				return fmt.Errorf("list item %d: %w", i, err)
			}
		}
		return nil
	case Typed:
		if !validKeyword(x.Type, true) {
			return fmt.Errorf("invalid type %q", x.Type)
		}
		if err := validateValue(x.Value, ids, depth+1); err != nil {
			return fmt.Errorf("typed %s: %w", x.Type, err)
		}
		return nil
	case Null, Omitted:
		return nil
	default:
		return fmt.Errorf("nil or unknown STEP value %T", v)
	}
}

func validKeyword(s string, allowUserDefined bool) bool {
	if allowUserDefined && strings.HasPrefix(s, "!") {
		s = s[1:]
	}
	if len(s) == 0 || s[0] < 'A' || s[0] > 'Z' {
		return false
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if c < 'A' || c > 'Z' {
			if c < '0' || c > '9' {
				if c != '_' {
					return false
				}
			}
		}
	}
	return true
}

func validEnumeration(s string) bool {
	if len(s) == 0 || s[0] < 'A' || s[0] > 'Z' {
		return false
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

func validBinary(s string) bool {
	if len(s) == 0 || s[0] < '0' || s[0] > '3' {
		return false
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if (c < '0' || c > '9') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}

// quoteString escapes an effective string into the edition-1-compatible ASCII
// spelling. Part 21 string length is limited on the encoded octets.
func quoteString(s string) (string, error) {
	if !utf8.ValidString(s) {
		return "", fmt.Errorf("string is not valid UTF-8")
	}
	var b strings.Builder
	b.WriteByte('\'')
	for _, r := range s {
		switch {
		case r == '\'':
			b.WriteString("''")
		case r == '\\':
			b.WriteString("\\\\")
		case r < 0x20 || r == 0x7f:
			fmt.Fprintf(&b, "\\X\\%02X", r)
		case r <= 0x7e:
			b.WriteRune(r)
		case r <= 0xffff:
			fmt.Fprintf(&b, "\\X2\\%04X\\X0\\", r)
		default:
			fmt.Fprintf(&b, "\\X4\\%08X\\X0\\", r)
		}
		if b.Len()+1 > 32769 {
			return "", fmt.Errorf("encoded string exceeds 32769 octets")
		}
	}
	b.WriteByte('\'')
	return b.String(), nil
}

func realToken(v Real) string {
	s := strings.ToUpper(strconv.FormatFloat(float64(v), 'g', -1, 64))
	if e := strings.IndexByte(s, 'E'); e >= 0 {
		if !strings.Contains(s[:e], ".") {
			s = s[:e] + "." + s[e:]
		}
		return s
	}
	if !strings.Contains(s, ".") {
		s += "."
	}
	return s
}
