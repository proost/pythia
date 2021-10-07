package evaluator

import (
	"fmt"
	"pythia/ast"
	"pythia/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.FunctionStatement:
		return evalFunctionStatement(node, env)
	case *ast.IfStatement:
		return evalIfStatement(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.InstructionStatement:
		return evalInstructionStatement(node)
	case *ast.LetStatement:
		return evalLetStatement(node, env)
	case *ast.ForStatement:
		return evalForStatement(node, env)
	case *ast.AssignmentExpression:
		return evalAssignmentExpression(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)

		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.MethodCallExpression:
		return evalMethodCallExpression(node, env)
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
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.NullLiteral:
		return NULL
	}

	return nil
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
