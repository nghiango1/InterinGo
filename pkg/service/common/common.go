package common

type EchoRequest struct {
	Data string `json:"data"`
}

type EchoResponse struct {
	Data string `json:"data"`
}

type EvalRequest struct {
	Data string `json:"data"`
}

// Interface can be tricky
// Status 200
type EvalResponseSuccess struct {
	Output string `json:"output,omitempty"`
}
