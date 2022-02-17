package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
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
	Date      string `json:"date"`
}

type RequestMessage struct {
	Title string `json:"title"`
	Text  string `json:"text"`
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

//This function returns all the active and completed todos
func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var message = ToDoMessage{
		Success:        true,
		Msg:            strconv.Itoa(len(openTodos)) + " active todos and " + strconv.Itoa(len(completedTodos)) + " completed todos",
		ActiveTodos:    openTodos,
		CompletedTodos: completedTodos,
	}
	json.NewEncoder(w).Encode(message)
}

//This function creates a todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todo ToDo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil || todo.Title == "" || todo.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		var message = NormalMessage{Success: false, Msg: "Bad Request"}
		json.NewEncoder(w).Encode(message)
		return
	}

	loc, locerr := time.LoadLocation("Asia/Kolkata")
	if locerr != nil {
		fmt.Println(locerr)
	}

	todo.Id = int32(rand.Intn(100000))
	todo.Date = time.Now().In(loc).Format(time.RFC1123)
	if todo.Completed {
		completedTodos = append(completedTodos, todo)
	} else {
		openTodos = append(openTodos, todo)
	}

	var message = NormalMessage{Success: true, Msg: "Todo created successfully"}
	json.NewEncoder(w).Encode(message)

}

//This function updates a todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tid int32
	params := mux.Vars(r)

	_, parseErr := fmt.Sscan(params["id"], &tid)
	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		var message = NormalMessage{Success: false, Msg: "Bad Request"}
		json.NewEncoder(w).Encode(message)
		return
	}

	var req RequestMessage
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var message = NormalMessage{Success: false, Msg: "Bad Request"}
		json.NewEncoder(w).Encode(message)
		return
	}

	for index, todo := range openTodos {
		if todo.Id == tid {
			newTodo := todo
			openTodos = append(openTodos[:index], openTodos[index+1:]...)
			newTodo.Title = req.Title
			newTodo.Text = req.Text
			openTodos = append(openTodos, newTodo)

			w.WriteHeader(http.StatusOK)
			var message = NormalMessage{Success: true, Msg: "Todo updated successfully"}
			json.NewEncoder(w).Encode(message)
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	var message = NormalMessage{Success: false, Msg: "Failed to update todo"}
	json.NewEncoder(w).Encode(message)

}

//This function deletes a todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tid int32
	params := mux.Vars(r)

	_, parseErr := fmt.Sscan(params["id"], &tid)
	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		var message = NormalMessage{Success: false, Msg: "Bad Request"}
		json.NewEncoder(w).Encode(message)
		return
	}

	for index, todo := range openTodos {
		if todo.Id == tid {
			openTodos = append(openTodos[:index], openTodos[index+1:]...)
			var message = NormalMessage{Success: true, Msg: "Todo deleted successfully"}
			json.NewEncoder(w).Encode(message)
			return
		}
	}

	for index, todo := range completedTodos {
		if todo.Id == tid {
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

//This function marks an active todo as completed
func MarkAsCompleted(w http.ResponseWriter, r *http.Request) {
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

//This function returns all the todos whose title or text match with the qiven query
func SearchTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type SearchRequest struct {
		Query string `json:"query"`
	}
	var querystring SearchRequest
	err := json.NewDecoder(r.Body).Decode(&querystring)
	if err != nil || querystring.Query == "" {
		w.WriteHeader(http.StatusBadRequest)
		var message = NormalMessage{Success: false, Msg: "Bad Request"}
		json.NewEncoder(w).Encode(message)
		return
	}
	fmt.Println(querystring)

	var resOpenTodos []ToDo
	var resCompletedTodos []ToDo

	for _, todo := range openTodos {
		titleMatched, _ := regexp.MatchString(querystring.Query, strings.ToLower(todo.Title))
		textMatched, _ := regexp.MatchString(querystring.Query, strings.ToLower(todo.Text))
		if titleMatched || textMatched {
			// fmt.Println(todo)
			resOpenTodos = append(resOpenTodos, todo)
		}
	}

	for _, todo := range completedTodos {
		titleMatched, _ := regexp.MatchString(querystring.Query, strings.ToLower(todo.Title))
		textMatched, _ := regexp.MatchString(querystring.Query, strings.ToLower(todo.Text))
		if titleMatched || textMatched {
			// fmt.Println(todo)
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

func main() {
	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	rand.Seed(time.Now().UnixNano())

	// Dummy Data

	openTodos = append(openTodos,
		ToDo{Id: int32(1), Title: "First Todo", Text: "Dummy text 1", Completed: false, Date: time.Now().Format(time.RFC1123)},
		ToDo{Id: int32(2), Title: "Second Todo", Text: "Dummy text 2", Completed: false, Date: time.Now().Format(time.RFC1123)},
	)

	completedTodos = append(completedTodos,
		ToDo{Id: int32(3), Title: "Completed Todo", Text: "Dummy text 3", Completed: true, Date: time.Now().Format(time.RFC1123)},
	)

	router.HandleFunc("/todos", GetAllTodos).Methods("GET")
	router.HandleFunc("/todos", CreateTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", DeleteTodo).Methods("DELETE")
	router.HandleFunc("/todos/markcompleted/{id}", MarkAsCompleted).Methods("PUT")
	router.HandleFunc("/todos/search", SearchTodo).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router)))
}
