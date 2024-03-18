package handlers

import (
	"errors"
	"fmt"
	"interingo/ast"
	"strings"

	_ "github.com/tliron/commonlog/simple"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// /**
//  * Value-object describing what options formatting should use.
//  */
// interface FormattingOptions {
// 	/** * Size of a tab in spaces. */
// 	tabSize: uinteger;
//
// 	/** * Prefer spaces over tabs. */
// 	insertSpaces: boolean;
//
// 	/** * Trim trailing whitespace on a line. */
// 	trimTrailingWhitespace?: boolean;
//
// 	/** * Insert a newline character at the end of the file if one does not exist. */
// 	insertFinalNewline?: boolean;
//
// 	/** * Trim all newlines after the final newline at the end of the file. */
// 	trimFinalNewlines?: boolean;
//
// 	/** * Signature for further properties. */
// 	[key: string]: boolean | integer | string;
// }

func indentPadding(option protocol.FormattingOptions, indent int) (string, error) {
	insertSpaces, ok := option[protocol.FormattingOptionInsertSpaces]
	if !ok {
		insertSpaces = true
	}

	boolInsertSpaces, ok := insertSpaces.(bool)
	if !ok {
		return "", errors.New(fmt.Sprintf("Formating option insertSpaces have wrong value, got = %v instead", insertSpaces))
	}

	if !boolInsertSpaces {
		return strings.Repeat("\t", indent), nil
	}

	var spaceTab string

	tabSize, ok := option[protocol.FormattingOptionTabSize]
	if !ok {
		tabSize = 4
	}

	intTabSize, ok := tabSize.(int)
	if !ok {
		return "", errors.New(fmt.Sprintf("Formating option insertSpaces have wrong value, got = %v instead", insertSpaces))
	}

	spaceTab = strings.Repeat(" ", intTabSize)

	return strings.Repeat(spaceTab, indent), nil
}

func FormatedFunctionLiteral(node *ast.FunctionLiteral, option protocol.FormattingOptions, indent int) string {
	indentPad, err := indentPadding(option, indent)
	if err != nil {
		indentPad = strings.Repeat("    ", indent)
	}

	formated := node.TokenLiteral()

	// Format space between "fn" token and params
	if false {
		formated += " "
	}
	formated += "("

	// Format space between "(" and the first param
	if false {
		formated += " "
	}

	for index, param := range node.Parameters {
		formated += param.Value
		if index < len(node.Parameters)-1 {
			formated += ", "
		}
	}

	// Format space between ")" and the last param
	if false {
		formated += " "
	}
	formated += ")"

	// Format space between ")" and function body, could be `\n`
	if false {
		formated += " "
	}

	formated += "{"

	// Format space between "{" and function body
	if true {
		formated += "\n"
	}

	stmtFormated := FormatedAST(node.Body, option, indent+1)
	formated += stmtFormated
	// Format space between function body and "}"
	if true {
		formated += "\n" + indentPad
	}
	formated += "}"

	return formated
}

func FormatedExpresionAST(node ast.Node, option protocol.FormattingOptions, indent int) string {
	var formated string
	indentPad, err := indentPadding(option, indent)
	if err != nil {
		indentPad = strings.Repeat("    ", indent)
	}
	switch node := node.(type) {
	case *ast.InfixExpression:
		formated += FormatedExpresionAST(node.Left, option, indent)
		formated += node.Operator
		formated += FormatedExpresionAST(node.Right, option, indent)

	case *ast.PrefixExpression:
		formated += node.Operator
		formated += FormatedExpresionAST(node.Right, option, indent)

	case *ast.FunctionLiteral:
		formated := FormatedFunctionLiteral(node, option, indent)
		return formated

	case *ast.IfExpression:
		formated = node.TokenLiteral()
		formated += " ("
		formated += FormatedExpresionAST(node.Condition, option, indent)
		formated += ") {\n"
		formated += FormatedAST(node.Consequence, option, indent+1)
		formated += "\n" + indentPad + "}"
		if node.Alternative != nil {
			formated += " else {\n"
			formated += FormatedAST(node.Alternative, option, indent+1)
			formated += "\n" + indentPad + "}"
		}

	case *ast.Identifier:
		formated = node.Value

	case *ast.Boolean:
		formated = node.TokenLiteral()

	case *ast.IntegerLiteral:
		formated = node.TokenLiteral()

	case *ast.CallExpression:
		formated = FormatedExpresionAST(node.Function, option, indent)

		formated += "("
		for index, args := range node.Arguments {
			formated += FormatedExpresionAST(args, option, indent)
			if index < len(node.Arguments)-1 {
				formated += ", "
			}
		}
		formated += ")"

	default:
		return node.String()
	}
	return formated
}

// Format statement
func FormatedAST(node ast.Node, option protocol.FormattingOptions, indent int) string {
	var formated string
	indentPad, err := indentPadding(option, indent)
	if err != nil {
		indentPad = strings.Repeat("    ", indent)
	}
	switch node := node.(type) {
	case *ast.Program:
		for index, statement := range node.Statements {
			stmtFormated := FormatedAST(statement, option, indent)
			formated += indentPad + stmtFormated

			if false {
				formated += fmt.Sprint(index)
			}

			if true {
				formated += ";\n"
			}
		}

	case *ast.BlockStatement:
		for index, statement := range node.Statements {
			stmtFormated := FormatedAST(statement, option, indent)
			formated += indentPad + stmtFormated

			if false {
				formated += fmt.Sprint(index)
			}

			if true {
				formated += ";"
			}

			if index < len(node.Statements)-1 {
				formated += "\n"
			}
		}

	case *ast.LetStatement:
		formated += node.TokenLiteral()
		formated += " "
		formated += node.Name.String()
		formated += " = "
		formated += FormatedExpresionAST(node.Value, option, indent)

	case *ast.ReturnStatement:
		formated += node.TokenLiteral()
		formated += " "
		formated += FormatedExpresionAST(node.ReturnValue, option, indent)

	case *ast.ExpressionStatement:
		formated += FormatedExpresionAST(node.Expression, option, indent)

	default:
		return node.String()
	}
	return formated
}
