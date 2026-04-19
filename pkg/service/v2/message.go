// Message type for API v1
package v2

type EvaluateRequest struct {
	Data string `json:"data"`
}

type EvaluateResponseSuccess struct {
	Output *string `json:"output,omitempty"`
}

type EvaluateResponseError struct {
	Type    int    `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type EvaluateResponseParserError struct {
	Type    int           `json:"type"`
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Errors  []ParserError `json:"error,omitempty"`
}

type EvaluateResponseEvalError struct {
	Type    int    `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ParserError struct {
	Message string `json:"message"`
}

type CreateReplRuntimeRequest struct{}

type CreateReplRuntimeResponseSuccess struct {
	RuntimeId string `json:"runtimeId"`
}
