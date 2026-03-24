package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"interingo/pkg/lsp/mappers"
	"interingo/pkg/lsp/store"
	"interingo/pkg/parser"

	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func handleTextDocumentSemanticTokensFull(uri protocol.DocumentUri) (*protocol.SemanticTokens, error) {
	data := []uint32{}

	ef, err := store.GetStore().Get(uri)
	if err != nil {
		return nil, err
	}

	var prevLine uint32 = 0
	var prevChar uint32 = 0
	for _, v := range ef.Parser.DocumentTokens {
		var deltaLine uint32 = uint32(v.Unwrap().Start.Line) - prevLine
		// New line will start with a absolute character position
		var deltaChar uint32 = uint32(v.Unwrap().Start.Character)

		// Same line will use delta character position
		if deltaLine == 0 {
			deltaChar -= prevChar
		}
		prevLine = uint32(v.Unwrap().Start.Line)
		prevChar = uint32(v.Unwrap().Start.Character)
		var length uint32 = uint32(len(v.Unwrap().Literal))
		var tokenType parser.SemanticTokenType = v.Kind
		var tokenModifiers uint32 = 0

		data = append(data, deltaLine,
			deltaChar,
			length,
			uint32(tokenType),
			tokenModifiers,
		)
	}

	return &protocol.SemanticTokens{
		ResultID: nil,
		Data:     data,
	}, nil

}

// Returns: SemanticTokens | SemanticTokensDelta | SemanticTokensDeltaPartialResult | nil
func HandleTextDocumentSemanticTokensFullDelta(context *glsp.Context, params *protocol.SemanticTokensDeltaParams) (any, error) {
	uri := params.TextDocument.URI
	return handleTextDocumentSemanticTokensFull(uri)
}

func HandleTextDocumentSemanticTokensFull(context *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	uri := params.TextDocument.URI
	return handleTextDocumentSemanticTokensFull(uri)
}

func getDiagnostic(errs []parser.ParserError) []protocol.Diagnostic {
	// Create diagnostic
	var diagnostics []protocol.Diagnostic

	for _, e := range errs {
		serverity := protocol.DiagnosticSeverityError
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range:              e.Range.ToProtocolRange(),
			Severity:           &serverity,
			Code:               nil,
			CodeDescription:    nil,
			Source:             nil,
			Message:            e.Message,
			Tags:               nil,
			RelatedInformation: nil,
			Data:               nil,
		})
	}
	return diagnostics
}

func HandleTextDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	ef := store.Wrap(&params.TextDocument)
	store.GetStore().Add(ef)
	found := getDiagnostic(ef.Parser.Errors)
	if len(found) != 0 {
		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
			URI:         params.TextDocument.URI,
			Diagnostics: found,
		})
	}
	return nil
}

func HandleTextDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
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
			return fmt.Errorf("ABORT: Can't following %d'th file change, get %v", index, contextChange)
		}
	}

	found := getDiagnostic(textDocObj.Parser.Errors)
	if len(found) != 0 {
		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
			URI:         params.TextDocument.URI,
			Diagnostics: found,
		})
	}
	return nil
}

func HandleDocumentFormatting(context *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	var formated []protocol.TextEdit

	uri := params.TextDocument.URI

	ef, err := store.GetStore().Get(uri)
	if err != nil {
		return nil, err
	}

	// Not format yet
	format := ef.Unwrap().Text


	var fo FormattingOptions
	d, err := json.Marshal(params.Options)
	if err != nil {
		return nil, err
	} 
	err = json.Unmarshal(d, &fo)

	if err != nil {
		return nil, err
	}

	if len(ef.Parser.Errors) == 0 {
		format = FormatedAST(ef.Parser.Program, fo, 0)
	} else {
		return nil, errors.New(ef.Parser.Errors[0].Message)
	}

	editAllFile := protocol.TextEdit{
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      protocol.UInteger(0),
				Character: protocol.UInteger(0),
			},
			End: protocol.Position{
				Line:      protocol.UInteger(ef.Parser.Lexer.Line),
				Character: protocol.UInteger(ef.Parser.Lexer.Character),
			},
		},
		NewText: format,
	}

	formated = append(formated, editAllFile)
	return formated, nil
}

func TextDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
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
