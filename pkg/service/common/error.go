package common

// Interface can be tricky
// Status 400 -> Message, Maybe have Param
// Status 500 -> Message is enough
type ErrorResponseInterface interface {
	GetType() int
	GetCode() string
	GetMessage() string
}

type ErrorResponseImpl struct {
	Type    int `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorResponseImpl) GetType() int    { return e.Type }
func (e *ErrorResponseImpl) GetCode() string    { return e.Code }
func (e *ErrorResponseImpl) GetMessage() string { return e.Message }

// This cover only for BadRequestErrorResponse
// Use ths if we know it possible to have Param as failed to validate field
type BadRequestErrorResponse struct {
	Type    int  `json:"type"`
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Param   *string `json:"param,omitempty"`
}

func (e *BadRequestErrorResponse) GetType() int    { return e.Type }
func (e *BadRequestErrorResponse) GetCode() string    { return e.Code }
func (e *BadRequestErrorResponse) GetMessage() string { return e.Message }

// 404
func NewNotFoundErrorResponse(message string) ErrorResponseInterface {
	if message == "" {
		return NewErrorResponse(404)
	}

	return &ErrorResponseImpl{
		Type:    404,
		Code:    "resource_not_found",
		Message: message,
	}
}

// 400
func NewBadRequestErrorResponse(message string, params *string) ErrorResponseInterface {
	if message == "" {
		if params == nil {
			return NewErrorResponse(400)
		}
		message = "The request was invalid"
	}

	return &BadRequestErrorResponse{
		Type:    400,
		Code:    "invalid_parameter",
		Message: message,
		Param:   params,
	}
}

// Use HTTP code and received default
func NewErrorResponse(status int) ErrorResponseInterface {
	switch status {
	case 400:
		return &ErrorResponseImpl{
			Type:    400,
			Code:    "invalid_parameter",
			Message: "The request was invalid",
			// Param:   "email",
		}

	case 401:
		return &ErrorResponseImpl{
			Type:    401,
			Code:    "unauthorized",
			Message: "Authentication required",
		}

	case 404:
		return &ErrorResponseImpl{
			Type:    404,
			Code:    "resource_not_found",
			Message: "The requested resource could not be found",
		}

	case 429:
		return &ErrorResponseImpl{
			Type:    429,
			Code:    "too_many_requests",
			Message: "Rate limit exceeded",
		}

	case 500:
		return &ErrorResponseImpl{
			Type:    500,
			Code:    "internal_error",
			Message: "Internal server error",
		}

	default:
		return &ErrorResponseImpl{
			Type:    500,
			Code:    "internal_error",
			Message: "Unknow error",
		}
	}
}
