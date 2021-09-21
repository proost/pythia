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
	case "exit":
		os.Exit(0)
	}

	return newError("unknown instruction: %s", instruction.Instruction)
}