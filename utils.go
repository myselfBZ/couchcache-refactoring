package main

import (
	"log"
	"math"
	"net/http"
	"time"
)

func ReturnTimeout(w http.ResponseWriter, k string) {
	log.Println(k + ": timeout")
	http.Error(w, k+": timeout", http.StatusRequestTimeout)
}

func TimeSpent(t0 int64) int64 {
	return int64(math.Floor(float64(time.Now().UnixNano()-t0)/1000000 + .5))
}

func DatastoreErrorToHTTPError(err error, w http.ResponseWriter) {
	switch err {
	case errNotFound:
		http.Error(w, "key not found", http.StatusNotFound)
	case errEmptyBody:
		http.Error(w, "empty value", http.StatusBadRequest)
	case errOversizedBody:
		http.Error(w, "oversized value", http.StatusBadRequest)
	case errInvalidKey:
		http.Error(w, "invalid key", http.StatusBadRequest)
	case errKeyExists:
		http.Error(w, "key exists", http.StatusBadRequest)
	default:
		http.Error(w, "cache server error", http.StatusInternalServerError)
	}
}
