package main

import (
	"github.com/gorilla/mux"

	"log"

	"net/http"

	"time"
)

var ds datastorer

var timeout = time.Millisecond * 100

func main() {
	if d, err := newDatastore(); err != nil {
		log.Fatalln(err)
	} else {
		ds = datastorer(d)
	}

	r := mux.NewRouter()
	kr := r.PathPrefix("/key/{key}").Subrouter()
	kr.Methods("GET").HandlerFunc(GetHandler)
	kr.Methods("POST").HandlerFunc(PostHandler)
	kr.Methods("DELETE").HandlerFunc(DeleteHandler)
	kr.Methods("PUT").HandlerFunc(PutHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
