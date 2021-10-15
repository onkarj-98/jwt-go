package main

import (
	hd "jwt-go/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/signin", hd.Signin)
	http.HandleFunc("/welcome", hd.Welcome)
	http.HandleFunc("/refresh", hd.Refresh)

	log.Println("Listening on localhot:8085")
	log.Fatal(http.ListenAndServe(":8085", nil))

}
