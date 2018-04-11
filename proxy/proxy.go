package main

import (
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"crypto/sha1"
	"encoding/base64"
	"bytes"
	"os"
)

var mapOfKeyToServer = make(map[string]string)
var numberOfServers = 2
var _ = numberOfServers
var server = "http://localhost"
var mapOfServerIdAndPort = make(map[int]string)
var proxyPort = 5555

func main() {

	for i, arg := range os.Args{
		if i == 0{
			continue
		}
		fmt.Println(arg)
		mapOfServerIdAndPort[i-1] = arg
	}

	for i := 0; i<= len(mapOfServerIdAndPort); i++ {
		fmt.Println(mapOfServerIdAndPort[i])
	}

	http.HandleFunc("/", proxyHandler)
	log.Printf("[INFO] Proxy listening for /get and /set requests on port: %v and its type is %T\n", proxyPort, proxyPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", proxyPort), nil))
}

func proxyHandler(w http.ResponseWriter, r *http.Request){
	endPoint := r.URL.Path
	_ = endPoint
	method := r.Method
	_ = method
	log.Println("endPoint: " + endPoint)
	if endPoint == "/set"{
		assignSetNode(w, r)
	}else if endPoint == "/get"{
		assignGetNode( w, r )
	}
}

func assignSetNode(w http.ResponseWriter, r *http.Request){
	var request setRequestOnNode
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	//var buff bytes.Buffer


	log.Printf("%T", r.Method)

	if err != nil {
		http.Error(w, "[ERROR] [PROXY] Bad request. Err decoding request on /set", http.StatusBadRequest)
		return
	}

	shaSumOfKey := getKeyForValue([]byte(request.Key))
	asciiSumOfShaSum := getAsciiSumOfIndividualCharactersInString(shaSumOfKey)
	serverId := asciiSumOfShaSum % numberOfServers
	serverToQuery := server + ":" + mapOfServerIdAndPort[serverId]
	log.Println("[INFO] [PROXY] asciiSumOfShaSum=", asciiSumOfShaSum, " and therefore, serverId: ", serverId,
		"\nand therefore, mapOfServerIdAndPort[serverId]: " + mapOfServerIdAndPort[serverId]);
	serverToQuery += "/set"
	log.Println("[INFO] [PROXY] serverToQuery=", serverToQuery)

	b, _ := json.Marshal(request)
	log.Println(string(b))

	relayedResponse, err := http.Post(serverToQuery, "application/json", bytes.NewReader( b ) )


	if err != nil {
		http.Error(w, "[ERROR] [PROXY] Bad request. Did not get response from server for /set", http.StatusBadRequest)
		return
	}


	//decode and parse the relayed response
	//TODO: find better way to simply relay the message
	var keyValuePair setKeyValueReponseOnNode
	decoder = json.NewDecoder(relayedResponse.Body)
	err = decoder.Decode(&keyValuePair)

	if err != nil {
		http.Error(w, "[ERROR] [PROXY] Bad request. Err decoding relay on /set", http.StatusBadRequest)
		return
	}

	log.Println("[INFO] [PROXY] response to relay...ooo : ", keyValuePair.Message)

	//w.Header().Add("contentType", "application/json")
	encoder := json.NewEncoder(w)

	encoder.Encode(keyValuePair)
}

func assignGetNode(w http.ResponseWriter, r *http.Request){
	var request getRequestOnNode
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	log.Printf("%T", r.Method)

	if err != nil {
		http.Error(w, "[ERROR] [PROXY] Bad request. Err decoding request on /set", http.StatusBadRequest)
		return
	}

	shaSumOfKey := getKeyForValue([]byte(request.Key))
	asciiSumOfShaSum := getAsciiSumOfIndividualCharactersInString(shaSumOfKey)
	serverId := asciiSumOfShaSum % numberOfServers
	serverToQuery := server + ":" + mapOfServerIdAndPort[serverId]
	log.Println("[INFO] [PROXY] asciiSumOfShaSum=", asciiSumOfShaSum, " and therefore, serverId: ", serverId,
		"\nand therefore, mapOfServerIdAndPort[serverId]: " + mapOfServerIdAndPort[serverId]);
	serverToQuery += "/get"
	log.Println("[INFO] [PROXY] serverToQuery=", serverToQuery)

	b, _ := json.Marshal(request)
	log.Println(string(b))

	relayedResponse, err := http.Post(serverToQuery, "application/json", bytes.NewReader( b ) )

	if err != nil {
		http.Error(w, "[ERROR] [PROXY] Bad request. Did not get response from server for /set", http.StatusBadRequest)
		return
	}


	//decode and parse the relayed response
	//TODO: find better way to simply relay the message
	var keyValuePair getResponse
	decoder = json.NewDecoder(relayedResponse.Body)
	err = decoder.Decode(&keyValuePair)

	if err != nil {
		http.Error(w, "[ERROR] [PROXY] Bad request. Err decoding relay on /set", http.StatusBadRequest)
		return
	}

	log.Println("[INFO] [PROXY] response to relay...ooo : ", keyValuePair.Value)

	//w.Header().Add("contentType", "application/json")
	encoder := json.NewEncoder(w)

	encoder.Encode(keyValuePair)

}

func getAsciiSumOfIndividualCharactersInString(inputString string) int{

	asciiSumOfString := 0
	for _, c :=  range inputString{
		characterInInt := int(c)
		asciiSumOfString += characterInInt
	}
	return asciiSumOfString
}

type setKeyValueReponseOnNode struct {
	// change the output field to be "message"
	Message string `json:"message"`
	// do not output this field
	Author string `json:"-"`
	// do not output the field if the value is empty
	Date string `json: ",omitempty"`
	//Convert the output to string and rename to "id"
	Id int `json:"id, string"`
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