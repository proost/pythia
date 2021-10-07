package object

import (
	"bytes"
	"strings"
)

type Array struct {
	Elements []Object
	offset   int // This is for for-loop
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

func (arr *Array) HasNext() bool {
	if arr.offset >= len(arr.Elements) {
		return false
	} else {
		return true
	}
}
func (arr *Array) Next() (Object, Object, bool) {
	if arr.HasNext() {
		idx := &Integer{Value: int64(arr.offset)}
		val := arr.Elements[arr.offset]

		arr.offset++

		return val, idx, true
	}

	return &Null{}, &Null{}, false
}
func (arr *Array) Reset() {
	arr.offset = 0
}

func (arr *Array) Apply(method string, env *Environment, args ...Object) (Object, bool) {
	switch method {
	case "isEmpty":
		return &Boolean{Value: arr.IsEmpty()}, true
	case "last":
		return arr.Last()
	}

	return nil, false
}
func (arr *Array) IsEmpty() bool {
	if len(arr.Elements) == 0 {
		return true
	}

	return false
}
func (arr *Array) Last() (Object, bool) {
	if arr.IsEmpty() {
		return nil, false
	}

	return arr.Elements[len(arr.Elements)-1], true
}
