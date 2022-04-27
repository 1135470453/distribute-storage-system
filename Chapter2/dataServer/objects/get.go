package objects

import (
	"distributed_storage_system/Chapter2/dataServer/locate"
	"distributed_storage_system/utils/headutils"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	log.Println("object get start")
	file := getFile(strings.Split(r.URL.EscapedPath(), "/")[2])
	if file == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sendFile(w, file)
	log.Println("object get end")
}

func getFile(hash string) string {
	log.Println("getFile start")
	file := os.Getenv("STORAGE_ROOT") + "/objects/" + hash
	f, _ := os.Open(file)
	d := url.PathEscape(headutils.CalculateHash(f))
	f.Close()
	if d != hash {
		log.Println("object hash mismatch, remove", file)
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	log.Println("getFile end")
	return file
}

func sendFile(w io.Writer, file string) {
	log.Println("sendFile start")
	f, _ := os.Open(file)
	defer f.Close()
	io.Copy(w, f)
	log.Println("sendFile end")
}
