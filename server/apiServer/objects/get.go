package objects

import (
	"compress/gzip"
	"distributed_storage_system/server/apiServer/heartbeat"
	"distributed_storage_system/server/apiServer/locate"
	"distributed_storage_system/utils/elasticSearch"
	"distributed_storage_system/utils/headutils"
	"distributed_storage_system/utils/rs"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	//获取url中的fileName和version
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
	//获取该文件和id对应的元数据
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
	//获取取得的元数据中的hash值
	hash := url.PathEscape(meta.Hash)
	//见函数注释
	stream, e := GetStream(hash, meta.Size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//获取请求中要求读取的偏移量
	offset := headutils.GetOffsetFromHeader(r.Header)
	//如果存在偏移量
	if offset != 0 {
		//读取到偏移量位置
		stream.Seek(offset, io.SeekCurrent)
		//在返回信息中记录偏移量数据
		w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", offset, meta.Size-1, meta.Size))
		w.WriteHeader(http.StatusPartialContent)
	}
	//表示是否需要gzip压缩
	acceptGzip := false
	//获取头部信息,判断是否需要gzip压缩
	encoding := r.Header["Accept-Encoding"]
	for i := range encoding {
		if encoding[i] == "gzip" {
			acceptGzip = true
			break
		}
	}
	//io.copy调用stream.read完成,即decoder.read函数
	if acceptGzip {
		//获取gzip压缩的文件
		w.Header().Set("content-encoding", "gzip")
		w2 := gzip.NewWriter(w)
		io.Copy(w2, stream)
		w2.Close()
	} else {
		//获取源文件
		io.Copy(w, stream)
	}
	stream.Close()
}

///返回RSGetStream,RSGetStream内嵌decoder,decoder中reader保存正确的分片文件,writer用于处理不正确的分片
func GetStream(hash string, size int64) (*rs.RSGetStream, error) {
	log.Println("apiServer GetStream start")
	//获取对应hash值的文件所在的dataServer的地址
	locateInfo := locate.Locate(hash)
	//返回节点数小于四,无法修复,报错
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateInfo)
	}
	//用于保存修复数据应该放在的dataServer地址
	dataServers := make([]string, 0)
	//返回节点数不等于6,选择新节点修复
	if len(locateInfo) != rs.ALL_SHARDS {
		//选择新的节点,用于保存修复数据
		dataServers = heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	log.Println("apiServer GetStream end")
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}
