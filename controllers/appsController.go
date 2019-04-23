package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

//AppsController
type AppsController struct {
}

//NewAppsController
func NewAppsController() *AppsController {
	return &AppsController{}
}

//RegisterRoutes creates the routes for the application that will be available via web api
func (c *AppsController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/apps", c.GetApps).Methods("GET")
	router.HandleFunc("/apps/{id}", c.GetApp).Methods("GET")
	router.HandleFunc("/apps", c.CreateApp).Methods("POST")
	router.HandleFunc("/apps/{id}", c.UpdateApp).Methods("PUT")
	router.HandleFunc("/apps/{id}", c.DeleteApp).Methods("DELETE")
}

//GetApps
func (c *AppsController) GetApps(w http.ResponseWriter, r *http.Request) {
	//json.NewEncoder(w).Encode() //fetch shit here)
	w.Write([]byte("Hello, World!!!\n"))
}

//GetApp
func (c *AppsController) GetApp(w http.ResponseWriter, r *http.Request) {}

//CreateApp
func (c *AppsController) CreateApp(w http.ResponseWriter, r *http.Request) {
}

//DeleteApp
func (c *AppsController) DeleteApp(w http.ResponseWriter, r *http.Request) {
}

//UpdateApp
func (c *AppsController) UpdateApp(w http.ResponseWriter, r *http.Request) {
}
