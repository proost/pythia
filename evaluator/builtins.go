package evaluator

import (
	"fmt"
	"pythia/object"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Hash:
				return &object.Integer{Value: int64(len(arg.Pairs))}
			default:
				return newError("argument to len not supported, got %s", args[0].Type())
			}
		},
	},
	"append": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to append must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
	"print": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				msg := strings.ReplaceAll(arg.Inspect(), `\n`, "\n") // escaped version of line breaking to real line breaking
				fmt.Printf(msg)
			}

			fmt.Println()

			return nil
		},
	},
	"type": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0] == nil {
				return &object.Type{InstanceType: NULL.Type()}
			}

			return &object.Type{InstanceType: args[0].Type()}
		},
	},
	"range": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if !(len(args) == 2 || len(args) == 3) {
				return newError("wrong number of arguments. got=%d, want= 2 or 3", len(args))
			}
			for _, arg := range args {
				if arg.Type() != object.INTEGER_OBJ {
					return newError("argument %v must be Integer, got %s", arg, arg.Type())
				}
			}

			left := args[0].(*object.Integer).Value
			right := args[1].(*object.Integer).Value
			stepVal := int64(1)
			if len(args) == 3 {
				stepVal = args[2].(*object.Integer).Value
			}

			if left <= right {
				if stepVal < 0 {
					return newError("start can't be smaller than end, when step is %d", stepVal)
				}

				arr := make([]object.Object, 0)
				for left < right {
					arr = append(arr, &object.Integer{Value: left})
					left += stepVal
				}

				return &object.Array{Elements: arr}
			} else {
				if stepVal > 0 {
					return newError("start can't be bigger than end, when step is %d", stepVal)
				}

				arr := make([]object.Object, 0)
				for left > right {
					arr = append(arr, &object.Integer{Value: left})
					left += stepVal
				}

				return &object.Array{Elements: arr}
			}
		},
	},
	"delete": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.HASH_OBJ {
				return newError("first argument of delete must be HASH, got %s", args[0].Type())
			}

			hash := args[0].(*object.Hash)
			index, ok := args[1].(object.Hashable)
			if !ok {
				return newError("unusable as hash key: %s", args[1].Type())
			}

			delete(hash.Pairs, index.HashKey())

			return nil
		},
	},
}
