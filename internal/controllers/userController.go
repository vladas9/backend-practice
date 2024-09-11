package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vladas9/backend-practice/internal/models"
)

// Login handles the HTTP request for user sign-in.
// It reads and decodes the request body into a UserModel instance.
// It then calls a service with the user information and returns an appropriate response based on the service result.
// If decoding fails, it returns a 400 Bad Request status with an error message.
// If successful, it returns a 200 OK status with a success message.
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) *ApiError {

	var user models.UserModel

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return &ApiError{fmt.Sprintf("SignIn failed: %s", err), http.StatusBadRequest}
	}

	// TODO: Call the service to handle sign-in with the user data

	return writeJSON(w, http.StatusOK, "Sign-in successful")
}

// Register	handles the HTTP request for user sign-up.
// It reads and decodes the request body into a UserModel instance.
// It then calls a service with the user information and returns an appropriate response based on the service result.
// If decoding fails, it returns a 400 Bad Request status with an error message.
// If successful, it returns a 200 OK status with a success message.
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) *ApiError {

	var user models.UserModel

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return &ApiError{fmt.Sprintf("SignUp failed: %s", err), http.StatusBadRequest}
	}

	// TODO: Call the service to handle sign-up with the user data

	return writeJSON(w, http.StatusOK, "Sign-up successful")
}