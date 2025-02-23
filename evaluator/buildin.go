package evaluator

import (
	"github.com/sachinaralapura/shoebill/object"
)

func lenBuildIn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newErrorObject("wrong number of arguments. got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newErrorObject("argument to `len` not supported, got %s", args[0].Type())
	}
}

func firstBuildIn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newErrorObject("wrong number of arguments. got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *object.Array:
		if len(arg.Elements) > 0 {
			return arg.Elements[0]
		}
		return NULL
	default:
		return newErrorObject("argument to `first` must be ARRAY, got %s", args[0].Type())
	}
}

func lastBuildIn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newErrorObject("wrong number of arguments. got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *object.Array:
		length := len(arg.Elements)
		if length > 0 {
			return arg.Elements[length-1]
		}
		return NULL
	default:
		return newErrorObject("argument to `last` must be ARRAY or STRING, got %s", args[0].Type())
	}
}

func restBuildIn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newErrorObject("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	switch arg := args[0].(type) {
	case *object.Array:
		length := len(arg.Elements)
		if length > 0 {
			newElements := make([]object.Object, length-1)
			copy(newElements, arg.Elements[1:length])
			return &object.Array{Elements: newElements}
		}
		return NULL
	default:
		return newErrorObject("argument to `rest` must be ARRAY of STRING, got %s", args[0].Type())
	}
}

func pushBuildIn(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newErrorObject("wrong number of arguments. got=%d, want=2", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newErrorObject("argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arg := args[0].(*object.Array)
	length := len(arg.Elements)
	newElements := make([]object.Object, length+1)
	copy(newElements, arg.Elements)
	newElements[length] = args[1]

	return &object.Array{Elements: newElements}
}

func printBuildIn(args ...object.Object) object.Object {
	return NULL
}

var BuildIns = map[string]*object.BuildIn{
	"len":   {Value: lenBuildIn},
	"print": {Value: printBuildIn},
	"first": {Value: firstBuildIn},
	"last":  {Value: lastBuildIn},
	"rest":  {Value: restBuildIn},
	"push":  {Value: pushBuildIn},
}
