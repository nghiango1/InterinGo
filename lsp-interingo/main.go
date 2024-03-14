package main

import (
	"interingo-lsp/handlers"
	"interingo-lsp/store"

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
		TextDocumentCompletion: handlers.TextDocumentCompletion,
		TextDocumentFormatting: handlers.HandleDocumentFormatting,
		TextDocumentDidOpen:    handlers.HandleTextDocumentDidOpen,
		TextDocumentDidChange:  handlers.HandleTextDocumentDidChange,
	}

	server := server.NewServer(&handler, lsName, true)

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	commonlog.NewInfoMessage(0, "Initializing server...")

	capabilities := handler.CreateServerCapabilities()

	capabilities.CompletionProvider = &protocol.CompletionOptions{}

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
