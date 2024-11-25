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
			default:
				return newError("`len` builtin function doesn't support argument of type %s", arg.Type())
			}
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
