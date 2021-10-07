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
	_, ok := ae.Left.(*ast.IndexExpression)
	if ok {
		result, ok := evalAssignmentWithIndexExpression(ae, env)
		if !ok {
			return result
		}

		return nil
	}

	newObj := Eval(ae.Value, env)
	if isError(newObj) {
		return newObj
	}

	ident, ok := ae.Left.(*ast.Identifier)
	currObj, ok := env.Get(ident.Value)
	if !ok {
		return newError("%s is not defined identifier", ident.Value)
	}

	res, ok := evalAssignmentOperationHelper(ae.Operator, currObj, newObj)
	if !ok {
		return res
	}

	env.Set(ident.Value, res)

	return nil
}

func evalAssignmentWithIndexExpression(ae *ast.AssignmentExpression, env *object.Environment) (object.Object, bool) {
	ie := ae.Left.(*ast.IndexExpression)
	ident := ie.Left.(*ast.Identifier)
	index := Eval(ie.Index, env)
	if isError(index) {
		return index, false
	}

	newObj := Eval(ae.Value, env)
	if isError(newObj) {
		return newObj, false
	}

	currObj, ok := env.Get(ident.Value)
	if !ok {
		return newError("%s is not defined identifier", ident.Value), false
	}

	switch {
	case currObj.Type() == object.ARRAY_OBJ:
		arr := currObj.(*object.Array)
		idx := index.(*object.Integer).Value
		max := int64(len(arr.Elements) - 1)
		if idx < 0 || idx > max {
			return newError("array index out of bound: %d", idx), false
		}

		res, ok := evalAssignmentOperationHelper(ae.Operator, arr.Elements[idx], newObj)
		if !ok {
			return res, false
		}

		arr.Elements[idx] = res
	case currObj.Type() == object.HASH_OBJ:
		hash := currObj.(*object.Hash)

		idx, ok := index.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", index.Type()), false
		}

		pair, ok := hash.Pairs[idx.HashKey()]
		if !ok {
			// It means key doesn't exist in hash. so add new key,value to hash if assign operator
			if ae.Operator == "=" {
				hash.Pairs[idx.HashKey()] = object.HashPair{Key: index, Value: newObj}
				return nil, true
			}
			return newError("%+v is not exist in hash", index), false
		}

		res, ok := evalAssignmentOperationHelper(ae.Operator, pair.Value, newObj)
		if !ok {
			return res, false
		}

		hash.Pairs[idx.HashKey()] = object.HashPair{Key: index, Value: res}
	default:
		return newError("%s is unknown index type, %T", ident.Value, ident), false
	}

	return nil, true
}

func evalAssignmentOperationHelper(op string, curr, rightOperand object.Object) (object.Object, bool) {
	switch op {
	case "=":
		return rightOperand, true
	case "+=":
		res := evalInfixExpression("+", curr, rightOperand)
		if isError(res) {
			return newError("+= operation is not supported for %s, %s", curr.Type(), rightOperand.Type()), false
		}
		return res, true
	case "-=":
		res := evalInfixExpression("-", curr, rightOperand)
		if isError(res) {
			return newError("-= operation is not supported for %s, %s", curr.Type(), rightOperand.Type()), false
		}
		return res, true
	case "*=":
		res := evalInfixExpression("*", curr, rightOperand)
		if isError(res) {
			return newError("* operation is not supported for %s, %s", curr.Type(), rightOperand.Type()), false
		}
		return res, true
	case "/=":
		res := evalInfixExpression("/", curr, rightOperand)
		if isError(res) {
			return newError("/ operation is not supported for %s, %s", curr.Type(), rightOperand.Type()), false
		}
		return res, true
	case "%=":
		res := evalInfixExpression("%", curr, rightOperand)
		if isError(res) {
			return newError("% operation is not supported for %s, %s", curr.Type(), rightOperand.Type()), false
		}
		return res, true
	default:
		return newError("%s is unknown assignment operator", op), false
	}
}

func evalCallExpression(ce *ast.CallExpression, env *object.Environment) object.Object {
	funcName := Eval(ce.Function, env)
	if isError(funcName) {
		return funcName
	}
	args := evalExpressions(ce.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	return applyFunction(funcName, args)
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

func evalMethodCallExpression(mce *ast.MethodCallExpression, env *object.Environment) object.Object {
	obj := Eval(mce.Object, env)
	if isError(obj) {
		return obj
	}

	method, ok := mce.Call.(*ast.CallExpression)
	if !ok {
		return newError("wrong type method: %s", method.Function.String())
	}

	args := evalExpressions(method.Arguments, env)

	callable, ok := obj.(object.Callable)
	if !ok {
		return newError("%s is not callable object", obj.Type())
	}

	result, ok := callable.Apply(method.Function.String(), env, args...)
	if !ok {
		return newError("%s is unknown method, %s", method.Function.String(), obj.Type())
	}

	return result
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
