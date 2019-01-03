package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/soajs/soajs.golang"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Response struct {
	Message string `json:"message"`
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	resp := Response{}
	resp.Message = fmt.Sprintf("heartbeat")
	respJson, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func SayHello(w http.ResponseWriter, r *http.Request) {
	soajs := r.Context().Value("soajs").(soajsGo.SOAJSObject)
	respJson, err := json.Marshal(soajs)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func SayHelloPost(w http.ResponseWriter, r *http.Request) {
	soajs := r.Context().Value("soajs").(soajsGo.SOAJSObject)
	controller := soajs.Awareness.GetHost()
	log.Println(controller)
	response, err := json.Marshal(soajs)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

//main function
func main() {
	router := mux.NewRouter()

	jsonFile, err := os.Open("soa.json")
	if err != nil {
		panic(err)
	}
	log.Println("Successfully Opened soajs.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	soajsMiddleware := soajsGo.InitMiddleware(result)
	router.Use(soajsMiddleware)

	router.HandleFunc("/tidbit/hello", SayHello).Methods("GET")
	router.HandleFunc("/tidbit/hello", SayHelloPost).Methods("POST")

	router.HandleFunc("/heartbeat", Heartbeat)

	log.Println("starting")
	log.Fatal(http.ListenAndServe(":"+strconv.FormatFloat(result["servicePort"].(float64), 'f', -1, 64), router))
}
