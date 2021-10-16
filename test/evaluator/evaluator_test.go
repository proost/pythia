package evaluator

import (
	"go/types"
	"pythia/evaluator"
	"pythia/lexer"
	"pythia/object"
	"pythia/parser"
	"testing"
)

func TestEvalEqualExpression(t *testing.T) {
	tests := []struct {
		left     string
		right    string
		expected bool
	}{
		{"1", "1", true},
		{"1", "2", false},
		{"1.1", "1.1", true},
		{"1.5", "2.5", false},
		{"true", "true", true},
		{"true", "false", false},
		{"null", "null", true},
		{`"abc"`, `"abc"`, true},
		{`"abc"`, `"bcd"`, false},
		{"[1,2,3]", "[1,2,3]", true},
		{"[1]", "[2]", false},
		{"type(1)", "type(2)", true},
		{"type(1)", "type(1.0)", false},
		{"{}", "{}", true},
		{`{"a": 1, "b": 2}`, `{"b": 2, "a": 1}`, true},
		{`{"a": 2, "b": 1}`, `{"a": 1, "b": 2}`, false},
	}

	for _, tt := range tests {
		leftEvaluated := testEval(tt.left)
		rightEvaluated := testEval(tt.right)
		if leftEvaluated.Equals(rightEvaluated) != tt.expected {
			t.Errorf("expected is wrong. left - (%+v), right - (%+v)", leftEvaluated, rightEvaluated)
		}
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"5 % 2", 1},
		{"40 % (2+5)", 5},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"10 | 20", 30},
		{"15 & 21", 5},
		{"10 ^ 20", 30},
		{"100 >> 2", 25},
		{"10 << 2", 40},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.0", 5.0},
		{"10.1", 10.1},
		{"-5.2", -5.2},
		{"5.0 + 5.0 + 5.0 + 5.0 - 10.0", 10.0},
		{"2.0 * 2.0 * 2.0 * 2.0 * 2.0", 32.0},
		{"-50.1 + 100.2 + -50.1", 0.0},
		{"5.0 * 2.0 + 10.5", 20.5},
		{"5 + 2.0 * 10", 25.0},
		{"20 + 2 * -10.0", 0.0},
		{"2.5 % 5", 2.5},
		{"7.2 % 5.2", 2.0},
		{"50.0 / 2.0 * 2.0 + 10.0", 60.0},
		{"2 * (5.2 + 10.8)", 32.0},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1.0 < 1.0", false},
		{"1.0 > 2.0", false},
		{"1 <= 2", true},
		{"1 >= 2", false},
		{"1.9 <= 1.0", false},
		{"1.0 >= 2.0", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1.1 == 2.2", false},
		{"1.1 != 2.2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"(1 < 2) == true", true},
		{"(1 > 2) == false", true},
		{"true && true", true},
		{"true && false", false},
		{"true || false", true},
		{"false || false", false},
		{`"a" && false`, false},
		{`"a" && "b"`, true},
		{`"a" || false`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let a = 0; if (true) { a = 10 }; a;", 10},
		{"let a = 0; if (false) { a = 10 }; a;", 0},
		{"let a = 0; if (1) { a = 10 }; a;", 10},
		{"let a = 0; if (1 < 2) { a = 10 }; a;", 10},
		{"let a = 0; if (1 > 2) { a = 10 }; a;", 0},
		{"let a = 0; if (1 > 2) { a = 10 } else { a = 20 }; a;", 20},
		{"let a = 0; if (1 < 2) { a = 10 } else { a = 20 }; a;", 10},
		{`func min(a,b) { if (a>b) { b } else { a } }; min(1,2)`, nil},
		{`func max(a,b) { if (a>b) { return a } else { return b }}; max(1,2)`, 2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case types.Nil:
			if evaluated != nil {
				t.Errorf("object is not nil. got=%T (%+v)", evaluated, expected)
			}
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"return", nil},
		{
			`func f(x) {
				  	return x;
				  	x + 10;
					};
					f(10);`,
			10,
		},
		{
			` func f(x) {
					   let result = x + 10;
					   return result;
					   return 10;
					};
					f(10);`,
			20,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			if evaluated != nil {
				t.Errorf("object has wrong type. got=%T, want=%T", evaluated, expected)
			}
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			".name",
			"unknown instruction: name",
		},
		{
			"[1, 2, 3][3]",
			"array index out of bound: 3",
		},
		{
			"[1, 2, 3][-1]",
			"array index out of bound: -1",
		},
		{
			"{[1,2]: 3}",
			"unusable as hash key: ARRAY",
		},
		{
			`{{"a":1} : 2}`,
			"unusable as hash key: HASH",
		},
		{
			`let a = "abc"; a -= "bcd";`,
			"-= operation is not supported for STRING, STRING",
		},
		{
			"range(1,5,-1)",
			"start can't be smaller than end, when step is -1",
		},
		{
			"range(-1,-5,1)",
			"start can't be bigger than end, when step is 1",
		},
		{
			`let h = {}; h[{1: true}] = 2;`,
			"unusable as hash key: HASH",
		},
		{
			`let h = {}; h[[1,2]] = 3;`,
			"unusable as hash key: ARRAY",
		},
		{
			`let h = {}; func add(){}; h[add] = 1;`,
			"unusable as hash key: FUNCTION",
		},
		{
			`let foobar = 1; foobar.call()`,
			"INTEGER is not callable object",
		},
		{
			`{}.foo()`,
			"foo is unknown method, HASH",
		},
		{
			`[].bar()`,
			"bar is unknown method, ARRAY",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}

}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestAssignmentExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let a = 1; a = 2; a;", 2},
		{"let a = 1.5; a += 1; a;", 2.5},
		{`let a = "ab"; a += "cd"; a;`, "abcd"},
		{"let a = 2; a -= 1.5; a;", 0.5},
		{"let a = 1.0; a *= 2.0; a;", 2.0},
		{"let a = 4; a /= 2; a;", 2},
		{"let a = 5; a %= 4; a;", 1},
		{"let a = 1; if(1<2) { a = 2 }; a;", 2},
		{"let a = 1; func t() { a = 2 }; t(); a;", 2},
		{`let h = {"a": 1}; h["a"] = 2; h["a"];`, 2},
		{`let h1 = {"a": 1}; let h2 = {1: "a"}; h1[h2[1]];`, 1},
		{`let arr = [0,1,2]; arr[0] = 3; arr[0];`, 3},
		{`let h = {"a": 1}; if(1<2) { h["a"] += 2 }; h["a"];`, 3},
		{`let h = {"a": 1}; h["b"] = 2; h["b"];`, 2},
	}

	for _, tt := range tests {
		switch expected := tt.expected.(type) {
		case string:
			evaluated := testEval(tt.input)

			str, ok := evaluated.(*object.String)
			if !ok {
				t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
			}

			if str.Value != expected {
				t.Errorf("Identifier has wrong value. got=%q", str.Value)
			}
		case float64:
			testFloatObject(t, testEval(tt.input), expected)
		case float32:
			testFloatObject(t, testEval(tt.input), float64(expected))
		case int:
			testIntegerObject(t, testEval(tt.input), int64(expected))
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func a(x) { x + 2; }; a;"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("prameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"func identify(x) { return x; }; identify(5);", 5},
		{"func identify(x) { return x; }; identify(5);", 5},
		{"func double(x) { return x * 2; }; double(5);", 10},
		{"func add(x, y) { return x + y; }; add(5, 5);", 10},
		{"func add(x, y) { return x + y; }; add(5 + 5, add(5, 5));", 20},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to len not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`len({"a": 1})`, 1},
		{"append([], 1)", []int{1}},
		{"append([1], 2)", []int{1, 2}},
		{"append([1], 2, 3)", "wrong number of arguments. got=3, want=2"},
		{"append(1, 2)", "argument to append must be ARRAY, got INTEGER"},
		{"range(1,4)", []int{1, 2, 3}},
		{"range(1,4,2)", []int{1, 3}},
		{"range(5,1,-1)", []int{5, 4, 3, 2}},
		{"range(-1,-5,-2)", []int{-1, -3}},
		{`let h = {"a": 1}; delete(h, "a"); h`, map[object.HashKey]object.HashPair{}},
		{`let h = {"a": 1, "b": 2}; delete(h, "a", "b"); h`, "wrong number of arguments. got=3, want=2"},
		{`[].isEmpty()`, true},
		{`[1,2].isEmpty()`, false},
		{`[1,2,3].last()`, 3},
		{`{}.isEmpty()`, true},
		{`{1: true}.isEmpty()`, false},
		{`{true: "a", false: "b"}.keys()`, []bool{true, false}},
		{`{true: "a", false: "b"}.values()`, []string{"a", "b"}},
		{`string([1,2,3], 4)`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case bool:
			testBooleanObject(t, evaluated, expected)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		case *object.Null:
			testNullObject(t, evaluated)
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testIntegerObject(t, array.Elements[i], int64(expectedElem))
			}
		case map[object.HashKey]object.HashPair:
			hash, ok := evaluated.(*object.Hash)
			if !ok {
				t.Errorf("obj not Hash. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			for k, pair := range hash.Pairs {
				if expected[k].Value != pair.Value {
					t.Errorf("object has wrong value. got=%+v, want=%+v", pair.Value, expected[k].Value)
				}
			}
		case []bool:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Elements))
				continue
			}

			// order is not demanded
			result := make([]bool, len(expected))
			for i, evaluatedEl := range expected {
				result[i] = evaluatedEl
			}

			t1 := true
			for i := 0; i < len(result); i++ {
				if result[i] != expected[i] {
					t1 = false
				}
			}

			t2 := true
			for i := 0; i < len(result); i++ {
				if result[i] != expected[len(expected)-1-i] {
					t2 = false
				}
			}

			if !(t1 || t2) {
				t.Errorf("object has wrong value. got=%+v, want=%+v", result, expected)
			}
		case []string:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Elements))
				continue
			}

			// order is not demanded
			result := make([]string, len(expected))
			for i, evaluatedEl := range expected {
				result[i] = evaluatedEl
			}

			t1 := true
			for i := 0; i < len(result); i++ {
				if result[i] != expected[i] {
					t1 = false
				}
			}

			t2 := true
			for i := 0; i < len(result); i++ {
				if result[i] != expected[len(expected)-1-i] {
					t2 = false
				}
			}

			if !(t1 || t2) {
				t.Errorf("object has wrong value. got=%+v, want=%+v", result, expected)
			}
		}
	}
}

func TestBuiltinTypeFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"type(1)", "Type: INTEGER"},
		{"type(null)", "Type: NULL"},
		{`type("s")`, "Type: STRING"},
		{"type(type)", "Type: BUILTIN"},
		{"type({})", "Type: HASH"},
		{"func a(){}; type(a())", "Type: NULL"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		typeObj, ok := evaluated.(*object.Type)
		if !ok {
			t.Errorf("object is not Type. got=%T (%+v)", evaluated, evaluated)
		}
		if typeObj.Inspect() != tt.expected {
			t.Errorf("wrong type. expected=%s, got=%s", tt.expected, typeObj.Inspect())
		}
	}
}

func TestBuiltinStringFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"string(1)", "1"},
		{"string(2.3)", "2.300000"},
		{"string(true)", "true"},
		{`string("foo")`, "foo"},
		{"string(null)", "null"},
		{"string({1: true, 0: false})", "{1: true, 0: false}"},
		{"string([1,2,3])", "[1, 2, 3]"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		result, ok := evaluated.(*object.String)
		if !ok {
			t.Errorf("object is not String. got=%T (%+v)", evaluated, evaluated)
		}

		if result.Value != tt.expected {
			t.Errorf("object has wrong value. got=%+v, want=%+v", result, tt.expected)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T, (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}

}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		evaluator.TRUE.HashKey():                   5,
		evaluator.FALSE.HashKey():                  6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Fatalf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}

}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
		{
			`{4.999999: 5}[4.999999]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestNullLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"null", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNullObject(t, evaluated)
	}
}

func TestForLoopStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let result = 0; for v in [10,20,30] { result += v }; result;", 60},
		{`let result = 0; for i,v in [10,20,30] { result += v }; result;`, 60},
		{`let result = 0; for i,v in [10,20,30] { result += i }; result;`, 3},
		{`let result = 0; for k,v in {1: 10, 2: 20, 3: 30 } { result += v } ; result;`, 60},
		{`let result = 0; for k,v in {1: 10, 2: 20, 3: 30 } { result += k } ; result;`, 6},
		{`let result = 0; for k in {1: 10, 2: 20, 3: 30 } { result += k } ; result;`, 6},
		{`let result = ""; for c in "abc" { result += c }; result;`, "abc"},
		{`let result = ""; for i,c in "abc" { result += c }; result;`, "abc"},
		{`let result = 0; for i,c in "abc" { result += i }; result;`, 3},
	}

	for _, tt := range tests {
		switch expected := tt.expected.(type) {
		case string:
			evaluated := testEval(tt.input)

			str, ok := evaluated.(*object.String)
			if !ok {
				t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
			}

			if str.Value != expected {
				t.Errorf("Identifier has wrong value. got=%q", str.Value)
			}
		case int:
			testIntegerObject(t, testEval(tt.input), int64(expected))
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return evaluator.Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != evaluator.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
