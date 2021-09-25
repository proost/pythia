package object

import (
	"bytes"
	"fmt"
	"strings"
)

type ObjectType string

const (
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
	TYPE_OBJ         = "TYPE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	Equals(o Object) bool
}

type Number interface {
	// TODO Add Complex inteface
	Number()
}

type Real interface {
	Number
	ToFloat64() float64
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Equals(o Object) bool {
	obj, ok := o.(*Integer)
	if !ok {
		return false
	}

	return i.Value == obj.Value
}
func (i *Integer) Number()            {}
func (i *Integer) ToFloat64() float64 { return float64(i.Value) }

type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Equals(o Object) bool {
	obj, ok := o.(*Float)
	if !ok {
		return false
	}

	return f.Value == obj.Value
}
func (f *Float) Number()            {}
func (f *Float) ToFloat64() float64 { return f.Value }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Equals(o Object) bool {
	obj, ok := o.(*Boolean)
	if !ok {
		return false
	}

	return b.Value == obj.Value
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }
func (n *Null) Equals(o Object) bool {
	_, ok := o.(*Null)
	if !ok {
		return false
	} else {
		return true
	}
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Equals(o Object) bool {
	obj, ok := o.(*ReturnValue)
	if !ok {
		return false
	}

	return rv.Value.Equals(obj.Value)
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Equals(o Object) bool {
	obj, ok := o.(*Error)
	if !ok {
		return false
	}

	return e.Message == obj.Message
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) Equals(o Object) bool {
	obj, ok := o.(*String)
	if !ok {
		return false
	}

	return s.Value == obj.Value
}

type Array struct {
	Elements []Object
}

func (arr *Array) Type() ObjectType { return ARRAY_OBJ }
func (arr *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range arr.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (arr *Array) Equals(o Object) bool {
	obj, ok := o.(*Array)
	if !ok {
		return false
	}

	if len(arr.Elements) != len(obj.Elements) {
		return false
	}

	i := 0
	for i < len(arr.Elements) {
		if !arr.Elements[i].Equals(obj.Elements[i]) {
			return false
		}
		i++
	}

	return true
}

type Type struct {
	InstanceType ObjectType
}

func (t *Type) Type() ObjectType { return TYPE_OBJ }
func (t *Type) Inspect() string  { return "Type: " + string(t.InstanceType) }
func (t *Type) Equals(o Object) bool {
	obj, ok := o.(*Type)
	if !ok {
		return false
	}

	return t.InstanceType == obj.InstanceType
}
