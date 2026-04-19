package common

// Interface can be tricky
// Status 400 -> Message, Maybe have Param
// Status 500 -> Message is enough
type ErrorResponse interface {
	GetType() int
	GetCode() string
	GetMessage() string
}

type errorResponseImpl struct {
	Type    int    `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *errorResponseImpl) GetType() int       { return e.Type }
func (e *errorResponseImpl) GetCode() string    { return e.Code }
func (e *errorResponseImpl) GetMessage() string { return e.Message }

// 404
func NewNotFoundErrorResponse(message string) ErrorResponse {
	if message == "" {
		return NewErrorResponse(404)
	}

	return &errorResponseImpl{
		Type:    404,
		Code:    "resource_not_found",
		Message: message,
	}
}

// Use HTTP code and received default
func NewErrorResponse(status int) ErrorResponse {
	switch status {
	case 400:
		return &errorResponseImpl{
			Type:    400,
			Code:    "invalid_parameter",
			Message: "The request was invalid",
			// Param:   "email",
		}

	case 401:
		return &errorResponseImpl{
			Type:    401,
			Code:    "unauthorized",
			Message: "Authentication required",
		}

	case 404:
		return &errorResponseImpl{
			Type:    404,
			Code:    "resource_not_found",
			Message: "The requested resource could not be found",
		}

	case 429:
		return &errorResponseImpl{
			Type:    429,
			Code:    "too_many_requests",
			Message: "Rate limit exceeded",
		}

	case 500:
		return &errorResponseImpl{
			Type:    500,
			Code:    "internal_error",
			Message: "Internal server error",
		}

	default:
		return &errorResponseImpl{
			Type:    500,
			Code:    "internal_error",
			Message: "Unknow error",
		}
	}
}
