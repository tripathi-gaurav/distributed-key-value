package main

import (
	"net/http"
	"log"
	"fmt"
)



func main() {
	port := 8080
	//key_value_store = make(map[string]int)
	http.HandleFunc("/get", getKeyValueHandler)
	log.Printf("Node listening for /get requests on port: %v and its type is %T\n", port, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
	http.HandleFunc("/set", setKeyValueHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))

}

func getKeyValueHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Welcome to my node :D")

}

func setKeyValueHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Thank you for attempting to request to write data to my keystore :)")

}


