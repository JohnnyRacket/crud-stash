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
	//.Queries("size", "{size:[0-9]+}", "token", "{token:*+}")
	router.HandleFunc("/{app}/{entity}/{id}", c.UpdateEntity).Methods("PUT")
}

func (c *CrudController) GetEntities(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// vals := r.URL.Query()
	app := vars["app"]
	entity := vars["entity"]

	size := 10
	var token []byte
	pageSize := r.URL.Query().Get("size")
	pagingToken := r.URL.Query().Get("token")
	// if they pass in a custom size convert it to an int
	if pageSize != "" {
		var err error
		size, err = strconv.Atoi(pageSize)
		if err != nil {
			fmt.Println("Error with page size input")
		}
	}
	// if there is a token convert it to a byte array
	if pagingToken != "" {
		var err error
		token, err = base64.URLEncoding.DecodeString(pagingToken)
		if err != nil {
			fmt.Println("Error Decoding Paging Token")
		}
	}

	//TODO errors around getting app and entity
	query := c.session.Query("SELECT * FROM " + app + "." + entity + ";")

	paginated := query.PageState(token).PageSize(size)
	iter := paginated.Iter()
	res := []map[string]interface{}{}
	row := &map[string]interface{}{}
	//scan items into a map
	for iter.MapScan(*row) {
		//fmt.Println(m)
		res = append(res, *row)
		row = &map[string]interface{}{}
	}

	//get page state to return in header
	resPagingToken := iter.PageState()
	retPagingToken := base64.URLEncoding.EncodeToString(resPagingToken)

	// marshall result into json fomat
	jsonRes, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
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
