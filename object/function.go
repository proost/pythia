package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"pythia/ast"
	"strings"
)

type Function struct {
	Parameters []*ast.Identifier
	Name       *ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("func ")
	out.WriteString(f.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (f *Function) HashKey() HashKey {
	h := fnv.New64()
	h.Write([]byte(fmt.Sprintf("%p", f)))

	return HashKey{Type: f.Type(), Value: h.Sum64()}
}
func (f *Function) Equals(o Object) bool {
	obj, ok := o.(*Function)
	if !ok {
		return false
	}

	return f == obj
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) HashKey() HashKey {
	h := fnv.New64()
	h.Write([]byte(fmt.Sprintf("%p", b)))

	return HashKey{Type: b.Type(), Value: h.Sum64()}
}
func (b *Builtin) Equals(o Object) bool {
	obj, ok := o.(*Builtin)
	if !ok {
		return false
	}

	return &b.Fn == &obj.Fn
}
