package evaluator

import (
	"gomonkey/object"
	"os"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Name: "len",
		Fn: func(args ...object.Object) object.Object {
			if lenArgs := len(args); lenArgs != 1 {
				return newError("wrong number of arguments. got=%d, want=1", lenArgs)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("`len` builtin function doesn't support argument of type %s", arg.Type())
			}
		},
	},
	"first": {
		Name: "first",
		Fn: func(args ...object.Object) object.Object {
			if lenArgs := len(args); lenArgs != 1 {
				return newError("wrong number of arguments. got=%d, want=1", lenArgs)
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": {
		Name: "last",
		Fn: func(args ...object.Object) object.Object {
			if lenArgs := len(args); lenArgs != 1 {
				return newError("wrong number of arguments. got=%d, want=1", lenArgs)
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if length := len(arr.Elements); length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},
	"rest": {
		Name: "rest",
		Fn: func(args ...object.Object) object.Object {
			if lenArgs := len(args); lenArgs != 1 {
				return newError("wrong number of arguments. got=%d, want=1", lenArgs)
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},
	"exit": {
		Name: "exit",
		Fn: func(args ...object.Object) object.Object {
			lenArgs := len(args)

			if lenArgs == 0 {
				os.Exit(0)
			}

			if lenArgs > 1 {
				return newError("wrong number of arguments. got=%d, want 0 or 1", lenArgs)
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				os.Exit(int(arg.Value))
			default:
				return newError("`exit` builtin function doesn't support argument of type %s", arg.Type())
			}

			return nil
		},
	},
}
