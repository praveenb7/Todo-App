package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Seeding dummy data
	openTodos = append(openTodos,
		ToDo{Id: int32(1), Title: "First Todo", Text: "Dummy text 1", Completed: false, Date: "test-date"},
		ToDo{Id: int32(2), Title: "Second Todo", Text: "Dummy text 2", Completed: false, Date: "test-date"},
	)

	completedTodos = append(completedTodos,
		ToDo{Id: int32(3), Title: "Completed Todo", Text: "Dummy text 3", Completed: true, Date: "test-date"},
	)
}

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/todos", GetAllTodos).Methods("GET")
	router.HandleFunc("/todos", CreateTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", DeleteTodo).Methods("DELETE")
	router.HandleFunc("/todos/search", SearchTodo).Methods("GET")
	router.HandleFunc("/todos/markcompleted/{id}", MarkAsCompleted).Methods("PUT")

	return router
}

func TestGetAllTodos(t *testing.T) {
	request, _ := http.NewRequest("GET", "/todos", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	message := `{"success":true,"msg":"2 active todos and 1 completed todos","activetodos":[{"id":1,"title":"First Todo","text":"Dummy text 1","completed":false,"date":"test-date"},{"id":2,"title":"Second Todo","text":"Dummy text 2","completed":false,"date":"test-date"}],"completedtodos":[{"id":3,"title":"Completed Todo","text":"Dummy text 3","completed":true,"date":"test-date"}]}`
	assert.Equalf(t, message+"\n", response.Body.String(), "Invalid Json! Expected %v but got %v instead", message, response.Body.String())
}

func TestCreateTodo(t *testing.T) {
	var jsonStr = []byte(`{"title":"dummy title","text":"dummy text","completed":false}`)

	request, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	message := `{"success":true,"msg":"Todo created successfully"}`
	assert.Equalf(t, message+"\n", response.Body.String(), "Invalid Json! Expected %v but got %v instead", message, response.Body.String())
}

func TestUpdateTodo(t *testing.T) {
	var jsonStr = []byte(`{"title":"dummy title","text":"dummy text"}`)

	request, _ := http.NewRequest("PUT", "/todos/2", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	message := `{"success":true,"msg":"Todo updated successfully"}`
	assert.Equalf(t, message+"\n", response.Body.String(), "Invalid Json! Expected %v but got %v instead", message, response.Body.String())
}

func TestDeleteTodo(t *testing.T) {
	request, _ := http.NewRequest("DELETE", "/todos/2", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	message := `{"success":true,"msg":"Todo deleted successfully"}`
	assert.Equalf(t, message+"\n", response.Body.String(), "Invalid Json! Expected %v but got %v instead", message, response.Body.String())
}

func TestSearchTodo(t *testing.T) {
	req, _ := http.NewRequest("GET", "/todos/search", nil)

	q := req.URL.Query()
	q.Add("query", "first")

	req.URL.RawQuery = q.Encode()

	response := httptest.NewRecorder()
	Router().ServeHTTP(response, req)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	message := `{"success":true,"msg":"1 active todos and 0 completed todos found","activetodos":[{"id":1,"title":"First Todo","text":"Dummy text 1","completed":false,"date":"test-date"}],"completedtodos":null}`
	assert.Equalf(t, message+"\n", response.Body.String(), "Invalid Json! Expected %v but got %v instead", message, response.Body.String())
}

func TestMarkAsCompleted(t *testing.T) {
	request, _ := http.NewRequest("PUT", "/todos/markcompleted/1", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	message := `{"success":true,"msg":"Todo marked as complete"}`
	assert.Equalf(t, message+"\n", response.Body.String(), "Invalid Json! Expected %v but got %v instead", message, response.Body.String())
}

// {"success":true,"msg":"1 active todos and 0 completed todos found","activetodos":[{"id":1,"title":"First Todo","text":"Dummy text 1","completed":false,"date":"test-date"}],"completedtodos":null}
// {"success":true,"msg":"0 active todos and 0 completed todos found","activetodos":null,"completedtodos":null}
