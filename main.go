package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `jason:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

type allTask []task

var tasks = allTask{
	{
		ID:      1,
		Name:    "Task one",
		Content: "Some Content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Tipo de dato en formato json
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid Task")
	}

	json.Unmarshal(reqBody, &newTask)

	newTask.ID = len(tasks) + 1

	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json") // Tipo de dato
	w.WriteHeader(http.StatusCreated)                  //Si fue correcto
	json.NewEncoder(w).Encode(newTask)                 //Enviar dato agregado
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for _, task := range tasks { //Recorrer elemento por elemento de la lista
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json") // Tipo de dato en formato json
			json.NewEncoder(w).Encode(task)
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "invalid ID to delete..")
		return
	}

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "The task with ID %v has been remove succesfully", taskID)
		}
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask task
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please Enter Valid Date")
	}
	json.Unmarshal(reqBody, &updatedTask)

	for i, t := range tasks {
		if t.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			updatedTask.ID = taskID
			tasks = append(tasks, updatedTask)

			fmt.Fprintf(w, "The task with the ID %v has been updated succesfully", taskID)
		}
	}

}

func indexRoute(w http.ResponseWriter, r *http.Request) { //w respuesta al cliente desde el servidor, r peticion del usuario al servidor
	fmt.Fprintf(w, "Welcome to my API")
}

func main() {
	router := mux.NewRouter().StrictSlash(true) // Crear enrutador

	router.HandleFunc("/", indexRoute) // "/"-> URL PRINCIPAL
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	http.ListenAndServe(":3000", router)
	log.Fatal(http.ListenAndServe(":3000", router))

}

//Compile Daemon
// 1. Compile Daemon - Compilar codigo cada vez que guardo
// 2. CompileDaemon -command="CRUD.exe"
