package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocql/gocql"

	"github.com/gorilla/mux"
)

type EntitiesController struct {
	session gocql.Session
}

func NewEntitiesController() *EntitiesController {
	return &EntitiesController{}
}

func (c *EntitiesController) RegisterRoutes(router *mux.Router) {
	//system entity creation/deletion
	router.HandleFunc("/{app}/entities", c.GetEntities).Methods("GET")
	router.HandleFunc("/{app}/entities", c.CreateEntity).Methods("POST")
	router.HandleFunc("/{app}/entities/{id}", c.CreateEntity).Methods("PUT")
	router.HandleFunc("/{app}/entities/{id}", c.CreateEntity).Methods("DELETE")
}

func (c *EntitiesController) GetEntities(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	entity := vars["entity"]

	var res map[string]interface{}
	//TODO errors around getting app and entity
	if err := c.session.Query(`SELECT * FROM %s.%s`, app, entity).Consistency(gocql.One).MapScan(res); err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func (c *EntitiesController) CreateEntity(w http.ResponseWriter, r *http.Request) {
	err := c.session.Query("CREATE TABLE IF NOT EXISTS example.tweet(timeline text, id UUID, text text, PRIMARY KEY(id));").Exec()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *EntitiesController) DeleteEntity(w http.ResponseWriter, r *http.Request) {}

func (c *EntitiesController) UpdateEntity(w http.ResponseWriter, r *http.Request) {}
func (c *EntitiesController) GetEntity(w http.ResponseWriter, r *http.Request)    {}
