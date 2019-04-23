package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocql/gocql"

	"github.com/gorilla/mux"
)

type CrudController struct {
	session *gocql.Session
}

func NewCrudController(session *gocql.Session) *CrudController {

	return &CrudController{session: session}
}

func (c *CrudController) RegisterRoutes(router *mux.Router) {
	// individual entities
	router.HandleFunc("/{app}/{entity}", c.CreateEntity).Methods("POST")
	router.HandleFunc("/{app}/{entity}/{id}", c.DeleteEntity).Methods("DELETE")
	router.HandleFunc("/{app}/{entity}/{id}", c.GetEntity).Methods("GET")
	router.HandleFunc("/{app}/{entity}", c.GetEntities).Methods("GET")
	router.HandleFunc("/{app}/{entity}/{id}", c.UpdateEntity).Methods("PUT")
}

func (c *CrudController) GetEntities(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	entity := vars["entity"]

	//TODO errors around getting app and entity
	query := c.session.Query("SELECT * FROM " + app + "." + entity + ";")
	iter := query.Iter()
	res := []map[string]interface{}{}
	m := &map[string]interface{}{}
	//scan items into a map
	for iter.MapScan(*m) {
		//fmt.Println(m)
		res = append(res, *m)
		m = &map[string]interface{}{}
	}
	// marshall result into json fomat
	jsonRes, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

func (c *CrudController) CreateEntity(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	app := vars["app"]
	entity := vars["entity"]

	// fetch metadata about table to know how to input the data; i.e. if we need quotes or not
	meta := "select column_name,type from system_schema.columns where keyspace_name ='" + app + "' AND table_name='" + entity + "';"
	metaQuery := c.session.Query(meta)
	iter := metaQuery.Iter()
	var columnName string
	var columnType string
	metadata := map[string]string{}
	//scan items into a map
	for iter.Scan(&columnName, &columnType) {
		metadata[columnName] = columnType
	}
	fmt.Println(metadata)

	decoder := json.NewDecoder(r.Body)
	input := map[string]interface{}{}
	err := decoder.Decode(&input)
	if err != nil {
		panic(err)
	}
	// log.Println(input)
	// get the keys and values from the body
	keys := make([]string, 0, len(input))
	values := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
		if metadata[key] != "text" {
			values = append(values, input[key].(string))
		} else {
			values = append(values, "'"+input[key].(string)+"'")
		}
	}
	// build a query string using them
	query := "INSERT INTO example.tweet (" + strings.Join(keys, ", ") + ") VALUES (" + strings.Join(values, ", ") + ");"
	fmt.Println(query)

	err = c.session.Query(query).Exec()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *CrudController) DeleteEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	entity := vars["entity"]
	id := vars["id"]

	query := "DELETE FROM " + app + "." + entity + " WHERE id = " + id + ";"
	fmt.Println(query)

	err := c.session.Query(query).Exec()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *CrudController) UpdateEntity(w http.ResponseWriter, r *http.Request) {}
func (c *CrudController) GetEntity(w http.ResponseWriter, r *http.Request)    {}
