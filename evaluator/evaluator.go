package evaluator

import (
	"fmt"

	"github.com/sachinaralapura/shoebill/ast"
	"github.com/sachinaralapura/shoebill/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.Return{Value: val}

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.BoolenExpression:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.IfExpression:
		return evalIfExpressionObject(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.FunctionObject{Parameters: params, Body: body, Env: env}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}
	return nil
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var results []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		results = append(results, evaluated)
	}
	return results
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.FunctionObject:
		extendedEnv, ok := extendFunctionEnv(fn, args)
		if !ok {
			return newErrorObject("excepted no. of arguments not passed to function")
		}
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.BuildIn:
		return fn.Value(args...)

	default:
		return newErrorObject("not a function: %s", fn.Type())
	}

}

func extendFunctionEnv(fn *object.FunctionObject, args []object.Object) (*object.Environment, bool) {
	env := object.NewEnclosedEnvironment(fn.Env)
	if len(fn.Parameters) != len(args) {
		return nil, false
	}
	for id, params := range fn.Parameters {
		env.Set(params.Value, args[id])
	}
	return env, true
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.Return); ok {
		return returnValue.Value
	}
	return obj
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, left, right)
	}
	if leftType == object.BOOLEAN_OBJ && rightType == object.BOOLEAN_OBJ {
		return evalBooleanInfixExpression(operator, left, right)
	}
	if leftType == object.STRING_OBJ && rightType == object.STRING_OBJ {
		return evalStringInfixExpression(operator, left, right)
	}
	if leftType != rightType {
		return newErrorObject("type mismatch: %s %s %s", leftType, operator, rightType)
	}
	return newErrorObject("unknown operator: %s %s %s", leftType, operator, rightType)
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorPrefixExpression(right)
	case "-":
		return evalMinusPrefixExpression(right)
	default:
		return newErrorObject("unknown operator: %s%s", operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "%":
		return &object.Integer{Value: leftValue % rightValue}
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	default:
		return newErrorObject("unknown operator : %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newErrorObject("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value
	return &object.String{Value: leftValue + rightValue}
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left == right) // pointer comparsion
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	case "&&":
		leftObj := left.(*object.Boolean)
		rightObj := right.(*object.Boolean)
		return nativeBoolToBooleanObject(leftObj.Value && rightObj.Value)
	case "||":
		leftObj := left.(*object.Boolean)
		rightObj := right.(*object.Boolean)
		return nativeBoolToBooleanObject(leftObj.Value || rightObj.Value)
	default:
		return newErrorObject("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalMinusPrefixExpression(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}
	case *object.ErrorObject:
		return right
	default:
		return newErrorObject("unknown operator: -%s", right.Type())
	}
}

func evalBangOperatorPrefixExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalIfExpressionObject(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var obj object.Object
	for _, stmt := range statements {
		obj = Eval(stmt, env)

		switch obj := obj.(type) {
		case *object.Return:
			return obj.Value
		case *object.ErrorObject:
			return obj
		}
	}
	return obj
}

func evalBlockStatements(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range statements {
		result = Eval(stmt, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if buildin, ok := BuildIns[node.Value]; ok {
		return buildin
	}
	return newErrorObject("identifier not found: %s", node.Value)
}

func nativeBoolToBooleanObject(value bool) object.Object {
	if value {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Null:
		return false
	case *object.Boolean:
		if obj.Value == true {
			return true
		}
		return false
	case *object.Integer:
		if obj.Value == 0 {
			return false
		}
		return true
	default:
		return false
	}
}

func newErrorObject(format string, a ...any) *object.ErrorObject {
	return &object.ErrorObject{Message: fmt.Sprintf(format, a...)}
}

func isError(node object.Object) bool {
	if node != nil {
		return node.Type() == object.ERROR_OBJ
	}
	return false
}
