package objects

import (
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("dataServer get a handler")
	m := r.Method
	if m == http.MethodGet {
		log.Println("dataServer get a get")
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
