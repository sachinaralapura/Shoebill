package evaluator

import "github.com/sachinaralapura/shoebill/object"

func lenBuildIn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newErrorObject("wrong number of arguments. got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newErrorObject("argument to `len` not supported, got %s", args[0].Type())
	}
}

func printBuildIn(args ...object.Object) object.Object {
	return NULL
}

var BuildIns = map[string]*object.BuildIn{
	"len":   {Value: lenBuildIn},
	"print": {Value: printBuildIn},
}
