package objects

import (
	"distributed_storage_system/Chapter2/apiServer/heartbeat"
	"distributed_storage_system/utils/elasticSearch"
	"distributed_storage_system/utils/hashutils"
	"distributed_storage_system/utils/objectStream"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//用hash值存储在服务器
func put(w http.ResponseWriter, r *http.Request) {
	log.Println("apiServer start put")
	hash := hashutils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("hash is " + hash)
	c, e := storeObject(r.Body, url.PathEscape(hash))
	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		log.Println("apiServer end put")
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		log.Println("apiServer end put")
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size := hashutils.GetSizeFromHeader(r.Header)
	e = elasticSearch.AddVersion(name, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Println("apiServer end put")
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
