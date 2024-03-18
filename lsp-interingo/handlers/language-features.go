package handlers

import (
	"errors"
	"fmt"
	"interingo-lsp/mappers"
	"interingo-lsp/store"
	"interingo/lexer"
	"interingo/parser"

	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func HandleTextDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	ef := store.Wrap(&params.TextDocument)
	store.GetStore().Add(ef)
	return nil
}

func HandleTextDocumentDidChange (context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	uri := params.TextDocument.URI
	textDocObj, err := store.GetStore().Get(uri)
	if err != nil {
		return err
	}

	textDocObj.Unwrap().Version = params.TextDocument.Version
	contentChanges := params.ContentChanges // TextDocumentContentChangeEvent or TextDocumentContentChangeEventWhole

	for index, contextChange := range contentChanges {
		switch changeType := contextChange.(type) {
		case protocol.TextDocumentContentChangeEventWhole:
			textDocObj.UpdateWhole(changeType)
		case protocol.TextDocumentContentChangeEvent:
			textDocObj.Update(changeType)
		default:
			return errors.New(fmt.Sprintf("ABORT: Can't following %d'th file change, get %v", index, contextChange))
		}
	}

	return nil
}

func HandleDocumentFormatting(context *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	formated := make([]protocol.TextEdit, 0)

	uri := params.TextDocument.URI

	ef, err := store.GetStore().Get(uri)
	if err != nil {
		return nil, err
	}

	// Not format yet
	format := ef.Unwrap().Text

	l := lexer.New(ef.Unwrap().Text)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) == 0 {
		format = FormatedAST(program, params.Options, 0)
	} else {
		return nil, errors.New(p.Errors()[0])
	}

	editAllFile := protocol.TextEdit{
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
		NewText: format,
	}

	formated = append(formated, editAllFile)
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
