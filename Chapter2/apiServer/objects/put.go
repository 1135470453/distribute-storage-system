package objects

import (
	"distributed_storage_system/Chapter2/apiServer/heartbeat"
	"distributed_storage_system/Chapter2/apiServer/locate"
	"distributed_storage_system/utils/elasticSearch"
	"distributed_storage_system/utils/headutils"
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
	hash := headutils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("hash is " + hash)
	size := headutils.GetSizeFromHeader(r.Header)
	log.Printf("size if %d\n", size)
	c, e := storeObject(r.Body, hash, size)
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
	e = elasticSearch.AddVersion(name, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Println("apiServer end put")
}

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	log.Println("apiServer start storeObject")
	//url.PathEscape确保可以放到url中使用
	if locate.Exist(url.PathEscape(hash)) {
		log.Println("this is hash is exist")
		return http.StatusOK, nil
	}
	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		return http.StatusInternalServerError, e
	}
	//将r写入stream,返回stream内容用于之后的判断
	//使用temStream.Write方法
	reader := io.TeeReader(r, stream)
	d := headutils.CalculateHash(reader)
	if d != hash {
		log.Println("客户端hash错误")
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch,"+
			"calculated=%s, requested=%s", d, hash)
	}
	stream.Commit(true)
	log.Println("apiServer : storeObject end!")
	return http.StatusOK, nil
}

//随机选择一个dataserver,并返回用于传数据的putStream
func putStream(hash string, size int64) (*objectStream.TempPutStream, error) {
	log.Println("apiServer start putStream")
	server := heartbeat.ChooseRandomDataServer()
	log.Println("server is " + server)
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}
	log.Println("apiServer :putStream end!")
	return objectStream.NewTempPutStream(server, hash, size)
}
