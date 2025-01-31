package utils

import "github.com/go-playground/validator/v10"

type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

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

func ServerErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error:   "Internal Server Error",
		Details: map[string]string{"error": err.Error()},
	}
}

func BadRequestResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error:   "Bad Request",
		Details: map[string]string{"error": err.Error()},
	}
}
func GenericErrorResponse() ErrorResponse {
	return ErrorResponse{
		Error: "Internal Server Error",
	}
}
func NotFoundResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error:   "Not Found",
		Details: map[string]string{"error": err.Error()},
	}
}
