package objects

import (
	"distributed_storage_system/Chapter2/apiServer/locate"
	"distributed_storage_system/utils/objectStream"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, e := getStream(object)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fali", object)
	}
	return objectStream.NewGetSteam(server, object)
}
