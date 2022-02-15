package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/todos", getAllTodos).Methods("GET")
	router.HandleFunc("/todos", createTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")
	router.HandleFunc("/todos/search", searchTodo).Methods("GET")
	router.HandleFunc("/todos/markcompleted/{id}", markAsCompleted).Methods("PUT")

	return router
}

func TestGetAllTodos(t *testing.T) {
	request, _ := http.NewRequest("GET", "/todos", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
	var message = ToDoMessage{
		Success:        true,
		Msg:            "0 active todos and 0 completed todos",
		ActiveTodos:    openTodos,
		CompletedTodos: completedTodos,
	}
	fmt.Println(message)
	databyte, _ := json.Marshal(message)
	assert.Equal(t, string(databyte)+"\n", response.Body.String(), "Invalid Json")
}
