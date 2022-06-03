package main

import (
	"employee-api/configs"
	"employee-api/routes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		json.NewEncoder(rw).Encode(map[string]string{"data": "Hello Tony Stark"})
	}).Methods("GET")

	configs.ConnectDB()
	routes.EmployeeRoute(router)

	log.Fatal(http.ListenAndServe(":5050", router))
}
