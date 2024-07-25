package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jwtly10/googl-bye/internal/errors"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, errResponse ErrorResponse) {
	errResponse.Message = capitalizeFirstLetter(errResponse.Message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		fmt.Printf("THIS SHOULD NOT HAPPEN. FAILED TO PARSE ERROR RESPONSE JSON: %v", err)
		http.Error(w, `{"error": "INTERNAL_SERVER_ERROR", "message": "THIS SHOULD NOT HAPPEN. WE WERE UNABLE TO PARSE THE REAL ERROR MSG INTO JSON"}`, http.StatusInternalServerError)
	}
}
func HandleBadRequest(w http.ResponseWriter, err error) {
	statusCode := http.StatusBadRequest
	errorResponse := ErrorResponse{Error: "BAD_REQUEST_ERROR", Message: err.Error()}
	WriteErrorResponse(w, statusCode, errorResponse)
}

func HandleInternalError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	errorResponse := ErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: err.Error()}
	WriteErrorResponse(w, statusCode, errorResponse)
}

func HandleValidationError(w http.ResponseWriter, err error) {
	statusCode := http.StatusBadRequest
	errorResponse := ErrorResponse{Error: "BAD_REQUEST_ERROR", Message: err.Error()}
	WriteErrorResponse(w, statusCode, errorResponse)
}

func HandleCustomErrors(w http.ResponseWriter, err error) {
	var statusCode int
	var errorResponse ErrorResponse

	switch e := err.(type) {
	case *errors.NotFoundError:
		statusCode = http.StatusNotFound
		errorResponse = ErrorResponse{Error: "NOT_FOUND", Message: e.Error()}
	case *errors.BadRequestError:
		statusCode = http.StatusBadRequest
		errorResponse = ErrorResponse{Error: "BAD_REQUEST_ERROR", Message: e.Error()}
	case *errors.InternalError:
		statusCode = http.StatusInternalServerError
		errorResponse = ErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: e.Error()}
	default:
		statusCode = http.StatusInternalServerError
		errorResponse = ErrorResponse{Error: "UNKNOWN_ERROR", Message: e.Error()}
	}

	WriteErrorResponse(w, statusCode, errorResponse)
}

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
