package object

import (
	"fmt"
	"math"
)

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
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}
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
func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), Value: math.Float64bits(f.Value)}
}
func (f *Float) Equals(o Object) bool {
	obj, ok := o.(*Float)
	if !ok {
		return false
	}

	return f.Value == obj.Value
}
func (f *Float) Number()            {}
func (f *Float) ToFloat64() float64 { return f.Value }
