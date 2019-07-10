package main

import (
	"crud-stash/controllers"
	"crud-stash/data"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// our main function
func main() {

	fmt.Println("about to wait print")
	//time.Sleep(time.Duration(180) * time.Second)
	//test cassandra connection
	fmt.Println("done waiting gonna test demo thing")
	session := data.Test()

	//creates a new router
	router := mux.NewRouter()

	//have controllers register their routes
	appController := controllers.NewAppsController() //pass in a db interface here in the future
	appController.RegisterRoutes(router)

	crudController := controllers.NewCrudController(session)
	crudController.RegisterRoutes(router)

	//listen on 8080
	log.Fatal(http.ListenAndServe(":8080", router))
}
