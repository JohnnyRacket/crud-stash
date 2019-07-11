package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	// if there is a token convert it to a byte array
	token, err := base64.URLEncoding.DecodeString(r.URL.Query().Get("token"))
	if err != nil {
		w.Write([]byte("Error with paging token format."))
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// if they pass in a custom size convert it to an int
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		w.Write([]byte("Error with page size format."))
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	//TODO errors around getting app and entity
	query := c.session.Query("SELECT * FROM " + app + "." + entity + ";")
	paginated := query.PageState(token).PageSize(size)
	iter := paginated.Iter()
	res := []map[string]interface{}{}
	row := &map[string]interface{}{}

	//scan items into a map
	for iter.MapScan(*row) {
		res = append(res, *row)
		row = &map[string]interface{}{}
	}
	//make sure iter closes without error
	if err := iter.Close(); err != nil {
		log.Println(err)
	}

	//get page state to return in header
	retPagingToken := base64.URLEncoding.EncodeToString(iter.PageState())

	// marshall result into json fomat
	jsonRes, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Paging-Token", retPagingToken)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

func (c *CrudController) CreateEntity(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	app := vars["app"]
	entity := vars["entity"]

	// fetch metadata about table to know how to input the data; i.e. if we need quotes or not
	metadata := c.getMetadata(app, entity)

	decoder := json.NewDecoder(r.Body)
	input := map[string]interface{}{}
	err := decoder.Decode(&input)
	if err != nil {
		http.Error(w, "Error Decoding Body.", 500)
	}

	// get the keys and values from the body
	keys, values := c.formatForCQL(input, metadata)
	// build a query string using them
	query := "INSERT INTO " + app + "." + entity + " (" + strings.Join(keys, ", ") + ") VALUES (" + strings.Join(values, ", ") + ");"
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

func (c *CrudController) UpdateEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	entity := vars["entity"]
	id := vars["id"]

	// fetch metadata about table to know how to input the data; i.e. if we need quotes or not
	metadata := c.getMetadata(app, entity)

	decoder := json.NewDecoder(r.Body)
	input := map[string]interface{}{}
	err := decoder.Decode(&input)
	if err != nil {
		http.Error(w, "Error Decoding Body.", 500)
	}

	//remove id if it is passed in (cant update primary key)
	input = c.filterPrimaryKey(input)

	// get the keys and values from the body
	keys, values := c.formatForCQL(input, metadata)

	// build update string
	updates := make([]string, len(keys))
	fmt.Println(keys)
	for i := range keys {
		updates[i] = (keys[i] + " = " + values[i])
	}
	updatesQuery := strings.Join(updates, ", ")
	fmt.Println(updatesQuery)

	// build a query string using them
	query := "UPDATE " + app + "." + entity + " SET " + updatesQuery + "WHERE id = " + id + ";"
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

func (c *CrudController) GetEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	entity := vars["entity"]
	id := vars["id"]

	query := "SELECT * FROM " + app + "." + entity + " WHERE id = " + id + ";"
	fmt.Println(query)
	iter := c.session.Query(query).Iter()

	// First iteration
	res := make(map[string]interface{})
	if !iter.MapScan(res) {
		iter.Close()
		log.Println("Failed to get by id.")
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

func (c *CrudController) getMetadata(app string, entity string) map[string]string {
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
	return metadata
}

func (c *CrudController) formatForCQL(input map[string]interface{}, metadata map[string]string) ([]string, []string) {
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
	return keys, values
}

func (c *CrudController) filterPrimaryKey(input map[string]interface{}) map[string]interface{} {
	filteredInput := make(map[string]interface{})

	for key, value := range input {
		if key != "id" {
			filteredInput[key] = value
		}
	}
	return filteredInput
}
