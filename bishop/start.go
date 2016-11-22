package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
)

func handlePlayerClient(w http.ResponseWriter, r *http.Request) {
    test := r.FormValue("test")
    fmt.Println("hello")
	fmt.Println(test)
}

func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/player", handlePlayerClient).Methods("POST")

    fmt.Println("listening...")
    err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
    //err := http.ListenAndServe(":1234", router)
    if err != nil {
        panic(err)
    }
}