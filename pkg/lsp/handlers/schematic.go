package handlers

import protocol "github.com/tliron/glsp/protocol_3_16"

type SemanticTokenType uint32

const (
	SemanticTokenTypeNamespace = SemanticTokenType(iota)
	/**
	 * Represents a generic type. Acts as a fallback for types which
	 * can't be mapped to a specific type like class or enum.
	 */
	SemanticTokenTypeType
	SemanticTokenTypeClass
	SemanticTokenTypeEnum
	SemanticTokenTypeInterface
	SemanticTokenTypeStruct
	SemanticTokenTypeTypeParameter
	SemanticTokenTypeParameter
	SemanticTokenTypeVariable
	SemanticTokenTypeProperty
	SemanticTokenTypeEnumMember
	SemanticTokenTypeEvent
	SemanticTokenTypeFunction
	SemanticTokenTypeMethod
	SemanticTokenTypeMacro
	SemanticTokenTypeKeyword
	SemanticTokenTypeModifier
	SemanticTokenTypeComment
	SemanticTokenTypeString
	SemanticTokenTypeNumber
	SemanticTokenTypeRegexp
	SemanticTokenTypeOperator
)

var tokenName = map[SemanticTokenType]protocol.SemanticTokenType{
	SemanticTokenTypeNamespace: protocol.SemanticTokenTypeNamespace,
	/**
	 * Represents a generic type. Acts as a fallback for types which
	 * can't be mapped to a specific type like class or enum.
	 */
	SemanticTokenTypeType:          protocol.SemanticTokenTypeType,
	SemanticTokenTypeClass:         protocol.SemanticTokenTypeClass,
	SemanticTokenTypeEnum:          protocol.SemanticTokenTypeEnum,
	SemanticTokenTypeInterface:     protocol.SemanticTokenTypeInterface,
	SemanticTokenTypeStruct:        protocol.SemanticTokenTypeStruct,
	SemanticTokenTypeTypeParameter: protocol.SemanticTokenTypeTypeParameter,
	SemanticTokenTypeParameter:     protocol.SemanticTokenTypeParameter,
	SemanticTokenTypeVariable:      protocol.SemanticTokenTypeVariable,
	SemanticTokenTypeProperty:      protocol.SemanticTokenTypeProperty,
	SemanticTokenTypeEnumMember:    protocol.SemanticTokenTypeEnumMember,
	SemanticTokenTypeEvent:         protocol.SemanticTokenTypeEvent,
	SemanticTokenTypeFunction:      protocol.SemanticTokenTypeFunction,
	SemanticTokenTypeMethod:        protocol.SemanticTokenTypeMethod,
	SemanticTokenTypeMacro:         protocol.SemanticTokenTypeMacro,
	SemanticTokenTypeKeyword:       protocol.SemanticTokenTypeKeyword,
	SemanticTokenTypeModifier:      protocol.SemanticTokenTypeModifier,
	SemanticTokenTypeComment:       protocol.SemanticTokenTypeComment,
	SemanticTokenTypeString:        protocol.SemanticTokenTypeString,
	SemanticTokenTypeNumber:        protocol.SemanticTokenTypeNumber,
	SemanticTokenTypeRegexp:        protocol.SemanticTokenTypeRegexp,
	SemanticTokenTypeOperator:      protocol.SemanticTokenTypeOperator,
}

func (ss SemanticTokenType) String() protocol.SemanticTokenType {
	return tokenName[ss]
}

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
