package objects

import (
	"distributed_storage_system/utils/elasticSearch"
	"log"
	"net/http"
	"strings"
)

//通过增加新版本并且size和hash设置为0的方法删除
func del(w http.ResponseWriter, r *http.Request) {
	log.Println("apiServer start del")
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	log.Println("name is " + name)
	version, e := elasticSearch.SearchLatestVersion(name)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	e = elasticSearch.PutMetadata(name, version.Version+1, 0, "")
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("apiServer del success and end")
		return
	}
	log.Println("apiServer del fail and end")
}
