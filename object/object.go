package object

import (
	"fmt"
	"hash/fnv"
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
	HASH_KEYS_OBJ    = "HASH_KEYS"
	HASH_VALUES_OBJ  = "HASH_VALUES"
	TYPE_OBJ         = "TYPE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	Equals(o Object) bool
}

type Iterator interface {
	HasNext() bool
	Next() (Object, Object, bool) // required value, optional value, hasNext
	Reset()                       // Before new iteration start, reset offset information
}

type Callable interface {
	Apply(method string, env *Environment, args ...Object) (Object, bool)
}

type Collections interface {
	Callable
	IsEmpty() bool
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}
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
func (n *Null) HashKey() HashKey {
	return HashKey{Type: n.Type(), Value: 0}
}
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
	Value  string
	offset int // This is for for-loop
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (s *String) Equals(o Object) bool {
	obj, ok := o.(*String)
	if !ok {
		return false
	}

	return s.Value == obj.Value
}
func (s *String) HasNext() bool {
	if s.offset >= len(s.Value) {
		return false
	} else {
		return true
	}
}
func (s *String) Next() (Object, Object, bool) {
	if s.HasNext() {
		idx := &Integer{Value: int64(s.offset)}
		val := &String{Value: string(s.Value[s.offset])}

		s.offset++

		return val, idx, true
	}

	return &Null{}, &Null{}, false
}
func (s *String) Reset() {
	s.offset = 0
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
