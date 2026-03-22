package main

import (
	"interingo/pkg/lsp/handlers"
	"interingo/pkg/lsp/store"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	_ "github.com/tliron/commonlog/simple"
)

const lsName = "Interingo Language Server"

var version string = "0.0.1"
var handler protocol.Handler

func Init() {
	store.Init()
}

func main() {
	commonlog.Configure(2, nil)

	handler = protocol.Handler{
		Initialize:             initialize,
		Shutdown:               shutdown,
		TextDocumentSemanticTokensFull: handlers.HandleTextDocumentSemanticTokensFull,
		TextDocumentCompletion: handlers.TextDocumentCompletion,
		TextDocumentFormatting: handlers.HandleDocumentFormatting,
		TextDocumentDidOpen:    handlers.HandleTextDocumentDidOpen,
		TextDocumentDidChange:  handlers.HandleTextDocumentDidChange,
		TextDocumentDocumentSymbol:  handlers.HandleTextDocumentDocumentSymbol,
		TextDocumentDocumentHighlight:  handlers.HandleTextDocumentDocumentHighlight,
	}

	server := server.NewServer(&handler, lsName, true)

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	commonlog.NewInfoMessage(0, "Initializing server...")

	capabilities := handler.CreateServerCapabilities()

	capabilities.CompletionProvider = &protocol.CompletionOptions{}
	capabilities.DocumentHighlightProvider = &protocol.DocumentHighlightOptions{}
	capabilities.DocumentSymbolProvider = &protocol.DocumentSymbolOptions{}
	capabilities.DocumentFormattingProvider = &protocol.DocumentFormattingOptions{}
	capabilities.SemanticTokensProvider = &protocol.SemanticTokensOptions{
		Legend: protocol.SemanticTokensLegend{
			TokenTypes: handlers.SupportedSemanticTokenType,
			TokenModifiers: nil,
		}, 
		Full: true,
	}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func shutdown(context *glsp.Context) error {
	return nil
}
