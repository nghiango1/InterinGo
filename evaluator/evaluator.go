package evaluator

import (
	"main/ast"
	"main/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
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

func evalIntegerLiteral(node ast.Node) object.Object {
	intLiteral, _ := node.(*ast.IntegerLiteral)
	intObj := &object.Integer{
		Value: intLiteral.Value,
	}
	return intObj
}
