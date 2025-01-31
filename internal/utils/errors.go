package utils

import "github.com/go-playground/validator/v10"

type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

// ValidationErrorResponse returns a validation error response based on the given error.
func ValidationErrorResponse(err error) ErrorResponse {
	errorsMap := make(map[string]string)

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return ErrorResponse{
			Error:   "validation error",
			Details: map[string]string{"error": err.Error()},
		}
	}
	for _, e := range validationErrs {
		switch e.Tag() {
		case "required":
			errorsMap[e.Field()] = e.Field() + " is required"
		case "ip":
			errorsMap[e.Field()] = "ip_address must be a valid IPv4 or IPv6 address"
		default:
			errorsMap[e.Field()] = e.Field() + " is invalid"
		}
	}

	return ErrorResponse{
		Error:   "Invalid Request",
		Details: errorsMap,
	}
}

// ServerErrorResponse returns a server error response.
func ServerErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error:   "Internal Server Error",
		Details: map[string]string{"error": err.Error()},
	}
}

// BadRequestResponse returns a bad request error response.
func BadRequestResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error:   "Bad Request",
		Details: map[string]string{"error": err.Error()},
	}
}

// GenericErrorResponse returns a generic error response.
func GenericErrorResponse() ErrorResponse {
	return ErrorResponse{
		Error: "Internal Server Error",
	}
}

// NotFoundResponse returns a not found error response.
func NotFoundResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error:   "Not Found",
		Details: map[string]string{"error": err.Error()},
	}
}
