package evaluator

import (
	"main/ast"
	"main/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.Boolean:
		return evalBooleanLiteral(node)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}

func evalBooleanLiteral(node ast.Node) object.Object {
	boolLiteral, _ := node.(*ast.Boolean)
	if boolLiteral.Value {
		return TRUE
	}
	return FALSE
}

func evalIntegerLiteral(node ast.Node) object.Object {
	intLiteral, _ := node.(*ast.IntegerLiteral)
	intObj := &object.Integer{
		Value: intLiteral.Value,
	}
	return intObj
}
