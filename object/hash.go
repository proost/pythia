package object

import (
	"bytes"
	"fmt"
	"strings"
)

type HashPair struct {
	Key   Object
	Value Object
}
type Hash struct {
	Pairs map[HashKey]HashPair
	check map[HashKey]struct{} // This is for for-loop
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
func (h *Hash) Equals(o Object) bool {
	obj, ok := o.(*Hash)
	if !ok {
		return false
	}

	if len(h.Pairs) != len(obj.Pairs) {
		return false
	}

	for pairKey, pairValue := range h.Pairs {
		if !obj.Pairs[pairKey].Value.Equals(pairValue.Value) {
			return false
		}
	}
	for pairKey, pairValue := range obj.Pairs {
		if !h.Pairs[pairKey].Value.Equals(pairValue.Value) {
			return false
		}
	}

	return true
}

func (h *Hash) HasNext() bool {
	if len(h.check) >= len(h.Pairs) {
		return false
	} else {
		return true
	}
}
func (h *Hash) Next() (Object, Object, bool) {
	if h.HasNext() {
		for k, pair := range h.Pairs {
			if _, isContained := h.check[k]; !isContained {
				h.check[k] = struct{}{}

				return pair.Key, pair.Value, true
			}
		}
	}

	return &Null{}, &Null{}, false
}
func (h *Hash) Reset() {
	h.check = map[HashKey]struct{}{}
}

func (h *Hash) Apply(method string, env *Environment, args ...Object) (Object, bool) {
	switch method {
	case "isEmpty":
		return &Boolean{Value: h.IsEmpty()}, true
	case "keys":
		return h.Keys()
	case "values":
		return h.Values()
	}

	return nil, false
}
func (h *Hash) IsEmpty() bool {
	if len(h.Pairs) == 0 {
		return true
	}

	return false
}
func (h *Hash) Keys() (Object, bool) {
	elements := make([]Object, len(h.Pairs))

	i := 0
	for _, pair := range h.Pairs {
		elements[i] = pair.Key
		i++
	}

	return &Array{Elements: elements}, true
}
func (h *Hash) Values() (Object, bool) {
	elements := make([]Object, len(h.Pairs))

	i := 0
	for _, pair := range h.Pairs {
		elements[i] = pair.Value
		i++
	}

	return &Array{Elements: elements}, true
}
