package main

import (
    "log"
    "fmt"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/soajs/soajs.golang"
)

type Response struct {
    Message             string          `json:"message"`
}

func SayHello(w http.ResponseWriter, r *http.Request) {
    vars := r.URL.Query()

    resp := Response{}
    resp.Message = fmt.Sprintf("Hello DEMO, I am a GOLANG service, you are %v and your last name is: %v", vars["username"], vars["lastname"])

    respJson, err := json.Marshal(resp)
    if err != nil {
        panic(err)
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func ReturnSoajsData(w http.ResponseWriter, r *http.Request) {
    soajs := r.Context().Value("soajs").(soajsGo.SOAJSObject)
    soajs.Controller = soajs.Awareness.GetHost()

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
    soajsMiddleware := soajsGo.InitMiddleware(map[string]string{"serviceName":"golang"})
    router.Use(soajsMiddleware)

    router.HandleFunc("/tidbit/hello", SayHello).Methods("GET")
    router.HandleFunc("/tidbit/hello", ReturnSoajsData).Methods("POST")

    log.Fatal(http.ListenAndServe(":4383", router))
}
