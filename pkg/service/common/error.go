package common

// Interface can be tricky
// Status 400 -> Message, Maybe have Param
// Status 500 -> Message is enough
type ErrorResponseInterface interface {
	getErrorResponse() *ErrorResponse
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e *ErrorResponse) getErrorResponse() *ErrorResponse { return e }

// This cover only for BadRequestErrorResponse
// Use ths if we know it possible to have Param as failed to validate field
type BadRequestErrorResponse struct {
	Code    string  `json:"code"`
	Type    string  `json:"type"`
	Message string  `json:"message"`
	Param   *string `json:"param,omitempty"`
}

func (e *BadRequestErrorResponse) getErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Code:    e.Code,
		Type:    e.Type,
		Message: e.Message,
	}
}

const NOT_FOUND_DEFAULT string = "The requested resource could not be found"

// 404
func NewNotFoundErrorResponse(message string) ErrorResponseInterface {
	if message == "" {
		return NewErrorResponse(404)
	}

	return &ErrorResponse{
		Type:    "not_found",
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
		Type:    "validation",
		Code:    "invalid_parameter",
		Message: message,
		Param:   params,
	}
}

// Use HTTP code and received default
func NewErrorResponse(status int) ErrorResponseInterface {
	switch status {
	case 400:
		return &ErrorResponse{
			Type:    "validation",
			Code:    "invalid_parameter",
			Message: "The request was invalid",
			// Param:   "email",
		}

	case 401:
		return &ErrorResponse{
			Type:    "authentication",
			Code:    "unauthorized",
			Message: "Authentication required",
		}

	case 404:
		return &ErrorResponse{
			Type:    "not_found",
			Code:    "resource_not_found",
			Message: "The requested resource could not be found",
		}

	case 429:
		return &ErrorResponse{
			Type:    "rate_limit",
			Code:    "too_many_requests",
			Message: "Rate limit exceeded",
		}

	case 500:
		return &ErrorResponse{
			Type:    "internal",
			Code:    "internal_error",
			Message: "Internal server error",
		}

	default:
		return &ErrorResponse{
			Type:    "internal",
			Code:    "internal_error",
			Message: "Unknow error",
		}
	}
}
