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

	return result
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
			extendedEnv.Set(forStmt.Index.Value, idx)
		}

		extendedEnv.Set(forStmt.Value.Value, val)

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
