package temp

import (
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("dataServer get a handler")
	m := r.Method
	if m == http.MethodPut {
		log.Println("dataServer get a put")
		put(w, r)
		return
	}
	if m == http.MethodPatch {
		log.Println("dataServer get a patch")
		patch(w, r)
		return
	}
	if m == http.MethodPost {
		log.Println("dataServer get a post")
		post(w, r)
		return
	}
	if m == http.MethodDelete {
		log.Println("dataServer get a del")
		del(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
