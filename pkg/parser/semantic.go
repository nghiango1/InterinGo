package parser

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
