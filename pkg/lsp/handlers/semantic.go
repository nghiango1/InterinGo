package handlers

import (
	"interingo/pkg/parser"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

var tokenName = map[parser.SemanticTokenType]protocol.SemanticTokenType{
	parser.SemanticTokenTypeNamespace: protocol.SemanticTokenTypeNamespace,
	/**
	 * Represents a generic type. Acts as a fallback for types which
	 * can't be mapped to a specific type like class or enum.
	 */
	parser.SemanticTokenTypeType:          protocol.SemanticTokenTypeType,
	parser.SemanticTokenTypeClass:         protocol.SemanticTokenTypeClass,
	parser.SemanticTokenTypeEnum:          protocol.SemanticTokenTypeEnum,
	parser.SemanticTokenTypeInterface:     protocol.SemanticTokenTypeInterface,
	parser.SemanticTokenTypeStruct:        protocol.SemanticTokenTypeStruct,
	parser.SemanticTokenTypeTypeParameter: protocol.SemanticTokenTypeTypeParameter,
	parser.SemanticTokenTypeParameter:     protocol.SemanticTokenTypeParameter,
	parser.SemanticTokenTypeVariable:      protocol.SemanticTokenTypeVariable,
	parser.SemanticTokenTypeProperty:      protocol.SemanticTokenTypeProperty,
	parser.SemanticTokenTypeEnumMember:    protocol.SemanticTokenTypeEnumMember,
	parser.SemanticTokenTypeEvent:         protocol.SemanticTokenTypeEvent,
	parser.SemanticTokenTypeFunction:      protocol.SemanticTokenTypeFunction,
	parser.SemanticTokenTypeMethod:        protocol.SemanticTokenTypeMethod,
	parser.SemanticTokenTypeMacro:         protocol.SemanticTokenTypeMacro,
	parser.SemanticTokenTypeKeyword:       protocol.SemanticTokenTypeKeyword,
	parser.SemanticTokenTypeModifier:      protocol.SemanticTokenTypeModifier,
	parser.SemanticTokenTypeComment:       protocol.SemanticTokenTypeComment,
	parser.SemanticTokenTypeString:        protocol.SemanticTokenTypeString,
	parser.SemanticTokenTypeNumber:        protocol.SemanticTokenTypeNumber,
	parser.SemanticTokenTypeRegexp:        protocol.SemanticTokenTypeRegexp,
	parser.SemanticTokenTypeOperator:      protocol.SemanticTokenTypeOperator,
}

// nvimsupported
// tokenModifiers = { "declaration", "definition", "readonly", "static", "deprecated", "abstract", "async", "modification", "documentation", "defaultLibrary" },
// tokenTypes = { "namespace", "type", "class", "enum", "interface", "struct", "typeParameter", "parameter", "variable", "property", "enumMember", "event", "function", "method", "macro", "keyword", "modifier", "comment", "string", "number", "regexp", "operator", "decorator" }
// @lsp.type.class Identifiers that declare or reference a class type
// @lsp.type.comment Tokens that represent a comment
// @lsp.type.decorator Identifiers that declare or reference decorators and annotations
// @lsp.type.enum Identifiers that declare or reference an enumeration type
// @lsp.type.enumMember Identifiers that declare or reference an enumeration property, constant, or member
// @lsp.type.event Identifiers that declare an event property
// @lsp.type.function Identifiers that declare a function
// @lsp.type.interface Identifiers that declare or reference an interface type
// @lsp.type.keyword Tokens that represent a language keyword
// @lsp.type.macro Identifiers that declare a macro
// @lsp.type.method Identifiers that declare a member function or method
// @lsp.type.modifier Tokens that represent a modifier
// @lsp.type.namespace Identifiers that declare or reference a namespace, module, or package
// @lsp.type.number Tokens that represent a number literal
// @lsp.type.operator Tokens that represent an operator
// @lsp.type.parameter Identifiers that declare or reference a function or method parameters
// @lsp.type.property Identifiers that declare or reference a member property, member field, or member variable
// @lsp.type.regexp Tokens that represent a regular expression literal
// @lsp.type.string Tokens that represent a string literal
// @lsp.type.struct Identifiers that declare or reference a struct type
// @lsp.type.type Identifiers that declare or reference a type that is not covered above
// @lsp.type.typeParameter Identifiers that declare or reference a type parameter
// @lsp.type.variable Identifiers that declare or reference a local or global variable
//
// server_legend = {
// tokenTypes = { "namespace", "type", "class", "enum", "interface", "struct", "typeParameter", "parameter", "variable", "property", "enumMember", "event", "function", "method", "macro", "keyword", "modifier", "comment", "string", "number", "regexp", "operator" }
// }
var SupportedSemanticTokenType []string = []string{
	string(protocol.SemanticTokenTypeNamespace),
	string(protocol.SemanticTokenTypeType),
	string(protocol.SemanticTokenTypeClass),
	string(protocol.SemanticTokenTypeEnum),
	string(protocol.SemanticTokenTypeInterface),
	string(protocol.SemanticTokenTypeStruct),
	string(protocol.SemanticTokenTypeTypeParameter),
	string(protocol.SemanticTokenTypeParameter),
	string(protocol.SemanticTokenTypeVariable),
	string(protocol.SemanticTokenTypeProperty),
	string(protocol.SemanticTokenTypeEnumMember),
	string(protocol.SemanticTokenTypeEvent),
	string(protocol.SemanticTokenTypeFunction),
	string(protocol.SemanticTokenTypeMethod),
	string(protocol.SemanticTokenTypeMacro),
	string(protocol.SemanticTokenTypeKeyword),
	string(protocol.SemanticTokenTypeModifier),
	string(protocol.SemanticTokenTypeComment),
	string(protocol.SemanticTokenTypeString),
	string(protocol.SemanticTokenTypeNumber),
	string(protocol.SemanticTokenTypeRegexp),
	string(protocol.SemanticTokenTypeOperator),
}
