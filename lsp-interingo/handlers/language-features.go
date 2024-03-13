package handlers

import (
	"errors"
	"fmt"
	"interingo-lsp/mappers"
	"interingo/lexer"
	"interingo/parser"
	"net/url"
	"os"

	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func HandleDocumentFormatting(context *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	uri := params.TextDocument.URI
	filepath, err := url.Parse(uri)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Can't parse file uri provided", uri))
	}

	data, err := os.ReadFile(filepath.Path)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Can't read file at", filepath))
	}

	stringData := string(data)

	l := lexer.New(stringData)
	p := parser.New(l)
	program := p.ParseProgram()

	formated := make([]protocol.TextEdit, 0)
	allfile := protocol.TextEdit{
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      protocol.UInteger(0),
				Character: protocol.UInteger(0),
			},
			End: protocol.Position{
				Line:      protocol.UInteger(l.Line),
				Character: protocol.UInteger(l.Character),
			},
		},
		NewText: fmt.Sprintf("// Yah want formatting? My job here is done\n// %s\n", program.String()) + stringData,
	}

	formated = append(formated, allfile)
	return formated, nil
}

func TextDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	var completionItems []protocol.CompletionItem

	kindConstant := protocol.CompletionItemKindConstant
	for word, constant := range mappers.ConstraintMapper {
		constantCopy := constant
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      word,
			Detail:     &constantCopy,
			InsertText: &constantCopy,
			Kind:       &kindConstant,
		})
	}

	kindKeyword := protocol.CompletionItemKindKeyword
	for word, keyword := range mappers.KeywordMapper {
		keywordCopy := keyword
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      word,
			Detail:     &keywordCopy,
			InsertText: &keywordCopy,
			Kind:       &kindKeyword,
		})
	}
	return completionItems, nil
}
