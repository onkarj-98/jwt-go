package main

import (
	"github.com/gorilla/mux"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("ORM")

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/users",allUsers).Methods("GET")
	myRouter.HandleFunc("/user/{name}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{name}/{email}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/user/{name}/{email}", newUser).Methods("POST")
	log.Println("Starting server on 8082")
	log.Fatal(http.ListenAndServe(":8082", myRouter))


}