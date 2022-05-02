package objects

import (
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
	"strings"
)

/*
PUT /objects/fileName
head: "Digest: SHA-256=<hash>"

将数据直接存储到dataServer,在元数据中记录
*/
func put(w http.ResponseWriter, r *http.Request) {
	//获取hash和size
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
	//将文件使用rs编码方式存储到dataServer
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
	//将文件在元数据中保存
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	e = elasticSearch.AddVersion(name, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Println("apiServer end put")
}

/*
r: 上传的文件内容 hash:文件hash值 size:文件size
*/
func storeObject(r io.Reader, hash string, size int64) (int, error) {
	log.Println("apiServer start storeObject")
	//url.PathEscape确保可以放到url中使用
	//判断是否已经存储
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
	//如果hash错误,删除临时文件并报错
	if d != hash {
		log.Println("客户端hash错误")
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch,"+
			"calculated=%s, requested=%s", d, hash)
	}
	//无误后转正临时文件
	stream.Commit(true)
	log.Println("apiServer : storeObject end!")
	return http.StatusOK, nil
}

//随机选择dataserver,并返回用于传数据的putStream
func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	log.Println("apiServer start putStream")
	servers := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(servers) != rs.ALL_SHARDS {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}
	log.Println("apiServer :putStream end!")
	return rs.NewRSPutStream(servers, hash, size)
}
