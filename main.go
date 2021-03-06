package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/soajs/soajs.golang"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	soajs := r.Context().Value(soajsgo.SoajsKey).(soajsgo.ContextData)
	respJson, err := json.Marshal(soajs)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func SayHelloPost(w http.ResponseWriter, r *http.Request) {
	soajs := r.Context().Value(soajsgo.SoajsKey).(soajsgo.ContextData)
	controller := soajs.Awareness.Path()
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
		log.Println(err)
	}
	log.Println("Successfully Opened soa.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result soajsgo.Config
	json.Unmarshal([]byte(byteValue), &result)

	soajs, err := soajsgo.NewFromConfig(context.Background(), result)
	if err != nil {
		log.Fatal(err)
	}

	router.Use(soajs.Middleware)

	router.HandleFunc("/tidbit/hello", SayHello).Methods("GET")
	router.HandleFunc("/tidbit/hello", SayHelloPost).Methods("POST")

	router.HandleFunc("/heartbeat", Heartbeat)

	port := fmt.Sprintf(":%d", result.ServicePort)
	log.Fatal(http.ListenAndServe(port, router))
}
