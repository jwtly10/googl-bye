package utils_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jwtly10/googl-bye/internal/errors"
	"github.com/jwtly10/googl-bye/internal/utils"
)

func TestWriteErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		errResponse    utils.ErrorResponse
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "bad request",
			statusCode:     http.StatusBadRequest,
			errResponse:    utils.ErrorResponse{Error: "BAD_REQUEST", Message: "invalid input"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"BAD_REQUEST","message":"Invalid input"}`,
		},
		{
			name:           "internal server error",
			statusCode:     http.StatusInternalServerError,
			errResponse:    utils.ErrorResponse{Error: "INTERNAL_SERVER_ERROR", Message: "failure processing request"},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"INTERNAL_SERVER_ERROR","message":"Failure processing request"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			utils.WriteErrorResponse(recorder, tt.statusCode, tt.errResponse)

			result := recorder.Result()
			defer result.Body.Close()

			if result.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, result.StatusCode)
			}

			var body utils.ErrorResponse
			if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
				t.Fatal("could not decode response body")
			}

			if body.Error != tt.errResponse.Error || strings.ToLower(body.Message) != strings.ToLower(tt.errResponse.Message) {
				t.Errorf("expected body %v, got %v", tt.expectedBody, body)
			}
		})
	}
}

func TestHandleBadRequest(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := fmt.Errorf("bad request error")
	utils.HandleBadRequest(recorder, err)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, result.StatusCode)
	}

	var body utils.ErrorResponse
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatal("could not decode response body")
	}

	expectedError := "BAD_REQUEST_ERROR"
	if body.Error != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, body.Error)
	}

	expectedMessage := "Bad request error"
	if strings.ToLower(body.Message) != strings.ToLower(err.Error()) {
		t.Errorf("expected message '%s', got '%s'", expectedMessage, body.Message)
	}
}

// Add similar tests for HandleInternalError, HandleValidationError, and HandleCustomErrors

func TestHandleInternalError(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := fmt.Errorf("internal server failure")
	utils.HandleInternalError(recorder, err)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, result.StatusCode)
	}

	var body utils.ErrorResponse
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatal("could not decode response body")
	}

	expectedError := "INTERNAL_SERVER_ERROR"
	if body.Error != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, body.Error)
	}

	expectedMessage := "Internal server failure"
	if strings.ToLower(body.Message) != strings.ToLower(err.Error()) {
		t.Errorf("expected message '%s', got '%s'", expectedMessage, body.Message)
	}
}

func TestHandleValidationError(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := fmt.Errorf("validation failed")
	utils.HandleValidationError(recorder, err)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, result.StatusCode)
	}

	var body utils.ErrorResponse
	if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
		t.Fatal("could not decode response body")
	}

	expectedError := "BAD_REQUEST_ERROR"
	if body.Error != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, body.Error)
	}

	expectedMessage := "Validation failed"
	if strings.ToLower(body.Message) != strings.ToLower(err.Error()) {
		t.Errorf("expected message '%s', got '%s'", expectedMessage, body.Message)
	}
}

func TestHandleCustomErrors(t *testing.T) {
	tests := []struct {
		name            string
		inputError      error
		expectedCode    int
		expectedError   string
		expectedMessage string
	}{
		{
			name:            "Not found error",
			inputError:      &errors.NotFoundError{Message: "not found"},
			expectedCode:    http.StatusNotFound,
			expectedError:   "NOT_FOUND",
			expectedMessage: "not found",
		},
		{
			name:            "Validation error",
			inputError:      &errors.BadRequestError{Message: "validation error"},
			expectedCode:    http.StatusBadRequest,
			expectedError:   "BAD_REQUEST_ERROR",
			expectedMessage: "validation error",
		},
		{
			name:            "Internal error",
			inputError:      &errors.InternalError{Message: "internal error"},
			expectedCode:    http.StatusInternalServerError,
			expectedError:   "INTERNAL_SERVER_ERROR",
			expectedMessage: "internal error",
		},
		{
			name:            "Unknown error",
			inputError:      fmt.Errorf("unknown error"),
			expectedCode:    http.StatusInternalServerError,
			expectedError:   "UNKNOWN_ERROR",
			expectedMessage: "unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			utils.HandleCustomErrors(recorder, tt.inputError)

			result := recorder.Result()
			defer result.Body.Close()

			if result.StatusCode != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, result.StatusCode)
			}

			var body utils.ErrorResponse
			if err := json.NewDecoder(result.Body).Decode(&body); err != nil {
				t.Fatal("could not decode response body")
			}

			if body.Error != tt.expectedError {
				t.Errorf("expected error '%s', got '%s'", tt.expectedError, body.Error)
			}

			if strings.ToLower(body.Message) != strings.ToLower(tt.expectedMessage) {
				t.Errorf("expected message '%s', got '%s'", tt.expectedMessage, body.Message)
			}
		})
	}
}
