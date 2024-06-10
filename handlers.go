package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UnixNano()
	k := mux.Vars(r)["key"]

	if err := ds.validKey(k); err != nil {
		http.Error(w, k+": invalid key", http.StatusBadRequest)
		return
	}

	ch := make(chan []byte, 1)
	go func() {
		ch <- ds.get(k)
	}()

	select {
	case v := <-ch:
		if v != nil {
			log.Println("get ["+k+"] in", TimeSpent(t0), "ms")
			w.Write(v)
		} else {
			log.Println(k + ": not found")
			http.Error(w, k+": not found", http.StatusNotFound)
		}
	case <-time.After(timeout):
		ReturnTimeout(w, k)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UnixNano()
	k := mux.Vars(r)["key"]
	ttl, _ := strconv.Atoi(r.FormValue("ttl"))
	if v, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, k+": can't get value", http.StatusBadRequest)
		return
	} else {
		if err = ds.validKey(k); err == nil {
			if err = ds.validValue(v); err == nil {
				go func() {
					ds.set(k, []byte(v), ttl)
				}()

				log.Println("set ["+k+"] in", TimeSpent(t0), "ms")
				w.WriteHeader(http.StatusCreated)
				return
			}
		}
		DatastoreErrorToHTTPError(err, w)
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UnixNano()
	k := mux.Vars(r)["key"]

	if err := ds.delete(k); err == nil {
		log.Println("delete ["+k+"] in", TimeSpent(t0), "ms")
		w.WriteHeader(http.StatusNoContent)
	} else {
		DatastoreErrorToHTTPError(err, w)
	}
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UnixNano()
	k := mux.Vars(r)["key"]

	if v, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, k+": can't get value", http.StatusBadRequest)
		return
	} else {
		if err = ds.append(k, v); err == nil {
			log.Println("append ["+k+"] in", TimeSpent(t0), "ms")
			w.WriteHeader(http.StatusOK)
		} else {
			DatastoreErrorToHTTPError(err, w)
		}
	}
}
