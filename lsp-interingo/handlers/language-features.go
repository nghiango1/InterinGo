package handlers

import (
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"interingo-lsp/mappers"
)

func TextDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	var completionItems []protocol.CompletionItem

	kindConstant := protocol.CompletionItemKindConstant
	for word, constant := range mappers.ConstraintMapper {
		constantCopy := constant // Create a copy of emoji
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      word,
			Detail:     &constantCopy,
			InsertText: &constantCopy,
			Kind:       &kindConstant,
		})
	}

	kindKeyword := protocol.CompletionItemKindKeyword
	for word, keyword := range mappers.KeywordMapper {
		keywordCopy := keyword // Create a copy of emoji
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      word,
			Detail:     &keywordCopy,
			InsertText: &keywordCopy,
			Kind:       &kindKeyword,
		})
	}
	return completionItems, nil
}
