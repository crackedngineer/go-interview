package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func GenericError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errorMessages []string

	for _, fieldErr := range errs {
		switch fieldErr.ActualTag() {
		case "required":
			errorMessages = append(errorMessages, fieldErr.Field()+" is required")
		case "email":
			errorMessages = append(errorMessages, fieldErr.Field()+" must be a valid email")
		case "min":
			errorMessages = append(errorMessages, fieldErr.Field()+" must be greater than or equal to "+fieldErr.Param())
		default:
			errorMessages = append(errorMessages, fieldErr.Field()+" is invalid")
		}
	}

	return Response{
		Status: StatusError,
		Error:  "Validation failed: " + strings.Join(errorMessages, ", "),
	}
}
