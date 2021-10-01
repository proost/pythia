package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"math"
	"strings"
)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), Value: math.Float64bits(f.Value)}
}

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

type Hashable interface {
	HashKey() HashKey
}
