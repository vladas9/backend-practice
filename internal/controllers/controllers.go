// Package controllers provides HTTP handlers for user authentication operations.
// It includes handlers for user sign-in and sign-up processes.
// Each handler reads the request body, decodes it into a Model, and interacts with a service to process the request.
// Based on the service's response, the handlers return appropriate HTTP responses.
package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	s "github.com/vladas9/backend-practice/internal/services"
	u "github.com/vladas9/backend-practice/internal/utils"
	"net/http"
	"reflect"
)

type Controller struct {
	userService *s.UserService
}

func NewController(db *sql.DB) (*Controller, error) {

	controller := &Controller{}

	if service, err := s.NewUserService(db); err != nil {
		return nil, err
	} else {
		controller.userService = service
	}

	return controller, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) *ApiError {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		aErr := &ApiError{fmt.Sprintf("Encoding of object of type %v failed", reflect.TypeOf(v)), 500}
		u.Logger.Error(aErr)
		return aErr
	}
	return nil
}

type ApiError struct {
	ErrorMsg string `json:"error"`
	Status   int
}

func (e ApiError) Error() string {
	return e.ErrorMsg
}
