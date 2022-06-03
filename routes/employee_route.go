package routes

import (
	"employee-api/controllers"

	"github.com/gorilla/mux"
)

func EmployeeRoute(router *mux.Router) {
	router.HandleFunc("/employee", controllers.CreateEmployee()).Methods("POST")
	router.HandleFunc("/employee/{employeeId}", controllers.GetAEmployee()).Methods("GET")
	router.HandleFunc("/employee/{employeeId}", controllers.EditAEmployee()).Methods("PUT")
	router.HandleFunc("/employee/{employeeId}", controllers.DeleteAEmployee()).Methods("DELETE")
	router.HandleFunc("/employees", controllers.GetAllEmployee()).Methods("GET")
}
