package controllers

import (
	"context"
	"employee-api/configs"
	"employee-api/models"
	"employee-api/responses"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var employeeCollection *mongo.Collection = configs.GetCollection(configs.DB, "employees")
var validate = validator.New()

func CreateEmployee() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var employee models.Employee
		defer cancel()

		//validate the request body
		if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.EmployeeResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&employee); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.EmployeeResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		newEmployee := models.Employee{
			Id:       primitive.NewObjectID(),
			Name:     employee.Name,
			Location: employee.Location,
			Title:    employee.Title,
			Email:    employee.Email,
			
		}
		result, err := employeeCollection.InsertOne(ctx, newEmployee)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		response := responses.EmployeeResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)
	}
}
func GetAEmployee() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		employeeId := params["employeeId"]
		var user models.Employee
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(employeeId)

		err := employeeCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&employeeId)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.EmployeeResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}}
		json.NewEncoder(rw).Encode(response)
	}
}

func EditAEmployee() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		employeeId := params["employeeId"]
		var employee models.Employee
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(employeeId)

		if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.EmployeeResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&employee); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.EmployeeResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		update := bson.M{"name": employee.Name, "location": employee.Location, "title": employee.Title}
		result, err := employeeCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//get updated user details
		var updatedUser models.Employee
		if result.MatchedCount == 1 {
			err := employeeCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
				return
			}
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.EmployeeResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}}
		json.NewEncoder(rw).Encode(response)
	}
}

func DeleteAEmployee() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		employeeId := params["employeeId"]
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(employeeId)

		result, err := employeeCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if result.DeletedCount < 1 {
			rw.WriteHeader(http.StatusNotFound)
			response := responses.EmployeeResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.EmployeeResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}}
		json.NewEncoder(rw).Encode(response)
	}
}

func GetAllEmployee() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var employees []models.Employee
		defer cancel()

		results, err := employeeCollection.Find(ctx, bson.M{})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleEmployee models.Employee
			if err = results.Decode(&singleEmployee); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
			}
			employees = append(employees, singleEmployee)
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.EmployeeResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": employees}}
		json.NewEncoder(rw).Encode(response)
	}
}
