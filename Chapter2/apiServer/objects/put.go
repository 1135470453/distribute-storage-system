package objects

import (
	"distributed_storage_system/Chapter2/apiServer/heartbeat"
	"distributed_storage_system/utils/objectStream"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	log.Println("apiServer start put")
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	log.Println("object :" + object)
	c, e := storeObject(r.Body, object)
	if e != nil {
		log.Println(e)
	}
	w.WriteHeader(c)
	log.Println("apiServer : put end!")
}

func storeObject(r io.Reader, object string) (int, error) {
	log.Println("apiServer start storeObject")
	stream, e := putStream(object)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}
	//putStream实现了write接口,这里应该是直接调用了write方法
	io.Copy(stream, r)
	e = stream.Close()
	if e != nil {
		return http.StatusInternalServerError, e
	}
	log.Println("apiServer : storeObject end!")
	return http.StatusOK, nil
}

//随机选择一个dataserver,并返回用于传数据的putStream
func putStream(object string) (*objectStream.PutStream, error) {
	log.Println("apiServer start putStream")
	server := heartbeat.ChooseRandomDataServer()
	log.Println("server is " + server)
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}
	log.Println("apiServer :putStream end!")
	return objectStream.NewPutStream(server, object), nil
}
