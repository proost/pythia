package evaluator

import (
	"fmt"
	"pythia/object"
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
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to first must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to last must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
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
				fmt.Println(arg.Inspect())
			}

			return nil
		},
	},
	"type": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
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
}
