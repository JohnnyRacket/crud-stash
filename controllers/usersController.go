package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

type UsersController struct {
}

// NewUsersController creates a new instance of a user controller
func NewUsersController() *UsersController {
	return &UsersController{}
}

// RegisterRoutes implements the controller interface
func (c *UsersController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users", c.GetUsers).Methods("GET")
}

// GetUsers returns all users
func (c *UsersController) GetUsers(w http.ResponseWriter, r *http.Request) {

}
