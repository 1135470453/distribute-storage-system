package objects

import (
	"distributed_storage_system/Chapter2/apiServer/heartbeat"
	"distributed_storage_system/Chapter2/apiServer/locate"
	"distributed_storage_system/utils/elasticSearch"
	"distributed_storage_system/utils/rs"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//func get(w http.ResponseWriter, r *http.Request) {
//	log.Println("apiServer start get")
//	name := strings.Split(r.URL.EscapedPath(), "/")[2]
//	versionId := r.URL.Query()["version"]
//	version := 0
//	var e error
//	if len(versionId) != 0 {
//		version, e = strconv.Atoi(versionId[0])
//		if e != nil {
//			log.Println(e)
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//	}
//	log.Println("name:" + name)
//	log.Printf("version:%d", version)
//	meta, e := elasticSearch.GetMetadata(name, version)
//	if e != nil {
//		log.Println(e)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	if meta.Hash == "" {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//	object := url.PathEscape(meta.Hash)
//	stream, e := getStream(object)
//	if e != nil {
//		log.Println(e)
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//	io.Copy(w, stream)
//}
//func getStream(object string) (io.Reader, error) {
//	server := locate.Locate(object)
//	if server == "" {
//		return nil, fmt.Errorf("object %s locate fali", object)
//	}
//	return objectStream.NewGetSteam(server, object)
//}

func get(w http.ResponseWriter, r *http.Request) {
	log.Println("apiServer get start")
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionId := r.URL.Query()["version"]
	version := 0
	var e error
	if len(versionId) != 0 {
		version, e = strconv.Atoi(versionId[0])
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	meta, e := elasticSearch.GetMetadata(name, version)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	hash := url.PathEscape(meta.Hash)
	//RS码会对没有写满的片进行填充,所以必须指明size
	stream, e := GetStream(hash, meta.Size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, e = io.Copy(w, stream)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	stream.Close()
	log.Println("apiServer get end")
}

func GetStream(hash string, size int64) (*rs.RSGetStream, error) {
	log.Println("apiServer GetStream start")
	locateInfo := locate.Locate(hash)
	//返回节点数小于四,无法修复,报错
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateInfo)
	}
	dataServers := make([]string, 0)
	//返回节点数不等于6,选择新节点修复
	if len(locateInfo) != rs.ALL_SHARDS {
		dataServers = heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	log.Println("apiServer GetStream end")
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}
