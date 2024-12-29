package evaluator

import (
	"gomonkey/lexer"
	"gomonkey/object"
	"gomonkey/parser"
	"testing"
)

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
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
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
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{`"Hello" == "World!"`, false},
		{`"Hello" == "Hello"`, true},
		{`"Hello" != "World!"`, true},
		{`"Hello" != "Hello"`, false},
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

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		expected interface{}
		input    string
	}{
		{input: "if (true) { 10 }", expected: 10},
		{input: "if (false) { 10 }", expected: nil},
		{input: "if (1) { 10 }", expected: 10},
		{input: "if (1 < 2) { 10 }", expected: 10},
		{input: "if (1 > 2) { 10 }", expected: nil},
		{input: "if (1 > 2) { 10 } else { 20 }", expected: 20},
		{input: "if (1 < 2) { 10 } else { 20 }", expected: 10},
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

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{input: "return 10;", expected: 10},
		{input: "return 10; 9;", expected: 10},
		{input: "return 2*5; 9;", expected: 10},
		{input: "9; return 2*5; 9;", expected: 10},
		{
			input: `if (10 > 1) {
        if (10 > 1) {
          return 10;
        }

        return 1;
      }`, expected: 10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			input:           "5 + true;",
			expectedMessage: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			input:           "5 + true; 5;",
			expectedMessage: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			input:           "-true",
			expectedMessage: "unknown operator: -BOOLEAN",
		},
		{
			input:           "true + false;",
			expectedMessage: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input:           "5; true + false; 5",
			expectedMessage: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input:           "if (10 > 1) { true + false; }",
			expectedMessage: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input: `if (10 > 1) {
        if (10 > 1) {
          return true + false;
        }

        return 1;
      }`,
			expectedMessage: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input:           "foobar",
			expectedMessage: "identifier not found: foobar",
		},
		{
			input:           `"Hello" - "World!";`,
			expectedMessage: "unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("test: %s. no error object returned. got=%T(%+v)", tt.input, evaluated, evaluated)
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
		{input: "let a = 5; a;", expected: 5},
		{input: "let a = 5 * 5; a;", expected: 25},
		{input: "let a = 5; let b = a; b;", expected: 5},
		{input: "let a = 5; let b = a; let c = a + b + 5; c;", expected: 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if bodyString := fn.Body.String(); bodyString != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, bodyString)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{input: "let identity = fn(x) { x; }; identity(5);", expected: 5},
		{input: "let identity = fn(x) { return x; }; identity(5);", expected: 5},
		{input: "let double = fn(x) { x * 2; }; double(5);", expected: 10},
		{input: "let add = fn(x, y) { x + y; }; add(5, 5);", expected: 10},
		{input: "let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", expected: 20},
		{input: "fn(x) { x; }(5)", expected: 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!";`

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
	input := `"Hello" + " " + "World!";`

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
		expected interface{}
		input    string
	}{
		{input: `len("")`, expected: 0},
		{input: `len("four")`, expected: 4},
		{input: `len("hello world")`, expected: 11},
		{input: `len([1, 2])`, expected: 2},
		{input: `len([])`, expected: 0},
		{input: `len(1)`, expected: "`len` builtin function doesn't support argument of type INTEGER"},
		{input: `len("one", "two")`, expected: "wrong number of arguments. got=2, want=1"},
		{input: `first([1, 2])`, expected: 1},
		{input: `first([])`, expected: nil},
		{input: `last([1, 2])`, expected: 2},
		{input: `last([])`, expected: nil},
		{input: `rest([1, 2, 3])`, expected: []int64{2, 3}},
		{input: `rest([])`, expected: nil},
		{input: `push([1, 2], 3)`, expected: []int64{1, 2, 3}},
		{input: `push([], 1)`, expected: []int64{1}},
		{input: `exit(0, 1)`, expected: "wrong number of arguments. got=2, want 0 or 1"},
		{input: `exit("1")`, expected: "`exit` builtin function doesn't support argument of type STRING"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		case []int64:
			arrObj, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("object is not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			for idx, v := range arrObj.Elements {
				intValue := int64(v.(*object.Integer).Value)
				if expected[idx] != intValue {
					t.Errorf("value at idx %d don't match. got=%d expecting %d", idx, intValue, expected[idx])
				}
			}

		}

	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)

	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
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
		expected interface{}
		input    string
	}{
		{
			1,
			"[1, 2, 3][0]",
		},
		{
			2,
			"[1, 2, 3][1]",
		},
		{
			3,
			"[1, 2, 3][2]",
		},
		{
			1,
			"let i = 0; [1][i];",
		},
		{
			3,
			"[1, 2, 3][1 + 1];",
		},
		{
			3,
			"let myArray = [1, 2, 3]; myArray[2];",
		},
		{
			6,
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
		},
		{
			2,
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
		},
		{
			3,
			"let myArray = [1, 2, 3]; myArray[-1];",
		},
		{
			1,
			"let myArray = [1, 2, 3]; myArray[-3];",
		},
		{
			nil,
			"let myArray = [1, 2, 3]; myArray[-4];",
		},
		{
			nil,
			"[1, 2, 3][3]",
		},
		{
			3,
			"[1, 2, 3][-1]",
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

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	env := object.NewEnvironment()

	return Eval(program, env)
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
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}
