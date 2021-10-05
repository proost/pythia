package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
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
func (h *Hash) HashKey() HashKey {
	b := make([]byte, len(h.Pairs)*2)
	i := 0
	for _, pair := range h.Pairs {
		b[i] = byte(pair.Key.HashKey().Value)
		b[i+1] = byte(pair.Value.HashKey().Value)
		i += 2
	}

	hs := fnv.New64()
	hs.Write(b)

	return HashKey{Type: h.Type(), Value: hs.Sum64()}
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
