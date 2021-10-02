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

func evalIfStatement(is *ast.IfStatement, env *object.Environment) object.Object {
	condition := Eval(is.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		Eval(is.Consequence, extendIfElseEnv(env))
	} else if is.Alternative != nil {
		Eval(is.Alternative, extendIfElseEnv(env))
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

	idx, val, ok := iter.Next()

	for ok {

		if forStmt.Index != nil {
			extendedEnv.SetInner(forStmt.Index.Value, idx)
		}

		extendedEnv.SetInner(forStmt.Value.Value, val)

		body := Eval(forStmt.Body, extendedEnv)
		if body != nil && body.Type() == object.RETURN_VALUE_OBJ {
			return nil
		}
		if isError(body) {
			return body
		}

		idx, val, ok = iter.Next()
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
