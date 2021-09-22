package object

import (
	"bytes"
	"fmt"
	"strings"
)

type ObjectType string

const (
	NUMBER_OBJ       = "NUMBER"
	READ_NUMBER_OBJ  = "REAL"
	INTEGER_OBJ      = "INTEGER"
	FLOAT_OBJ        = "FLOAT"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Number interface {
	// TODO Add Complex inteface
	Number()
}

type Real interface {
	ToFloat64() float64
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string    { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType   { return INTEGER_OBJ }
func (i *Integer) Number()            {}
func (i *Integer) ToFloat64() float64 { return float64(i.Value) }

type Float struct {
	Value float64
}

func (f *Float) Inspect() string    { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Type() ObjectType   { return FLOAT_OBJ }
func (f *Float) Number()            {}
func (f *Float) ToFloat64() float64 { return f.Value }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
