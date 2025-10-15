package internal

import (
	"fmt"
	"strconv"

	"marktstammdatenregister.dev/internal/spec"
)

var (
	unknownXsdType = "don't know how to handle XSD type '%s'"
	xsd2sqliteType = map[string]string{
		"nonNegativeInteger": "integer",
		"boolean":            "integer",
		"decimal":            "real",
		"date":               "text",
		"dateTime":           "text",
		"":                   "text",
	}
)

func Xsd2SqliteType(xsd string) (string, bool) {
	typ, ok := xsd2sqliteType[xsd]
	return typ, ok
}

type Fields struct {
	order []string
	pos   map[string]int
	typ   map[string]string
}

func NewFields(fields []spec.Field) (*Fields, error) {
	order := make([]string, len(fields))
	pos := make(map[string]int)
	typ := make(map[string]string)
	for i, field := range fields {
		t, ok := Xsd2SqliteType(field.Xsd)
		if !ok {
			return nil, fmt.Errorf(unknownXsdType, field.Xsd)
		}
		order[i] = field.Name
		pos[field.Name] = i
		typ[field.Name] = t
	}
	return &Fields{order: order, pos: pos, typ: typ}, nil
}

func (f *Fields) Header() []string {
	result := make([]string, len(f.order))
	copy(result, f.order)
	return result
}

func (f *Fields) EnsureField(name, sqlType string) (bool, string) {
	if _, ok := f.pos[name]; ok {
		if sqlType != "" && f.typ[name] == "" {
			f.typ[name] = sqlType
		}
		return false, f.typ[name]
	}
	typ := sqlType
	if typ == "" {
		typ = "text"
	}
	f.pos[name] = len(f.order)
	f.order = append(f.order, name)
	f.typ[name] = typ
	return true, typ
}

func (f *Fields) Record(item map[string]string) ([]interface{}, error) {
	result := make([]interface{}, len(f.order))
	for idx, name := range f.order {
		value, ok := item[name]
		if !ok || value == "" {
			result[idx] = nil
			continue
		}
		switch f.typ[name] {
		case "integer":
			v, err := strconv.Atoi(value)
			if err != nil {
				return result, err
			}
			result[idx] = v
		case "real":
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return result, err
			}
			result[idx] = v
		case "text", "":
			result[idx] = value
		default:
			return nil, fmt.Errorf(unknownXsdType, f.typ[name])
		}
	}
	return result, nil
}
