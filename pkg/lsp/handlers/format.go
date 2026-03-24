package handlers

import (
	"encoding/json"
	"fmt"
	"interingo/pkg/ast"
	"interingo/pkg/token"
	"io"
	"log/slog"
	"strings"

	_ "github.com/tliron/commonlog/simple"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type FormattingOptions struct {
	// Size of a tab in spaces.
	TabSize *int `json:"tabSize,omitempty"`
	// 	Prefer spaces over tabs.
	InsertSpaces *bool `json:"insertSpaces,omitempty"`
	// 	Trim trailing whitespace on a line.
	TrimTrailingWhitespace *bool `json:"trimTrailingWhitespace,omitempty"`
	// 	Insert a newline character at the end of the file if one does not exist.
	InsertFinalNewline *bool `json:"insertFinalNewline,omitempty"`
	// 	Trim all newlines after the final newline at the end of the file.
	TrimFinalNewlines *bool `json:"trimFinalNewlines,omitempty"`
}

func indentPadding(option protocol.FormattingOptions, indent int) (string, error) {
	var fo FormattingOptions
	d, err := json.Marshal(option)
	if err != nil {
	} else {
		err = json.Unmarshal(d, &fo)
	}

	if fo.InsertSpaces != nil {
		if *(fo.InsertSpaces) == true {
			return strings.Repeat("\t", indent), nil
		}
	}

	var spaceTab string

	tabSize := 4
	if fo.TabSize != nil {
		tabSize = *(fo.TabSize)
	}
	spaceTab = strings.Repeat(" ", tabSize)

	return strings.Repeat(spaceTab, indent), nil
}

var comments []token.Token
var curr_comments int

func FormatedFunctionLiteral(node *ast.FunctionLiteral, option protocol.FormattingOptions, indent int) string {
	indentPad, err := indentPadding(option, indent)
	if err != nil {
		indentPad = strings.Repeat("    ", indent)
	}

	var formated strings.Builder
	formated.WriteString(node.TokenLiteral())

	formated.WriteString("(")

	for index, param := range node.Parameters {
		formated.WriteString(param.Value)
		if index < len(node.Parameters)-1 {
			formated.WriteString(", ")
		}
	}

	formated.WriteString(")")

	// Format space between ")" and "{" function body, could be `\n`
	formated.WriteString(" ")

	formated.WriteString("{")

	// Format space between "{" and function body
	formated.WriteString("\n")

	// Each statement already cmoe with a new line \n
	stmtFormated := FormatedAST(node.Body, option, indent+1)

	formated.WriteString(stmtFormated)
	// Format space between function body and "}"
	formated.WriteString(indentPad)
	formated.WriteString("}")

	return formated.String()
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

func checkInjectComment(currentLine int, out io.Writer, option protocol.FormattingOptions, indent int) {
	if curr_comments >= len(comments) {
		return
	}
	slog.Debug(fmt.Sprintf("Try inject comment (%v 'th) %v, at Line %v", curr_comments, comments[curr_comments], currentLine))
	if currentLine > comments[curr_comments].Start.Line {
		pad, err := indentPadding(option, indent)
		if err != nil {
			slog.Debug(fmt.Sprintf("When do indent padding, got %v and have to skip", err.Error()))
		} else {
			fmt.Fprintf(out, "%s", pad)
		}
		fmt.Fprintf(out, "%s\n", comments[curr_comments].Literal)
		curr_comments += 1
	}
}

// Format statement
func FormatedAST(node ast.Node, option protocol.FormattingOptions, indent int) string {
	var formated strings.Builder
	indentPad, err := indentPadding(option, indent)
	if err != nil {
		indentPad = strings.Repeat("    ", indent)
	}
	switch node := node.(type) {
	case *ast.Program:
		comments = node.Comments
		curr_comments = 0
		for _, statement := range node.Statements {
			stmtFormated := FormatedAST(statement, option, indent)
			formated.WriteString(stmtFormated)

			formated.WriteString(";\n")
		}
		// Cover the remain comment, this doesn't do while loop as there can be possible
		// inf loop error, this however can lost comment
		for i := curr_comments; i < len(comments); i++ {
			checkInjectComment(node.GetRange().End.Line+1, &formated, option, indent)
		}

		if curr_comments != len(comments) {
			slog.Error("Can't inject all comments found in program!")
		}

	case *ast.BlockStatement:
		for index, statement := range node.Statements {
			stmtFormated := FormatedAST(statement, option, indent)
			formated.WriteString(stmtFormated)
			formated.WriteString(";")

			if index < len(node.Statements)-1 {
				formated.WriteString("\n")
			}
		}

		// Cover the remain comment, this doesn't do while loop as there can be possible
		// inf loop error, this however can lost comment
		for i := curr_comments; i < len(comments); i++ {
			currentLine := node.GetRange().End.Line + 1
			if currentLine <= comments[curr_comments].Start.Line {
				break
			}
			// We are sure that we can inject new comment base from previous check
			// New line to seperate the comment if it not originally goes along with the line
			formated.WriteString("\n")
			checkInjectComment(currentLine, &formated, option, indent)
		}

	case *ast.LetStatement:
		checkInjectComment(node.GetRange().Start.Line, &formated, option, indent)
		formated.WriteString(indentPad)
		formated.WriteString(node.TokenLiteral())
		formated.WriteString(" ")
		formated.WriteString(node.Name.String())
		formated.WriteString(" = ")
		formated.WriteString(FormatedExpresionAST(node.Value, option, indent))

	case *ast.ReturnStatement:
		checkInjectComment(node.GetRange().Start.Line, &formated, option, indent)
		formated.WriteString(indentPad)
		formated.WriteString(node.TokenLiteral())
		formated.WriteString(" ")
		formated.WriteString(FormatedExpresionAST(node.ReturnValue, option, indent))

	case *ast.ExpressionStatement:
		checkInjectComment(node.GetRange().Start.Line, &formated, option, indent)
		formated.WriteString(indentPad)
		formated.WriteString(FormatedExpresionAST(node.Expression, option, indent))

	default:
		return node.String()
	}
	return formated.String()
}
