package objects

import (
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("apiServer get a handler")
	m := r.Method
	if m == http.MethodPut {
		log.Println("apiServer get a put handler")
		put(w, r)
		return
	}
	if m == http.MethodGet {
		log.Println("apiServer get a get handler")
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
