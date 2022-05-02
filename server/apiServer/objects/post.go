package objects

import (
	"distributed_storage_system/server/apiServer/heartbeat"
	"distributed_storage_system/server/apiServer/locate"
	"distributed_storage_system/utils/elasticSearch"
	"distributed_storage_system/utils/headutils"
	"distributed_storage_system/utils/rs"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

/*
	POST /objects/filename
	Digest: SHA-256=<hash>
	size:<size>
	返回对应的token，用于后续上传
*/
func post(w http.ResponseWriter, r *http.Request) {
	//获取对象的name、size、hash
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	hash := headutils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//若该hash已正确存在
	if locate.Exist(url.PathEscape(hash)) {
		//增加版本号就返回
		e = elasticSearch.AddVersion(name, hash, size)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	//选择要存分片的数据节点
	ds := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(ds) != rs.ALL_SHARDS {
		log.Println("cannot find enough dataServer")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	//返回一个RSResumablePutStream(保存编码器和写入数据、writer、token)
	stream, e := rs.NewRSResumablePutStream(ds, name, url.PathEscape(hash), size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//将token写入返回头部
	w.Header().Set("location", "/temp/"+url.PathEscape(stream.ToToken()))
	w.WriteHeader(http.StatusCreated)
}
