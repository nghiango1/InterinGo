package common

// This cover only for invalidParamsErrorResponse
// Use ths if we know it possible to have Param as failed to validate field
type invalidParamsErrorResponse struct {
	Message string  `json:"message"`
	Param   *string `json:"param,omitempty"`
}

func (e *invalidParamsErrorResponse) GetType() int       { return 400 }
func (e *invalidParamsErrorResponse) GetCode() string    { return "invalid_parameter" }
func (e *invalidParamsErrorResponse) GetMessage() string { return e.Message }

// 400
func NewInvalidParamsErrorResponse(message string, params *string) ErrorResponse {
	if message == "" {
		if params == nil {
			return NewErrorResponse(400)
		}
		message = "The request was invalid"
	}

	return &invalidParamsErrorResponse{
		Message: message,
		Param:   params,
	}
}
