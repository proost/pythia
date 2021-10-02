package evaluator

import (
	"math"
	"pythia/ast"
	"pythia/object"
)

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalAssignmentExpression(ae *ast.AssignmentExpression, env *object.Environment) object.Object {
	evaluated := Eval(ae.Value, env)
	if isError(evaluated) {
		return evaluated
	}

	switch ae.Operator {
	case "=":
		_, ok := env.Get(ae.Name.String())
		if !ok {
			return newError("%s is not defined identifier", ae.Name.String())
		}

		env.Set(ae.Name.String(), evaluated)
	case "+=":
		curr, ok := env.Get(ae.Name.String())
		if !ok {
			return newError("%s is not defined identifier", ae.Name.String())
		}

		res := evalInfixExpression("+", curr, evaluated)
		if isError(res) {
			return newError("+= operation is not supported for %s, %s", curr.Type(), evaluated.Type())
		}

		env.Set(ae.Name.String(), res)
	case "-=":
		curr, ok := env.Get(ae.Name.String())
		if !ok {
			return newError("%s is not defined identifier", ae.Name.String())
		}

		res := evalInfixExpression("-", curr, evaluated)
		if isError(res) {
			return newError("-= operation is not supported for %s, %s", curr.Type(), evaluated.Type())
		}

		env.Set(ae.Name.String(), res)
	case "*=":
		curr, ok := env.Get(ae.Name.String())
		if !ok {
			return newError("%s is not defined identifier", ae.Name.String())
		}

		res := evalInfixExpression("*", curr, evaluated)
		if isError(res) {
			return newError("* operation is not supported for %s, %s", curr.Type(), evaluated.Type())
		}

		env.Set(ae.Name.String(), res)
	case "/=":
		curr, ok := env.Get(ae.Name.String())
		if !ok {
			return newError("%s is not defined identifier", ae.Name.String())
		}

		res := evalInfixExpression("/", curr, evaluated)
		if isError(res) {
			return newError("/ operation is not supported for %s, %s", curr.Type(), evaluated.Type())
		}

		env.Set(ae.Name.String(), res)
	case "%=":
		curr, ok := env.Get(ae.Name.String())
		if !ok {
			return newError("%s is not defined identifier", ae.Name.String())
		}

		res := evalInfixExpression("%", curr, evaluated)
		if isError(res) {
			return newError("% operation is not supported for %s, %s", curr.Type(), evaluated.Type())
		}

		env.Set(ae.Name.String(), res)
	default:
		return newError("%s is unknown assignment operator", ae.Operator)
	}

	return nil
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		env := extendIfElseEnv(env)
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		env := extendIfElseEnv(env)
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func extendIfElseEnv(env *object.Environment) *object.Environment {
	return object.NewEnclosedEnvironment(env)
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.SetInner(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	return nativeBoolToBooleanObject(!isTruthy(right))
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	value, ok := right.(object.Real)
	if !ok {
		return newError("unknown operator: -%s", right.Type())
	}

	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	} else {
		return &object.Float{Value: -value.ToFloat64()}
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case areBothRealNumber(left, right):
		return evalRealNumberInfixExpression(operator, left, right)
	case operator == "&&":
		return evalLogicalAndExpression(left, right)
	case operator == "||":
		return evalLogicalOrExpression(left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left.Equals(right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!left.Equals(right))
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "|":
		return &object.Integer{Value: leftVal | rightVal}
	case "&":
		return &object.Integer{Value: leftVal & rightVal}
	case "^":
		return &object.Integer{Value: leftVal ^ rightVal}
	case ">>":
		return &object.Integer{Value: leftVal >> rightVal}
	case "<<":
		return &object.Integer{Value: leftVal << rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func areBothRealNumber(left, right object.Object) bool {
	_, ok := left.(object.Real)
	if !ok {
		return false
	}

	_, ok = right.(object.Real)
	if !ok {
		return false
	}
	return true
}

func evalRealNumberInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(object.Real).ToFloat64()
	rightVal := right.(object.Real).ToFloat64()
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "%":
		return &object.Float{Value: math.Mod(leftVal, rightVal)}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	switch operator {
	case "+":
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.String{Value: leftVal + rightVal}
	case "==":
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.Boolean{Value: leftVal == rightVal}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return newError("array index out of bound: %d", idx)
	}
	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalLogicalAndExpression(left, right object.Object) object.Object {
	leftVal := isTruthy(left)

	if !leftVal {
		return nativeBoolToBooleanObject(false)
	}

	rightVal := isTruthy(right)

	return nativeBoolToBooleanObject(leftVal && rightVal)
}

func evalLogicalOrExpression(left, right object.Object) object.Object {
	leftVal := isTruthy(left)

	if leftVal {
		return nativeBoolToBooleanObject(true)
	}

	rightVal := isTruthy(right)

	return nativeBoolToBooleanObject(leftVal || rightVal)
}
