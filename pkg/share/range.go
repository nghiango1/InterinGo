package share

import protocol "github.com/tliron/glsp/protocol_3_16"

type Position struct {
	/**
	 * Line position in a document (zero-based).
	 */
	Line int

	/**
	 * Character offset on a line in a document (zero-based). The meaning of this
	 * offset is determined by the negotiated `PositionEncodingKind`.
	 *
	 * If the character value is greater than the line length it defaults back
	 * to the line length.
	 */
	Character int
}

type Range struct {
	Start Position
	End   Position
}

// Support glsp
func (r *Range) ToProtocolRange() protocol.Range {
	return protocol.Range{
		Start: protocol.Position{
			Line:      protocol.UInteger(r.Start.Line),
			Character: protocol.UInteger(r.Start.Character),
		},
		End: protocol.Position{
			Line:      protocol.UInteger(r.End.Line),
			Character: protocol.UInteger(r.End.Character),
		},
	}
}
