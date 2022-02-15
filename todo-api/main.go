package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ToDo struct {
	Id        int32  `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

type ToDoMessage struct {
	Success        bool   `json:"success"`
	Msg            string `json:"msg"`
	ActiveTodos    []ToDo `json:"activetodos"`
	CompletedTodos []ToDo `json:"completedtodos"`
}

type NormalMessage struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

var openTodos []ToDo
var completedTodos []ToDo

func getAllTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var message = ToDoMessage{
		Success:        true,
		Msg:            strconv.Itoa(len(openTodos)) + " active todos and " + strconv.Itoa(len(completedTodos)) + " completed todos",
		ActiveTodos:    openTodos,
		CompletedTodos: completedTodos,
	}
	json.NewEncoder(w).Encode(message)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todo ToDo
	err := json.NewDecoder(r.Body).Decode(&todo)
	fmt.Println("POST:", todo)
	if err != nil || todo.Title == "" || todo.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		var message = NormalMessage{Success: false, Msg: "Bad Request"}
		json.NewEncoder(w).Encode(message)
		return
	}

	todo.Id = int32(rand.Intn(100000))
	if todo.Completed {
		completedTodos = append(completedTodos, todo)
	} else {
		openTodos = append(openTodos, todo)
	}
	fmt.Println("Final POST:", todo)

	var message = NormalMessage{Success: true, Msg: "Todo created successfully"}
	json.NewEncoder(w).Encode(message)

}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var id int32
	params := mux.Vars(r)

	var tempTodo ToDo
	err := json.NewDecoder(r.Body).Decode(&tempTodo)
	if err != nil {
		fmt.Println(err)
	}

	for index, todo := range openTodos {
		fmt.Sscan(params["id"], &id)
		if todo.Id == id {
			openTodos = append(openTodos[:index], openTodos[index+1:]...)
			tempTodo.Id = id
			if tempTodo.Completed {
				completedTodos = append(completedTodos, tempTodo)
			} else {
				openTodos = append(openTodos, tempTodo)
			}

			fmt.Println("PUT:", tempTodo)

			var message = NormalMessage{Success: true, Msg: "Todo updated successfully"}
			json.NewEncoder(w).Encode(message)
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	var message = NormalMessage{Success: false, Msg: "Failed to update todo"}
	json.NewEncoder(w).Encode(message)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var id int32
	params := mux.Vars(r)
	fmt.Sscan(params["id"], &id)

	for index, todo := range openTodos {
		if todo.Id == id {
			openTodos = append(openTodos[:index], openTodos[index+1:]...)
			var message = NormalMessage{Success: true, Msg: "Todo deleted successfully"}
			json.NewEncoder(w).Encode(message)
			return
		}
	}

	for index, todo := range completedTodos {
		if todo.Id == id {
			completedTodos = append(completedTodos[:index], completedTodos[index+1:]...)
			var message = NormalMessage{Success: true, Msg: "Todo deleted successfully"}
			json.NewEncoder(w).Encode(message)
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	var message = NormalMessage{Success: false, Msg: "Failed to delete todo"}
	json.NewEncoder(w).Encode(message)
}

func searchTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := strings.ToLower(r.URL.Query().Get("query"))
	fmt.Println(query)

	var resOpenTodos []ToDo
	var resCompletedTodos []ToDo

	for _, todo := range openTodos {
		titleMatched, _ := regexp.MatchString(query, strings.ToLower(todo.Title))
		textMatched, _ := regexp.MatchString(query, strings.ToLower(todo.Text))
		if titleMatched || textMatched {
			fmt.Println(todo)
			resOpenTodos = append(resOpenTodos, todo)
		}
	}

	for _, todo := range completedTodos {
		titleMatched, _ := regexp.MatchString(query, strings.ToLower(todo.Title))
		textMatched, _ := regexp.MatchString(query, strings.ToLower(todo.Text))
		if titleMatched || textMatched {
			fmt.Println(todo)
			resCompletedTodos = append(resCompletedTodos, todo)
		}
	}

	var message = ToDoMessage{
		Success:        true,
		Msg:            strconv.Itoa(len(resOpenTodos)) + " active todos and " + strconv.Itoa(len(resCompletedTodos)) + " completed todos found",
		ActiveTodos:    resOpenTodos,
		CompletedTodos: resCompletedTodos,
	}

	json.NewEncoder(w).Encode(message)
}

func markAsCompleted(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var id int32
	params := mux.Vars(r)

	for index, todo := range openTodos {
		fmt.Sscan(params["id"], &id)
		if todo.Id == id {
			todo.Completed = true
			openTodos = append(openTodos[:index], openTodos[index+1:]...)
			completedTodos = append(completedTodos, todo)

			var message = NormalMessage{Success: true, Msg: "Todo marked as complete"}
			json.NewEncoder(w).Encode(message)
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	var message = NormalMessage{Success: false, Msg: "Failed to update todo"}
	json.NewEncoder(w).Encode(message)

}

func main() {
	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	rand.Seed(time.Now().UnixNano())
	openTodos = []ToDo{}
	completedTodos = []ToDo{}

	// Dummy Data

	// openTodos = append(openTodos,
	// 	ToDo{Id: int32(rand.Intn(100000)), Title: "First Todo", Text: "Dummy text 1", Completed: false},
	// 	ToDo{Id: int32(rand.Intn(100000)), Title: "Second Todo", Text: "Dummy text 2", Completed: false},
	// )

	// completedTodos = append(completedTodos,
	// 	ToDo{Id: int32(rand.Intn(100000)), Title: "Completed Todo", Text: "Dummy text 3", Completed: true},
	// )

	router.HandleFunc("/todos", getAllTodos).Methods("GET")
	router.HandleFunc("/todos", createTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")
	router.HandleFunc("/todos/search", searchTodo).Methods("GET")
	router.HandleFunc("/todos/markcompleted/{id}", markAsCompleted).Methods("PUT")

	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(headers, methods, origins)(router)))
}
