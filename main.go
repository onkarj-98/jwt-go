package main

import (
	"log"
	"net/http"
	hd"jwt-go/handlers"
)

func main() {
	http.HandleFunc("/signin", hd.Signin)
	http.HandleFunc("/welcome",hd.Welcome)
	http.HandleFunc("/refresh",hd.Refresh)
	log.Println("Server listening on 8082")
	http.ListenAndServe(":8022", nil)
}