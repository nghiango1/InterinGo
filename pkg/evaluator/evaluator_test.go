package evaluator_test

import (
	"interingo/pkg/ast"
	"interingo/pkg/lexer"
	"interingo/pkg/object"
	"interingo/pkg/parser"
	"interingo/pkg/evaluator"
	"interingo/pkg/test"
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

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return evaluator.Eval(program, env)
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
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}
func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj,
			obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
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
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
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

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != evaluator.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`if (10 > 1) {
if (10 > 1) {
  return 10;
}
return 1;
}
`,
			10,
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
			`if (10 > 1) {
	if (10 > 1) {
		true + false;
		return 10;
	}
	return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
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
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"2/0",
			"divide by zero: 2 / 0",
		},
		{
			"if (2/0 > 0) { 1 } else { 2 }",
			"divide by zero: 2 / 0",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("%v\n", tt.input)
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
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
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated,
			evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expectedBody := "{ (x + 2) }"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody,
			fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { return x; }; identity(identity(5));", 5},
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
		{"return 5;", 5},
	}
	for _, tt := range tests {
		result := testEval(tt.input)
		testIntegerObject(t, result, tt.expected)
	}
}

func Test_evalBuiltInObject(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		b    object.BuiltIn
		args []ast.Expression
		want map[string]object.Object
	}{
		{
			"Standard",
			&test.MockBuiltinImpl{MockParameters: object.Parameters{
				Standard: []*ast.Identifier{{Value: "code"}},
			}},
			[]ast.Expression{
				&ast.IntegerLiteral{Value: 1},
			},
			map[string]object.Object{
				"code": &object.Integer{Value: 1},
			},
		},
		{
			"Default",
			&test.MockBuiltinImpl{MockParameters: object.Parameters{
				Default: []object.DefaultParameter{
					{
						Key:   &ast.Identifier{Value: "code"},
						Value: &object.Integer{Value: 1000},
					},
				},
			}},
			[]ast.Expression{
				&ast.IntegerLiteral{Value: 1},
			},
			map[string]object.Object{
				"code": &object.Integer{Value: 1},
			},
		},
		{
			"Default 2",
			&test.MockBuiltinImpl{MockParameters: object.Parameters{
				Default: []object.DefaultParameter{
					{
						Key:   &ast.Identifier{Value: "code"},
						Value: &object.Integer{Value: 1000},
					},
				},
			}},
			nil,
			map[string]object.Object{
				"code": &object.Integer{Value: 1000},
			},
		},
		{
			"Standard and Default",
			&test.MockBuiltinImpl{MockParameters: object.Parameters{
				Standard: []*ast.Identifier{{Value: "code"}},
				Default: []object.DefaultParameter{
					{
						Key:   &ast.Identifier{Value: "default"},
						Value: &object.Integer{Value: 1000},
					},
				},
			}},
			[]ast.Expression{
				&ast.IntegerLiteral{Value: 1},
			},
			map[string]object.Object{
				"code":    &object.Integer{Value: 1},
				"default": &object.Integer{Value: 1000},
			},
		},
		{
			"Standard and Default 2",
			&test.MockBuiltinImpl{MockParameters: object.Parameters{
				Standard: []*ast.Identifier{{Value: "code"}},
				Default: []object.DefaultParameter{
					{
						Key:   &ast.Identifier{Value: "default"},
						Value: &object.Integer{Value: 1000},
					},
				},
			}},
			[]ast.Expression{
				&ast.IntegerLiteral{Value: 1},
				&ast.IntegerLiteral{Value: 1},
			},
			map[string]object.Object{
				"code":    &object.Integer{Value: 1},
				"default": &object.Integer{Value: 1},
			},
		},
		{
			"Rest",
			&test.MockBuiltinImpl{MockParameters: object.Parameters{
				Rest: true,
			}},
			[]ast.Expression{
				&ast.IntegerLiteral{Value: 1},
				&ast.IntegerLiteral{Value: 1},
			},
			map[string]object.Object{
				evaluator.REST_ARGS: &object.Array{
					Value: []object.Object{
						&object.Integer{Value: 1},
						&object.Integer{Value: 1},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator.EvalBuiltInObject(tt.b, tt.args)
			// TODO: update the condition below to compare got with tt.want.
			for k, v := range tt.want {
				got, ok := tt.b.Env().Get(k)
				if !ok {
					t.Errorf("evalBuiltInObject() = %v, want %v", got, tt.want)
				}
				if got.Inspect() != v.Inspect() || got.Type() != v.Type() {
					t.Errorf("evalBuiltInObject() = %v, want %v", got, tt.want)
				}

			}
		})
	}
}
