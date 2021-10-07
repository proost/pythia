package evaluator

import (
	"os"
	"pythia/ast"
	"pythia/object"
)

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalLetStatement(ls *ast.LetStatement, env *object.Environment) object.Object {
	val := Eval(ls.Value, env)
	if isError(val) {
		return val
	}
	env.Set(ls.Name.Value, val)

	return nil
}

func evalIfStatement(is *ast.IfStatement, env *object.Environment) object.Object {
	condition := Eval(is.Condition, env)
	if isError(condition) {
		return condition
	}

	var result object.Object

	if isTruthy(condition) {
		result = Eval(is.Consequence, extendIfElseEnv(env))
	} else if is.Alternative != nil {
		result = Eval(is.Alternative, extendIfElseEnv(env))
	}

	if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
		return result
	}

	return nil
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

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return nil
}

func evalFunctionStatement(fn *ast.FunctionStatement, env *object.Environment) object.Object {
	params := fn.Parameters
	name := fn.Name
	body := fn.Body

	funcObj := &object.Function{Parameters: params, Name: name, Body: body}
	env.Set(name.Value, funcObj)
	funcObj.Env = env

	return nil
}

func evalReturnStatement(ret *ast.ReturnStatement, env *object.Environment) object.Object {
	val := Eval(ret.ReturnValue, env)
	if isError(val) {
		return val
	}
	if ret.ReturnValue == nil && val == nil {
		return &object.ReturnValue{Value: nil} // void return
	}

	return &object.ReturnValue{Value: val}
}

func evalInstructionStatement(instruction *ast.InstructionStatement) object.Object {
	switch instruction.Instruction {
	case "quit":
		os.Exit(0)
	}

	return newError("unknown instruction: %s", instruction.Instruction)
}

func evalForStatement(forStmt *ast.ForStatement, env *object.Environment) object.Object {

	container := Eval(forStmt.Container, env)

	iter, ok := container.(object.Iterator)
	if !ok {
		return newError("%s object doesn't implement the Iterator interface", container.Type())
	}

	// Initialize index, value in for-loop
	var variables []*ast.Identifier
	variables = append(variables, forStmt.Value)
	if forStmt.Index != nil {
		variables = append(variables, forStmt.Index)
	}
	extendedEnv := extendForLoopEnv(variables, env)

	// loop iterator object
	iter.Reset()

	required, optional, ok := iter.Next()

	for ok {

		extendedEnv.SetInner(forStmt.Value.Value, required)

		if forStmt.Index != nil {
			if container.Type() == object.HASH_OBJ {
				extendedEnv.SetInner(forStmt.Index.Value, required)
				extendedEnv.SetInner(forStmt.Value.Value, optional)
			} else {
				extendedEnv.SetInner(forStmt.Index.Value, optional)
			}
		}

		body := Eval(forStmt.Body, extendedEnv)
		if body != nil && body.Type() == object.RETURN_VALUE_OBJ {
			return nil
		}
		if isError(body) {
			return body
		}

		required, optional, ok = iter.Next()
	}

	return nil
}

func extendForLoopEnv(variables []*ast.Identifier, env *object.Environment) *object.Environment {
	extendedEnv := object.NewEnclosedEnvironment(env)

	for _, v := range variables {
		extendedEnv.Set(v.Value, NULL)
	}

	return extendedEnv
}
