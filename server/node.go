package main

import (
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"encoding/base64"
	"crypto/sha1"
	"os"
	"strconv"
)

var key_value_store = make(map[string]string)
var port = 8080
var selfNode = "127.0.0.1" + string(port)
var selfNodeInByteArray = []byte(selfNode)
var nodeId = getKeyForValue(selfNodeInByteArray)

func main() {

	//key_value_store = make(map[string]int)
	http.HandleFunc("/get", getKeyValueHandlerOnNode)
	http.HandleFunc("/set", setKeyValueHandlerOnNode)
	log.Println(nodeId)
	port, _	 := strconv.ParseInt(os.Args[1], 10, 64)
	log.Printf("[INFO] Node listening for /get and /set requests on port: %v and its type is %T\n", port, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))


}

func getKeyValueHandlerOnNode(w http.ResponseWriter, r *http.Request){
	//fmt.Fprint(w, "Welcome to my node :D")


	var request getRequestOnNode
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil {
		http.Error(w, "[ERROR] Bad request. Error decoding request on /get", http.StatusBadRequest)
		return
	}

	var resp = ""

	log.Printf("[INFO] Fetching: key: %s\n", request.Key)
	resp += key_value_store[request.Key]
	log.Printf("[INFO]%s\n and len of k-v store:%d", resp, len(key_value_store))
	for i, n := range key_value_store{
		fmt.Println("i " + i + " v:" + n)
	}
	response := getResponse{Value:resp}

	/*
	The encoding/json package has a function called NewEncoder
	this returns us an Encoder object that can be used to write JSON straight to an open writer
	func NewEncoder(w io.Writer) *Encoder
	*/

	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}


func setKeyValueHandlerOnNode(w http.ResponseWriter, r *http.Request){

	/*body, err := ioutil.ReadAll(r.Body)
	if ( err != nil ){
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	fmt.Printf(string(body))
	var request setRequestOnNode
	err = json.Unmarshal(body, &request)
	if ( err != nil ){
		http.Error(w, "Failed to unmarshal", http.StatusBadRequest)
		return
	}*/


	var request setRequestOnNode
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	log.Printf("%T in set", r.Method)

	/*if r.Method == "PUT" {
		//update

	}else if r.Method == "POST"{
		//insert

	}*/

	//Current method capable of handling GET/POST/PUT

	if err != nil {

		http.Error(w, "[ERROR] Bad request. Err decoding request on /set", http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] Storing: key: %s, data: %s\n", request.Key, request.Value)
	key_value_store[request.Key] = request.Value

	log.Printf( "[INFO] Len of k-v store: %d\n", len(key_value_store) )
	response := setKeyValueReponseOnNode{Message:"Success"}
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

type setKeyValueReponseOnNode struct {
	// change the output field to be "message"
	Message string `json:"message"`
}


type setRequestOnNode struct{
	Key string `"json:string"`
	Value string `"json:string"`
}

type getRequestOnNode struct{
	Key string `"json:string"`
}

type getResponse struct{
	Value string `"json:string"`
}

func /* (ms *MapServer) */ getKeyForValue(byteArray []byte) string{
	hasher := sha1.New()
	hasher.Write(byteArray)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(byteArray))
	return sha
}